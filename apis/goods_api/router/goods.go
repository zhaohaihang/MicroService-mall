package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/goods_api/api/goods"
	"github.com/zhaohaihang/goods_api/middlewares"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("goods", middlewares.Trace())
	{
		GoodsRouter.GET("", goods.List)
		GoodsRouter.POST("", goods.New)
		GoodsRouter.GET("/:id", goods.Detail)
		GoodsRouter.DELETE("/:id", goods.Delete)
		GoodsRouter.PUT("/:id", goods.Update)
		GoodsRouter.PATCH("/:id", goods.UpdateStatus)
	}
}
