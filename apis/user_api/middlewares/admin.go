package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/user_api/models"
)

// AdminAuth 登录管理员权限验证
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		currentUser := claims.(*models.CustomClaims)

		if currentUser.AuthorityId != 2 {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "role has not Authorited",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
