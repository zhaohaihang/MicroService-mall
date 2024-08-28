package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/order_api/api/shop_cart"
	"github.com/zhaohaihang/order_api/middlewares"
)

func InitShopCartRouter(Router *gin.RouterGroup) {
	ShopCartRouter := Router.Group("shopcarts", middlewares.Trace())
	{
		ShopCartRouter.GET("", middlewares.JWTAuth(), shop_cart.ListGoodsInCart)          //购物车列表
		ShopCartRouter.DELETE("/:id", middlewares.JWTAuth(), shop_cart.DeleteGoodsInCart) //删除条目
		ShopCartRouter.POST("", middlewares.JWTAuth(), shop_cart.AddGoodsToCart)          //添加商品到购物车
		ShopCartRouter.PATCH("/:id", middlewares.JWTAuth(), shop_cart.UpdateGoodsInCart)  //修改条目
	}
}
