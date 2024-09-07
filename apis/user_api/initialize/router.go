package initialize

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/user_api/global"
	"github.com/zhaohaihang/user_api/middlewares"
	"github.com/zhaohaihang/user_api/router"
	"github.com/zhaohaihang/user_api/swagger"
	"go.uber.org/zap"
)

// InitRouters 初始化路由
func InitRouters() {
	Router := gin.New()
	Router.Use(gin.LoggerWithConfig(
		gin.LoggerConfig{
			SkipPaths: []string{"/v1/health"},
		},
	),gin.Recovery())

	// ApiGroup := Router.Group("/user/v1")
	ApiGroup := Router.Group("/v1")
	swagger.InitSwaggarRoute(ApiGroup)

	ApiGroup.Use(middlewares.Cors())
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
