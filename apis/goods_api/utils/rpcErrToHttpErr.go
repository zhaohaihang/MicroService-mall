package utils

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

// HandleGrpcErrorToHttpError
func HandleGrpcErrorToHttpError(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": "can not find",
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "internal server error",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "Bad Request",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "other internal server error",
				})
			}
		}
	}
}
