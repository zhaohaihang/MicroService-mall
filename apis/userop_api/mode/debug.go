package mode

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/zhaohaihang/userop_api/global"
	"github.com/zhaohaihang/userop_api/proto"
)

func DebugMode() {
	zap.S().Warnf("start debug mode")
	target := "127.0.0.1:8000"
	useropConn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	global.UserFavoriteClient = proto.NewUserFavoriteClient(useropConn)
	global.AddressClient = proto.NewAddressClient(useropConn)
	global.MessageClient = proto.NewMessageClient(useropConn)
	zap.S().Infof("conn success")
}
