package handler

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/zhaohaihang/order_service/global"
	"github.com/zhaohaihang/order_service/model"
	"github.com/zhaohaihang/order_service/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CartItemList 获取用户购物车列表
func (s *OrderService) CartItemList(ctx context.Context, request *proto.UserInfo) (*proto.CartItemListResponse, error) {
	zap.S().Infow("Info", "service", serviceName, "method", "CartItemList", "request", request)

	parentSpan := opentracing.SpanFromContext(ctx)
	cartItemListSpan := opentracing.GlobalTracer().StartSpan("CartItemList", opentracing.ChildOf(parentSpan.Context()))
	defer cartItemListSpan.Finish()

	response := &proto.CartItemListResponse{}
	// 根据UserId 查询购物车
	var shopCarts []model.ShoppingCart
	result := global.DB.Where(&model.ShoppingCart{UserId: request.Id}).Find(&shopCarts)
	if result.Error != nil {
		zap.S().Errorw("CartItemList failed", "err", result.Error)
		return nil, result.Error
	}
	response.Total = int32(result.RowsAffected)
	// 组装返回数据
	for _, shopCart := range shopCarts {
		response.Data = append(response.Data, &proto.ShopCartInfoResponse{
			Id:      int32(shopCart.Model.ID),
			UserId:  shopCart.UserId,
			GoodsId: shopCart.GoodsId,
			Nums:    shopCart.Nums,
			Checked: shopCart.Checked,
		})
	}

	return response, nil
}

// CreateCartItem 商品加入购物车
func (s *OrderService) CreateCartItem(ctx context.Context, request *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	zap.S().Infow("Info", "service", serviceName, "method", "CreateCartItem", "request", request)

	parentSpan := opentracing.SpanFromContext(ctx)
	createCartItemSpan := opentracing.GlobalTracer().StartSpan("CreateCartItem", opentracing.ChildOf(parentSpan.Context()))
	defer createCartItemSpan.Finish()

	response := &proto.ShopCartInfoResponse{}

	var shopCart model.ShoppingCart
	result := global.DB.Where(&model.ShoppingCart{GoodsId: request.GoodsId, UserId: request.UserId}).First(&shopCart)
	if result.RowsAffected == 1 {
		shopCart.Nums += request.Nums // 如果购物车中已经有，则直接更新数量
	} else {
		shopCart.UserId = request.UserId
		shopCart.GoodsId = request.GoodsId
		shopCart.Nums = request.Nums
		shopCart.Checked = false
	}
	result = global.DB.Save(&shopCart)
	if result.Error != nil {
		zap.S().Errorw("CreateCartItem save failed", "err", result.Error)
		return nil, result.Error
	}
	response.Id = int32(shopCart.ID)

	return response, nil
}

// UpdateCartItem 更新购物车 是否选择/数量改变
func (s *OrderService) UpdateCartItem(ctx context.Context, request *proto.CartItemRequest) (*emptypb.Empty, error) {
	zap.S().Infow("Info", "service", serviceName, "method", "UpdateCartItem", "request", request)

	parentSpan := opentracing.SpanFromContext(ctx)
	updateCartItemSpan := opentracing.GlobalTracer().StartSpan("UpdateCartItem", opentracing.ChildOf(parentSpan.Context()))
	defer updateCartItemSpan.Finish()

	var shopCart model.ShoppingCart
	result := global.DB.Where(&model.ShoppingCart{GoodsId: request.GoodsId, UserId: request.UserId}).First(&shopCart)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "ShoppingCart is not exist")
	}
	shopCart.Checked = request.Checked
	if request.Nums > 0 {
		shopCart.Nums = request.Nums
	}
	result = global.DB.Save(&shopCart)
	if result.Error != nil {
		zap.S().Errorw("UpdateCartItem save failed", "err", result.Error)
		return nil, result.Error
	}

	return &emptypb.Empty{}, nil
}

// DeleteCartItem  删除购物车记录
func (s *OrderService) DeleteCartItem(ctx context.Context, request *proto.CartItemRequest) (*emptypb.Empty, error) {
	zap.S().Infow("Info", "service", serviceName, "method", "DeleteCartItem", "request", request)

	parentSpan := opentracing.SpanFromContext(ctx)
	deleteCartItemSpan := opentracing.GlobalTracer().StartSpan("DeleteCartItem", opentracing.ChildOf(parentSpan.Context()))
	defer deleteCartItemSpan.Finish()

	var shopCart model.ShoppingCart
	result := global.DB.Where(&model.ShoppingCart{UserId: request.UserId, GoodsId: request.GoodsId}).Delete(&shopCart)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "ShoppingCart is not exist")
	}

	return &emptypb.Empty{}, nil
}
