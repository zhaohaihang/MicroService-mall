package router

import (
	"github.com/zhaohaihang/user_api/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/user_api/api"
)

func InitBaseRouter(Router *gin.RouterGroup) {
	
	BaseRouter := Router.Group("base", middlewares.Trace())
	{
		BaseRouter.GET("/captcha", api.GetCaptcha)
		BaseRouter.POST("/note_code", api.SendNoteCode)
	}

}
