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

// BrandList 获取品牌列表
func (g *GoodsServer) BrandList(ctx context.Context, request *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "BrandList", "request", request)
	
	parentSpan := opentracing.SpanFromContext(ctx)
	brandListSpan := opentracing.GlobalTracer().StartSpan("BrandList", opentracing.ChildOf(parentSpan.Context()))
	response := &proto.BrandListResponse{}

	var brands []model.Brand
	result := global.DB.Scopes(utils.Paginate(int(request.Pages), int(request.PagePerNums))).Find(&brands)
	if result.Error != nil {
		return nil, result.Error
	}
	var total int64
	global.DB.Model(&model.Brand{}).Count(&total)
	response.Total = int32(total)
	var brandList []*proto.BrandInfoResponse
	for _, brand := range brands {
		brandResponse := proto.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Name,
		}
		brandList = append(brandList, &brandResponse)
	}
	response.Data = brandList
	
	brandListSpan.Finish()
	return response, nil
}

// CreateBrand 创建品牌
func (g *GoodsServer) CreateBrand(ctx context.Context, request *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "CreateBrand", "request", request)
	
	parentSpan := opentracing.SpanFromContext(ctx)
	createBrand := opentracing.GlobalTracer().StartSpan("CreateBrand", opentracing.ChildOf(parentSpan.Context()))
	
	result := global.DB.Where("name=?", request.Name).First(&model.Brand{})
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.InvalidArgument, "brand has exists")
	}
	brand := model.Brand{
		Name: request.Name,
		Logo: request.Logo,
	}
	zap.S().Infof("create brand %#v", brand)
	
	global.DB.Create(&brand)
	createBrand.Finish()
	response := &proto.BrandInfoResponse{
		Id:   brand.ID,
		Name: brand.Name,
		Logo: brand.Logo,
	}
	
	return response, nil
}

// DeleteBrand 删除品牌
func (g *GoodsServer) DeleteBrand(ctx context.Context, request *proto.BrandRequest) (*proto.OperationResult, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "DeleteBrand", "request", request)
	parentSpan := opentracing.SpanFromContext(ctx)
	deleteBrandSPan := opentracing.GlobalTracer().StartSpan("DeleteBrand", opentracing.ChildOf(parentSpan.Context()))
	response := &proto.OperationResult{}

	var brand model.Brand
	result := global.DB.Delete(&brand, request.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "brand is not exists")
	}
	deleteBrandSPan.Context()
	response.Success = true
	return response, nil
}

// UpdateBrand 更新品牌信息
func (g *GoodsServer) UpdateBrand(ctx context.Context, request *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "UpdateBrand", "request", request)
	parentSpan := opentracing.SpanFromContext(ctx)
	updateBrandSpan := opentracing.GlobalTracer().StartSpan("UpdateBrand", opentracing.ChildOf(parentSpan.Context()))
	response := &proto.BrandInfoResponse{}

	var brand model.Brand
	result := global.DB.First(&brand, request.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "brand is not exists")
	}
	if request.Name != "" {
		brand.Name = request.Name
	}
	if request.Logo != "" {
		brand.Logo = request.Logo
	}
	global.DB.Save(&brand)
	updateBrandSpan.Finish()
	response.Id = brand.ID
	response.Name = brand.Name
	response.Logo = brand.Logo
	return response, nil
}
