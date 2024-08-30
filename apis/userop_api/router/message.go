package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/userop_api/api/message"
	"github.com/zhaohaihang/userop_api/middlewares"
)

func InitMessageRouter(Router *gin.RouterGroup) {
	MessageRouter := Router.Group("message", middlewares.Trace())
	{
		MessageRouter.GET("", middlewares.JWTAuth(), message.ListMessage)
		MessageRouter.POST("", middlewares.JWTAuth(), message.CreateMessage)
	}
}
