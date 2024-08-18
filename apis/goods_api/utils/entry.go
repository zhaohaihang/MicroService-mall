package utils

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
)

func SentinelEntry(ctx *gin.Context) (*base.SentinelEntry, *base.BlockError) {
	entry, blockError := sentinel.Entry("goods_api", sentinel.WithTrafficType(base.Inbound))
	return entry, blockError
}
