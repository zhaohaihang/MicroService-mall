package mode

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/zhaohaihang/order_api/global"
	"github.com/zhaohaihang/order_api/proto"
)

func DebugMode() {
	target := "127.0.0.1:8000"
	orderConn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	global.OrderClient = proto.NewOrderClient(orderConn)
	zap.S().Infof("RPC debugg模式 服务连接成功")
}
