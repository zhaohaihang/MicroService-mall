package initialize

import (
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/zhaohaihang/inventory_service/global"
	"github.com/zhaohaihang/inventory_service/handler"
)

func InitRocketMq() {
	//监听库存归还topic
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{fmt.Sprintf("%s:%d",
			global.ServiceConfig.RocketMQInfo.Host,
			global.ServiceConfig.RocketMQInfo.Port)}),
		consumer.WithGroupName("mq-inventory"),
	)

	if err := c.Subscribe("order_reback", consumer.MessageSelector{}, handler.AutoReback); err != nil {
		fmt.Println("读取消息失败")
	}
	_ = c.Start()

}
