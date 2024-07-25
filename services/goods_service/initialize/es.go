package initialize

import (
	"context"
	"fmt"

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
		zap.S().Errorw("Error", "message", "es init failed", "err", err.Error())
		panic(err)
	}
	EsGoods := model.EsGoods{}
	// 判断Index是否存在
	isExists, err := global.EsClient.IndexExists(EsGoods.GetIndexName()).Do(context.Background())
	if err != nil {
		zap.S().Errorw("Error", "message", "check es index exists failed ", "err", err.Error())
		panic(err)
	}
	if !isExists {
		_, err := global.EsClient.CreateIndex(EsGoods.GetIndexName()).BodyString(EsGoods.GetMapping()).Do(context.Background())
		if err != nil {
			zap.S().Errorw("Error", "message", "create index failed ", "err", err.Error())
			panic(err)
		}
	}
	zap.S().Infof("elastic search init success")
}
