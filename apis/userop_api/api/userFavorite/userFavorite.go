package userFavorite

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhaohaihang/userop_api/forms"
	"github.com/zhaohaihang/userop_api/global"
	"github.com/zhaohaihang/userop_api/proto"
	"github.com/zhaohaihang/userop_api/utils"
	"go.uber.org/zap"
)

func ListFavorite(ctx *gin.Context) {
	entry, blockError := utils.SentinelEntry(ctx)
	if blockError != nil {
		zap.S().Errorw("Error", "message", "Request too frequent")
		utils.HandleRequestFrequentError(ctx)
		return
	}

	// 获取用户的收藏的商品列表
	userId, _ := ctx.Get("userId")
	response, err := global.UserFavoriteClient.GetFavoriteList(context.WithValue(context.Background(), "ginContext", ctx), &proto.UserFavoriteRequest{
		UserId: int32(userId.(uint)),
	})
	if err != nil {
		zap.S().Errorw("Error", "message", "get favorite list failed", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	ids := make([]int32, 0)
	for _, item := range response.Data {
		ids = append(ids, item.GoodsId)
	}
	if len(ids) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"total": 0,
		})
		return
	}

	// 获取商品详细信息
	goodsResponse, err := global.GoodsClient.BatchGetGoods(context.WithValue(context.Background(), "ginContext", ctx), &proto.BatchGoodsIdInfo{Id: ids})
	if err != nil {
		zap.S().Errorw("Error", "message", "get goods info failed", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	responseMap := gin.H{
		"total": goodsResponse.Total,
	}
	goodsList := make([]interface{}, 0)
	for _, item := range response.Data {
		data := gin.H{
			"id": item.GoodsId,
		}

		for _, goods := range goodsResponse.Data {
			if item.GoodsId == goods.Id {
				data["name"] = goods.Name
				data["shop_price"] = goods.ShopPrice
			}
		}
		goodsList = append(goodsList, data)
	}
	responseMap["data"] = goodsList
	ctx.JSON(http.StatusOK, responseMap)
	entry.Exit()
}

func CreateFavorite(ctx *gin.Context) {
	entry, blockError := utils.SentinelEntry(ctx)
	if blockError != nil {
		zap.S().Errorw("Error", "message", "Request too frequent")
		utils.HandleRequestFrequentError(ctx)
		return
	}

	userFavFrom := forms.UserFavForm{}
	err := ctx.ShouldBind(&userFavFrom)
	if err != nil {
		zap.S().Errorw("Error", "message", "favorite bind failed", "err", err.Error())
		utils.HandleValidatorError(ctx, err)
		return
	}

	userId, _ := ctx.Get("userId")
	_, err = global.UserFavoriteClient.AddUserFavorite(context.WithValue(context.Background(), "ginContext", ctx), &proto.UserFavoriteRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: userFavFrom.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("Error", "message", "add fauorite failed", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
	entry.Exit()
}

func DeleteFavorite(ctx *gin.Context) {
	entry, blockError := utils.SentinelEntry(ctx)
	if blockError != nil {
		zap.S().Errorw("Error", "message", "Request too frequent")
		utils.HandleRequestFrequentError(ctx)
		return
	}

	id := ctx.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	userId, _ := ctx.Get("userId")
	_, err = global.UserFavoriteClient.DeleteUserFavorite(context.WithValue(context.Background(), "ginContext", ctx), &proto.UserFavoriteRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(idInt),
	})
	if err != nil {
		zap.S().Errorw("Error", "message", "delete favorite failed", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
	entry.Exit()
}

func DetailFavorite(ctx *gin.Context) {
	entry, blockError := utils.SentinelEntry(ctx)
	if blockError != nil {
		zap.S().Errorw("Error", "message", "Request too frequent")
		utils.HandleRequestFrequentError(ctx)
		return
	}

	goodsId := ctx.Param("id")
	goodsIdInt, err := strconv.ParseInt(goodsId, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	userId, _ := ctx.Get("userId")
	_, err = global.UserFavoriteClient.GetUserFavoriteDetail(context.WithValue(context.Background(), "ginContext", ctx), &proto.UserFavoriteRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(goodsIdInt),
	})
	if err != nil {
		zap.S().Errorw("Error", "message", "get UserFavoriteDetail failed", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
	entry.Exit()
}
