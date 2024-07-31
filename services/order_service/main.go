package main

import (
	"flag"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/zhaohaihang/order_service/handler"
	"github.com/zhaohaihang/order_service/initialize"
	"github.com/zhaohaihang/order_service/mode"
	"github.com/zhaohaihang/order_service/proto"
	otgrpc "github.com/zhaohaihang/order_service/util/otgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip address")
	Port := flag.Int("port", 8000, "port ")
	Mode := flag.String("mode", "release", "mode debug / release ")
	flag.Parse()

	initialize.InitFileAbsPath()
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	initialize.InitOtherService()

	tracer, closer := initialize.InitTracer()
	defer func(closer io.Closer) {
		err := closer.Close()
		if err != nil {
			panic(err)
		}
	}(closer)
	opentracing.SetGlobalTracer(tracer)
	server := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))
	proto.RegisterOrderServer(server, &handler.OrderService{})
	
	if *Mode == "debug" {
		zap.S().Warnf("start debug mode  \n")
		mode.DebugMode(server, *IP, *Port)
	} else if *Mode == "release" {
		zap.S().Warnf("start release mode\n")
		mode.ReleaseMode(server, *IP)
	}
}
