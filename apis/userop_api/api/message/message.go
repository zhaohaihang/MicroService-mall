package message

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/userop_web/forms"
	"github.com/zhaohaihang/userop_web/global"
	"github.com/zhaohaihang/userop_web/models"
	"github.com/zhaohaihang/userop_web/proto"
	"github.com/zhaohaihang/userop_web/utils"
	"go.uber.org/zap"
)

func ListMessage(ctx *gin.Context) {
	entry, blockError := utils.SentinelEntry(ctx)
	if blockError != nil {
		zap.S().Errorw("Error", "message", "Request too frequent")
		utils.HandleRequestFrequentError(ctx)
		return
	}

	request := &proto.MessageRequest{}
	userId, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	if currentUser.AuthorityId == 1 {
		request.UserId = int32(userId.(uint))
	}

	response, err := global.MessageClient.MessageList(context.WithValue(context.Background(), "ginContext", ctx), request)
	if err != nil {
		zap.S().Errorw("Error", "message", "get message list failed", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	responseMap := gin.H{
		"total": response.Total,
	}

	result := make([]interface{}, 0)
	for _, item := range response.Data {
		temp := map[string]interface{}{}
		temp["id"] = item.Id
		temp["user_id"] = item.UserId
		temp["type"] = item.MessageType
		temp["subject"] = item.Subject
		temp["message"] = item.Message
		temp["file"] = item.File
		result = append(result, temp)
	}
	responseMap["data"] = result
	ctx.JSON(http.StatusOK, responseMap)
	entry.Exit()
}

func CreateMessage(ctx *gin.Context) {
	entry, blockError := utils.SentinelEntry(ctx)
	if blockError != nil {
		zap.S().Errorw("Error", "message", "Request too frequent")
		utils.HandleRequestFrequentError(ctx)
		return
	}

	userId, _ := ctx.Get("userId")

	messageForm := forms.MessageForm{}
	err := ctx.ShouldBind(&messageForm)
	if err != nil {
		zap.S().Errorw("Error", "message", "bind form failed", "err", err.Error())
		utils.HandleValidatorError(ctx, err)
		return
	}

	response, err := global.MessageClient.CreateMessage(context.WithValue(context.Background(), "ginContext", ctx), &proto.MessageRequest{
		UserId:      int32(userId.(uint)),
		MessageType: messageForm.MessageType,
		Subject:     messageForm.Subject,
		Message:     messageForm.Message,
		File:        messageForm.File,
	})
	if err != nil {
		zap.S().Errorw("Error", "message", "create meaages failed", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Id": response.Id,
	})
	entry.Exit()
}
