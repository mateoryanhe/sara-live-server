package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/rechargeorderdto"
	"xr-game-server/module/recharge"
)

// RechargeOrderAppPublicController App充值订单公开接口(无需鉴权)
type RechargeOrderAppPublicController struct{}

func initRechargeOrderAppPublicController() {
	httpserver.RegNonAuthAPI(RechargeOrderAppUrl, &RechargeOrderAppPublicController{})
}

// CreateChannelRechargeOrder App渠道充值建单(无需鉴权,userId由App上报)
func (c *RechargeOrderAppPublicController) CreateChannelRechargeOrder(ctx context.Context, req *rechargeorderdto.AppCreateChannelRechargeOrderReq) (res *rechargeorderdto.AppCreateChannelRechargeOrderRes, err error) {
	return recharge.CreateChannelRechargeOrder(ctx, req)
}

// GetChannelPayUserProfile App查询渠道付款人资料(无需鉴权)
func (c *RechargeOrderAppPublicController) GetChannelPayUserProfile(ctx context.Context, req *rechargeorderdto.AppGetChannelPayUserProfileReq) (*rechargeorderdto.AppGetChannelPayUserProfileRes, error) {
	return recharge.GetChannelPayUserProfile(ctx, req)
}

// SaveChannelPayUserProfile App保存渠道付款人资料(无需鉴权)
func (c *RechargeOrderAppPublicController) SaveChannelPayUserProfile(ctx context.Context, req *rechargeorderdto.AppSaveChannelPayUserProfileReq) (*rechargeorderdto.AppSaveChannelPayUserProfileRes, error) {
	return recharge.SaveChannelPayUserProfile(ctx, req)
}
