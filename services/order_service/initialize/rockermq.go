package initialize

import (
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/zhaohaihang/order_service/global"
	"github.com/zhaohaihang/order_service/handler"
	"go.uber.org/zap"
)

func InitRocketMq() {
	//监听订单支付超时
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{fmt.Sprintf("%s:%d",
			global.ServiceConfig.RocketMQInfo.Host,
			global.ServiceConfig.RocketMQInfo.Port)}),
		consumer.WithGroupName("mq-order"),
	)

	if err := c.Subscribe("order_timeout", consumer.MessageSelector{}, handler.OrderTimeout); err != nil {
		zap.S().Errorf("Subscribe  order_timeout msg error\n")
	}
	_ = c.Start()

}
