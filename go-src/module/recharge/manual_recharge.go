package recharge

import (
	"context"
	"strconv"
	"time"
	"xr-game-server/constants/currency"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/rechargeorderdao"
	"xr-game-server/dto/rechargeorderdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

// ManualRecharge 后台手动给玩家充值到账(完成待支付订单并发放金币)
func ManualRecharge(ctx context.Context, req *rechargeorderdto.CMSManualRechargeReq) (*rechargeorderdto.CMSManualRechargeRes, error) {
	orderId, err := strconv.ParseUint(req.OrderId, 10, 64)
	if err != nil || orderId == 0 {
		return nil, errercode.CreateCode(errercode.RechargeOrderNonExist)
	}
	order := rechargeorderdao.GetById(orderId)
	if order == nil {
		return nil, errercode.CreateCode(errercode.RechargeOrderNonExist)
	}

	operatorId := httpserver.GetAuthId(ctx)
	order.SetOperatorId(operatorId)

	after, err := completeOrder(order, currency.ReasonGmAdjust)
	if err != nil {
		// 发放失败,标记订单为已取消,避免遗留"待支付"脏数据
		order.SetStatus(entity.RechargeOrderStatusCancelled)
		order.SetUpdatedAt(time.Now())
		rechargeorderdao.FlushOrderCache(order)
		return nil, err
	}

	return &rechargeorderdto.CMSManualRechargeRes{
		OrderId: strconv.FormatUint(order.ID, 10),
		Gold:    0,
		After:   after,
		Success: true,
	}, nil
}
