package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/zhaohaihang/order_api/config"
	"github.com/zhaohaihang/order_api/proto"
)

var (
	FilePathConfig       *config.FilePathConfig
	NacosConfig      *config.NacosConfig
	ApiConfig 		 *config.ApiConfig
	Translator       ut.Translator
	OrderClient      proto.OrderClient
	GoodsClient      proto.GoodsClient
	InventoryClient  proto.InventoryClient
	Port        int
)
