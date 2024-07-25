package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/user_api/api"
)

func InitHealthRoute(Router *gin.RouterGroup) {
	Router.GET("/health", api.HealthCheck)
}
