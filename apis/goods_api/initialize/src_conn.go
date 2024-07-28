package initialize

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // 负载均衡
	"github.com/opentracing/opentracing-go"
	"github.com/zhaohaihang/goods_api/global"
	"github.com/zhaohaihang/goods_api/proto"
	otgrpc "github.com/zhaohaihang/goods_api/utils/otgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitGoodsServiceConn() {
	cfg := api.DefaultConfig()
	consulConfig := global.ApiConfig.ConsulInfo
	cfg.Address = fmt.Sprintf("%s:%d", consulConfig.Host, consulConfig.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Errorw("create consul client failed", "err", err.Error())
		return
	}

	data, err := client.Agent().ServicesWithFilter(`Service == "goods_service"`)
	if err != nil {
		zap.S().Errorw("search  goods-service failed", "err", err.Error())
		return
	}
	zap.S().Info(data)

	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",
			global.ApiConfig.ConsulInfo.Host,
			global.ApiConfig.ConsulInfo.Port,
			global.ApiConfig.GoodsServiceInfo.Name),

		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)

	if err != nil {
		zap.S().Errorw("grpc Dial failed", "err", err.Error())
		return
	}
	global.GoodsClient = proto.NewGoodsClient(goodsConn)
	zap.S().Infow("RPC release mode conn service success")
}
