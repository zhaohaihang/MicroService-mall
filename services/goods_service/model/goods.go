package model

import (
	"context"
	"strconv"

	"github.com/zhaohaihang/goods_service/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Category  目录表结构
type Category struct {
	gorm.Model
	Name             string      `gorm:"type:varchar(20);not null" json:"name"`
	ParentCategoryID uint       `json:"parent" gorm:"default:null"`   // default:null 解决顶级分类外键无法插入null的问题
	ParentCategory   *Category   `json:"-"`
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"`
	Level            int32       `gorm:"type:int;not null;default:1" json:"level"`
	IsTab            bool        `gorm:"default:false;not null" json:"is_tab"`
}

// TableName 自定义表名
func (Category) TableName() string {
	return "category"
}

// Brands 品牌表结构
type Brand struct {
	gorm.Model
	Name string `gorm:"type:varchar(50);not null"`
	Logo string `gorm:"type:varchar(200);default:'';not null"`
}

func (Brand) TableName() string {
	return "brands"
}

// GoodsCategoryBrand 商品目录表结构
type GoodsCategoryBrand struct {
	gorm.Model
	CategoryID uint `gorm:"type:int;index:idx_category_brand,unique"`
	Category   Category

	BrandID uint `gorm:"type:int;index:idx_category_brand,unique"`
	Brand   Brand
}

// TableName 自定义表名
func (GoodsCategoryBrand) TableName() string {
	return "goodscategorybrand"
}

// Banner 横幅表结构
type Banner struct {
	gorm.Model
	Image string `gorm:"type:varchar(200);not null"`
	Url   string `gorm:"type:varchar(200);not null"`
	Index int32  `gorm:"type:int;default:1;not null"`
}

// TableName 自定义表名
func (Banner) TableName() string {
	return "banner"
}

// Goods 商品表结构
type Goods struct {
	gorm.Model

	CategoryID uint `gorm:"type:int;not null"`
	Category   Category
	BrandID    uint `gorm:"type:int;not null;column:brand_id"`
	Brand      Brand

	OnSale   bool `gorm:"default:false;not null"`
	ShipFree bool `gorm:"default:false;not null"`
	IsNew    bool `gorm:"default:false;not null"`
	IsHot    bool `gorm:"default:false;not null"`

	Name            string   `gorm:"type:varchar(100);not null"`
	GoodsSn         string   `gorm:"type:varchar(50);not null"`
	ClickNum        int32    `gorm:"type:int;default:0;not null"`
	SoldNum         int32    `gorm:"type:int;default:0;not null"`
	FavNum          int32    `gorm:"type:int;default:0;not null"`
	MarketPrice     float32  `gorm:"not null"`
	ShopPrice       float32  `gorm:"not null"`
	GoodsBrief      string   `gorm:"type:varchar(100);not null"`
	Images          GormList `gorm:"type:varchar(1000);not null"`
	DescImages      GormList `gorm:"not null"`
	GoodsFrontImage string   `gorm:"type:varchar(200);not null"`
	Stocks          int32    `gorm:"type:int;default:0;not null;column:stocks"`
}

func (g *Goods) AfterCreate(tx *gorm.DB) (err error) {
	esModel := EsGoods{
		ID:          int32(g.ID),
		CategoryID:  int32(g.CategoryID),
		BrandsID:    int32(g.BrandID),
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,
		MarketPrice: g.MarketPrice,
		GoodsBrief:  g.GoodsBrief,
		ShopPrice:   g.ShopPrice,
	}
	_, err = global.EsClient.Index().Index(esModel.GetIndexName()).BodyJson(esModel).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		zap.S().Errorw("Error", "message", "sync es data failed", "err", err.Error())
		return err
	}
	return nil
}

func (g *Goods) AfterUpdate(tx *gorm.DB) (err error) {
	esModel := EsGoods{
		ID:          int32(g.ID),
		CategoryID:  int32(g.CategoryID),
		BrandsID:    int32(g.BrandID),
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,
		MarketPrice: g.MarketPrice,
		GoodsBrief:  g.GoodsBrief,
		ShopPrice:   g.ShopPrice,
	}

	_, err = global.EsClient.Update().Index(esModel.GetIndexName()).Doc(esModel).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		zap.S().Errorw("Error", "message", "sync es data failed", "err", err.Error())
		return err
	}
	return nil
}

func (g *Goods) AfterDelete(tx *gorm.DB) (err error) {
	_, err = global.EsClient.Delete().Index(EsGoods{}.GetIndexName()).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		zap.S().Errorw("Error", "message", "sync es data failed", "err", err.Error())
		return err
	}
	return nil
}
