package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/goods_api/api/brand"
	"github.com/zhaohaihang/goods_api/api/categoryBrand"
	"github.com/zhaohaihang/goods_api/middlewares"
)

func InitBrandRouter(Router *gin.RouterGroup) {
	BrandRouter := Router.Group("brands", middlewares.Trace())
	{
		BrandRouter.GET("", brand.List)
		BrandRouter.DELETE("/:id", brand.Delete)
		BrandRouter.POST("", brand.New)
		BrandRouter.PUT("/:id", brand.Update)
	}
	CategoryBrandRouter := Router.Group("categorybrands")
	{
		CategoryBrandRouter.GET("", categoryBrand.List)
		CategoryBrandRouter.DELETE("/:id", categoryBrand.Delete)
		CategoryBrandRouter.POST("", categoryBrand.New)
		CategoryBrandRouter.PUT("/:id", categoryBrand.Update)
		CategoryBrandRouter.GET("/:id", categoryBrand.Detail)
	}
}
