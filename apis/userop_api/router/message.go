package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/userop_web/api/message"
	"github.com/zhaohaihang/userop_web/middlewares"
)

func InitMessageRouter(Router *gin.RouterGroup) {
	MessageRouter := Router.Group("message", middlewares.Trace())
	{
		MessageRouter.GET("", middlewares.JWTAuth(), message.List)
		MessageRouter.POST("", middlewares.JWTAuth(), message.New)
	}
}