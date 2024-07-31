package handler
import "github.com/zhaohaihang/order_service/proto"

type OrderService struct {
	proto.UnimplementedOrderServer
}

const(
	SERVICE_NAME = "[order_service]"
)  
