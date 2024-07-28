package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis/v9"
	"github.com/zhaohaihang/user_api/config"
	"github.com/zhaohaihang/user_api/proto"
)

var (
	Translator  ut.Translator    // 翻译器
	UserClient  proto.UserClient // grpc客户端
	RedisClient *redis.Client
	FileConfig  *config.FileConfig // 文件配置
	NacosConfig *config.NacosConfig
	ApiConfig   *config.ApiConfig
	Port        int
)
