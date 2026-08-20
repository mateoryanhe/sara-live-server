package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/activitydto"
	"xr-game-server/module/activity"
)

const FirstRechargeActivityAppUrl = "/firstRechargeActivity"

type FirstRechargeActivityAppController struct{}

func initFirstRechargeActivityAppController() {
	httpserver.RegAPI(FirstRechargeActivityAppUrl, &FirstRechargeActivityAppController{})
}

func (c *FirstRechargeActivityAppController) FirstRechargeActivityCfgForApp(ctx context.Context, req *activitydto.AppFirstRechargeActivityCfgReq) (*activitydto.AppFirstRechargeActivityCfgRes, error) {
	return activity.GetAppFirstRechargeActivityCfg(ctx, req)
}
