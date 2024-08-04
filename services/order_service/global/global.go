package global

import (
	"github.com/hashicorp/consul/api"
	"github.com/zhaohaihang/order_service/config"
	"github.com/zhaohaihang/order_service/proto"
	"gorm.io/gorm"
)

var (
	FilePath               *config.FilePathConfig
	NacosConfig            *config.NacosConfig
	ServiceConfig          *config.ServiceConfig
	DB                     *gorm.DB
	Client                 *api.Client
	FreePort               int
	ServiceID              string
	GoodsServiceClient     proto.GoodsClient
	InventoryServiceClient proto.InventoryClient
)

const (
	// 订单状态
	PAYING         = "PAYING"         //待支付,
	TRADE_SUCCESS  = "TRADE_SUCCESS"  //成功,
	TRADE_CLOSED   = "TRADE_CLOSED"   //超时关闭
	WAIT_BUYER_PAY = "WAIT_BUYER_PAY" //(交易创建
	TRADE_FINISHED = "TRADE_FINISHED" //交易结束)
)
