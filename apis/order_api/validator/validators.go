package validator

import (
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"regexp"
)

// ValidateMobile 验证手机号是否符合规则
func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	ok, _ := regexp.MatchString(`^1(3[0-9]|5[0-3,5-9]|7[1-3,5-8]|8[0-9])\d{8}$`, mobile)
	if !ok {
		zap.S().Errorw("mobile Validate failed ", "mobile:", mobile)
		return false
	}
	return true
}
