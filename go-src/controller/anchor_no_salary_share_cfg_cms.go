package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/anchornosalarysharecfgdto"
	"xr-game-server/module/anchornosalarysharecfg"
)

const AnchorNoSalaryShareCfgCMSURL = "/anchorNoSalaryShareCfg"

type AnchorNoSalaryShareCfgCMSController struct{}

func initAnchorNoSalaryShareCfgCMSController() {
	httpserver.RegCMS(AnchorNoSalaryShareCfgCMSURL, &AnchorNoSalaryShareCfgCMSController{})
}

func (c *AnchorNoSalaryShareCfgCMSController) AnchorNoSalaryShareCfgList(ctx context.Context, req *anchornosalarysharecfgdto.AnchorNoSalaryShareCfgListReq) (*httpserver.CMSQueryResp, error) {
	return anchornosalarysharecfg.GetList(ctx, req)
}

func (c *AnchorNoSalaryShareCfgCMSController) CreateAnchorNoSalaryShareCfg(ctx context.Context, req *anchornosalarysharecfgdto.CreateAnchorNoSalaryShareCfgReq) (*anchornosalarysharecfgdto.CreateAnchorNoSalaryShareCfgRes, error) {
	return anchornosalarysharecfg.Create(ctx, req)
}

func (c *AnchorNoSalaryShareCfgCMSController) UpdateAnchorNoSalaryShareCfg(ctx context.Context, req *anchornosalarysharecfgdto.UpdateAnchorNoSalaryShareCfgReq) (*anchornosalarysharecfgdto.UpdateAnchorNoSalaryShareCfgRes, error) {
	return anchornosalarysharecfg.Update(ctx, req)
}

func (c *AnchorNoSalaryShareCfgCMSController) DeleteAnchorNoSalaryShareCfg(ctx context.Context, req *anchornosalarysharecfgdto.DeleteAnchorNoSalaryShareCfgReq) (*anchornosalarysharecfgdto.DeleteAnchorNoSalaryShareCfgRes, error) {
	return anchornosalarysharecfg.Delete(ctx, req)
}
