package anchornosalarysharecfg

import (
	"context"
	"math"
	"strconv"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/anchornosalarysharecfgdao"
	"xr-game-server/dto/anchornosalarysharecfgdto"
	"xr-game-server/entity/live"
	"xr-game-server/errercode"
)

const (
	DefaultAnchorSocialSharePercent = 10
	DefaultGuildSocialSharePercent  = 10
)

func Init() {
	seedDefaultIfEmpty()
	normalizeLegacySingleton()
}

func seedDefaultIfEmpty() {
	if anchornosalarysharecfgdao.CountAll() > 0 {
		return
	}
	_ = anchornosalarysharecfgdao.Create(&entity.AnchorNoSalaryShareCfg{
		Level:                     1,
		SocialTotalDiamondRevenue: 0,
		AnchorSocialSharePercent:  DefaultAnchorSocialSharePercent,
		GuildSocialSharePercent:   DefaultGuildSocialSharePercent,
	})
}

// normalizeLegacySingleton 保留旧版单行配置的比例，并把新增的档位字段补成可编辑的首档。
func normalizeLegacySingleton() {
	rows := anchornosalarysharecfgdao.ListAllOrderByThresholdDesc()
	if len(rows) != 1 || rows[0] == nil || rows[0].Level != 0 {
		return
	}
	rows[0].Level = 1
	_ = anchornosalarysharecfgdao.Update(rows[0])
}

func GetList(_ context.Context, req *anchornosalarysharecfgdto.AnchorNoSalaryShareCfgListReq) (*httpserver.CMSQueryResp, error) {
	total, list := anchornosalarysharecfgdao.GetList(req)
	return httpserver.NewCMSQueryResp(total, list), nil
}

func Create(_ context.Context, req *anchornosalarysharecfgdto.CreateAnchorNoSalaryShareCfgReq) (*anchornosalarysharecfgdto.CreateAnchorNoSalaryShareCfgRes, error) {
	if req == nil || req.Level == 0 || !validRevenue(req.SocialTotalDiamondRevenue) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if !validPercent(req.AnchorSocialSharePercent) || !validPercent(req.GuildSocialSharePercent) || anchornosalarysharecfgdao.Exists(req.Level, 0) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := &entity.AnchorNoSalaryShareCfg{
		Level:                     req.Level,
		SocialTotalDiamondRevenue: roundRevenue(req.SocialTotalDiamondRevenue),
		AnchorSocialSharePercent:  roundPercent(req.AnchorSocialSharePercent),
		GuildSocialSharePercent:   roundPercent(req.GuildSocialSharePercent),
	}
	if err := anchornosalarysharecfgdao.Create(row); err != nil {
		return nil, err
	}
	return &anchornosalarysharecfgdto.CreateAnchorNoSalaryShareCfgRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func Update(_ context.Context, req *anchornosalarysharecfgdto.UpdateAnchorNoSalaryShareCfgReq) (*anchornosalarysharecfgdto.UpdateAnchorNoSalaryShareCfgRes, error) {
	if req == nil || req.Level == 0 || !validRevenue(req.SocialTotalDiamondRevenue) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if !validPercent(req.AnchorSocialSharePercent) || !validPercent(req.GuildSocialSharePercent) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := anchornosalarysharecfgdao.GetByID(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if anchornosalarysharecfgdao.Exists(req.Level, row.ID) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row.Level = req.Level
	row.SocialTotalDiamondRevenue = roundRevenue(req.SocialTotalDiamondRevenue)
	row.AnchorSocialSharePercent = roundPercent(req.AnchorSocialSharePercent)
	row.GuildSocialSharePercent = roundPercent(req.GuildSocialSharePercent)
	if err := anchornosalarysharecfgdao.Update(row); err != nil {
		return nil, err
	}
	return &anchornosalarysharecfgdto.UpdateAnchorNoSalaryShareCfgRes{Success: true}, nil
}

func Delete(_ context.Context, req *anchornosalarysharecfgdto.DeleteAnchorNoSalaryShareCfgReq) (*anchornosalarysharecfgdto.DeleteAnchorNoSalaryShareCfgRes, error) {
	if req == nil || anchornosalarysharecfgdao.GetByID(req.ID) == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := anchornosalarysharecfgdao.Delete(req.ID); err != nil {
		return nil, err
	}
	return &anchornosalarysharecfgdto.DeleteAnchorNoSalaryShareCfgRes{Success: true}, nil
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
