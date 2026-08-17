package anchorsalarycfg

import (
	"context"
	"strconv"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/anchorsalarycfgdao"
	"xr-game-server/dto/anchorsalarycfgdto"
	"xr-game-server/entity/live"
	"xr-game-server/errercode"
)

func GetList(_ context.Context, req *anchorsalarycfgdto.AnchorSalaryCfgListReq) (*httpserver.CMSQueryResp, error) {
	total, list := anchorsalarycfgdao.GetList(req)
	return httpserver.NewCMSQueryResp(total, list), nil
}

func Create(_ context.Context, req *anchorsalarycfgdto.CreateAnchorSalaryCfgReq) (*anchorsalarycfgdto.CreateAnchorSalaryCfgRes, error) {
	row := &entity.AnchorSalaryCfg{
		DailyEffectiveLiveCount:  req.DailyEffectiveLiveCount,
		WeeklyEffectiveLiveCount: req.WeeklyEffectiveLiveCount,
		SalaryAmount:             req.SalaryAmount,
		Sort:                     req.Sort,
	}
	if err := anchorsalarycfgdao.Create(row); err != nil {
		return nil, err
	}
	return &anchorsalarycfgdto.CreateAnchorSalaryCfgRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func Update(_ context.Context, req *anchorsalarycfgdto.UpdateAnchorSalaryCfgReq) (*anchorsalarycfgdto.UpdateAnchorSalaryCfgRes, error) {
	row := anchorsalarycfgdao.GetById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row.DailyEffectiveLiveCount = req.DailyEffectiveLiveCount
	row.WeeklyEffectiveLiveCount = req.WeeklyEffectiveLiveCount
	row.SalaryAmount = req.SalaryAmount
	row.Sort = req.Sort
	if err := anchorsalarycfgdao.Update(row); err != nil {
		return nil, err
	}
	return &anchorsalarycfgdto.UpdateAnchorSalaryCfgRes{Success: true}, nil
}

func Delete(_ context.Context, req *anchorsalarycfgdto.DeleteAnchorSalaryCfgReq) (*anchorsalarycfgdto.DeleteAnchorSalaryCfgRes, error) {
	if anchorsalarycfgdao.GetById(req.ID) == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := anchorsalarycfgdao.Delete(req.ID); err != nil {
		return nil, err
	}
	return &anchorsalarycfgdto.DeleteAnchorSalaryCfgRes{Success: true}, nil
}
