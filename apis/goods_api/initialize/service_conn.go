package initialize

import (
	"fmt"

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
