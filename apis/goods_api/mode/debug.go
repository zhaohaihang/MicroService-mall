package mode

import (
	"github.com/zhaohaihang/goods_api/global"
	"github.com/zhaohaihang/goods_api/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func DebugMode() {
	zap.S().Warnf("start debug mode")
	target := "127.0.0.1:8000"
	goodsConn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	global.GoodsClient = proto.NewGoodsClient(goodsConn)

	zap.S().Infof("conn success")
}
