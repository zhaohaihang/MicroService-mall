package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
	"github.com/zhaohaihang/user_api/global"
)

// HandleValidatorError
func HandleValidatorError(c *gin.Context, err error) {
	errors, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": "--" + err.Error(),
		})
		return 
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errors.Translate(global.Translator)),
	})
}

func removeTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for filed, err := range fields {
		rsp[filed[strings.Index(filed, ".")+1:]] = err
	}
	return rsp
}
