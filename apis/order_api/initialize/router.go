package initialize

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/order_api/global"
	"github.com/zhaohaihang/order_api/middlewares"
	"github.com/zhaohaihang/order_api/router"
	"go.uber.org/zap"
)

func InitRouters() {
	Router := gin.Default()
	ApiGroup := Router.Group("/order/v1")
	router.InitHealthRoute(ApiGroup)

	ApiGroup.Use(middlewares.Cors())
	router.InitOrderRouter(ApiGroup)
	router.InitShopCartRouter(ApiGroup)
	zap.S().Infow("init order_api router success")

	zap.S().Infof("start order_api")
	err := Router.Run(fmt.Sprintf("0.0.0.0:%d", global.Port))
	if err != nil {
		zap.S().Errorw("run order api router failed", "err", err.Error())
	}

}
