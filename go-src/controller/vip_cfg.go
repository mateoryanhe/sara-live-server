package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/vipcfgdto"
	"xr-game-server/module/vipcfg"
)

const VipCfgUrl = "/vipCfg"

type VipCfgController struct{}

func initVipCfgController() {
	httpserver.RegCMS(VipCfgUrl, &VipCfgController{})
}

func (c *VipCfgController) VipCfgList(ctx context.Context, req *vipcfgdto.VipCfgListReq) (*httpserver.CMSQueryResp, error) {
	return vipcfg.GetList(ctx, req)
}

func (c *VipCfgController) CreateVipCfg(ctx context.Context, req *vipcfgdto.CreateVipCfgReq) (*vipcfgdto.CreateVipCfgRes, error) {
	return vipcfg.Create(ctx, req)
}

func (c *VipCfgController) UpdateVipCfg(ctx context.Context, req *vipcfgdto.UpdateVipCfgReq) (*vipcfgdto.UpdateVipCfgRes, error) {
	return vipcfg.Update(ctx, req)
}

func (c *VipCfgController) DeleteVipCfg(ctx context.Context, req *vipcfgdto.DeleteVipCfgReq) (*vipcfgdto.DeleteVipCfgRes, error) {
	return vipcfg.Delete(ctx, req)
}
