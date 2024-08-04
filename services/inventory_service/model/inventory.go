package model

import (
	"database/sql/driver"
	"encoding/json"

	"gorm.io/gorm"
)

type Inventory struct {
	gorm.Model
	Goods   int32 `gorm:"type:int;index"`
	Stocks  int32 `gorm:"type:int"`
	Version int32 `gorm:"type:int"`
}

func (i Inventory) TableName() string {
	return "inventory"
}

type GoodsDetail struct {
	GoodsId int32
	Num int32
}

type GoodsDetailList []GoodsDetail

func (g GoodsDetailList) Value() (driver.Value,error) {
	return json.Marshal(g)
}

func (g GoodsDetailList)Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte),&g)
}

type StockSellDetail struct {
	OrderSn string          `gorm:"type:varchar(200);index:idx_order_sn,unique;"`
	Status  int32           `gorm:"type:varchar(200)"` //1 表示已扣减 2. 表示已归还
	Details  GoodsDetailList `gorm:"type:varchar(200)"`
}

func (StockSellDetail) TableName() string {
	return "stockselldetail"
}