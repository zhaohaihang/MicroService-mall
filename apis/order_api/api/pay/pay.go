package pay

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"github.com/zhaohaihang/order_api/global"
	"github.com/zhaohaihang/order_api/proto"
	"github.com/zhaohaihang/order_api/utils"
	"go.uber.org/zap"
)

func Notify(ctx *gin.Context) {
	entry, blockError := utils.SentinelEntry(ctx)
	if blockError != nil {
		zap.S().Errorw("Error", "message", "Request too frequent")
		utils.HandleRequestFrequentError(ctx)
		return
	}
	
	aliPayInfo := global.ApiConfig.AlipayInfo
	client, err := alipay.New(aliPayInfo.AppID, aliPayInfo.PrivateKey, false)
	if err != nil {
		zap.S().Errorw("Error", "message", "creat alipay client failed", "err", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = client.LoadAliPayPublicKey(aliPayInfo.AliPublicKey)
	if err != nil {
		zap.S().Errorw("Error", "message", "LoadAliPayPublicKey failed", "err", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	notification, err := client.GetTradeNotification(ctx.Request)
	if err != nil {
		zap.S().Errorw("Error", "message", "get Notification failed", "err", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	_, err = global.OrderClient.UpdateOrderStatus(context.WithValue(context.Background(), "ginContext", ctx), &proto.OrderStatus{
		OrderSn: notification.OutTradeNo,
		Status:  string(notification.TradeStatus),
	})
	if err != nil {
		zap.S().Errorw("Error", "message", "update order status failed", "err", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.String(http.StatusOK, "success")
	entry.Exit()
}
