package global

import (
	"github.com/hashicorp/consul/api"
	"gorm.io/gorm"
	"github.com/zhaohaihang/order_service/config"
	"github.com/zhaohaihang/order_service/proto"
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
