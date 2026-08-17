package guildsalarycfg

import (
	"context"
	"strconv"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/guildsalarycfgdao"
	"xr-game-server/dto/guildsalarycfgdto"
	"xr-game-server/entity/live"
	"xr-game-server/errercode"
)

func GetList(_ context.Context, req *guildsalarycfgdto.GuildSalaryCfgListReq) (*httpserver.CMSQueryResp, error) {
	total, list := guildsalarycfgdao.GetList(req)
	return httpserver.NewCMSQueryResp(total, list), nil
}

func Create(_ context.Context, req *guildsalarycfgdto.CreateGuildSalaryCfgReq) (*guildsalarycfgdto.CreateGuildSalaryCfgRes, error) {
	row := &entity.GuildSalaryCfg{
		DailyEffectiveLiveCount:  req.DailyEffectiveLiveCount,
		WeeklyEffectiveLiveCount: req.WeeklyEffectiveLiveCount,
		SalaryAmount:             req.SalaryAmount,
		Sort:                     req.Sort,
	}
	if err := guildsalarycfgdao.Create(row); err != nil {
		return nil, err
	}
	return &guildsalarycfgdto.CreateGuildSalaryCfgRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func Update(_ context.Context, req *guildsalarycfgdto.UpdateGuildSalaryCfgReq) (*guildsalarycfgdto.UpdateGuildSalaryCfgRes, error) {
	row := guildsalarycfgdao.GetById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row.DailyEffectiveLiveCount = req.DailyEffectiveLiveCount
	row.WeeklyEffectiveLiveCount = req.WeeklyEffectiveLiveCount
	row.SalaryAmount = req.SalaryAmount
	row.Sort = req.Sort
	if err := guildsalarycfgdao.Update(row); err != nil {
		return nil, err
	}
	return &guildsalarycfgdto.UpdateGuildSalaryCfgRes{Success: true}, nil
}

func Delete(_ context.Context, req *guildsalarycfgdto.DeleteGuildSalaryCfgReq) (*guildsalarycfgdto.DeleteGuildSalaryCfgRes, error) {
	if guildsalarycfgdao.GetById(req.ID) == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := guildsalarycfgdao.Delete(req.ID); err != nil {
		return nil, err
	}
	return &guildsalarycfgdto.DeleteGuildSalaryCfgRes{Success: true}, nil
}
