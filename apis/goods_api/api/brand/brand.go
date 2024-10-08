package brand

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

// List  获取品牌列表
func List(ctx *gin.Context) {
	entry, blockError := utils.SentinelEntry(ctx)
	if blockError != nil {
		zap.S().Errorw("Error", "message", "Request too frequent")
		utils.HandleRequestFrequentError(ctx)
		return 
	}

	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)
	response, err := global.GoodsClient.BrandList(context.WithValue(context.Background(), "ginContext", ctx), &proto.BrandFilterRequest{
		Pages:       int32(pnInt),
		PagePerNums: int32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	responseMap := make(map[string]interface{})
	brandList := make([]interface{}, 0)
	responseMap["total"] = response.Total
	for _, value := range response.Data {
		brandList = append(brandList, map[string]interface{}{
			"id":   value.Id,
			"name": value.Name,
			"logo": value.Logo,
		})
	}
	responseMap["data"] = brandList
	ctx.JSON(http.StatusOK, responseMap)
	entry.Exit()
}


// New godoc
// @Summary 创建品牌
// @Description 创建品牌
// @Tags Brand
// @ID  /goods/v1/brands  
// @Accept  json
// @Produce  json
// @Param data body forms.BrandForm true "body"
// @Success 200 {string} string "ok"
// @Router /goods/v1/brands   [post]
func New(ctx *gin.Context) {
	entry, blockError := utils.SentinelEntry(ctx)
	if blockError != nil {
		zap.S().Errorw("Error", "message", "Request too frequent")
		utils.HandleRequestFrequentError(ctx)
		return 
	}

	brandForm := forms.BrandForm{}
	err := ctx.ShouldBind(&brandForm)
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())
		utils.HandleValidatorError(ctx, err)
		return
	}

	response, err := global.GoodsClient.CreateBrand(context.WithValue(context.Background(), "ginContext", ctx), &proto.BrandRequest{
		Name: brandForm.Name,
		Logo: brandForm.Logo,
	})
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	responseMap := make(map[string]interface{})
	responseMap["id"] = response.Id
	responseMap["name"] = response.Name
	responseMap["logo"] = response.Logo
	ctx.JSON(http.StatusOK, responseMap)
	entry.Exit()
}

func Delete(ctx *gin.Context) {
	entry, blockError := utils.SentinelEntry(ctx)
	if blockError != nil {
		zap.S().Errorw("Error", "message", "Request too frequent")
		utils.HandleRequestFrequentError(ctx)
		return 
	}

	id := ctx.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())
		ctx.Status(http.StatusNotFound)
		return
	}

	response, err := global.GoodsClient.DeleteBrand(context.WithValue(context.Background(), "ginContext", ctx), &proto.BrandRequest{
		Id: int32(idInt),
	})
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

func Update(ctx *gin.Context) {
	entry, blockError := utils.SentinelEntry(ctx)
	if blockError != nil {
		zap.S().Errorw("Error", "message", "Request too frequent")
		utils.HandleRequestFrequentError(ctx)
		return 
	}

	brandForm := forms.BrandForm{}
	err := ctx.ShouldBind(&brandForm)
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())

		utils.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())
		ctx.Status(http.StatusNotFound)
		return
	}

	response, err := global.GoodsClient.UpdateBrand(context.WithValue(context.Background(), "ginContext", ctx), &proto.BrandRequest{
		Id:   int32(idInt),
		Name: brandForm.Name,
		Logo: brandForm.Logo,
	})
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id":   response.Id,
		"name": response.Name,
		"logo": response.Logo,
	})
	entry.Exit()
}
