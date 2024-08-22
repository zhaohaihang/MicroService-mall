package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/olivere/elastic/v7"
	"github.com/opentracing/opentracing-go"
	"github.com/zhaohaihang/goods_service/global"
	"github.com/zhaohaihang/goods_service/model"
	"github.com/zhaohaihang/goods_service/proto"
	"github.com/zhaohaihang/goods_service/utils"
	"go.uber.org/zap"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Result struct {
	ID int32
}

// GoodsList 获取商品列表
func (g *GoodsServer) GoodsList(ctx context.Context, request *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "GoodsList", "request", request)

	parentSpan := opentracing.SpanFromContext(ctx)

	esSpan := opentracing.GlobalTracer().StartSpan("elasticsearch-goods", opentracing.ChildOf(parentSpan.Context()))
	// es条件查询
	q := elastic.NewBoolQuery()
	if request.KeyWords != "" {
		q = q.Must(elastic.NewMultiMatchQuery(request.KeyWords, "name", "desc"))
	}
	if request.IsHot {
		q = q.Filter(elastic.NewTermQuery("is_hot", request.IsHot))
	}
	if request.IsNew {
		q = q.Filter(elastic.NewTermQuery("is_new", request.IsNew))
	}
	if request.PriceMin > 0 {
		q = q.Filter(elastic.NewRangeQuery("shop_price").Gte(request.PriceMin))
	}
	if request.PriceMax > 0 {
		q = q.Filter(elastic.NewRangeQuery("shop_price").Lte(request.PriceMax))
	}
	if request.Brand > 0 {
		q = q.Filter(elastic.NewTermQuery("brands_id", request.Brand))
	}
	if request.TopCategory > 0 {
		categoryIds := make([]interface{}, 0)
		var subQuery string
		var category model.Category
		result := global.DB.First(&category, request.TopCategory)
		if result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "cateogry not found")
		}
		if category.Level == 1 {
			subQuery = fmt.Sprintf("select id from category where parent_category_id in (select id from category WHERE parent_category_id=%d)", request.TopCategory)
		} else if category.Level == 2 {
			subQuery = fmt.Sprintf("select id from category WHERE parent_category_id=%d", request.TopCategory)
		} else if category.Level == 3 {
			subQuery = fmt.Sprintf("select id from category WHERE id=%d", request.TopCategory)
		}

		var results []Result
		global.DB.Model(model.Category{}).Raw(subQuery).Scan(&results)
		for _, re := range results {
			categoryIds = append(categoryIds, re.ID)
		}
		q = q.Filter(elastic.NewTermsQuery("category_id", categoryIds...))
	}

	// 分页
	if request.Pages == 0 {
		request.Pages = 1
	}
	if request.PagePerNums > 100 {
		request.PagePerNums = 100
	}else if request.PagePerNums <= 0 {
		request.PagePerNums = 10
	}
	result, err := global.EsClient.Search().Index(model.EsGoods{}.GetIndexName()).Query(q).From(int(request.Pages)).Size(int(request.PagePerNums)).Do(context.Background())
	if err != nil {
		zap.S().Errorw("Error", "message", "es 查询goods失败", "err", err.Error())
	}
	
	// 获取es中查询出来的所有商品Id
	goodsIds := make([]int32, 0)
	for _, value := range result.Hits.Hits {
		goods := model.EsGoods{}
		_ = json.Unmarshal(value.Source, &goods)
		goodsIds = append(goodsIds, goods.ID)
	}
	esSpan.Finish()

	// 通过数据库补充字段
	goodListSpan := opentracing.GlobalTracer().StartSpan("goods_list", opentracing.ChildOf(parentSpan.Context()))
	
	response := &proto.GoodsListResponse{}
	var goods []model.Goods
	localResult :=  global.DB.Model(&model.Goods{}).Preload("Category").Preload("Brand").Find(&goods, goodsIds)
	response.Total = int32(localResult.RowsAffected)

	var goodsListResponse []*proto.GoodsInfoResponse
	for _, goods := range goods {
		goodsResponse := utils.GoodsToGoodsInfoResponse(&goods)
		goodsListResponse = append(goodsListResponse, &goodsResponse)
	}
	response.Data = goodsListResponse

	goodListSpan.Finish()
	return response, nil
}

// BatchGetGoods 批量获取商品信息
func (g *GoodsServer) BatchGetGoods(ctx context.Context, request *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "BatchGetGoods", "request", request)
	parentSpan := opentracing.SpanFromContext(ctx)
	batchGetGoodsSpan := opentracing.GlobalTracer().StartSpan("BatchGetGoods", opentracing.ChildOf(parentSpan.Context()))
	defer batchGetGoodsSpan.Finish()

	response := &proto.GoodsListResponse{}
	
	var goodsList []model.Goods
	result := global.DB.Where(request.Id).Find(&goodsList)
	response.Total = int32(result.RowsAffected)

	var goodsListResponse []*proto.GoodsInfoResponse
	for _, goods := range goodsList {
		goodsInfoResponse := utils.GoodsToGoodsInfoResponse(&goods)
		goodsListResponse = append(goodsListResponse, &goodsInfoResponse)
	}
	
	response.Data = goodsListResponse
	return response, nil
}

// CreateGoods 创建商品
func (g *GoodsServer) CreateGoods(ctx context.Context, request *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "CreateGoods", "request", request)
	parentSpan := opentracing.SpanFromContext(ctx)
	createGoodsSpan := opentracing.GlobalTracer().StartSpan("CreateGoods", opentracing.ChildOf(parentSpan.Context()))
	defer createGoodsSpan.Finish()

	var category model.Category
	if result := global.DB.First(&category, request.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "category is not exist")
	}

	var brand model.Brand
	if result := global.DB.First(&brand, request.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "brands is not exist")
	}

	goods := model.Goods{
		Brand:           brand,
		BrandID:         int32(brand.ID),
		Category:        category,
		CategoryID:      int32(category.ID),
		Name:            request.Name,
		GoodsSn:         request.GoodsSn,
		MarketPrice:     request.MarketPrice,
		ShopPrice:       request.ShopPrice,
		GoodsBrief:      request.GoodsBrief,
		ShipFree:        request.ShipFree,
		Images:          request.Images,
		DescImages:      request.DescImages,
		GoodsFrontImage: request.GoodsFrontImage,
		IsNew:           request.IsNew,
		IsHot:           request.IsHot,
		OnSale:          request.OnSale,
		Stocks:          request.Stocks,
	}
	global.DB.Create(&goods)
	
	response := utils.GoodsToGoodsInfoResponse(&goods)
	return &response, nil
}

// DeleteGoods 删除商品
func (g *GoodsServer) DeleteGoods(ctx context.Context, request *proto.DeleteGoodsInfo) (*proto.OperationResult, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "DeleteGoods", "request", request)
	parentSpan := opentracing.SpanFromContext(ctx)
	deleteGoodsSpan := opentracing.GlobalTracer().StartSpan("DeleteGoods", opentracing.ChildOf(parentSpan.Context()))
	defer deleteGoodsSpan.Finish()

	response := &proto.OperationResult{}
	result := global.DB.Delete(&model.Goods{}, request.Id)
	if result.RowsAffected == 0 {
		response.Success = false
		return nil, status.Errorf(codes.NotFound, "goods is not exist")
	}

	response.Success = true
	return response, nil
}

// UpdateGoods  更新商品信息
func (g GoodsServer) UpdateGoods(ctx context.Context, request *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "UpdateGoods", "request", request)
	parentSpan := opentracing.SpanFromContext(ctx)
	updateGoodsSpan := opentracing.GlobalTracer().StartSpan("UpdateGoods", opentracing.ChildOf(parentSpan.Context()))
	defer 	updateGoodsSpan.Finish()

	var goods model.Goods
	if result := global.DB.First(&goods, request.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "goods is not exist")
	}
	var category model.Category
	if result := global.DB.First(&category, request.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "category is not exist")
	}
	var brand model.Brand
	if result := global.DB.First(&brand, request.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "brand is not exist")
	}

	goods.Brand = brand
	goods.BrandID = int32(brand.ID)
	goods.Category = category
	goods.CategoryID = int32(category.ID)
	goods.Name = request.Name
	goods.GoodsSn = request.GoodsSn
	goods.MarketPrice = request.MarketPrice
	goods.ShopPrice = request.ShopPrice
	goods.GoodsBrief = request.GoodsBrief
	goods.ShipFree = request.ShipFree
	goods.Images = request.Images
	goods.DescImages = request.DescImages
	goods.GoodsFrontImage = request.GoodsFrontImage
	goods.IsNew = request.IsNew
	goods.IsHot = request.IsHot
	goods.OnSale = request.OnSale
	global.DB.Save(&goods)

	response := utils.GoodsToGoodsInfoResponse(&goods)
	return &response, nil
}

// GetGoodsDetail 获取商品详细信息
func (g *GoodsServer) GetGoodsDetail(ctx context.Context, request *proto.GoodsInfoRequest) (*proto.GoodsInfoResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "GetGoodsDetail", "request", request)
	parentSpan := opentracing.SpanFromContext(ctx)
	getGoodsDetailSpan := opentracing.GlobalTracer().StartSpan("GetGoodsDetail", opentracing.ChildOf(parentSpan.Context()))
	defer getGoodsDetailSpan.Finish()
	
	var goods model.Goods
	result := global.DB.Preload("Category").Preload("Brand").First(&goods, request.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "goods is not exist")
	}
	
	response := utils.GoodsToGoodsInfoResponse(&goods)
	return &response, nil
}

// UpdateGoodsStatus 更新商品状态
func (g *GoodsServer) UpdateGoodsStatus(ctx context.Context, request *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "UpdateGoodsStatus", "request", request)
	
	parentSpan := opentracing.SpanFromContext(ctx)
	updateGoodsStatusSpan := opentracing.GlobalTracer().StartSpan("UpdateGoodsStatus", opentracing.ChildOf(parentSpan.Context()))
	defer updateGoodsStatusSpan.Finish()

	var goods model.Goods
	result := global.DB.Preload("Category").Preload("Brand").First(&goods, request.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "goods is not exist")
	}
	goods.IsHot = request.IsHot
	goods.IsNew = request.IsNew
	goods.OnSale = request.OnSale
	result = global.DB.Save(&goods)
	if result.Error != nil {
		return nil, result.Error
	}
	
	response := utils.GoodsToGoodsInfoResponse(&goods)
	return &response, nil
}
