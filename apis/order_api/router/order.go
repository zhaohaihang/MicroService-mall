package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/order_api/api/order"
	"github.com/zhaohaihang/order_api/api/pay"
	"github.com/zhaohaihang/order_api/middlewares"
)

func InitOrderRouter(Router *gin.RouterGroup) {
	OrderRouter := Router.Group("orders", middlewares.Trace())
	{
		OrderRouter.GET("", middlewares.JWTAuth(), order.List)
		OrderRouter.POST("", middlewares.JWTAuth(), order.New)
		OrderRouter.GET("/:id", middlewares.JWTAuth(), order.Detail)
	}
	PayRouter := Router.Group("pay")
	{
		PayRouter.POST("alipy/notify", pay.Notify)
	}
}
