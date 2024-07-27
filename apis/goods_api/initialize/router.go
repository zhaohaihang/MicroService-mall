package initialize

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/goods_api/global"
	"github.com/zhaohaihang/goods_api/router"
	"go.uber.org/zap"
)

func InitRouter() {
	Router := gin.Default()
	ApiGroup := Router.Group("/goods/v1")
	router.InitHealthRouter(ApiGroup)
	router.InitBannerRouter(ApiGroup)
	router.InitGoodsRouter(ApiGroup)
	router.InitCategoryRouter(ApiGroup)
	router.InitBrandRouter(ApiGroup)

	zap.S().Info("init router success")

	zap.S().Info("start run goods api router")
	err := Router.Run(fmt.Sprintf("0.0.0.0:%d", global.Port))
	if err != nil {
		zap.S().Errorw("run goods api router failed", "err", err.Error())
	}

}
