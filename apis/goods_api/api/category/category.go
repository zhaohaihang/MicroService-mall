package category

import (
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/zhaohaihang/goods_api/forms"
	"github.com/zhaohaihang/goods_api/global"
	"github.com/zhaohaihang/goods_api/proto"
	"github.com/zhaohaihang/goods_api/utils"
	"go.uber.org/zap"

	"net/http"
	"strconv"
)

// List 获取商品目录列表
func List(ctx *gin.Context) {
	entry, blockError := utils.SentinelEntry(ctx)
	if blockError != nil {
		zap.S().Errorw("Error", "message", "Request too frequent")
		utils.HandleRequestFrequentError(ctx)
		return 
	}

	response, err := global.GoodsClient.GetAllCategoriesList(context.WithValue(context.Background(), "ginContext", ctx), &empty.Empty{})
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	data := make([]interface{}, 0)
	err = json.Unmarshal([]byte(response.JsonData), &data)
	if err != nil {
		zap.S().Errorw("nmarshal data from rpc failed", "err", err.Error())
		return
	}
	ctx.JSON(http.StatusOK, data)
	entry.Exit()
}
 
// Detail 获取商品目录详情信息
func Detail(ctx *gin.Context) {
	entry, blockError := utils.SentinelEntry(ctx)
	if blockError != nil {
		zap.S().Errorw("Error", "message", "Request too frequent")
		utils.HandleRequestFrequentError(ctx)
		return 
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())
		ctx.Status(http.StatusNotFound)
		return
	}

	responseMap := make(map[string]interface{})
	subCategorys := make([]interface{}, 0)
	response, err := global.GoodsClient.GetSubCategory(context.WithValue(context.Background(), "ginContext", ctx), &proto.CategoryListRequest{
		Id: int32(i),
	})
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	for _, value := range response.SubCategory {
		subCategorys = append(subCategorys, map[string]interface{}{
			"id":              value.Id,
			"name":            value.Name,
			"level":           value.Level,
			"parent_category": value.ParentCategory,
			"is_tab":          value.IsTab,
		})
	}
	responseMap["id"] = response.Info.Id
	responseMap["name"] = response.Info.Name
	responseMap["level"] = response.Info.Level
	responseMap["parent_category"] = response.Info.ParentCategory
	responseMap["is_tab"] = response.Info.IsTab
	responseMap["sub_categorys"] = subCategorys
	ctx.JSON(http.StatusOK, responseMap)
	entry.Exit()
}

// New godoc
// @Summary 创建分类
// @Description 创建分类
// @Tags Category
// @ID  /goods/v1/categorys
// @Accept  json
// @Produce  json
// @Param data body forms.CategoryForm true "body"
// @Success 200 {string} string "ok"
// @Router /goods/v1/categorys [post]
func New(ctx *gin.Context) {
	entry, blockError := utils.SentinelEntry(ctx)
	if blockError != nil {
		zap.S().Errorw("Error", "message", "Request too frequent")
		utils.HandleRequestFrequentError(ctx)
		return 
	}

	categoryForm := forms.CategoryForm{}
	err := ctx.ShouldBind(&categoryForm)
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())
		utils.HandleValidatorError(ctx, err)
		return
	}

	response, err := global.GoodsClient.CreateCategory(context.WithValue(context.Background(), "ginContext", ctx), &proto.CategoryInfoRequest{
		Name:           categoryForm.Name,
		ParentCategory: categoryForm.ParentCategory,
		Level:          categoryForm.Level,
		IsTab:          *categoryForm.IsTab,
	})
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	responseMap := map[string]interface{}{
		"id":     response.Id,
		"name":   response.Name,
		"parent": response.ParentCategory,
		"level":  response.Level,
		"is_tab": response.IsTab,
	}
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
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())
		ctx.Status(http.StatusNotFound)
		return
	}

	response, err := global.GoodsClient.DeleteCategory(context.WithValue(context.Background(), "ginContext", ctx), &proto.DeleteCategoryRequest{Id: int32(i)})
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

// Update 更新目录信息
func Update(ctx *gin.Context) {
	entry, blockError := utils.SentinelEntry(ctx)
	if blockError != nil {
		zap.S().Errorw("Error", "message", "Request too frequent")
		utils.HandleRequestFrequentError(ctx)
		return 
	}

	categoryForm := forms.UpdateCategoryForm{}
	err := ctx.ShouldBind(&categoryForm)
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())
		utils.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())
		ctx.Status(http.StatusNotFound)
		return
	}

	request := &proto.CategoryInfoRequest{
		Id:   int32(i),
		Name: categoryForm.Name,
	}
	if categoryForm.IsTab != nil {
		request.Name = categoryForm.Name
		request.IsTab = *categoryForm.IsTab
	}
	response, err := global.GoodsClient.UpdateCategory(context.WithValue(context.Background(), "ginContext", ctx), request)
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"id":     response.Id,
		"name":   response.Name,
		"is_tab": response.IsTab,
		"level":  response.Level,
	})
	entry.Exit()
}
