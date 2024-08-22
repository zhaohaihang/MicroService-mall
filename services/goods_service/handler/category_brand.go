package handler

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/zhaohaihang/goods_service/global"
	"github.com/zhaohaihang/goods_service/model"
	"github.com/zhaohaihang/goods_service/proto"
	"github.com/zhaohaihang/goods_service/utils"
	"go.uber.org/zap"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CategoryBrandList 获取分类的品牌列表
func (g *GoodsServer) CategoryBrandList(ctx context.Context, request *proto.CategoryBrandFilterRequest) (*proto.CategoryBrandListResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "CategoryBrandList", "request", request)
	
	parentSpan := opentracing.SpanFromContext(ctx)
	CategoryBrandList := opentracing.GlobalTracer().StartSpan("CategoryBrandList", opentracing.ChildOf(parentSpan.Context()))
	defer CategoryBrandList.Finish()

	response := &proto.CategoryBrandListResponse{}
	// 查询总数
	var total int64
	global.DB.Find(&model.GoodsCategoryBrand{}).Count(&total)
	response.Total = int32(total)

	// 连表分页查询
	var categoryBrands []model.GoodsCategoryBrand
	global.DB.Preload("Category").Preload("Brand").Scopes(utils.Paginate(int(request.Pages), int(request.PagePerNums))).Find(&categoryBrands)
	var categroyBrandsResponse []*proto.CategoryBrandResponse
	for _, categoryBrand := range categoryBrands {
		categroyBrandsResponse = append(categroyBrandsResponse, &proto.CategoryBrandResponse{
			Category: &proto.CategoryInfoResponse{
				Id:             int32(categoryBrand.Category.ID),
				Name:           categoryBrand.Category.Name,
				ParentCategory: categoryBrand.Category.ParentCategoryID,
				Level:          categoryBrand.Category.Level,
				IsTab:          categoryBrand.Category.IsTab,
			},
			Brand: &proto.BrandInfoResponse{
				Id:   int32(categoryBrand.Brand.ID),
				Name: categoryBrand.Brand.Name,
				Logo: categoryBrand.Brand.Logo,
			},
		})
	}
	response.Data = categroyBrandsResponse
	return response, nil
}
 
// GetCategoryBrandList 获取某一个分类的品牌列表
func (g GoodsServer) GetCategoryBrandList(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.BrandListResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "GetCategoryBrandList", "request", request)
	
	parentSpan := opentracing.SpanFromContext(ctx)
	GetCategoryBrandListSpan := opentracing.GlobalTracer().StartSpan("GetCategoryBrandList", opentracing.ChildOf(parentSpan.Context()))
	defer GetCategoryBrandListSpan.Finish()

	response := &proto.BrandListResponse{} // 查询该商品分类是否存在
	var category model.Category
	result := global.DB.Find(&category, request.Id).First(&category)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "category is not exists")
	}
	// 返回该分类下的 所有品牌
	var categoryBrands []model.GoodsCategoryBrand
	result = global.DB.Preload("Brands").Where(&model.GoodsCategoryBrand{CategoryID: request.Id}).Find(&categoryBrands)
	if result.RowsAffected > 0 {
		response.Total = int32(result.RowsAffected)
	}
	var brandInfoResponse []*proto.BrandInfoResponse
	for _, categoryBrand := range categoryBrands {
		brandInfoResponse = append(brandInfoResponse, &proto.BrandInfoResponse{
			Id:   int32(categoryBrand.Brand.ID),
			Name: categoryBrand.Brand.Name,
			Logo: categoryBrand.Brand.Logo,
		})
	}
	response.Data = brandInfoResponse
	return response, nil
}

// CreateCategoryBrand 创建目录下的品牌
func (g *GoodsServer) CreateCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "CreateCategoryBrand", "request", request)
	
	parentSpan := opentracing.SpanFromContext(ctx)
	CreateCategoryBrandSpan := opentracing.GlobalTracer().StartSpan("CreateCategoryBrand", opentracing.ChildOf(parentSpan.Context()))
	defer CreateCategoryBrandSpan.Finish()

	var category model.Category
	result := global.DB.First(&category, request.CategoryId)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brand
	brandResult := global.DB.First(&brand, request.BrandId)
	if brandResult.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "brand is not exists")
	}

	categoryBrand := model.GoodsCategoryBrand{
		CategoryID: request.CategoryId,
		BrandID:    request.BrandId,
	}
	
	global.DB.Save(&categoryBrand)
	response := &proto.CategoryBrandResponse{
		Id: int32(categoryBrand.ID),
		Category: &proto.CategoryInfoResponse{
			Id:             int32(categoryBrand.Category.ID),
			Name:           categoryBrand.Category.Name,
			ParentCategory: categoryBrand.Category.ParentCategoryID,
			Level:          categoryBrand.Category.Level,
			IsTab:          categoryBrand.Category.IsTab,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   int32(categoryBrand.Brand.ID),
			Name: categoryBrand.Brand.Name,
			Logo: categoryBrand.Brand.Logo,
		},
	}
	
	return response, nil
}

// DeleteCategoryBrand 删除类别下的品牌
func (g *GoodsServer) DeleteCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*proto.OperationResult, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "DeleteCategoryBrand", "request", request)
	parentSpan := opentracing.SpanFromContext(ctx)
	DeleteCategoryBrandSpan := opentracing.GlobalTracer().StartSpan("DeleteCategoryBrand", opentracing.ChildOf(parentSpan.Context()))
	defer DeleteCategoryBrandSpan.Finish()
	
	response := &proto.OperationResult{}
	result := global.DB.Delete(&model.GoodsCategoryBrand{}, request.Id)
	if result.RowsAffected == 0 {
		response.Success = false
		return response, status.Errorf(codes.InvalidArgument, "brand is not exists")
	}
	
	response.Success = true
	return response, nil
}

// UpdateCategoryBrand 更新目录下的品牌信息
func (g *GoodsServer) UpdateCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*proto.OperationResult, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "UpdateCategoryBrand", "request", request)
	
	parentSpan := opentracing.SpanFromContext(ctx)
	UpdateCategoryBrandSpan := opentracing.GlobalTracer().StartSpan("UpdateCategoryBrand", opentracing.ChildOf(parentSpan.Context()))
	defer UpdateCategoryBrandSpan.Finish()
	
	response := &proto.OperationResult{
		Success: true,
	}
	var categoryBrand model.GoodsCategoryBrand
	result := global.DB.First(&categoryBrand, request.Id)
	if result.RowsAffected == 0 {
		response.Success = false
		return response, status.Errorf(codes.InvalidArgument, "CategoryBrand is not exists")
	}

	result = global.DB.Find(&model.Category{}, request.CategoryId)
	if result.RowsAffected == 0 {
		response.Success = false
		return response, status.Errorf(codes.InvalidArgument, "Category is not exists")
	}

	result = global.DB.Find(&model.Brand{}, request.BrandId)
	if result.RowsAffected == 0 {
		response.Success = false
		return response, status.Errorf(codes.InvalidArgument, "Brand is not exists")
	}

	categoryBrand.CategoryID = request.CategoryId
	categoryBrand.BrandID = request.BrandId

	global.DB.Save(&categoryBrand)
	return response, nil
}
