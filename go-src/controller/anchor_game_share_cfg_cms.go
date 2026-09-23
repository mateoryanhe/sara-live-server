package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/anchorgamesharecfgdto"
	"xr-game-server/module/anchorgamesharecfg"
)

const AnchorGameShareCfgCMSURL = "/anchorGameShareCfg"

type AnchorGameShareCfgCMSController struct{}

func initAnchorGameShareCfgCMSController() {
	httpserver.RegCMS(AnchorGameShareCfgCMSURL, &AnchorGameShareCfgCMSController{})
}

func (c *AnchorGameShareCfgCMSController) AnchorGameShareCfgList(ctx context.Context, req *anchorgamesharecfgdto.AnchorGameShareCfgListReq) (*httpserver.CMSQueryResp, error) {
	return anchorgamesharecfg.GetList(ctx, req)
}

func (c *AnchorGameShareCfgCMSController) CreateAnchorGameShareCfg(ctx context.Context, req *anchorgamesharecfgdto.CreateAnchorGameShareCfgReq) (*anchorgamesharecfgdto.CreateAnchorGameShareCfgRes, error) {
	return anchorgamesharecfg.Create(ctx, req)
}

func (c *AnchorGameShareCfgCMSController) UpdateAnchorGameShareCfg(ctx context.Context, req *anchorgamesharecfgdto.UpdateAnchorGameShareCfgReq) (*anchorgamesharecfgdto.UpdateAnchorGameShareCfgRes, error) {
	return anchorgamesharecfg.Update(ctx, req)
}

func (c *AnchorGameShareCfgCMSController) DeleteAnchorGameShareCfg(ctx context.Context, req *anchorgamesharecfgdto.DeleteAnchorGameShareCfgReq) (*anchorgamesharecfgdto.DeleteAnchorGameShareCfgRes, error) {
	return anchorgamesharecfg.Delete(ctx, req)
}
