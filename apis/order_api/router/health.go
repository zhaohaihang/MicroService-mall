package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/order_api/api/health"
)

func InitHealthRoute(Router *gin.RouterGroup) {
	Router.GET("/health", health.HealthCheck)
}
