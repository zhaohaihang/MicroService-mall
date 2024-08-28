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
	IP := flag.String("ip", "0.0.0.0", "ip address")
	Port := flag.Int("port", 8000, "port ")
	Mode := flag.String("mode", "release", "mode debug / release ")
	nacosConfig := flag.String("nacosConfig", "home", "home / com")
	flag.Parse()

	initialize.InitFileAbsPath(*nacosConfig)
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
