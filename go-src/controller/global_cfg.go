package controller

import (
	"context"
	"time"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/snowflake"
	"xr-game-server/dao/globalcfgdao"
	"xr-game-server/dto/globalcfgdto"
)

type GlobalCfgController struct {
}

func initGlobalCfgController() {
	httpserver.RegCMS("/globalCfg", &GlobalCfgController{})
}

func (a *GlobalCfgController) GetGlobalCfg(ctx context.Context, req *globalcfgdto.GetGlobalCfgReq) (*httpserver.CMSQueryResp, error) {
	ret := globalcfgdao.GetCfgList(req.Module, req.ModuleName)
	data := make([]*globalcfgdto.GlobalCfgDto, 0)
	for _, v := range ret {
		data = append(data, globalcfgdto.NewGlobalCfgDto(v))
	}
	return httpserver.NewCMSQueryResp(len(ret), data), nil
}

func (a *GlobalCfgController) SaveGlobalCfg(ctx context.Context, req *globalcfgdto.SaveGlobalCfgReq) (bool, error) {
	if req.ID == 0 {
		req.ID = snowflake.GetId()
		req.UpdatedAt = time.Now()
		req.CreatedAt = time.Now()
	}
	globalcfgdao.Save(req.GlobalCfg)
	return true, nil
}

func (a *GlobalCfgController) DelGlobalCfg(ctx context.Context, req *globalcfgdto.DelGlobalCfgReq) (bool, error) {
	globalcfgdao.DelById(req.ID)
	return true, nil
}
