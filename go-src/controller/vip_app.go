package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/vipdto"
	"xr-game-server/module/vip"
)

const VipAppUrl = "/vip"

type VipAppController struct{}

func initVipAppController() {
	httpserver.RegAPI(VipAppUrl, &VipAppController{})
}

// GetVipDetail App端查询当前用户VIP详情
func (c *VipAppController) GetVipDetail(ctx context.Context, req *vipdto.AppVipDetailReq) (*vipdto.AppVipDetailRes, error) {
	return vip.GetAppVipDetail(ctx, req)
}
