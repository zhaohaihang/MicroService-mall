package router

import (
	"github.com/zhaohaihang/goods_api/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/goods_api/api/banner"
)

func InitBannerRouter(Router *gin.RouterGroup) {
	BannerRouter := Router.Group("banners", middlewares.Trace())
	{
		BannerRouter.GET("", banner.List)
		BannerRouter.DELETE("/:id", banner.Delete)
		BannerRouter.POST("", banner.New)
		BannerRouter.PUT("/:id", banner.Update)
	}
}
