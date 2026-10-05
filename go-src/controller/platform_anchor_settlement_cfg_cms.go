package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/platformanchorsettlementcfgdto"
	"xr-game-server/module/platformanchorsettlementcfg"
)

const PlatformAnchorSettlementCfgCMSUrl = "/platformAnchorSettlementCfg"

type PlatformAnchorSettlementCfgCMSController struct{}

func initPlatformAnchorSettlementCfgCMSController() {
	httpserver.RegCMS(PlatformAnchorSettlementCfgCMSUrl, &PlatformAnchorSettlementCfgCMSController{})
}

func (c *PlatformAnchorSettlementCfgCMSController) GetPlatformAnchorSettlementCfg(ctx context.Context, req *platformanchorsettlementcfgdto.GetPlatformAnchorSettlementCfgReq) (*platformanchorsettlementcfgdto.GetPlatformAnchorSettlementCfgRes, error) {
	return platformanchorsettlementcfg.GetPlatformAnchorSettlementCfg(ctx, req)
}

func (c *PlatformAnchorSettlementCfgCMSController) SavePlatformAnchorSettlementCfg(ctx context.Context, req *platformanchorsettlementcfgdto.SavePlatformAnchorSettlementCfgReq) (*platformanchorsettlementcfgdto.SavePlatformAnchorSettlementCfgRes, error) {
	return platformanchorsettlementcfg.SavePlatformAnchorSettlementCfg(ctx, req)
}
