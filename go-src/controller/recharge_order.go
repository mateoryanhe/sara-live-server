package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/rechargeorderdto"
	"xr-game-server/module/recharge"
)

const (
	RechargeOrderUrl = "/rechargeOrder"
)

type RechargeOrderController struct{}

func initRechargeOrderController() {
	httpserver.RegCMS(RechargeOrderUrl, &RechargeOrderController{})
}

// RechargeOrderList 分页查询充值订单
func (c *RechargeOrderController) RechargeOrderList(ctx context.Context, req *rechargeorderdto.CMSRechargeOrderListReq) (res *httpserver.CMSQueryResp, err error) {
	return recharge.GetCMSList(ctx, req)
}

// ManualRecharge 后台手动给玩家充值到账
func (c *RechargeOrderController) ManualRecharge(ctx context.Context, req *rechargeorderdto.CMSManualRechargeReq) (res *rechargeorderdto.CMSManualRechargeRes, err error) {
	return recharge.ManualRecharge(ctx, req)
}

// ManualCreateOrder 后台人工创建充值订单
func (c *RechargeOrderController) ManualCreateOrder(ctx context.Context, req *rechargeorderdto.CMSCreateRechargeOrderReq) (res *rechargeorderdto.CMSCreateRechargeOrderRes, err error) {
	return recharge.CMSCreateOrder(ctx, req)
}
