package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/zhaohaihang/goods_api/config"
	"github.com/zhaohaihang/goods_api/proto"
)

var (
	Port             int
	FilePath         *config.FilePathConfig
	NacosConfig      *config.NacosConfig
	ApiConfig *config.ApiConfig
	Translator       ut.Translator
	GoodsClient      proto.GoodsClient
)
