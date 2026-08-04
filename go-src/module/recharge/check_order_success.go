package recharge

import (
	"context"
	"strconv"
	"strings"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/rechargeorderdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/rechargeorderdto"
	"xr-game-server/entity"
)

// CheckOrderRechargeSuccess App 端查询订单是否充值成功(优先缓存,未命中查库)
func CheckOrderRechargeSuccess(ctx context.Context, req *rechargeorderdto.AppCheckRechargeOrderSuccessReq) (*rechargeorderdto.AppCheckRechargeOrderSuccessRes, error) {
	userId := httpserver.GetAuthId(ctx)
	orderId, _ := strconv.ParseUint(strings.TrimSpace(req.OrderId), 10, 64)
	if orderId == 0 {
		return &rechargeorderdto.AppCheckRechargeOrderSuccessRes{Success: false}, nil
	}
	order := rechargeorderdao.GetById(orderId)
	if order == nil || order.UserId != userId {
		return &rechargeorderdto.AppCheckRechargeOrderSuccessRes{Success: false}, nil
	}
	res := &rechargeorderdto.AppCheckRechargeOrderSuccessRes{
		Success: order.Status == entity.RechargeOrderStatusCompleted,
		Price:   order.Price,
		Gold:    order.Gold,
	}
	if userInfo := userinfodao.GetUserInfoByUserId(userId); userInfo != nil {
		res.GoldBalance = userInfo.Gold
	}
	return res, nil
}
