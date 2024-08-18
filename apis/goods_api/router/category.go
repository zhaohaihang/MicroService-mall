package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/goods_api/api/category"
	"github.com/zhaohaihang/goods_api/middlewares"
)

func InitCategoryRouter(Router *gin.RouterGroup) {
	CategoryRouter := Router.Group("categorys", middlewares.Trace())
	{
		CategoryRouter.GET("", category.List)
		CategoryRouter.DELETE("/:id", category.Delete)
		CategoryRouter.GET("/:id", category.Detail)
		CategoryRouter.POST("", category.New)
		CategoryRouter.PUT("/:id", category.Update)
	}
}
