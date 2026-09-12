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

// CreateChannelRechargeOrder App渠道充值建单(需登录)
func (c *RechargeOrderAppController) CreateChannelRechargeOrder(ctx context.Context, req *rechargeorderdto.AppCreateChannelRechargeOrderReq) (res *rechargeorderdto.AppCreateChannelRechargeOrderRes, err error) {
	return recharge.CreateChannelRechargeOrder(ctx, req)
}

// CreateCoinMerchantChannelRechargeOrder 币商App渠道建单(档位来自币商充值配置)
func (c *RechargeOrderAppController) CreateCoinMerchantChannelRechargeOrder(ctx context.Context, req *rechargeorderdto.AppCreateCoinMerchantChannelRechargeOrderReq) (*rechargeorderdto.AppCreateChannelRechargeOrderRes, error) {
	return recharge.CreateCoinMerchantChannelRechargeOrder(ctx, req)
}

// GetChannelPayUserProfile App查询渠道付款人资料(需登录)
func (c *RechargeOrderAppController) GetChannelPayUserProfile(ctx context.Context, req *rechargeorderdto.AppGetChannelPayUserProfileReq) (*rechargeorderdto.AppGetChannelPayUserProfileRes, error) {
	return recharge.GetChannelPayUserProfile(ctx, req)
}

// SaveChannelPayUserProfile App保存渠道付款人资料(需登录)
func (c *RechargeOrderAppController) SaveChannelPayUserProfile(ctx context.Context, req *rechargeorderdto.AppSaveChannelPayUserProfileReq) (*rechargeorderdto.AppSaveChannelPayUserProfileRes, error) {
	return recharge.SaveChannelPayUserProfile(ctx, req)
}

// MyRechargeOrderList App端查询本人充值订单分页列表
func (c *RechargeOrderAppController) MyRechargeOrderList(ctx context.Context, req *rechargeorderdto.AppMyRechargeOrderListReq) (res *rechargeorderdto.AppMyRechargeOrderListRes, err error) {
	return recharge.GetMyOrderList(ctx, req)
}

// CheckRechargeOrderSuccess App端查询订单是否充值成功
func (c *RechargeOrderAppController) CheckRechargeOrderSuccess(ctx context.Context, req *rechargeorderdto.AppCheckRechargeOrderSuccessReq) (res *rechargeorderdto.AppCheckRechargeOrderSuccessRes, err error) {
	return recharge.CheckOrderRechargeSuccess(ctx, req)
}
