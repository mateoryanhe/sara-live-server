package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/anchorsalarycfgdto"
	"xr-game-server/module/anchorsalarycfg"
)

const AnchorSalaryCfgCMSUrl = "/anchorSalaryCfg"

type AnchorSalaryCfgCMSController struct{}

func initAnchorSalaryCfgCMSController() {
	httpserver.RegCMS(AnchorSalaryCfgCMSUrl, &AnchorSalaryCfgCMSController{})
}

func (c *AnchorSalaryCfgCMSController) AnchorSalaryCfgList(ctx context.Context, req *anchorsalarycfgdto.AnchorSalaryCfgListReq) (*httpserver.CMSQueryResp, error) {
	return anchorsalarycfg.GetList(ctx, req)
}

func (c *AnchorSalaryCfgCMSController) CreateAnchorSalaryCfg(ctx context.Context, req *anchorsalarycfgdto.CreateAnchorSalaryCfgReq) (*anchorsalarycfgdto.CreateAnchorSalaryCfgRes, error) {
	return anchorsalarycfg.Create(ctx, req)
}

func (c *AnchorSalaryCfgCMSController) UpdateAnchorSalaryCfg(ctx context.Context, req *anchorsalarycfgdto.UpdateAnchorSalaryCfgReq) (*anchorsalarycfgdto.UpdateAnchorSalaryCfgRes, error) {
	return anchorsalarycfg.Update(ctx, req)
}

func (c *AnchorSalaryCfgCMSController) DeleteAnchorSalaryCfg(ctx context.Context, req *anchorsalarycfgdto.DeleteAnchorSalaryCfgReq) (*anchorsalarycfgdto.DeleteAnchorSalaryCfgRes, error) {
	return anchorsalarycfg.Delete(ctx, req)
}
