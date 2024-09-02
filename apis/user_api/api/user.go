package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/zhaohaihang/user_api/forms"
	"github.com/zhaohaihang/user_api/global"
	"github.com/zhaohaihang/user_api/proto"
	"github.com/zhaohaihang/user_api/response"
	"go.uber.org/zap"

	"github.com/zhaohaihang/user_api/utils"
)

// GetUserList
// @Tags      User
// @Summary   获取用户列表
// @Security  ApiKeyAuth
// @Produce  application/json
// @Param     page    query   string  false "页码"
// @Param     size    query   string  false "大小"
// @Success 200 {string} string "ok"
// @Router   /user_api/v1/user/list [get]
func GetUserList(c *gin.Context) {
	// 1.获取参数
	pageNum := c.DefaultQuery("page", "0")
	pageNumInt, _ := strconv.Atoi(pageNum)
	pageSize := c.DefaultQuery("size", "10")
	pageSizeInt, _ := strconv.Atoi(pageSize)
	// 2.调用rpc服务
	resp, err := global.UserClient.GetUserList(context.WithValue(context.Background(), "ginContext", c), &proto.PageInfoRequest{
		PageNum:  uint32(pageNumInt),
		PageSize: uint32(pageSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] get user list failed", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, c)
		return
	}
	// 3.返回查询结果
	result := make([]interface{}, 0)
	for _, value := range resp.Data {
		user := response.UserResponse{
			Id:       value.Id,
			Name:     value.NickName,
			Gender:   value.Gender,
			Mobile:   value.Mobile,
			Birthday: time.Time(time.Unix(int64(value.Birthday), 0)),
		}
		result = append(result, user)
	}
	c.JSON(http.StatusOK, result)
}

// PasswordLogin godoc
// @Summary 手机密码登录
// @Description 手机密码登录
// @Tags User
// @ID /user/login
// @Accept  json
// @Produce  json
// @Param data body forms.PasswordLoginForm true "body"
// @Success 200 {string} string "ok"
// @Router /user_api/v1/user/pwd_login [post]
func PasswordLogin(c *gin.Context) {
	// 1.实例化验证对象
	passwordLoginForm := forms.PasswordLoginForm{}
	// 2.判断是否有错误
	if err := c.ShouldBind(&passwordLoginForm); err != nil {
		utils.HandleValidatorError(c, err)
		return
	}

	verify := store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, true)
	if !verify {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "smscode is not right",
		})
		return
	}

	// 3.登录
	// 3.1获取用户加密后的密码
	userInfoResponse, err := global.UserClient.GetUserByMobile(context.WithValue(context.Background(), "ginContext", c), &proto.MobileRequest{Mobile: passwordLoginForm.Mobile})
	if err != nil {
		zap.S().Errorw("get user by mobile failed", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, c)
	}
	// 4.密码进行验证比对
	checkPasswordResponse, err := global.UserClient.CheckPassword(context.WithValue(context.Background(), "ginContext", c), &proto.CheckPasswordRequest{
		Password:          passwordLoginForm.Password,
		EncryptedPassword: userInfoResponse.Password,
	})
	if err != nil {
		zap.S().Errorw("check passwd failed", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, c)
	}
	// 5.根据获取的结果返回
	if checkPasswordResponse.Success {
		token, err := utils.GenerateToken(uint(userInfoResponse.Id), userInfoResponse.NickName, uint(userInfoResponse.Role))
		if err != nil {
			zap.S().Errorw("generate token failed", "err:", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "generate token failed",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"id":        userInfoResponse.Id,
			"nickName":  userInfoResponse.NickName,
			"token":     token,
			"expiresAt": (time.Now().Unix() + 60*60*24*30) * 1000,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "login failed",
		})
	}
}

// Register godoc
// @Summary 注册用户
// @Description 注册用户
// @Tags User
// @ID  /user_api/v1/user/register
// @Accept  json
// @Produce  json
// @Param data body forms.RegisterForm true "body"
// @Success 200 {string} string "ok"
// @Router /user_api/v1/user/register [post]
func Register(c *gin.Context) {
	// 1.表单认证
	registerForm := forms.RegisterForm{}
	err := c.ShouldBind(&registerForm)
	if err != nil {
		zap.S().Errorw("bind error", "err", err.Error())
		utils.HandleValidatorError(c, err)
		return
	}
	// 2.通过redis 验证 验证码是否正确
	value, err := global.RedisClient.Get(context.Background(), registerForm.Mobile).Result()
	if err == redis.Nil { // redis中没有验证码
		zap.S().Warnw("can not find code in redis", "mobile num", registerForm.Mobile)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "code error",
		})
		return
	} else { // 验证码错误
		if value != registerForm.Code {
			zap.S().Warnw("code not match", "mobile:", registerForm.Mobile, "rediscode:", value, "code:", " registerForm.Code ")
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "code error",
			})
			return
		}
	}
	userResponse, err := global.UserClient.CreateUser(context.Background(), &proto.CreateUserInfoRequest{
		NickName: registerForm.Mobile,
		Password: registerForm.Password,
		Mobile:   registerForm.Mobile,
	})
	if err != nil {
		zap.S().Errorw("create user failed", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, c)
		return
	}
	token, err := utils.GenerateToken(uint(userResponse.Id), userResponse.NickName, uint(userResponse.Role))
	if err != nil {
		zap.S().Errorw("generate token failed", "err:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "generate token failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":        userResponse.Id,
		"nickName":  userResponse.NickName,
		"token":     token,
		"expiresAt": (time.Now().Unix() + 60*60*24*30) * 1000,
	})
}
