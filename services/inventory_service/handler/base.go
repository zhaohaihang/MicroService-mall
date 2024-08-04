package handler

import (
	"context"
	"encoding/json"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/zhaohaihang/inventory_service/global"
	"github.com/zhaohaihang/inventory_service/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func AutoReback(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	type OrderInfo struct {
		OrderSn string
	}

	for i := range msgs {
		// 获取需要归还库存的订单号
		var orderInfo OrderInfo
		err := json.Unmarshal(msgs[i].Body, &orderInfo)
		if err != nil {
			zap.S().Errorf("unmarshal json failed")
			return consumer.ConsumeSuccess, nil
		}

		// 根据订单号获取订单中每个商品的信息
		tx := global.DB.Begin()
		var sellDetail model.StockSellDetail
		if result := tx.Model(&model.StockSellDetail{}).Where(&model.StockSellDetail{OrderSn: orderInfo.OrderSn, Status: 1}).First(&sellDetail); result.RowsAffected == 0 {
			return consumer.ConsumeSuccess, nil
		}

		// 归还每个商品的库存
		for _, orderGoods := range sellDetail.Details {
			if result := tx.Model(&model.Inventory{}).Where(&model.Inventory{Goods: orderGoods.GoodsId}).Update("stocks", gorm.Expr("stocks+?", orderGoods.Num)); result.RowsAffected == 0 {
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}
		}

		// 更新订单归还表
		if result := tx.Model(&model.StockSellDetail{}).Where(&model.StockSellDetail{OrderSn: orderInfo.OrderSn}).Update("status", 2); result.RowsAffected == 0 {
			tx.Rollback()
			return consumer.ConsumeRetryLater, nil
		}
		tx.Commit()
		return consumer.ConsumeSuccess, nil
	}
	return consumer.ConsumeSuccess, nil
}
