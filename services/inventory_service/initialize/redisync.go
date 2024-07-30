package initialize

import (
	"fmt"
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"go.uber.org/zap"
	"github.com/zhaohaihang/inventory_service/global"
)

func InitRedis() {
	redisConfig := global.ServiceConfig.RedisInfo
	client := goredislib.NewClient(&goredislib.Options{
		Addr: fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port), 
		Password: redisConfig.Password, PoolSize: redisConfig.PoolSize})
	pool := goredis.NewPool(client)
	global.Redsync = redsync.New(pool)
	zap.S().Infof("redsync init success \n")
}
