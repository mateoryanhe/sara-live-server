package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/activitydto"
	"xr-game-server/module/activity"
)

const FirstRechargeActivityCMSUrl = "/firstRechargeActivity"

type FirstRechargeActivityCMSController struct{}

func initFirstRechargeActivityCMSController() {
	httpserver.RegCMS(FirstRechargeActivityCMSUrl, &FirstRechargeActivityCMSController{})
}

func (c *FirstRechargeActivityCMSController) GetFirstRechargeActivityCfg(ctx context.Context, req *activitydto.GetFirstRechargeActivityCfgReq) (*activitydto.GetFirstRechargeActivityCfgRes, error) {
	return activity.GetFirstRechargeActivityCfg(ctx, req)
}

func (c *FirstRechargeActivityCMSController) SaveFirstRechargeActivityCfg(ctx context.Context, req *activitydto.SaveFirstRechargeActivityCfgReq) (*activitydto.SaveFirstRechargeActivityCfgRes, error) {
	return activity.SaveFirstRechargeActivityCfg(ctx, req)
}
