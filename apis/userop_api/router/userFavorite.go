package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/userop_web/api/userFavorite"
	"github.com/zhaohaihang/userop_web/middlewares"
)

func InitUserFavoriteRouter(Router *gin.RouterGroup) {
	userFavoriteRouter := Router.Group("userfavs", middlewares.Trace())
	{
		userFavoriteRouter.GET("", middlewares.JWTAuth(), userFavorite.ListFavorite)
		userFavoriteRouter.POST("", middlewares.JWTAuth(), userFavorite.CreateFavorite)
		userFavoriteRouter.GET("/:id", middlewares.JWTAuth(), userFavorite.DetailFavorite)
		userFavoriteRouter.DELETE("/:id", middlewares.JWTAuth(), userFavorite.DeleteFavorite)
	}
}
