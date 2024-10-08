package goods

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/goods_api/forms"
	"github.com/zhaohaihang/goods_api/global"
	"github.com/zhaohaihang/goods_api/proto"
	"github.com/zhaohaihang/goods_api/utils"
	"go.uber.org/zap"

	"net/http"
	"strconv"
)

func List(ctx *gin.Context) {
	entry, blockError := utils.SentinelEntry(ctx)
	if blockError != nil {
		zap.S().Errorw("Error", "message", "Request too frequent")
		utils.HandleRequestFrequentError(ctx)
		return 
	}

	zap.S().Infof("goods List request:%v", ctx.Request.Host)
	request := &proto.GoodsFilterRequest{}

	priceMin := ctx.DefaultQuery("pmin", "0")
	priceMinInt, _ := strconv.Atoi(priceMin)
	request.PriceMin = int32(priceMinInt)
	priceMax := ctx.DefaultQuery("pmax", "0")
	priceMaxInt, _ := strconv.Atoi(priceMax)
	request.PriceMax = int32(priceMaxInt)
	isHot := ctx.DefaultQuery("ih", "0")
	if isHot == "1" {
		request.IsHot = true
	}
	isNew := ctx.DefaultQuery("in", "0")
	if isNew == "1" {
		request.IsNew = true
	}
	isTab := ctx.DefaultQuery("it", "0")
	if isTab == "1" {
		request.IsTab = true
	}
	categoryId := ctx.DefaultQuery("c", "0")
	categoryIdInt, _ := strconv.Atoi(categoryId)
	request.TopCategory = int32(categoryIdInt)
	pages := ctx.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)
	request.Pages = int32(pagesInt)
	perNums := ctx.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	request.PagePerNums = int32(perNumsInt)
	keywords := ctx.DefaultQuery("q", "")
	request.KeyWords = keywords
	brandId := ctx.DefaultQuery("b", "0")
	brandIdInt, _ := strconv.Atoi(brandId)
	request.Brand = int32(brandIdInt)
	ginContext := context.WithValue(context.Background(), "ginContext", ctx)
	response, err := global.GoodsClient.GoodsList(ginContext, request)
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	responseMap := map[string]interface{}{
		"total": response.Total,
	}
	goodsList := make([]interface{}, 0)
	for _, goods := range response.Data {
		goodsMap := map[string]interface{}{
			"id":          goods.Id,
			"name":        goods.Name,
			"goods_brief": goods.GoodsBrief,
			"desc":        goods.GoodsDesc,
			"ship_free":   goods.ShipFree,
			"desc_image":  goods.DescImages,
			"front_image": goods.GoodsFrontImage,
			"shop_price":  goods.ShopPrice,
			"category": map[string]interface{}{
				"id":   goods.Category.Id,
				"name": goods.Category.Name,
			},
			"brand": map[string]interface{}{
				"id":   goods.Brand.Id,
				"name": goods.Brand.Name,
				"logo": goods.Brand.Logo,
			},
			"is_host": goods.IsHot,
			"is_new":  goods.IsNew,
			"on_sale": goods.OnSale, 
		}
		goodsList = append(goodsList, goodsMap)
	}
	responseMap["data"] = goodsList
	ctx.JSON(http.StatusOK, responseMap)
	entry.Exit()
}

// New godoc
// @Summary 创建商品
// @Description 创建商品
// @Tags Goods
// @ID  /goods/v1/goods
// @Accept  json
// @Produce  json
// @Param data body forms.GoodsForm true "body"
// @Success 200 {string} string "ok"
// @Router /goods/v1/goods [post]
func New(ctx *gin.Context) {
	entry, blockError := utils.SentinelEntry(ctx)
	if blockError != nil {
		zap.S().Errorw("Error", "message", "Request too frequent")
		utils.HandleRequestFrequentError(ctx)
		return 
	}

	zap.S().Infof("goods New request:%v", ctx.Request.Host)
	goodsForm := forms.GoodsForm{}
	if err := ctx.ShouldBindJSON(&goodsForm); err != nil {
		zap.S().Errorw("Error", "err", err.Error())
		utils.HandleValidatorError(ctx, err)
		return
	}

	goodsClient := global.GoodsClient
	rsp, err := goodsClient.CreateGoods(context.WithValue(context.Background(), "ginContext", ctx), &proto.CreateGoodsInfo{
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	})
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, rsp)
	entry.Exit()
}

// Detail  获取商品详情
func Detail(ctx *gin.Context) {
	entry, blockError := utils.SentinelEntry(ctx)
	if blockError != nil {
		zap.S().Errorw("Error", "message", "Request too frequent")
		utils.HandleRequestFrequentError(ctx)
		return 
	}
	
	zap.S().Infof("goods Detail request:%v", ctx.Request.Host)
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	response, err := global.GoodsClient.GetGoodsDetail(context.WithValue(context.Background(), "ginContext", ctx), &proto.GoodsInfoRequest{Id: int32(i)})
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	rep := map[string]interface{}{
		"id":          response.Id,
		"name":        response.Name,
		"goods_brief": response.GoodsBrief,
		"desc":        response.GoodsDesc,
		"ship_free":   response.ShipFree,
		"images":      response.Images,
		"desc_images": response.DescImages,
		"front_image": response.GoodsFrontImage,
		"shop_price":  response.ShopPrice,
		"category": map[string]interface{}{
			"id":   response.Category.Id,
			"name": response.Category.Name,
		},
		"brand": map[string]interface{}{
			"id":   response.Brand.Id,
			"name": response.Brand.Name,
			"logo": response.Brand.Logo,
		},
		"is_hot":  response.IsHot,
		"is_new":  response.IsNew,
		"on_sale": response.OnSale,
	}
	ctx.JSON(http.StatusOK, rep)
	entry.Exit()
}

func Delete(ctx *gin.Context) {
	entry, blockError := utils.SentinelEntry(ctx)
	if blockError != nil {
		zap.S().Errorw("Error", "message", "Request too frequent")
		utils.HandleRequestFrequentError(ctx)
		return 
	}
	zap.S().Infof("goods Delelte request:%v", ctx.Request.Host)
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())

		ctx.Status(http.StatusNotFound)
		return
	}
	response, err := global.GoodsClient.DeleteGoods(context.WithValue(context.Background(), "ginContext", ctx), &proto.DeleteGoodsInfo{Id: int32(i)})
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": response.Success,
	})

	entry.Exit()
}

func UpdateStatus(ctx *gin.Context) {
	entry, blockError := utils.SentinelEntry(ctx)
	if blockError != nil {
		zap.S().Errorw("Error", "message", "Request too frequent")
		utils.HandleRequestFrequentError(ctx)
		return 
	}
	zap.S().Infof("goods Update request:%v", ctx.Request.Host)
	goodsStatusForm := forms.GoodsStatusForm{}
	err := ctx.ShouldBind(&goodsStatusForm)
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())
		utils.HandleValidatorError(ctx, err)
		return
	}
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	response, err := global.GoodsClient.UpdateGoodsStatus(context.WithValue(context.Background(), "ginContext", ctx), &proto.CreateGoodsInfo{
		Id:     int32(i),
		IsHot:  *goodsStatusForm.IsHot,
		IsNew:  *goodsStatusForm.IsNew,
		OnSale: *goodsStatusForm.OnSale,
	})
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id":          response.Id,
		"name":        response.Name,
		"goods_brief": response.GoodsBrief,
		"desc":        response.GoodsDesc,
		"ship_free":   response.ShipFree,
		"images":      response.Images,
		"desc_images": response.DescImages,
		"front_image": response.GoodsFrontImage,
		"shop_price":  response.ShopPrice,
		"category": map[string]interface{}{
			"id":   response.Category.Id,
			"name": response.Category.Name,
		},
		"brand": map[string]interface{}{
			"id":   response.Brand.Id,
			"name": response.Brand.Name,
			"logo": response.Brand.Logo,
		},
		"is_hot":  response.IsHot,
		"is_new":  response.IsNew,
		"on_sale": response.OnSale,
	})
	entry.Exit()
}

func Update(ctx *gin.Context) {
	entry, blockError := utils.SentinelEntry(ctx)
	if blockError != nil {
		zap.S().Errorw("Error", "message", "Request too frequent")
		utils.HandleRequestFrequentError(ctx)
		return 
	}
	goodsForm := forms.GoodsForm{}
	err := ctx.ShouldBind(&goodsForm)
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())
		utils.HandleValidatorError(ctx, err)
		return
	}
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	response, err := global.GoodsClient.UpdateGoods(context.WithValue(context.Background(), "ginContext", ctx), &proto.CreateGoodsInfo{
		Id:              int32(i),
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	})
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id":          response.Id,
		"name":        response.Name,
		"goods_brief": response.GoodsBrief,
		"desc":        response.GoodsDesc,
		"ship_free":   response.ShipFree,
		"images":      response.Images,
		"desc_images": response.DescImages,
		"front_image": response.GoodsFrontImage,
		"shop_price":  response.ShopPrice,
		"category": map[string]interface{}{
			"id":   response.Category.Id,
			"name": response.Category.Name,
		},
		"brand": map[string]interface{}{
			"id":   response.Brand.Id,
			"name": response.Brand.Name,
			"logo": response.Brand.Logo,
		},
		"is_hot":  response.IsHot,
		"is_new":  response.IsNew,
		"on_sale": response.OnSale,
	})
	entry.Exit()
}
