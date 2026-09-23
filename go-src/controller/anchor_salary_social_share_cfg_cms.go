package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/anchorsalarysocialsharecfgdto"
	"xr-game-server/module/anchorsalarysocialsharecfg"
)

const AnchorSalarySocialShareCfgCMSURL = "/anchorSalarySocialShareCfg"

type AnchorSalarySocialShareCfgCMSController struct{}

func initAnchorSalarySocialShareCfgCMSController() {
	httpserver.RegCMS(AnchorSalarySocialShareCfgCMSURL, &AnchorSalarySocialShareCfgCMSController{})
}

func (c *AnchorSalarySocialShareCfgCMSController) AnchorSalarySocialShareCfgList(ctx context.Context, req *anchorsalarysocialsharecfgdto.AnchorSalarySocialShareCfgListReq) (*httpserver.CMSQueryResp, error) {
	return anchorsalarysocialsharecfg.GetList(ctx, req)
}

func (c *AnchorSalarySocialShareCfgCMSController) CreateAnchorSalarySocialShareCfg(ctx context.Context, req *anchorsalarysocialsharecfgdto.CreateAnchorSalarySocialShareCfgReq) (*anchorsalarysocialsharecfgdto.CreateAnchorSalarySocialShareCfgRes, error) {
	return anchorsalarysocialsharecfg.Create(ctx, req)
}

func (c *AnchorSalarySocialShareCfgCMSController) UpdateAnchorSalarySocialShareCfg(ctx context.Context, req *anchorsalarysocialsharecfgdto.UpdateAnchorSalarySocialShareCfgReq) (*anchorsalarysocialsharecfgdto.UpdateAnchorSalarySocialShareCfgRes, error) {
	return anchorsalarysocialsharecfg.Update(ctx, req)
}

func (c *AnchorSalarySocialShareCfgCMSController) DeleteAnchorSalarySocialShareCfg(ctx context.Context, req *anchorsalarysocialsharecfgdto.DeleteAnchorSalarySocialShareCfgReq) (*anchorsalarysocialsharecfgdto.DeleteAnchorSalarySocialShareCfgRes, error) {
	return anchorsalarysocialsharecfg.Delete(ctx, req)
}
