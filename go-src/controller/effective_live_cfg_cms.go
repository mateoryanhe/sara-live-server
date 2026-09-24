package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/effectivelivecfgdto"
	"xr-game-server/module/effectivelivecfg"
)

const EffectiveLiveCfgCMSUrl = "/effectiveLiveCfg"

type EffectiveLiveCfgCMSController struct{}

func initEffectiveLiveCfgCMSController() {
	httpserver.RegCMS(EffectiveLiveCfgCMSUrl, &EffectiveLiveCfgCMSController{})
}

func (c *EffectiveLiveCfgCMSController) GetEffectiveLiveCfg(ctx context.Context, req *effectivelivecfgdto.GetEffectiveLiveCfgReq) (*effectivelivecfgdto.GetEffectiveLiveCfgRes, error) {
	return effectivelivecfg.GetEffectiveLiveCfg(ctx, req)
}

func (c *EffectiveLiveCfgCMSController) SaveEffectiveLiveCfg(ctx context.Context, req *effectivelivecfgdto.SaveEffectiveLiveCfgReq) (*effectivelivecfgdto.SaveEffectiveLiveCfgRes, error) {
	return effectivelivecfg.SaveEffectiveLiveCfg(ctx, req)
}
