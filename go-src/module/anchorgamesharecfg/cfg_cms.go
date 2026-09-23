package anchorgamesharecfg

import (
	"context"
	"math"
	"strconv"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/anchorgamesharecfgdao"
	"xr-game-server/dto/anchorgamesharecfgdto"
	"xr-game-server/entity/live"
	"xr-game-server/errercode"
)

func GetList(_ context.Context, req *anchorgamesharecfgdto.AnchorGameShareCfgListReq) (*httpserver.CMSQueryResp, error) {
	if req == nil || !validSalaryType(req.SalaryType) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	total, list := anchorgamesharecfgdao.GetList(req)
	return httpserver.NewCMSQueryResp(total, list), nil
}

func Create(_ context.Context, req *anchorgamesharecfgdto.CreateAnchorGameShareCfgReq) (*anchorgamesharecfgdto.CreateAnchorGameShareCfgRes, error) {
	if req == nil || !validSalaryType(req.SalaryType) || req.Level == 0 || !validGoldRevenue(req.GameTotalGoldRevenue) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if !validPercent(req.AnchorGameSharePercent) || !validPercent(req.GuildGameSharePercent) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if anchorgamesharecfgdao.Exists(req.SalaryType, req.Level, 0) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := &entity.AnchorGameShareCfg{
		SalaryType:             req.SalaryType,
		Level:                  req.Level,
		GameTotalGoldRevenue:   roundGoldRevenue(req.GameTotalGoldRevenue),
		AnchorGameSharePercent: roundPercent(req.AnchorGameSharePercent),
		GuildGameSharePercent:  roundPercent(req.GuildGameSharePercent),
	}
	if err := anchorgamesharecfgdao.Create(row); err != nil {
		return nil, err
	}
	return &anchorgamesharecfgdto.CreateAnchorGameShareCfgRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func Update(_ context.Context, req *anchorgamesharecfgdto.UpdateAnchorGameShareCfgReq) (*anchorgamesharecfgdto.UpdateAnchorGameShareCfgRes, error) {
	if req == nil || req.Level == 0 || !validGoldRevenue(req.GameTotalGoldRevenue) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if !validPercent(req.AnchorGameSharePercent) || !validPercent(req.GuildGameSharePercent) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := anchorgamesharecfgdao.GetByID(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if anchorgamesharecfgdao.Exists(row.SalaryType, req.Level, row.ID) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row.Level = req.Level
	row.GameTotalGoldRevenue = roundGoldRevenue(req.GameTotalGoldRevenue)
	row.AnchorGameSharePercent = roundPercent(req.AnchorGameSharePercent)
	row.GuildGameSharePercent = roundPercent(req.GuildGameSharePercent)
	if err := anchorgamesharecfgdao.Update(row); err != nil {
		return nil, err
	}
	return &anchorgamesharecfgdto.UpdateAnchorGameShareCfgRes{Success: true}, nil
}

func Delete(_ context.Context, req *anchorgamesharecfgdto.DeleteAnchorGameShareCfgReq) (*anchorgamesharecfgdto.DeleteAnchorGameShareCfgRes, error) {
	if req == nil || anchorgamesharecfgdao.GetByID(req.ID) == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := anchorgamesharecfgdao.Delete(req.ID); err != nil {
		return nil, err
	}
	return &anchorgamesharecfgdto.DeleteAnchorGameShareCfgRes{Success: true}, nil
}

func validSalaryType(value uint32) bool {
	return value == entity.AnchorGameShareSalaryTypeWithSalary || value == entity.AnchorGameShareSalaryTypeNoSalary
}

func validPercent(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0) && value >= 0 && value <= 100
}

func validGoldRevenue(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0) && value >= 0
}

func roundGoldRevenue(value float64) float64 {
	return math.Round(value*10000) / 10000
}

func roundPercent(value float64) float64 {
	return math.Round(value*100) / 100
}
