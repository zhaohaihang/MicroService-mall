package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"github.com/zhaohaihang/user_api/middlewares"
	"github.com/zhaohaihang/user_api/models"
)

// GenerateToken 生成Token
func GenerateToken(Id uint, NickName string, Role uint) (string, error) {
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          Id,
		NickName:    NickName,
		AuthorityId: Role,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               // 签名的生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, // 设置30天过期
			Issuer:    "microservice-mall",
		},
	}
	token, err := j.CreateToken(claims)
	return token, err
}
