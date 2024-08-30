package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/zhaohaihang/userop_api/config"
	"github.com/zhaohaihang/userop_api/proto"
)

var (
	FilePath           *config.FilePathConfig
	NacosConfig        *config.NacosConfig
	ApiConfig          *config.ApiConfig
	Translator         ut.Translator
	UserFavoriteClient proto.UserFavoriteClient
	MessageClient      proto.MessageClient
	AddressClient      proto.AddressClient
	GoodsClient        proto.GoodsClient
	Port               int
)
