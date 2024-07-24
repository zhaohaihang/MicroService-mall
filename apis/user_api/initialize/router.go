package initialize

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/user_api/api"
	"github.com/zhaohaihang/user_api/global"
	"github.com/zhaohaihang/user_api/router"
	"go.uber.org/zap"
)

// InitRouters 初始化路由
func InitRouters()  {
	Router := gin.Default()
	Router.GET("/health", api.HealthCheck)
	ApiGroup := Router.Group("/user/v1")
	router.InitBaseRouter(ApiGroup)
	router.InitUserRouter(ApiGroup)
	zap.S().Infow("路由启动成功")

	zap.S().Infof("start user api")
	err := Router.Run(fmt.Sprintf("0.0.0.0:%d", global.Port))
	if err != nil {
		panic(err)
	}
}
