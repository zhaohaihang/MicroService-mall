package api

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/user_api/forms"
	"github.com/zhaohaihang/user_api/global"
	"github.com/zhaohaihang/user_api/utils"
	"go.uber.org/zap"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	dysmsapi "github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

// SendNoteCode 发送短信验证码
func SendNoteCode(c *gin.Context) {
	// 表单验证
	sendSmsForm := forms.SendSmsForm{}
	err := c.ShouldBind(&sendSmsForm)
	if err != nil {
		zap.S().Errorw("Error", "method", "bind sms form failed", "err", err.Error())
		utils.HandleValidatorError(c, err)
		return
	}

	config := sdk.NewConfig()
	credential := credentials.NewAccessKeyCredential(global.ApiConfig.AliSmsInfo.ApiKey, global.ApiConfig.AliSmsInfo.ApiSecret)
	client, err := dysmsapi.NewClientWithOptions(global.ApiConfig.AliSmsInfo.RegionId, config, credential)
	if err != nil {
		zap.S().Fatal("Error", "method", "create sms client failed", "err", err.Error())
	}

	smsCode := generateNoteCode(5)
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.SignName = global.ApiConfig.AliSmsInfo.SignName
	request.TemplateCode = global.ApiConfig.AliSmsInfo.TemplateCode
	request.PhoneNumbers = sendSmsForm.Mobile
	request.TemplateParam = "{\"code\":\"" + smsCode + "\"}"

	response, err := client.SendSms(request)
	if err != nil {
		zap.S().Errorw("Error", "method", "send sms failed", "err", err.Error())
	}
	zap.S().Infof("response is %#v\n", response)

	global.RedisClient.Set(context.WithValue(context.Background(), "ginContext", c), sendSmsForm.Mobile, smsCode, 300*time.Second)
	c.JSON(http.StatusOK, gin.H{
		"msg": "send success",
	})
}

// generateNoteCode 生成随机验证码
func generateNoteCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().Unix())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		_, _ = fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}
