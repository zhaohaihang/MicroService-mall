package initialize

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/userop_web/global"
	"github.com/zhaohaihang/userop_web/middlewares"
	"github.com/zhaohaihang/userop_web/router"
	"go.uber.org/zap"
)

func InitRouters() {
	Router := gin.Default()
	ApiGroup := Router.Group("/userop/v1")
	router.InitHealthRoute(ApiGroup)

	Router.Use(middlewares.Cors())
	router.InitMessageRouter(ApiGroup)
	router.InitAddressRouter(ApiGroup)
	router.InitUserFavoriteRouter(ApiGroup)
	zap.S().Infow("init userop router success")

	zap.S().Infof("start userop api")
	err := Router.Run(fmt.Sprintf("0.0.0.0:%d", global.Port))
	if err != nil {
		zap.S().Fatalw("start userop api failed", "err", err)
	}
}
