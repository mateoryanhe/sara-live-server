package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/globalcfgdto"
	"xr-game-server/module/globalcfg"
)

type GlobalCfgController struct {
}

func initGlobalCfgController() {
	httpserver.RegCMS("/globalCfg", &GlobalCfgController{})
}

func (a *GlobalCfgController) GetGlobalCfg(ctx context.Context, req *globalcfgdto.GetGlobalCfgReq) (*httpserver.CMSQueryResp, error) {
	return globalcfg.GetGlobalCfg(ctx, req)
}

func (a *GlobalCfgController) SaveGlobalCfg(ctx context.Context, req *globalcfgdto.SaveGlobalCfgReq) (bool, error) {
	return globalcfg.SaveGlobalCfg(ctx, req)
}

func (a *GlobalCfgController) DelGlobalCfg(ctx context.Context, req *globalcfgdto.DelGlobalCfgReq) (bool, error) {
	return globalcfg.DelGlobalCfg(ctx, req)
}
