package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/vipcfgdto"
	"xr-game-server/module/vipcfg"
)

const VipCfgAppUrl = "/vipCfg"

type VipCfgAppController struct{}

func initVipCfgAppController() {
	httpserver.RegAPI(VipCfgAppUrl, &VipCfgAppController{})
}

// GetVipCfgByLevel App端按等级查询VIP配置
func (c *VipCfgAppController) GetVipCfgByLevel(ctx context.Context, req *vipcfgdto.AppVipCfgByLevelReq) (*vipcfgdto.AppVipCfgByLevelRes, error) {
	return vipcfg.GetAppVipCfgByLevel(ctx, req)
}

// VipCfgListForApp App端查询全部VIP配置
func (c *VipCfgAppController) VipCfgListForApp(ctx context.Context, req *vipcfgdto.AppVipCfgListReq) (*vipcfgdto.AppVipCfgListRes, error) {
	return vipcfg.GetAppVipCfgList(ctx, req)
}
