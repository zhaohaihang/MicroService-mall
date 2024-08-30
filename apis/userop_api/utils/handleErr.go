package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"github.com/zhaohaihang/userop_api/global"
)

func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": RemoveTopStruct(errs.Translate(global.Translator)),
	})
}

func HandleRequestFrequentError(c *gin.Context) {
	c.JSON(http.StatusTooManyRequests, gin.H{
		"msg": "Request too frequent",
	})
}