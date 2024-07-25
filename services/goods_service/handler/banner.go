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
	"google.golang.org/protobuf/types/known/emptypb"
)

// BannerList 返回轮播图列表
func (g *GoodsServer) BannerList(ctx context.Context, request *emptypb.Empty) (*proto.BannerListResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "BannerList", "request", request)
	
	parentSpan := opentracing.SpanFromContext(ctx)
	bannerListSpan := opentracing.GlobalTracer().StartSpan("BannerList", opentracing.ChildOf(parentSpan.Context()))
	banners := []model.Banner{}
	result := global.DB.Find(&banners)
	bannerListSpan.Finish()
	
	bannerListResponse := &proto.BannerListResponse{}
	bannerListResponse.Total = int32(result.RowsAffected)
	for _, banner := range banners {
		bannerListResponse.Data = append(bannerListResponse.Data, utils.BannerToBannerResponse(&banner))
	}
	
	return bannerListResponse, nil
}

// CreateBanner 创建轮播图
func (g *GoodsServer) CreateBanner(ctx context.Context, request *proto.BannerRequest) (*proto.BannerResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "CreateBanner", "request", request)
	
	banner :=  model.Banner{
		Image: request.Image,
		Index: request.Index,
		Url: request.Url,
	}

	parentSpan := opentracing.SpanFromContext(ctx)
	createBannerSpan := opentracing.GlobalTracer().StartSpan("CreateBanner", opentracing.ChildOf(parentSpan.Context()))
	result := global.DB.Create(&banner)
	if result.Error != nil {
		zap.S().Errorw("create banner failed", "err", result.Error.Error())
		return nil, result.Error
	}
	createBannerSpan.Finish()

	bannerResponse :=utils.BannerToBannerResponse(&banner)
	return bannerResponse, nil
}

// DeleteBanner 删除轮播图
func (g *GoodsServer) DeleteBanner(ctx context.Context, request *proto.BannerRequest) (*proto.OperationResult, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "DeleteBanner", "request", request)
	
	parentSpan := opentracing.SpanFromContext(ctx)
	deleteBannerSpan := opentracing.GlobalTracer().StartSpan("DeleteBanner", opentracing.ChildOf(parentSpan.Context()))
	result := global.DB.Delete(&model.Banner{}, request.Id)
	deleteBannerSpan.Finish()

	response := &proto.OperationResult{
		Success: true,
	}
	if result.RowsAffected == 0 {
		response.Success = false
		return response, status.Errorf(codes.NotFound, "banner is not exists")
	}

	return response, nil
}

// UpdateBanner 更新轮播图
func (g *GoodsServer) UpdateBanner(ctx context.Context, request *proto.BannerRequest) (*proto.BannerResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "UpdateBanner", "request", request)
	
	parentSpan := opentracing.SpanFromContext(ctx)
	updateBannerSpan := opentracing.GlobalTracer().StartSpan("UpdateBanner", opentracing.ChildOf(parentSpan.Context()))

	var banner model.Banner
	result := global.DB.First(&banner, request.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "banner is not exists")
	}
	if request.Url != "" {
		banner.Url = request.Url
	}
	if request.Image != "" {
		banner.Image = request.Image
	}
	if request.Index != 0 {
		banner.Index = request.Index
	}
	result = global.DB.Save(&banner)
	if result.RowsAffected != 1 {
		return nil, result.Error
	}
	updateBannerSpan.Finish()

	bannerResponse :=utils.BannerToBannerResponse(&banner)
	return bannerResponse, nil
}

