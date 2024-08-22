package handler

import (
	"context"
	"fmt"

	"github.com/go-redsync/redsync/v4"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/opentracing/opentracing-go"
	"github.com/zhaohaihang/inventory_service/global"
	"github.com/zhaohaihang/inventory_service/model"
	"github.com/zhaohaihang/inventory_service/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"time"
)

type InventoryService struct {
	proto.UnimplementedInventoryServer
}

// 设置商品库存
func (i *InventoryService) SetInv(ctx context.Context, request *proto.GoodsInvInfo) (*empty.Empty, error) {
	zap.S().Infow("Info", "service", serviceName, "method", "SetInv", "request", request)
	parentSpan := opentracing.SpanFromContext(ctx)
	setInventorySpan := opentracing.GlobalTracer().StartSpan("setInventory", opentracing.ChildOf(parentSpan.Context()))
	defer setInventorySpan.Finish()

	var inventory model.Inventory
	global.DB.Where(&model.Inventory{GoodsId: request.GoodsId}).First(&inventory)
	inventory.GoodsId = request.GoodsId
	inventory.Stocks = request.Num
	global.DB.Save(&inventory)

	return &empty.Empty{}, nil
}

// 查询商品的库存
func (i InventoryService) InvDetail(ctx context.Context, request *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	zap.S().Infow("Info", "service", serviceName, "method", "InvDetail", "request", request)
	
	parentSpan := opentracing.SpanFromContext(ctx)
	inventoryDetailSpan := opentracing.GlobalTracer().StartSpan("inventoryDetailSpan", opentracing.ChildOf(parentSpan.Context()))
	defer inventoryDetailSpan.Finish()

	var inventory model.Inventory
	result := global.DB.Where(&model.Inventory{
		GoodsId: request.GoodsId,
	}).First(&inventory)
	if result.RowsAffected == 0 {
		zap.S().Errorw("global.DB.First result = 0", "err", "goods Inventory info not exists")
		return nil, status.Errorf(codes.NotFound, "goods Inventory info not exists")
	}
	if result.Error != nil {
		zap.S().Errorw("global.DB.First result error", "err", result.Error)
		return nil, status.Errorf(codes.Internal, "database error")
	}
	response := &proto.GoodsInvInfo{
		Num: inventory.Stocks,
		GoodsId: inventory.GoodsId,
	}

	return response, nil
}

func (i *InventoryService) Sell(ctx context.Context, request *proto.SellInfo) (*empty.Empty, error) {
	zap.S().Infow("Info", "service", serviceName, "method", "Sell", "request", request)
	parentSpan := opentracing.SpanFromContext(ctx)
	sellSpan := opentracing.GlobalTracer().StartSpan("sell", opentracing.ChildOf(parentSpan.Context()))
	tx := global.DB.Begin()

	var details []model.GoodsDetail

	for _, goodInfo := range request.GoodsInfo {

		details = append(details, model.GoodsDetail{
			GoodsId: goodInfo.GoodsId,
			Num:     goodInfo.Num,
		})

		// 加锁
		mutex := global.Redsync.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId), redsync.WithTries(100), redsync.WithExpiry(time.Second*20))
		err := mutex.Lock()
		if err != nil {
			zap.S().Errorw("redisync锁错误", "goods_id", goodInfo.GoodsId, "err", err)
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}

		var inventory model.Inventory
		result := global.DB.Where(&model.Inventory{
			GoodsId: goodInfo.GoodsId,
		}).First(&inventory)
		if result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品库存信息不存在")
		}
		if inventory.Stocks < goodInfo.Num {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.ResourceExhausted, "商品库存不足")
		}
		// 扣减
		inventory.Stocks -= goodInfo.Num
		result = tx.Save(&inventory)
		if result.Error != nil {
			return nil, status.Errorf(codes.Internal, "更新库存失败")
		}

		// 解锁
		ok, err := mutex.Unlock()
		if !ok || err != nil {
			zap.S().Errorw("redisync解锁失败", "goods_id", goodInfo.GoodsId, "err", err.Error())
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}
	}

	sellDetail := model.StockSellDetail{
		OrderSn: request.OrderSn,
		Status:  1,
		Details: details,
	}
	if result := tx.Create(&sellDetail); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "保存库存扣减历史失败")
	}

	tx.Commit()
	sellSpan.Finish()
	return &empty.Empty{}, nil
}

func (i *InventoryService) ReBack(ctx context.Context, request *proto.SellInfo) (*empty.Empty, error) {
	zap.S().Infow("Info", "service", serviceName, "method", "ReBack", "request", request)
	parentSpan := opentracing.SpanFromContext(ctx)
	rebackSpan := opentracing.GlobalTracer().StartSpan("reback", opentracing.ChildOf(parentSpan.Context()))

	tx := global.DB
	for _, goodsInvInfo := range request.GoodsInfo {
		var inventory model.Inventory
		result := global.DB.Where(&model.Inventory{
			GoodsId: goodsInvInfo.GoodsId,
		}).First(&inventory)
		if result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.NotFound, "商品库存信息不存在")
		}
		inventory.Stocks += goodsInvInfo.Num
		tx.Save(&inventory)
	}
	tx.Commit()
	rebackSpan.Finish()
	return &empty.Empty{}, nil
}
