package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/userop_api/api/address"
	"github.com/zhaohaihang/userop_api/middlewares"
)

func InitAddressRouter(Router *gin.RouterGroup) {
	AddressRouter := Router.Group("address", middlewares.Trace())
	{
		AddressRouter.GET("", middlewares.JWTAuth(), address.ListAddress)
		AddressRouter.DELETE("/:id", middlewares.JWTAuth(), address.DeleteAddress)
		AddressRouter.POST("", middlewares.JWTAuth(), address.CreateAddress)
		AddressRouter.PUT("/:id", middlewares.JWTAuth(), address.UpdateAddress)
	}
}
