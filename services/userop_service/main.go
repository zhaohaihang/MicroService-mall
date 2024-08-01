package main

import (
	"flag"
	"io"
	"github.com/zhaohaihang/userop_service/handler"
	"github.com/zhaohaihang/userop_service/initialize"
	"github.com/zhaohaihang/userop_service/mode"
	"github.com/zhaohaihang/userop_service/proto"
	"github.com/zhaohaihang/userop_service/util/otgrpc"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// 命令行参数
	IP := flag.String("ip", "0.0.0.0", "ip地址：服务启动ip地址")
	Port := flag.Int("port", 8000, "port端口号：服务启动端口号")
	Mode := flag.String("mode", "release", "mode启动模式：debug 本地调试/release 服务注册")
	flag.Parse()

	initialize.InitFileAbsPath()
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	tracer, closer := initialize.InitTracer()
	defer func(closer io.Closer) {
		err := closer.Close()
		if err != nil {
			panic(err)
		}
	}(closer)
	opentracing.SetGlobalTracer(tracer)
	server := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))
	proto.RegisterMessageServer(server, &handler.UserOpService{})
	proto.RegisterAddressServer(server, &handler.UserOpService{})
	proto.RegisterUserFavoriteServer(server, &handler.UserOpService{})
	if *Mode == "debug" {
		zap.S().Warnf("start debug mode  \n")
		mode.DebugMode(server, *IP, *Port)
	} else if *Mode == "release" {
		zap.S().Warnf("start release mode\n")
		mode.ReleaseMode(server, *IP)
	}
}
