package initialize

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/opentracing/opentracing-go"
	"github.com/zhaohaihang/user_api/global"
	"github.com/zhaohaihang/user_api/proto"
	"github.com/zhaohaihang/user_api/utils/otgrpc"
	_ "github.com/mbobakov/grpc-consul-resolver" // 负载均衡
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitUserService() {
	// 创建conusl 客户端
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ApiConfig.ConsulInfo.Host, global.ApiConfig.ConsulInfo.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Errorw("create consul client failed", "err", err.Error())
		return
	}

	// 获取服务信息
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.ApiConfig.UserServiceInfo.Name))
	if err != nil {
		zap.S().Errorw("search user-service failed", "err", err.Error())
		return
	}
	zap.S().Info(data)

	userConn, err := grpc.Dial(
			fmt.Sprintf("consul://%s:%d/%s?wait=14s", 
			global.ApiConfig.ConsulInfo.Host,
			global.ApiConfig.ConsulInfo.Port, 
			global.ApiConfig.UserServiceInfo.Name),

		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)

	if err != nil {
		zap.S().Errorw("grpc Dial failed", "err", err.Error())
		return
	}
	global.UserClient = proto.NewUserClient(userConn)
	zap.S().Infof("RPC release mode conn service success")
}
