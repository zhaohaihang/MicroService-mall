package shop_cart

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"github.com/zhaohaihang/order_api/forms"
	"github.com/zhaohaihang/order_api/global"
	"github.com/zhaohaihang/order_api/proto"
	"github.com/zhaohaihang/order_api/utils"
	"strconv"
)

// ListGoods 获取购物车内商品列表
func ListGoodsInCart(ctx *gin.Context) {
	entry, blockError := utils.SentinelEntry(ctx)
	if blockError != nil {
		zap.S().Errorw("Error", "message", "Request too frequent")
		utils.HandleRequestFrequentError(ctx)
		return
	}
	// 获取购物车中的商品id
	userId, _ := ctx.Get("userId")
	response, err := global.OrderClient.CartItemList(context.WithValue(context.Background(), "ginContext", ctx), &proto.UserInfo{Id: int32(userId.(uint))})
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())
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
	goodsListResponse, err := global.GoodsClient.BatchGetGoods(context.WithValue(context.Background(), "ginContext", ctx), &proto.BatchGoodsIdInfo{Id: ids})
	if err != nil {
		zap.S().Errorw("Error", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	reMap := gin.H{
		"total": response.Total,
	}
	goodsList := make([]interface{}, 0)
	for _, item := range response.Data { // 遍历购物车 获取商品ID
		for _, good := range goodsListResponse.Data { // 遍历商品 获得商品的详细信息 将两者进行组装
			if good.Id == item.GoodsId {
				tmpMap := map[string]interface{}{}
				tmpMap["id"] = item.Id
				tmpMap["goods_id"] = item.GoodsId
				tmpMap["good_name"] = good.Name
				tmpMap["good_image"] = good.GoodsFrontImage
				tmpMap["good_price"] = good.ShopPrice
				tmpMap["nums"] = item.Nums
				tmpMap["checked"] = item.Checked

				goodsList = append(goodsList, tmpMap)
			}
		}
	}
	reMap["data"] = goodsList
	ctx.JSON(http.StatusOK, reMap)
	entry.Exit()
}

// AddGoodsToCart 添加商品到购物车
func AddGoodsToCart(ctx *gin.Context) {
	entry, blockError := utils.SentinelEntry(ctx)
	if blockError != nil {
		zap.S().Errorw("Error", "message", "Request too frequent")
		utils.HandleRequestFrequentError(ctx)
		return
	}

	// 绑定参数
	itemForm := forms.ShopCartItemForm{}
	err := ctx.ShouldBind(&itemForm)
	if err != nil {
		zap.S().Errorw("Error", "message", "cart bind form failed", "err", err.Error())
		utils.HandleValidatorError(ctx, err)
		return
	}

	// 检查商品是否存在
	_, err = global.GoodsClient.GetGoodsDetail(context.WithValue(context.Background(), "ginContext", ctx), &proto.GoodsInfoRequest{
		Id: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("Error")
		utils.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	// 校验商品库存
	inventoryResp, err := global.InventoryClient.InvDetail(context.WithValue(context.Background(), "ginContext", ctx), &proto.GoodsInvInfo{
		GoodsId: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("Error", "message", "inventory is not exists", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	if inventoryResp.Num < itemForm.Nums {
		ctx.JSON(http.StatusOK, gin.H{
			"nums": "inventory less than nums",
		})
		return
	}

	userId, _ := ctx.Get("userId")
	orderResp, err := global.OrderClient.CreateCartItem(context.WithValue(context.Background(), "ginContext", ctx), &proto.CartItemRequest{
		GoodsId: itemForm.GoodsId,
		UserId:  int32(userId.(uint)),
		Nums:    itemForm.Nums,
	})
	if err != nil {
		zap.S().Errorw("Error", "message", "add to cart failed", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": orderResp.Id,
	})
	entry.Exit()
}
 
// Update 更新购物车商品信息
func UpdateGoodsInCart(ctx *gin.Context) {
	entry, blockError := utils.SentinelEntry(ctx)
	if blockError != nil {
		zap.S().Errorw("Error", "message", "Request too frequent")
		utils.HandleRequestFrequentError(ctx)
		return
	} 

	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		zap.S().Errorw("Error", "message", "param bind failed", "err", err.Error())
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "url param error",
		})
		return
	}

	itemForm := forms.ShopCartItemUpdateForm{}
	err = ctx.ShouldBind(&itemForm)
	if err != nil {
		zap.S().Errorw("Error", "message", "updateForm  bind failed", "err", err.Error())
		utils.HandleValidatorError(ctx, err)
		return
	}

	// 更新商品数量 和 选中状态
	userId, _ := ctx.Get("userId")
	_, err = global.OrderClient.UpdateCartItem(context.WithValue(context.Background(), "ginContext", ctx), &proto.CartItemRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(idInt),
		Nums:    itemForm.Nums,
		Checked: *itemForm.Checked,
	})
	if err != nil {
		zap.S().Errorw("Error", "message", "update failed", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	entry.Exit()
}

// Delete 删除购物车内商品
func DeleteGoodsInCart(ctx *gin.Context) {
	entry, blockError := utils.SentinelEntry(ctx)
	if blockError != nil {
		zap.S().Errorw("Error", "message", "Request too frequent")
		utils.HandleRequestFrequentError(ctx)
		return
	}

	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		zap.S().Errorw("Error", "message", "param bind failed", "err", err.Error())
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "url param error",
		})
		return
	}

	userId, _ := ctx.Get("userId")
	_, err = global.OrderClient.DeleteCartItem(context.WithValue(context.Background(), "ginContext", ctx), &proto.CartItemRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(idInt),
	})
	if err != nil {
		zap.S().Errorw("Error", "message", "delete cart failed", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
	entry.Exit()
}
