package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/zhaohaihang/order_service/global"
	"github.com/zhaohaihang/order_service/model"
	"github.com/zhaohaihang/order_service/proto"
	"go.uber.org/zap"
)

type OrderService struct {
	proto.UnimplementedOrderServer
}

const (
	SERVICE_NAME = "[order_service]"
)

func OrderTimeout(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

	type OrderInfo struct {
		OrderSn string
	}

	for i := range msgs {
		var orderInfo OrderInfo
		_ = json.Unmarshal(msgs[i].Body, &orderInfo)
		zap.S().Info("get order timeout msg")

		var order model.OrderInfo
		if result := global.DB.Model(model.OrderInfo{}).Where(model.OrderInfo{OrderSn: orderInfo.OrderSn}).First(&order); result.RowsAffected == 0 {
			return consumer.ConsumeSuccess, nil
		}
		if order.Status != "TRADE_SUCCESS" {
			tx := global.DB.Begin()
		
			//修改订单的状态
			order.Status = "TRADE_CLOSED"
			tx.Save(&order)

			//归还库存
			p, err := rocketmq.NewProducer(producer.WithNameServer([]string{
				fmt.Sprintf("%s:%d",
				global.ServiceConfig.RocketMQInfo.Host,
				global.ServiceConfig.RocketMQInfo.Port)}))
			if err != nil {
				panic("生成producer失败")
			}
			if err = p.Start(); err != nil {
				panic("启动producer失败")
			}
			_, err = p.SendSync(context.Background(), primitive.NewMessage("order_reback", msgs[i].Body))
			if err != nil {
				tx.Rollback()
				fmt.Printf("发送失败: %s\n", err)
				return consumer.ConsumeRetryLater, nil
			}
			return consumer.ConsumeSuccess, nil
		}
	}

	return consumer.ConsumeSuccess, nil
}
