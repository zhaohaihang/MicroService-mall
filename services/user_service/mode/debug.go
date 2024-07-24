package mode

import (
	"fmt"
	"net"

	"github.com/zhaohaihang/user_service/handler"
	"github.com/zhaohaihang/user_service/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// DebugMode debug开发模式
func DebugMode(server *grpc.Server, ip string, port int) {
	proto.RegisterUserServer(server, &handler.UserService{})
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
