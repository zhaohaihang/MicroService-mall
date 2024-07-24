package models

import (
	"github.com/dgrijalva/jwt-go"
)

// CustomClaims JWT中携带的信息数据
type CustomClaims struct {
	ID          uint
	NickName    string
	AuthorityId uint
	jwt.StandardClaims
}
