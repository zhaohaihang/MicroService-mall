package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"github.com/zhaohaihang/goods_service/global"
	"github.com/zhaohaihang/goods_service/model"
	"github.com/zhaohaihang/goods_service/proto"
	"go.uber.org/zap"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GetAllCategoriesList 获取目录列表
func (g *GoodsServer) GetAllCategoriesList(ctx context.Context, request *emptypb.Empty) (*proto.CategoryListResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "GetAllCategoriesList", "request", request)
	parentSpan := opentracing.SpanFromContext(ctx)
	getAllCategoriesListSpan := opentracing.GlobalTracer().StartSpan("GetAllCategoriesList", opentracing.ChildOf(parentSpan.Context()))
	response := &proto.CategoryListResponse{}

	var categorys []model.Category
	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys)
	getAllCategoriesListSpan.Finish()
	b, err := json.Marshal(&categorys)
	if err != nil {
		zap.S().Errorw("json转换failed", "err", err.Error())
		return nil, err
	}
	response.JsonData = string(b)
	return response, nil
}

// GetSubCategory 获取二级目录
func (g *GoodsServer) GetSubCategory(ctx context.Context, request *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "GetSubCategory", "request", request)
	parentSpan := opentracing.SpanFromContext(ctx)
	getSubCategorySpan := opentracing.GlobalTracer().StartSpan("GetSubCategory", opentracing.ChildOf(parentSpan.Context()))
	response := &proto.SubCategoryListResponse{}

	var category model.Category
	result := global.DB.First(&category, request.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	response.Info = &proto.CategoryInfoResponse{
		Id:             int32(category.ID),
		Name:           category.Name,
		ParentCategory: category.ParentCategoryID,
		Level:          category.Level,
		IsTab:          category.IsTab,
	}

	var subCategorys []model.Category
	var subCategorysResponse []*proto.CategoryInfoResponse
	global.DB.Where(&model.Category{ParentCategoryID: request.Id}).Find(&subCategorys)
	getSubCategorySpan.Finish()
	for _, subCategory := range subCategorys {
		subCategorysResponse = append(subCategorysResponse, &proto.CategoryInfoResponse{
			Id:             int32(subCategory.ID),
			Name:           subCategory.Name,
			ParentCategory: subCategory.ParentCategoryID,
			Level:          subCategory.Level,
			IsTab:          subCategory.IsTab,
		})
	}
	response.SubCategory = subCategorysResponse
	return response, nil
}

// CreateCategory 创建目录
func (g *GoodsServer) CreateCategory(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "CreateCategory", "request", request)
	parentSpan := opentracing.SpanFromContext(ctx)
	createCategorySpan := opentracing.GlobalTracer().StartSpan("CreateCategory", opentracing.ChildOf(parentSpan.Context()))
	response := &proto.CategoryInfoResponse{}

	var category model.Category
	category.Name = request.Name
	if request.Level != 1 {
		category.ParentCategoryID = request.ParentCategory
	}
	result := global.DB.Create(&category)
	if result.RowsAffected == 0 {
		fmt.Println(result.Error)
		zap.S().Errorw("创建目录失败", "err", result.Error)
		return nil, result.Error
	}
	createCategorySpan.Finish()
	response.Id = int32(category.ID)
	response.IsTab = category.IsTab
	response.Level = category.Level
	response.Name = category.Name
	response.ParentCategory = category.ParentCategoryID

	return response, nil
}

// DeleteCategory  删除目录
func (g *GoodsServer) DeleteCategory(ctx context.Context, request *proto.DeleteCategoryRequest) (*proto.OperationResult, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "DeleteCategory", "request", request)
	parentSpan := opentracing.SpanFromContext(ctx)
	deleteCategorySpan := opentracing.GlobalTracer().StartSpan("DeleteCategory", opentracing.ChildOf(parentSpan.Context()))
	response := &proto.OperationResult{}

	result := global.DB.Delete(&model.Category{}, request.Id)
	if result.RowsAffected == 0 {
		response.Success = false
		return response, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	deleteCategorySpan.Finish()
	response.Success = true
	return response, nil
}

// UpdateCategory 更新目录信息
func (g *GoodsServer) UpdateCategory(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "UpdateCategory", "request", request)
	parentSpan := opentracing.SpanFromContext(ctx)
	UpdateCategorySpan := opentracing.GlobalTracer().StartSpan("UpdateCategory", opentracing.ChildOf(parentSpan.Context()))
	response := &proto.CategoryInfoResponse{}
	var category model.Category
	result := global.DB.First(&category, request.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品")
	}
	if request.Name != "" {
		category.Name = request.Name
	}
	if request.Level != 0 {
		category.Level = request.Level
	}
	if request.ParentCategory != 0 {
		category.ParentCategoryID = request.ParentCategory
	}
	if request.IsTab {
		category.IsTab = request.IsTab
	}
	result = global.DB.Save(&category)
	if result.Error != nil {
		zap.S().Errorw("更新目录失败", "err", result.Error)
		fmt.Println(result.Error)
		return nil, result.Error
	}
	UpdateCategorySpan.Finish()
	response.Id = int32(category.ID)
	response.Name = category.Name
	response.Level = category.Level
	response.IsTab = category.IsTab
	return response, nil
}
