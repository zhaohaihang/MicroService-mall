package handler

import (
	"context"
	"encoding/json"

	"github.com/opentracing/opentracing-go"
	"github.com/zhaohaihang/goods_service/global"
	"github.com/zhaohaihang/goods_service/model"
	"github.com/zhaohaihang/goods_service/proto"
	"go.uber.org/zap"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GetAllCategoriesList 获取分类列表
func (g *GoodsServer) GetAllCategoriesList(ctx context.Context, request *emptypb.Empty) (*proto.CategoryListResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "GetAllCategoriesList", "request", request)
	parentSpan := opentracing.SpanFromContext(ctx)
	getAllCategoriesListSpan := opentracing.GlobalTracer().StartSpan("GetAllCategoriesList", opentracing.ChildOf(parentSpan.Context()))
	defer getAllCategoriesListSpan.Finish()

	var categorys []model.Category
	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys)

	b, err := json.Marshal(&categorys)
	if err != nil {
		zap.S().Errorw("json Marshal failed", "err", err.Error())
		return nil, status.Errorf(codes.Internal, "category is not exists")
	}
	response := &proto.CategoryListResponse{
		JsonData: string(b),
	}
	return response, nil
}

// GetSubCategory 获取二级目录
func (g *GoodsServer) GetSubCategory(ctx context.Context, request *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "GetSubCategory", "request", request)

	parentSpan := opentracing.SpanFromContext(ctx)
	getSubCategorySpan := opentracing.GlobalTracer().StartSpan("GetSubCategory", opentracing.ChildOf(parentSpan.Context()))
	defer getSubCategorySpan.Finish()

	response := &proto.SubCategoryListResponse{}

	// 先获取一级目录
	var category model.Category
	result := global.DB.First(&category, request.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	response.Info = &proto.CategoryInfoResponse{
		Id:             int32(category.ID),
		Name:           category.Name,
		ParentCategory: int32(category.ParentCategoryID),
		Level:          category.Level,
		IsTab:          category.IsTab,
	}

	// 获取二级目录
	var subCategorys []model.Category
	var subCategorysResponse []*proto.CategoryInfoResponse
	global.DB.Where(&model.Category{ParentCategoryID: uint(request.Id)}).Find(&subCategorys)
	for _, subCategory := range subCategorys {
		subCategorysResponse = append(subCategorysResponse, &proto.CategoryInfoResponse{
			Id:             int32(subCategory.ID),
			Name:           subCategory.Name,
			ParentCategory: int32(subCategory.ParentCategoryID),
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
	defer createCategorySpan.Finish()

	var category model.Category
	category.Name = request.Name
	if request.Level != 1 {
		category.ParentCategoryID = uint(request.ParentCategory)
	}
	result := global.DB.Create(&category)
	if result.RowsAffected == 0 {
		zap.S().Errorw("create Category error ", "err", result.Error)
		return nil, result.Error
	}

	response := &proto.CategoryInfoResponse{
		Id:             int32(category.ID),
		IsTab:          category.IsTab,
		Level:          category.Level,
		Name:           category.Name,
		ParentCategory: int32(category.ParentCategoryID),
	}
	return response, nil
}

// DeleteCategory  删除目录
func (g *GoodsServer) DeleteCategory(ctx context.Context, request *proto.DeleteCategoryRequest) (*proto.OperationResult, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "DeleteCategory", "request", request)
	parentSpan := opentracing.SpanFromContext(ctx)
	deleteCategorySpan := opentracing.GlobalTracer().StartSpan("DeleteCategory", opentracing.ChildOf(parentSpan.Context()))
	defer deleteCategorySpan.Finish()

	response := &proto.OperationResult{}

	result := global.DB.Delete(&model.Category{}, request.Id)
	if result.RowsAffected == 0 {
		response.Success = false
		return response, status.Errorf(codes.NotFound, "Category is not exists")
	}

	response.Success = true
	return response, nil
}

// UpdateCategory 更新目录信息
func (g *GoodsServer) UpdateCategory(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "UpdateCategory", "request", request)
	parentSpan := opentracing.SpanFromContext(ctx)
	UpdateCategorySpan := opentracing.GlobalTracer().StartSpan("UpdateCategory", opentracing.ChildOf(parentSpan.Context()))
	defer UpdateCategorySpan.Finish()

	var category model.Category
	result := global.DB.First(&category, request.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "category is not exists")
	}
	if request.Name != "" {
		category.Name = request.Name
	}
	if request.Level != 0 {
		category.Level = request.Level
	}
	if request.ParentCategory != 0 {
		category.ParentCategoryID = uint(request.ParentCategory)
	}
	if request.IsTab {
		category.IsTab = request.IsTab
	}
	result = global.DB.Save(&category)
	if result.Error != nil {
		zap.S().Errorw("update category failed", "err", result.Error)
		return nil, result.Error
	}

	response := &proto.CategoryInfoResponse{
		Id:             int32(category.ID),
		IsTab:          category.IsTab,
		Level:          category.Level,
		Name:           category.Name,
		ParentCategory: int32(category.ParentCategoryID),
	}
	return response, nil
}
