package globalcfg

import (
	"context"
	"time"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/snowflake"
	"xr-game-server/dao/globalcfgdao"
	"xr-game-server/dto/globalcfgdto"
)

func GetGlobalCfg(ctx context.Context, req *globalcfgdto.GetGlobalCfgReq) (*httpserver.CMSQueryResp, error) {
	ret := globalcfgdao.GetCfgByModule(req.Module)
	data := make([]*globalcfgdto.GlobalCfgDto, 0)
	for _, v := range ret {
		data = append(data, globalcfgdto.NewGlobalCfgDto(v))
	}
	return httpserver.NewCMSQueryResp(len(ret), data), nil
}

func SaveGlobalCfg(ctx context.Context, req *globalcfgdto.SaveGlobalCfgReq) (bool, error) {
	if req.ID == 0 {
		req.ID = snowflake.GetId()
		req.UpdatedAt = time.Now()
		req.CreatedAt = time.Now()
	}
	globalcfgdao.Save(req.GlobalCfg)
	return true, nil
}

func DelGlobalCfg(ctx context.Context, req *globalcfgdto.DelGlobalCfgReq) (bool, error) {
	globalcfgdao.DelById(req.ID)
	return true, nil
}
