package initialize

import (
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/zhaohaihang/user_api/global"
	myValidate "github.com/zhaohaihang/user_api/validator"
	"go.uber.org/zap"
)

var validate *validator.Validate
var ok bool

func InitValidator() {
	validate, ok = binding.Validator.Engine().(*validator.Validate)
	if !ok {
		zap.S().Errorw("binding validator failed")
	}
	initMobileValidator()
}

func initMobileValidator() {
	err := validate.RegisterValidation("mobile", myValidate.ValidateMobile)
	if err != nil {
		zap.S().Errorw("mobile bind failed", "error", err.Error())
	}
	err = validate.RegisterTranslation("mobile", global.Translator, func(ut ut.Translator) error {
		return ut.Add("mobile", "invalid mobile number", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, err := ut.T("mobile", fe.Field())
		if err != nil {
			zap.S().Errorw("mobile  Trans failed", "err", err.Error())
		}
		return t
	})
	if err != nil {
		zap.S().Errorw("mobile register translation", "err", err.Error())
	}
	zap.S().Infow("init validator")
}
