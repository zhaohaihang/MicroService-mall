package model

import "gorm.io/gorm"

type Inventory struct {
	gorm.Model
	Goods   int32 `gorm:"type:int;index"`
	Stocks  int32 `gorm:"type:int"`
	Version int32 `gorm:"type:int"`
}

func (i Inventory) TableName() string {
	return "inventory"
}
