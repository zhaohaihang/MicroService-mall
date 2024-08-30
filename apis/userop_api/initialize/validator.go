package initialize

import (
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/zhaohaihang/userop_api/global"
	"go.uber.org/zap"

	myValidate "github.com/zhaohaihang/userop_api/validator"
)

var validate *validator.Validate
var ok bool

func InitValidator() {
	validate, ok = binding.Validator.Engine().(*validator.Validate)
	if !ok {
		zap.S().Errorw("bind consum validator failed")
	}
	initMobileValidator()
}

func initMobileValidator() {
	err := validate.RegisterValidation("mobile", myValidate.ValidateMobile)
	if err != nil {
		zap.S().Errorw("mobile bind failed", "err", err.Error())
	}
	err = validate.RegisterTranslation("mobile", global.Translator, func(ut ut.Translator) error {
		return ut.Add("mobile", "{0} 非法的手机号码", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, err := ut.T("mobile", fe.Field())
		if err != nil {
			zap.S().Errorw("mobile register translation failed", "err", err.Error())
		}
		return t
	})
	if err != nil {
		zap.S().Errorw("mobile register translation failed", "err", err.Error())
	}
	zap.S().Infow("mobile load validation success")
}
