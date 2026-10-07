package anchorsalarysocialsharecfg

import (
	"context"
	"math"
	"strconv"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/anchorsalarysocialsharecfgdao"
	"xr-game-server/dto/anchorsalarysocialsharecfgdto"
	"xr-game-server/entity/live"
	"xr-game-server/errercode"
)

func GetList(_ context.Context, req *anchorsalarysocialsharecfgdto.AnchorSalarySocialShareCfgListReq) (*httpserver.CMSQueryResp, error) {
	total, list := anchorsalarysocialsharecfgdao.GetList(req)
	return httpserver.NewCMSQueryResp(total, list), nil
}

func Create(_ context.Context, req *anchorsalarysocialsharecfgdto.CreateAnchorSalarySocialShareCfgReq) (*anchorsalarysocialsharecfgdto.CreateAnchorSalarySocialShareCfgRes, error) {
	if req == nil || req.Level == 0 || !validRevenue(req.SocialTotalDiamondRevenue) || !validRevenue(req.LiveBaseSalaryDiamond) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if !validPercent(req.AnchorSocialSharePercent) || !validPercent(req.GuildSocialSharePercent) || anchorsalarysocialsharecfgdao.Exists(req.Level, 0) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := &entity.AnchorSalarySocialShareCfg{
		Level:                     req.Level,
		EffectiveLiveDays:         req.EffectiveLiveDays,
		LiveBaseSalaryDiamond:     roundRevenue(req.LiveBaseSalaryDiamond),
		SocialTotalDiamondRevenue: roundRevenue(req.SocialTotalDiamondRevenue),
		AnchorSocialSharePercent:  roundPercent(req.AnchorSocialSharePercent),
		GuildSocialSharePercent:   roundPercent(req.GuildSocialSharePercent),
	}
	if err := anchorsalarysocialsharecfgdao.Create(row); err != nil {
		return nil, err
	}
	return &anchorsalarysocialsharecfgdto.CreateAnchorSalarySocialShareCfgRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func Update(_ context.Context, req *anchorsalarysocialsharecfgdto.UpdateAnchorSalarySocialShareCfgReq) (*anchorsalarysocialsharecfgdto.UpdateAnchorSalarySocialShareCfgRes, error) {
	if req == nil || req.Level == 0 || !validRevenue(req.SocialTotalDiamondRevenue) || !validRevenue(req.LiveBaseSalaryDiamond) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if !validPercent(req.AnchorSocialSharePercent) || !validPercent(req.GuildSocialSharePercent) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := anchorsalarysocialsharecfgdao.GetByID(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if anchorsalarysocialsharecfgdao.Exists(req.Level, row.ID) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row.Level = req.Level
	row.EffectiveLiveDays = req.EffectiveLiveDays
	row.LiveBaseSalaryDiamond = roundRevenue(req.LiveBaseSalaryDiamond)
	row.SocialTotalDiamondRevenue = roundRevenue(req.SocialTotalDiamondRevenue)
	row.AnchorSocialSharePercent = roundPercent(req.AnchorSocialSharePercent)
	row.GuildSocialSharePercent = roundPercent(req.GuildSocialSharePercent)
	if err := anchorsalarysocialsharecfgdao.Update(row); err != nil {
		return nil, err
	}
	return &anchorsalarysocialsharecfgdto.UpdateAnchorSalarySocialShareCfgRes{Success: true}, nil
}

func Delete(_ context.Context, req *anchorsalarysocialsharecfgdto.DeleteAnchorSalarySocialShareCfgReq) (*anchorsalarysocialsharecfgdto.DeleteAnchorSalarySocialShareCfgRes, error) {
	if req == nil || anchorsalarysocialsharecfgdao.GetByID(req.ID) == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := anchorsalarysocialsharecfgdao.Delete(req.ID); err != nil {
		return nil, err
	}
	return &anchorsalarysocialsharecfgdto.DeleteAnchorSalarySocialShareCfgRes{Success: true}, nil
}

func validPercent(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0) && value >= 0 && value <= 100
}

func validRevenue(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0) && value >= 0 && math.Trunc(value) == value
}

func roundRevenue(value float64) float64 {
	return math.Trunc(value)
}

func roundPercent(value float64) float64 {
	return math.Round(value*100) / 100
}
