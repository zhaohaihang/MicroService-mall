package router

import (
	"github.com/zhaohaihang/userop_web/api/health"

	"github.com/gin-gonic/gin"
)

func InitHealthRoute(Router *gin.RouterGroup) {
	Router.GET("/health", health.HealthCheck)
}