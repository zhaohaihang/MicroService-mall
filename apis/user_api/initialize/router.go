package initialize

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/user_api/global"
	"github.com/zhaohaihang/user_api/router"
	"go.uber.org/zap"
)

// InitRouters 初始化路由
func InitRouters() {
	Router := gin.New()
	Router.Use(gin.LoggerWithConfig(
		gin.LoggerConfig{
			SkipPaths: []string{"/user/v1/health"},
		},
	),gin.Recovery())

	ApiGroup := Router.Group("/user/v1")
	router.InitHealthRoute(ApiGroup)
	router.InitBaseRouter(ApiGroup)
	router.InitUserRouter(ApiGroup)
	zap.S().Infow("init user_api router success")

	zap.S().Infof("start user api")
	err := Router.Run(fmt.Sprintf("0.0.0.0:%d", global.Port))
	if err != nil {
		panic(err)
	}
}
