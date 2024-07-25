package handler
import "github.com/zhaohaihang/goods_service/proto"

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

const(
	SERVICE_NAME = "[Goods_Service]"
)  
