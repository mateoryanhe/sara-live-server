package liveroomdao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity/live"
)

const recentUnsettledDailyGuildEffectiveLiveLimit = 8

var dailyGuildEffectiveLiveCacheMgr *cache.CacheMgr

func initDailyGuildEffectiveLiveDao() {
	dailyGuildEffectiveLiveCacheMgr = cache.NewCacheMgr()
}

func GetDailyGuildEffectiveLive(date string, guildId uint64) *entity.DailyGuildEffectiveLive {
	if date == "" || guildId == 0 || dailyGuildEffectiveLiveCacheMgr == nil {
		return nil
	}
	id := entity.BuildDailyGuildEffectiveLiveId(date, guildId)
	v := dailyGuildEffectiveLiveCacheMgr.GetData(id, func(ctx context.Context) (interface{}, error) {
		var row *entity.DailyGuildEffectiveLive
		_ = g.Model(string(entity.TbDailyGuildEffectiveLive)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewDailyGuildEffectiveLive(date, guildId), nil
		}
		return row, nil
	})
	if v == nil {
		return nil
	}
	row, _ := v.(*entity.DailyGuildEffectiveLive)
	return row
}

// ListRecentUnsettledDailyGuildEffectiveLives 直查DB:某工会最近N条未结算日直播时长记录
func ListRecentUnsettledDailyGuildEffectiveLives(guildId uint64) []*entity.DailyGuildEffectiveLive {
	if guildId == 0 {
		return nil
	}
	var rows []*entity.DailyGuildEffectiveLive
	_ = g.Model(string(entity.TbDailyGuildEffectiveLive)).
		Where(string(entity.DailyGuildEffectiveLiveGuildId)+" = ? AND "+string(entity.DailyGuildEffectiveLiveSettled)+" = ?", guildId, false).
		Order(string(db.CreatedAtName) + " desc").
		Limit(recentUnsettledDailyGuildEffectiveLiveLimit).
		Scan(&rows)
	return rows
}

// ClearRecentUnsettledDailyGuildLiveDuration 工会下架:直查DB取最近8条未结算日表,直播时长置0
func ClearRecentUnsettledDailyGuildLiveDuration(guildId uint64) {
	rows := ListRecentUnsettledDailyGuildEffectiveLives(guildId)
	for _, row := range rows {
		if row == nil || row.ID == "" {
			continue
		}
		target := resolveDailyGuildEffectiveLiveTarget(row)
		target.SetLiveDuration(0)
	}
}

// MarkDailyGuildEffectiveLivesSettled 将工会日直播时长记录标记为已结算
func MarkDailyGuildEffectiveLivesSettled(rows []*entity.DailyGuildEffectiveLive) {
	for _, row := range rows {
		if row == nil || row.ID == "" {
			continue
		}
		target := resolveDailyGuildEffectiveLiveTarget(row)
		if target.Settled {
			continue
		}
		target.SetSettled(true)
	}
}

func resolveDailyGuildEffectiveLiveTarget(row *entity.DailyGuildEffectiveLive) *entity.DailyGuildEffectiveLive {
	if row == nil {
		return nil
	}
	if dailyGuildEffectiveLiveCacheMgr != nil {
		if v := dailyGuildEffectiveLiveCacheMgr.GetFromCache(row.ID); v != nil {
			if cached, ok := v.(*entity.DailyGuildEffectiveLive); ok && cached != nil {
				return cached
			}
		}
	}
	return row
}

// AddDailyGuildLiveDuration 累加工会当日直播时长(主播下播且单场>30分钟时同步调用)
func AddDailyGuildLiveDuration(guildId uint64, at time.Time, durationSec float64) {
	if guildId == 0 || durationSec <= 0 {
		return
	}
	date := entity.FormatDailyAnchorEffectiveLiveDate(at)
	row := GetDailyGuildEffectiveLive(date, guildId)
	if row == nil {
		return
	}
	row.AddLiveDuration(durationSec)
}
