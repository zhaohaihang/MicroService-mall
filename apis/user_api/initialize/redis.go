package initialize

import (
	"fmt"

	"github.com/go-redis/redis/v9"
	"github.com/zhaohaihang/user_api/global"
	"go.uber.org/zap"
)

func InitRedis() {
	addr := fmt.Sprintf("%s:%d", global.ApiConfig.RedisInfo.Host, global.ApiConfig.RedisInfo.Port)
	global.RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: global.ApiConfig.RedisInfo.Password,
	})

	zap.S().Infof("redis client init success \n")
}
