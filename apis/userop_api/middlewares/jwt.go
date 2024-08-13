package middlewares

import (
	"errors"
	"net/http"

	"github.com/zhaohaihang/userop_web/global"
	"go.uber.org/zap"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/userop_web/models"
)

type JWT struct {
	SigningKey []byte
}

var (
	ErrTokenExpired     = errors.New("Token is expired")
	ErrTokenNotValidYet = errors.New("Token not active yet")
	ErrTokenMalformed   = errors.New("That's not even a token")
	ErrTokenInvalid     = errors.New("Couldn't handle this token:")
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 判断是否携带token
		token := c.Request.Header.Get("x-token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "please login first",
			})
			c.Abort()
			return
		}

		// 解析token
		j := NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == ErrTokenExpired {
				c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "authorize has expire",
				})
				c.Abort()
				return
			}
			zap.S().Errorw("parse token failed", "err", err)
			c.JSON(http.StatusUnauthorized, "has not login")
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Set("userId", claims.ID)
		c.Next()
	}
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.ApiConfig.JWTInfo.SigningKey),
	}
}

func (j *JWT) ParseToken(tokenString string) (*models.CustomClaims, error) {

	keyFunc := func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, keyFunc)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			} else {
				return nil, ErrTokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, ErrTokenInvalid
	} else {
		return nil, ErrTokenInvalid
	}
}
