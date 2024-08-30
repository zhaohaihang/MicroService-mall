package initialize

import (
	"fmt"

	"github.com/opentracing/opentracing-go"
	"github.com/zhaohaihang/userop_api/global"
	"github.com/zhaohaihang/userop_api/proto"
	"github.com/zhaohaihang/userop_api/utils/otgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitService() {
	initUseropService()
	initGoodsService()
}

func initUseropService() {
	consulConfig := global.ApiConfig.ConsulInfo
	useropConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",
			consulConfig.Host,
			consulConfig.Port,
			global.ApiConfig.UseropService.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		zap.S().Fatalw("userop_service conn failed", "err", err)
	}
	zap.S().Infof("userop_service conn success")

	global.UserFavoriteClient = proto.NewUserFavoriteClient(useropConn)
	global.AddressClient = proto.NewAddressClient(useropConn)
	global.MessageClient = proto.NewMessageClient(useropConn)
}

func initGoodsService() {
	consulConfig := global.ApiConfig.ConsulInfo
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",
			consulConfig.Host,
			consulConfig.Port,
			global.ApiConfig.GoodsService.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		zap.S().Fatalw("goods_service conn failed", "err", err)
	}
	zap.S().Infof("goods_service conn success")
	global.GoodsClient = proto.NewGoodsClient(goodsConn)
}
