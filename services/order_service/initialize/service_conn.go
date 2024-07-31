package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/zhaohaihang/order_service/global"
	"github.com/zhaohaihang/order_service/proto"
)

func InitOtherService() {
	initGoodsService()
	initInventoryService()
}

func initGoodsService() {
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",
			global.ServiceConfig.ConsulInfo.Host, 
			global.ServiceConfig.ConsulInfo.Port, 
			global.ServiceConfig.GoodsServiceInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		zap.S().Fatalw("goods_service conn failed", "err", err)
	}
	zap.S().Infof("goods_service conn success")
	global.GoodsServiceClient = proto.NewGoodsClient(goodsConn)
}

func initInventoryService() {
	inventoryConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",
			global.ServiceConfig.ConsulInfo.Host, 
			global.ServiceConfig.ConsulInfo.Port, 
			global.ServiceConfig .InventoryServiceInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		zap.S().Fatalw("inventory_service conn failed", "err", err)
	}
	zap.S().Infof("inventory_service conn success")
	global.InventoryServiceClient = proto.NewInventoryClient(inventoryConn)
}
