package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"github.com/zhaohaihang/order_service/global"
	"github.com/zhaohaihang/order_service/model"
	"github.com/zhaohaihang/order_service/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

var serviceName = "[order_service]"

func (s *OrderService) OrderList(ctx context.Context, request *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	zap.S().Infow("Info", "service", serviceName, "method", "OrderList", "request", request)
	parentSpan := opentracing.SpanFromContext(ctx)
	orderListSpan := opentracing.GlobalTracer().StartSpan("OrderList", opentracing.ChildOf(parentSpan.Context()))
	response := &proto.OrderListResponse{}
	var orders []model.OrderInfo
	var total int64
	global.DB.Where(&model.OrderInfo{UserId: request.UserId}).Count(&total)
	response.Total = int32(total)
	result := global.DB.Scopes(model.Paginate(int(request.Pages), int(request.PagePerNums))).Find(&orders)
	if result.RowsAffected == 0 {
		zap.S().Warnw("[orderList] 订单分页查询失败")
		return nil, status.Errorf(codes.NotFound, "未查询到订单")
	}
	for _, order := range orders {
		response.Data = append(response.Data, &proto.OrderInfoResponse{
			Id:      int32(order.Model.ID),
			UserId:  order.UserId,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status:  order.Status,
			Post:    order.Post,
			Total:   order.OrderMount,
			Address: order.Address,
			Name:    order.SignerName,
			Mobile:  order.SingerMobile,
		})
	}
	orderListSpan.Finish()
	return response, nil
}

// OrderDetail  获取订单详情
func (s *OrderService) OrderDetail(ctx context.Context, request *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	zap.S().Infow("Info", "service", serviceName, "method", "OrderDetail", "request", request)
	parentSpan := opentracing.SpanFromContext(ctx)
	orderDetailSpan := opentracing.GlobalTracer().StartSpan("OrderDetail", opentracing.ChildOf(parentSpan.Context()))

	response := &proto.OrderInfoDetailResponse{}
	var order model.OrderInfo

	result := global.DB.Where(&model.OrderInfo{Model: gorm.Model{ID: order.Model.ID}, UserId: request.UserId}).First(&order)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	orderInfo := proto.OrderInfoResponse{}
	orderInfo.Id = int32(order.Model.ID)
	orderInfo.UserId = order.UserId
	orderInfo.OrderSn = order.OrderSn
	orderInfo.PayType = order.PayType
	orderInfo.Status = order.Status
	orderInfo.Post = order.Post
	orderInfo.Total = order.OrderMount
	orderInfo.Address = order.Address
	orderInfo.Name = order.SignerName
	orderInfo.Mobile = order.SingerMobile

	response.OrderInfo = &orderInfo

	var orderGoods []model.OrderGoods
	result = global.DB.Where(&model.OrderGoods{OrderId: int32(order.Model.ID)}).Find(&orderGoods)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, orderGood := range orderGoods {
		response.Goods = append(response.Goods, &proto.OrderItemResponse{
			GoodsId:    orderGood.GoodsId,
			GoodsName:  orderGood.GoodsName,
			GoodsImage: orderGood.GoodsImage,
			Nums:       orderGood.Nums,
		})
	}
	orderDetailSpan.Finish()
	return response, nil
}

/*
func (s *OrderService) CreateOrder(ctx context.Context, request *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	zap.S().Infow("Info", "service", serviceName, "method", "CreateOrder", "request", request)
	parentSpan := opentracing.SpanFromContext(ctx)

	response := &proto.OrderInfoResponse{}

	var goodsId []int32
	var shopCarts []model.ShoppingCart
	goodsNumsMap := make(map[int32]int32)
	// 获取用户购物车中 已选中的商品
	result := global.DB.Where(&model.ShoppingCart{User: request.UserId, Checked: true}).Find(&shopCarts)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "没有选中的结算商品")
	}
	fmt.Printf("购物车选中商品%v", shopCarts)

	for _, shopCart := range shopCarts {
		goodsId = append(goodsId, shopCart.Goods)
		goodsNumsMap[shopCart.Goods] = shopCart.Nums
	}
	fmt.Printf("goodsId %#v", goodsId)
	// 调用商品微服务 查询商品信息
	goodsServiceSpan := opentracing.GlobalTracer().StartSpan("goodsService", opentracing.ChildOf(parentSpan.Context()))
	goods, err := global.GoodsServiceClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{Id: goodsId})
	if err != nil {
		zap.S().Errorw("[goods_service]服务批量查询商品失败", "err", err)
		return nil, status.Errorf(codes.Internal, "批量查询商品信息失败")
	}
	var orderAmount float32
	var orderGoods []*model.OrderGoods
	var goodsInvInfo []*proto.GoodsInvInfo
	for _, goods := range goods.Data {
		orderAmount += goods.ShopPrice * float32(goodsNumsMap[goods.Id])
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      goods.Id,
			GoodsName:  goods.Name,
			GoodsImage: goods.GoodsFrontImage,
			GoodsPrice: goods.ShopPrice,
			Nums:       goodsNumsMap[goods.Id],
		})
		goodsInvInfo = append(goodsInvInfo, &proto.GoodsInvInfo{GoodsId: goods.Id, Num: goodsNumsMap[goods.Id]})
	}
	goodsServiceSpan.Finish()
	// 调用库存服务 扣减库存
	inventoryServiceSpan := opentracing.GlobalTracer().StartSpan("inventoryService", opentracing.ChildOf(parentSpan.Context()))
	_, err = global.InventoryServiceClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInfo: goodsInvInfo,
	})
	if err != nil {
		return nil, status.Errorf(codes.ResourceExhausted, "库存服务扣减失败")
	}
	// 生成订单表
	tx := global.DB.Begin()
	order := model.OrderInfo{
		OrderSn:      GenerateOrderSn(request.UserId),
		OrderMount:   orderAmount,
		Address:      request.Address,
		SignerName:   request.Name,
		SingerMobile: request.Mobile,
		Post:         request.Post,
	}
	result = tx.Save(&order)
	if result.Error != nil {
		zap.S().Errorw("保存订单失败", "err", result.Error)
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "创建订单失败")

	}
	for _, orderGood := range orderGoods {
		orderGood.Order = int32(order.Model.ID)
	}

	// 删除购物车记录
	result = tx.Where(&model.ShoppingCart{User: request.UserId, Checked: true}).Delete(&model.ShoppingCart{})
	if result.RowsAffected == 0 {
		zap.S().Errorw("删除购物车记录失败", "err", result.Error)
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "创建订单失败")
	}
	tx.Commit()
	response.Id = int32(order.Model.ID)
	response.OrderSn = order.OrderSn
	response.Total = order.OrderMount
	inventoryServiceSpan.Finish()
	return response, nil
}

*/

func (s *OrderService) UpdateOrderStatus(ctx context.Context, request *proto.OrderStatus) (*emptypb.Empty, error) {
	zap.S().Infow("Info", "service", serviceName, "method", "UpdateOrderStatus", "request", request)
	parentSpan := opentracing.SpanFromContext(ctx)
	updateOrderStatusSpan := opentracing.GlobalTracer().StartSpan("UpdateOrderStatus", opentracing.ChildOf(parentSpan.Context()))

	result := global.DB.Model(&model.OrderInfo{}).Where(&model.OrderInfo{OrderSn: request.OrderSn}).Update("status", request.Status)
	if result.RowsAffected == 0 || result.Error != nil {
		return nil, status.Errorf(codes.Internal, "更新订单状态失败")
	}
	updateOrderStatusSpan.Finish()
	return &emptypb.Empty{}, nil
}

func (s *OrderService) CreateOrder(ctx context.Context, request *proto.OrderRequest) (*proto.OrderInfoResponse, error){
	orderListener := OrderListener{Ctx :ctx}

	p,err := rocketmq.NewTransactionProducer(
		&orderListener,
		producer.WithNameServer([]string{
			fmt.Sprintf("%s:%d",
			global.ServiceConfig.RocketMQInfo.Host,
			global.ServiceConfig.RocketMQInfo.Port)}),
	)
	if err != nil {
		zap.S().Errorw("create transaction producer failed", "err", err)
		return nil, err
	}
	if err = p.Start(); err != nil {
		zap.S().Errorw("transaction producer start failed", "err", err)
		return nil, err
	}

	order := model.OrderInfo{
		OrderSn:      GenerateOrderSn(request.UserId),
		Address:      request.Address,
		SignerName:   request.Name,
		SingerMobile: request.Mobile,
		Post:         request.Post,
		UserId:         request.UserId,
	}

	orderJson,_ := json.Marshal(order)
	_,err  = p.SendMessageInTransaction(context.Background(), primitive.NewMessage("order_reback", orderJson))
	if err != nil {
		return nil, status.Error(codes.Internal, "send meg failed")
	}

	if orderListener.Code != codes.OK {
		return nil, status.Error(orderListener.Code, orderListener.Detail)
	}

	return &proto.OrderInfoResponse{Id: orderListener.ID, OrderSn: order.OrderSn, Total: orderListener.OrderAmount}, nil
}

type OrderListener struct {
	Code        codes.Code
	Detail      string
	ID          int32
	OrderAmount float32
	Ctx         context.Context
}

func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	parentSpan := opentracing.SpanFromContext(o.Ctx)

	var orderInfo model.OrderInfo
	json.Unmarshal(msg.Body,&orderInfo)

	// 获取购物车里被选中的商品
	shopCartSpan := opentracing.GlobalTracer().StartSpan("select_shopcart", opentracing.ChildOf(parentSpan.Context()))
	var shopCarts []model.ShoppingCart
	if result := global.DB.Where(&model.ShoppingCart{UserId: orderInfo.UserId,Checked: true}).Find(&shopCarts);result.RowsAffected == 0 {
		o.Code = codes.InvalidArgument
		o.Detail = "no checked goods in shopping cart"
		zap.S().Errorw("no checked goods in shopping cart")
		return primitive.RollbackMessageState
	}
	shopCartSpan.Finish()

	var goodsIds []int32
	goodsNumMap := make(map[int32]int32)
	for _,shopCart := range shopCarts {
		goodsIds = append(goodsIds, shopCart.GoodsId)
		goodsNumMap[shopCart.GoodsId] = shopCart.Nums
	}

	// 查询被选中的商品信息
	queryGoodsSpan := opentracing.GlobalTracer().StartSpan("query_goods", opentracing.ChildOf(parentSpan.Context()))
	goods,err := global.GoodsServiceClient.BatchGetGoods(context.Background(),&proto.BatchGoodsIdInfo{Id: goodsIds})
	if err != nil {
		o.Code = codes.Internal
		o.Detail = " batch get goods info from goods_service failed"
		zap.S().Errorw("goods_service get goods info failed", "err", err)
		return primitive.RollbackMessageState
	}
	queryGoodsSpan.Finish()

	// 计算总金额,购买数量
	orderAmount := float32(0)
	var orderGoods []*model.OrderGoods
	var goodsInvInfo []*proto.GoodsInvInfo
	for _,goods := range goods.Data {
		orderAmount += goods.ShopPrice * float32(goodsNumMap[goods.Id])
		orderGoods = append(orderGoods, &model.OrderGoods{
			GoodsId:    goods.Id,
			GoodsName:  goods.Name,
			GoodsImage: goods.GoodsFrontImage,
			GoodsPrice: goods.ShopPrice,
			Nums:       goodsNumMap[goods.Id],
		})
		goodsInvInfo = append(goodsInvInfo, &proto.GoodsInvInfo{
			GoodsId: goods.Id,
			Num:     goodsNumMap[goods.Id],
		})
	}

	// 远程调用，更新库存
	queryInvSpan := opentracing.GlobalTracer().StartSpan("query_inv", opentracing.ChildOf(parentSpan.Context()))
	if _, err = global.InventoryServiceClient.Sell(context.Background(), &proto.SellInfo{OrderSn: orderInfo.OrderSn, GoodsInfo: goodsInvInfo}); err != nil {
		o.Code = codes.ResourceExhausted
		o.Detail = "subject inventory failed"
		zap.S().Errorw("inventory_serviceget subject goods inventory failed", "err", err)
		return primitive.RollbackMessageState
	}
	queryInvSpan.Finish()

	tx := global.DB.Begin()
	// 创建订单
	saveOrderSpan := opentracing.GlobalTracer().StartSpan("save_order", opentracing.ChildOf(parentSpan.Context()))
	orderInfo.OrderMount = orderAmount
	if result := tx.Save(&orderInfo); result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "create orderinfo failed"
		zap.S().Errorw("create orderinfo failed")
		return primitive.CommitMessageState
	}
	saveOrderSpan.Finish()

	// 向orderlistener 回填数据
	o.OrderAmount = orderAmount
	o.ID = int32(orderInfo.ID)

	//填充orderGoods 的订单号，并插入
	saveOrderGoodsSpan := opentracing.GlobalTracer().StartSpan("save_order_goods", opentracing.ChildOf(parentSpan.Context()))
	for _, orderGood := range orderGoods {
		orderGood.OrderId = int32(orderInfo.ID)
	}
	if result := tx.CreateInBatches(orderGoods, 100); result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "create orderGoods failed"
		zap.S().Errorw("create orderGoods failed")
		return primitive.CommitMessageState
	}
	saveOrderGoodsSpan.Finish()

	// 删除购物车已经下单的商品
	deleteShopCartSpan := opentracing.GlobalTracer().StartSpan("delete_shopcart", opentracing.ChildOf(parentSpan.Context()))
	if result := tx.Where(&model.ShoppingCart{UserId: orderInfo.UserId, Checked: true}).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "delete check goods in shoppingcart failed"
		zap.S().Errorw("delete check goods in shoppingcart failed")
		return primitive.CommitMessageState
	}
	deleteShopCartSpan.Finish()

	//发送延时消息
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{
		fmt.Sprintf("%s:%d",
		global.ServiceConfig.RocketMQInfo.Host,
		global.ServiceConfig.RocketMQInfo.Port)}),)
	if err != nil {
		zap.S().Fatalw("gen name producer failed")
	}
	if err = p.Start(); err != nil {
		panic("启动producer失败")
	}

	msg = primitive.NewMessage("order_timeout", msg.Body)
	msg.WithDelayTimeLevel(3)
	_, err = p.SendSync(context.Background(), msg)
	if err != nil {
		zap.S().Errorf("发送延时消息失败: %v\n", err)
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "发送延时消息失败"
		return primitive.CommitMessageState
	}

	tx.Commit()
	o.Code = codes.OK
	return primitive.CommitMessageState
}

func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	var orderInfo model.OrderInfo
	_ = json.Unmarshal(msg.Body, &orderInfo)
	if result := global.DB.Where(model.OrderInfo{OrderSn: orderInfo.OrderSn}).First(&orderInfo); result.RowsAffected == 0 {
		return primitive.CommitMessageState 
	}

	return primitive.RollbackMessageState
}
