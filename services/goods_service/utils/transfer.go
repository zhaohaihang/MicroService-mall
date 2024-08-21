package utils

import (
	"github.com/zhaohaihang/goods_service/model"
	"github.com/zhaohaihang/goods_service/proto"
)

func BannerToBannerResponse(banner *model.Banner) *proto.BannerResponse {
	return &proto.BannerResponse{
		Id:    int32(banner.ID),
		Index: banner.Index,
		Image: banner.Image,
		Url:   banner.Url,
	}
}

func ModelToResponse(goods *model.Goods) proto.GoodsInfoResponse {
	return proto.GoodsInfoResponse{
		Id:              int32(goods.ID),
		CategoryId:      goods.CategoryID,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		DescImages:      goods.DescImages,
		Images:          goods.Images,
		Category: &proto.CategoryBriefInfoResponse{
			Id:   int32(goods.Category.ID),
			Name: goods.Category.Name,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   int32(goods.Brand.ID),
			Name: goods.Brand.Name,
			Logo: goods.Brand.Logo,
		},
	}
}