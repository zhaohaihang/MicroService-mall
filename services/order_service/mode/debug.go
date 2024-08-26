package mode

import (
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func DebugMode(server *grpc.Server, ip string, port int) {

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		zap.S().Errorw("net.Listen failed", "err", err.Error())
		return
	}
	zap.S().Infof("service start success ,port %s:%d", ip, port)
	err = server.Serve(listen)
	if err != nil {
		zap.S().Errorw("server.Serve failed", "err", err.Error())
		return
	}

}
