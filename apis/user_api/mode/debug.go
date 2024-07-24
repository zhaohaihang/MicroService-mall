package mode

import (
	"github.com/opentracing/opentracing-go"
	"github.com/zhaohaihang/user_api/global"
	"github.com/zhaohaihang/user_api/initialize"
	"github.com/zhaohaihang/user_api/proto"
	"github.com/zhaohaihang/user_api/utils/otgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func DebugMode() {
	target := "127.0.0.1:8000"
	global.ApiConfig.Host = "127.0.0.1"
	global.ApiConfig.Port = 8000
	userConn, err := grpc.Dial(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())))
	if err != nil {
		panic(err)
	}
	global.UserClient = proto.NewUserClient(userConn)
	zap.S().Infof("debugg mode conn grpc server success")

	initialize.InitUserService()
}
