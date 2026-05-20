package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/rechargeorderdto"
	"xr-game-server/module/recharge"
)

const (
	RechargeOrderAppUrl = "/rechargeOrder"
)

type RechargeOrderAppController struct{}

func initRechargeOrderAppController() {
	httpserver.RegAPI(RechargeOrderAppUrl, &RechargeOrderAppController{})
}

// CreateRechargeOrder App端创建充值订单(按 cfgId 取价格与金币,订单初始为待支付)
func (c *RechargeOrderAppController) CreateRechargeOrder(ctx context.Context, req *rechargeorderdto.AppCreateRechargeOrderReq) (res *rechargeorderdto.AppCreateRechargeOrderRes, err error) {
	return recharge.CreateOrder(ctx, req)
}

// MyRechargeOrderList App端查询本人充值订单分页列表
func (c *RechargeOrderAppController) MyRechargeOrderList(ctx context.Context, req *rechargeorderdto.AppMyRechargeOrderListReq) (res *rechargeorderdto.AppMyRechargeOrderListRes, err error) {
	return recharge.GetMyOrderList(ctx, req)
}
