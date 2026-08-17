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

// ListRecentUnsettledDailyGuildEffectiveLives 直查DB:某工会最近N条未结算日有效直播记录
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

// ClearRecentUnsettledDailyGuildEffectiveLiveCount 工会下架:直查DB取最近8条未结算日表,EffectiveLiveCount置0
func ClearRecentUnsettledDailyGuildEffectiveLiveCount(guildId uint64) {
	rows := ListRecentUnsettledDailyGuildEffectiveLives(guildId)
	for _, row := range rows {
		if row == nil || row.ID == "" {
			continue
		}
		target := resolveDailyGuildEffectiveLiveTarget(row)
		target.SetEffectiveLiveCount(0)
	}
}

// MarkDailyGuildEffectiveLivesSettled 将工会日有效直播记录标记为已结算(不清零次数)
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

// AddDailyGuildEffectiveLiveCount 工会当日有效直播次数+1
func AddDailyGuildEffectiveLiveCount(guildId uint64, at time.Time) {
	if guildId == 0 {
		return
	}
	date := entity.FormatDailyAnchorEffectiveLiveDate(at)
	row := GetDailyGuildEffectiveLive(date, guildId)
	if row == nil {
		return
	}
	row.AddEffectiveLiveCount(1)
}
