package main

import (
	"flag"

	"github.com/opentracing/opentracing-go"
	"github.com/zhaohaihang/goods_service/handler"
	"github.com/zhaohaihang/goods_service/initialize"
	"github.com/zhaohaihang/goods_service/mode"
	"github.com/zhaohaihang/goods_service/proto"
	"github.com/zhaohaihang/goods_service/utils/otgrpc"
	"go.uber.org/zap"

	"io"

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
	initialize.InitEs()
	tracer, closer := initialize.InitTracer()
	defer func(closer io.Closer) {
		err := closer.Close()
		if err != nil {
			panic(err)
		}
	}(closer)

	opentracing.SetGlobalTracer(tracer)
	server := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))
	proto.RegisterGoodsServer(server, &handler.GoodsServer{})

	if *Mode == "debug" {
		zap.S().Warnf("start debug mode  \n")
		mode.DebugMode(server, *IP, *Port)
	} else if *Mode == "release" {
		zap.S().Warnf("start release mode\n")
		mode.ReleaseMode(server, *IP)
	}
}
