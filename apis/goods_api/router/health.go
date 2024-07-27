package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/goods_api/api/health"
)

func InitHealthRouter(Router *gin.RouterGroup) {
	Router.GET("/health", health.HealthCheck)
}
