package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/activitydto"
	"xr-game-server/module/activity"
)

const InviteRechargeRewardCMSUrl = "/inviteRechargeReward"

type InviteRechargeRewardCMSController struct{}

func initInviteRechargeRewardCMSController() {
	httpserver.RegCMS(InviteRechargeRewardCMSUrl, &InviteRechargeRewardCMSController{})
}

func (c *InviteRechargeRewardCMSController) GetInviteRechargeRewardCfg(ctx context.Context, req *activitydto.GetInviteRechargeRewardCfgReq) (*activitydto.GetInviteRechargeRewardCfgRes, error) {
	return activity.GetInviteRechargeRewardCfg(ctx, req)
}

func (c *InviteRechargeRewardCMSController) SaveInviteRechargeRewardCfg(ctx context.Context, req *activitydto.SaveInviteRechargeRewardCfgReq) (*activitydto.SaveInviteRechargeRewardCfgRes, error) {
	return activity.SaveInviteRechargeRewardCfg(ctx, req)
}
