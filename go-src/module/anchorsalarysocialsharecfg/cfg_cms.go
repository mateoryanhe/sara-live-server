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

type defaultTier struct {
	level                     uint32
	socialTotalDiamondRevenue float64
	anchorSocialSharePercent  float64
}

var defaultTiers = []defaultTier{
	{level: 1, socialTotalDiamondRevenue: 10, anchorSocialSharePercent: 1},
	{level: 2, socialTotalDiamondRevenue: 11, anchorSocialSharePercent: 2},
	{level: 3, socialTotalDiamondRevenue: 12, anchorSocialSharePercent: 3},
	{level: 4, socialTotalDiamondRevenue: 13, anchorSocialSharePercent: 4},
	{level: 5, socialTotalDiamondRevenue: 14, anchorSocialSharePercent: 5},
	{level: 6, socialTotalDiamondRevenue: 15, anchorSocialSharePercent: 6},
	{level: 7, socialTotalDiamondRevenue: 16, anchorSocialSharePercent: 7},
	{level: 8, socialTotalDiamondRevenue: 17, anchorSocialSharePercent: 8},
	{level: 9, socialTotalDiamondRevenue: 18, anchorSocialSharePercent: 9},
	{level: 10, socialTotalDiamondRevenue: 19, anchorSocialSharePercent: 10},
	{level: 11, socialTotalDiamondRevenue: 20, anchorSocialSharePercent: 11},
	{level: 12, socialTotalDiamondRevenue: 21, anchorSocialSharePercent: 12},
}

func Init() {
	seedDefaultsIfEmpty()
}

func seedDefaultsIfEmpty() {
	if anchorsalarysocialsharecfgdao.CountAll() > 0 {
		return
	}
	for _, tier := range defaultTiers {
		_ = anchorsalarysocialsharecfgdao.Create(&entity.AnchorSalarySocialShareCfg{
			Level:                     tier.level,
			SocialTotalDiamondRevenue: tier.socialTotalDiamondRevenue,
			AnchorSocialSharePercent:  tier.anchorSocialSharePercent,
			GuildSocialSharePercent:   0,
		})
	}
}

func GetList(_ context.Context, req *anchorsalarysocialsharecfgdto.AnchorSalarySocialShareCfgListReq) (*httpserver.CMSQueryResp, error) {
	total, list := anchorsalarysocialsharecfgdao.GetList(req)
	return httpserver.NewCMSQueryResp(total, list), nil
}

func Create(_ context.Context, req *anchorsalarysocialsharecfgdto.CreateAnchorSalarySocialShareCfgReq) (*anchorsalarysocialsharecfgdto.CreateAnchorSalarySocialShareCfgRes, error) {
	if req == nil || req.Level == 0 || !validRevenue(req.SocialTotalDiamondRevenue) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if !validPercent(req.AnchorSocialSharePercent) || !validPercent(req.GuildSocialSharePercent) || anchorsalarysocialsharecfgdao.Exists(req.Level, 0) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := &entity.AnchorSalarySocialShareCfg{
		Level:                     req.Level,
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
	if req == nil || req.Level == 0 || !validRevenue(req.SocialTotalDiamondRevenue) {
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
	return !math.IsNaN(value) && !math.IsInf(value, 0) && value >= 0
}

func roundRevenue(value float64) float64 {
	return math.Round(value*10000) / 10000
}

func roundPercent(value float64) float64 {
	return math.Round(value*100) / 100
}
