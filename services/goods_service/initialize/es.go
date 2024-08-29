package initialize

import (
	"context"
	"fmt"
	"strconv"

	"github.com/olivere/elastic/v7"
	"github.com/zhaohaihang/goods_service/global"
	"github.com/zhaohaihang/goods_service/model"
	"go.uber.org/zap"
)

func InitEs() {
	var err error
	EsInfo := global.ServiceConfig.EsInfo
	url := fmt.Sprintf("http://%s:%d", EsInfo.Host, EsInfo.Port)
	global.EsClient, err = elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		zap.S().Fatalw("Error", "message", "es init failed", "err", err.Error())
	}
	EsGoods := model.EsGoods{}
	// 判断Index是否存在
	isExists, err := global.EsClient.IndexExists(EsGoods.GetIndexName()).Do(context.Background())
	if err != nil {
		zap.S().Fatalw("Error", "message", "check es index exists failed ", "err", err.Error())
	}
	if !isExists {
		_, err := global.EsClient.CreateIndex(EsGoods.GetIndexName()).BodyString(EsGoods.GetMapping()).Do(context.Background())
		if err != nil {
			zap.S().Fatalw("Error", "message", "create index failed ", "err", err.Error())
		}
	}
	zap.S().Infof("elastic search init success")
	zap.S().Infof("start sync data from mysql to es")
	Mysql2Es()
}

func Mysql2Es() {
	var goods []model.Goods
	global.DB.Find(&goods)
	for _, g := range goods {
		esModel := model.EsGoods{
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
		_, err := global.EsClient.Index().
			Index(esModel.GetIndexName()).
			BodyJson(esModel).
			Id(strconv.Itoa(int(g.ID))).
			Do(context.Background())
		if err != nil {
			zap.S().Errorw("Error", "message", "sync goods data to es", "err", err.Error())
		}
	}
}
