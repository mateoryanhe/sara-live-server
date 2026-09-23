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

func (c *AnchorNoSalaryShareCfgCMSController) GetAnchorNoSalaryShareCfg(ctx context.Context, req *anchornosalarysharecfgdto.GetAnchorNoSalaryShareCfgReq) (*anchornosalarysharecfgdto.GetAnchorNoSalaryShareCfgRes, error) {
	return anchornosalarysharecfg.Get(ctx, req)
}

func (c *AnchorNoSalaryShareCfgCMSController) SaveAnchorNoSalaryShareCfg(ctx context.Context, req *anchornosalarysharecfgdto.SaveAnchorNoSalaryShareCfgReq) (*anchornosalarysharecfgdto.SaveAnchorNoSalaryShareCfgRes, error) {
	return anchornosalarysharecfg.Save(ctx, req)
}
