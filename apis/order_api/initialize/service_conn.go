package initialize

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/opentracing/opentracing-go"
	"github.com/zhaohaihang/order_api/global"
	"github.com/zhaohaihang/order_api/proto"
	"github.com/zhaohaihang/order_api/utils/otgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitService() {
	initOrderService()
	initGoodsService()
	initInventoryService()
}

func initOrderService() {
	consulConfig := global.ApiConfig.ConsulInfo
	orderConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",
			consulConfig.Host,
			consulConfig.Port,
			global.ApiConfig.OrderService.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()), 
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		zap.S().Fatalw("order_service conn failed", "err", err)
	}
	zap.S().Infof("order_service conn success")
	global.OrderClient = proto.NewOrderClient(orderConn)
}

func initGoodsService() {
	consulConfig := global.ApiConfig.ConsulInfo
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",
			consulConfig.Host,
			consulConfig.Port,
			global.ApiConfig.GoodsService.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()), 
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		zap.S().Fatalw("goods_service conn failed", "err", err)
	}
	zap.S().Infof("goods_service conn success")
	global.GoodsClient = proto.NewGoodsClient(goodsConn)
}

func initInventoryService() {
	consulConfig := global.ApiConfig.ConsulInfo
	inventoryConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", 
			consulConfig.Host, 
			consulConfig.Port, 
			global.ApiConfig.InventoryService.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()), 
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		zap.S().Fatalw("inventory_service conn failed", "err", err)
	}
	zap.S().Infof("inventory_service conn success")
	global.InventoryClient = proto.NewInventoryClient(inventoryConn)
}
