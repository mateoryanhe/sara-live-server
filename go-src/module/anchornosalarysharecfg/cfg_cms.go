package anchornosalarysharecfg

import (
	"context"
	"math"
	"strconv"
	"time"

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
	if anchornosalarysharecfgdao.GetFirst() != nil {
		return
	}
	_ = anchornosalarysharecfgdao.Save(&entity.AnchorNoSalaryShareCfg{
		AnchorSocialSharePercent: DefaultAnchorSocialSharePercent,
		GuildSocialSharePercent:  DefaultGuildSocialSharePercent,
	})
}

func Get(_ context.Context, _ *anchornosalarysharecfgdto.GetAnchorNoSalaryShareCfgReq) (*anchornosalarysharecfgdto.GetAnchorNoSalaryShareCfgRes, error) {
	row := anchornosalarysharecfgdao.GetFirst()
	if row == nil {
		return &anchornosalarysharecfgdto.GetAnchorNoSalaryShareCfgRes{
			Cfg: &anchornosalarysharecfgdto.AnchorNoSalaryShareCfgItem{
				AnchorSocialSharePercent: DefaultAnchorSocialSharePercent,
				GuildSocialSharePercent:  DefaultGuildSocialSharePercent,
			},
		}, nil
	}
	return &anchornosalarysharecfgdto.GetAnchorNoSalaryShareCfgRes{Cfg: toItem(row)}, nil
}

func Save(_ context.Context, req *anchornosalarysharecfgdto.SaveAnchorNoSalaryShareCfgReq) (*anchornosalarysharecfgdto.SaveAnchorNoSalaryShareCfgRes, error) {
	if req == nil || !validPercent(req.AnchorSocialSharePercent) || !validPercent(req.GuildSocialSharePercent) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	existing := anchornosalarysharecfgdao.GetFirst()
	row := &entity.AnchorNoSalaryShareCfg{
		AnchorSocialSharePercent: roundPercent(req.AnchorSocialSharePercent),
		GuildSocialSharePercent:  roundPercent(req.GuildSocialSharePercent),
	}
	if req.ID > 0 {
		if existing == nil || existing.ID != req.ID {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		row.ID = existing.ID
		row.CreatedAt = existing.CreatedAt
	} else if existing != nil {
		row.ID = existing.ID
		row.CreatedAt = existing.CreatedAt
	}
	row.UpdatedAt = time.Now()
	if row.CreatedAt.IsZero() {
		row.CreatedAt = row.UpdatedAt
	}
	if err := anchornosalarysharecfgdao.Save(row); err != nil {
		return nil, err
	}
	return &anchornosalarysharecfgdto.SaveAnchorNoSalaryShareCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func validPercent(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0) && value >= 0 && value <= 100
}

func roundPercent(value float64) float64 {
	return math.Round(value*100) / 100
}

func toItem(row *entity.AnchorNoSalaryShareCfg) *anchornosalarysharecfgdto.AnchorNoSalaryShareCfgItem {
	if row == nil {
		return nil
	}
	return &anchornosalarysharecfgdto.AnchorNoSalaryShareCfgItem{
		ID:                       strconv.FormatUint(row.ID, 10),
		AnchorSocialSharePercent: row.AnchorSocialSharePercent,
		GuildSocialSharePercent:  row.GuildSocialSharePercent,
		CreatedAt:                formatTime(row.CreatedAt),
		UpdatedAt:                formatTime(row.UpdatedAt),
	}
}

func formatTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format("2006-01-02 15:04:05")
}
