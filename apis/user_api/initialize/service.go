package initialize

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/opentracing/opentracing-go"
	"github.com/zhaohaihang/user_api/global"
	"github.com/zhaohaihang/user_api/proto"
	"github.com/zhaohaihang/user_api/utils/otgrpc"
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
	userServiceHost := ""
	userServicePort := 0
	for _, value := range data {
		zap.S().Info(value.Address)
		zap.S().Info(value.Port)
		userServiceHost = value.Address
		userServicePort = value.Port
		break
	}
	if userServiceHost == "" || userServicePort == 0 {
		zap.S().Fatal("Init RPC failed")
		return
	}
	zap.S().Infof("find user-service %s:%d", userServiceHost, userServicePort)
	target := fmt.Sprintf("%s:%d", userServiceHost, userServicePort)

	// 连接服务端
	userConn, err := grpc.Dial(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())))
	if err != nil {
		zap.S().Errorw("grpc Dial failed", "err", err.Error())
		return
	}
	global.UserClient = proto.NewUserClient(userConn)
	zap.S().Infof("RPC release mode conn service success")
}
