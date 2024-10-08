package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/user_api/api"
	"github.com/zhaohaihang/user_api/middlewares"
)

func InitUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("user", middlewares.Trace())
	{
		userRouter.GET("list", middlewares.JWTAuth(), middlewares.AdminAuth(), api.GetUserList)
		userRouter.POST("pwd_login", api.PasswordLogin)
		userRouter.POST("register", api.Register)
	}
}
