package liveroomdao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity/live"
)

const recentUnsettledDailyEffectiveLiveLimit = 8

var dailyAnchorEffectiveLiveCacheMgr *cache.CacheMgr

func initDailyAnchorEffectiveLiveDao() {
	dailyAnchorEffectiveLiveCacheMgr = cache.NewCacheMgr()
}

// GetDailyAnchorEffectiveLive 按日期+直播间获取每日直播时长(缓存;无则新建缓冲对象)
func GetDailyAnchorEffectiveLive(date string, roomId uint64) *entity.DailyAnchorEffectiveLive {
	if date == "" || roomId == 0 || dailyAnchorEffectiveLiveCacheMgr == nil {
		return nil
	}
	id := entity.BuildDailyAnchorEffectiveLiveId(date, roomId)
	v := dailyAnchorEffectiveLiveCacheMgr.GetData(id, func(ctx context.Context) (interface{}, error) {
		var row *entity.DailyAnchorEffectiveLive
		_ = g.Model(string(entity.TbDailyAnchorEffectiveLive)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewDailyAnchorEffectiveLive(date, roomId), nil
		}
		return row, nil
	})
	if v == nil {
		return nil
	}
	row, _ := v.(*entity.DailyAnchorEffectiveLive)
	return row
}

// ClearRecentUnsettledDailyLiveDuration 主播下架:直查DB取最近8条未结算日表,直播时长置0
func ClearRecentUnsettledDailyLiveDuration(roomId uint64) {
	rows := ListRecentUnsettledDailyEffectiveLives(roomId)
	for _, row := range rows {
		if row == nil || row.ID == "" {
			continue
		}
		target := resolveDailyAnchorEffectiveLiveTarget(row)
		target.SetLiveDuration(0)
	}
}

// ListRecentUnsettledDailyEffectiveLives 直查DB:某主播最近N条未结算日直播时长记录
func ListRecentUnsettledDailyEffectiveLives(roomId uint64) []*entity.DailyAnchorEffectiveLive {
	if roomId == 0 {
		return nil
	}
	var rows []*entity.DailyAnchorEffectiveLive
	_ = g.Model(string(entity.TbDailyAnchorEffectiveLive)).
		Where(string(entity.DailyAnchorEffectiveLiveRoomId)+" = ? AND "+string(entity.DailyAnchorEffectiveLiveSettled)+" = ?", roomId, false).
		Order(string(db.CreatedAtName) + " desc").
		Limit(recentUnsettledDailyEffectiveLiveLimit).
		Scan(&rows)
	return rows
}

// MarkDailyEffectiveLivesSettled 将日直播时长记录标记为已结算(优先写缓存对象)
func MarkDailyEffectiveLivesSettled(rows []*entity.DailyAnchorEffectiveLive) {
	for _, row := range rows {
		if row == nil || row.ID == "" {
			continue
		}
		target := resolveDailyAnchorEffectiveLiveTarget(row)
		if target.Settled {
			continue
		}
		target.SetSettled(true)
	}
}

func resolveDailyAnchorEffectiveLiveTarget(row *entity.DailyAnchorEffectiveLive) *entity.DailyAnchorEffectiveLive {
	if row == nil {
		return nil
	}
	if dailyAnchorEffectiveLiveCacheMgr != nil {
		if v := dailyAnchorEffectiveLiveCacheMgr.GetFromCache(row.ID); v != nil {
			if cached, ok := v.(*entity.DailyAnchorEffectiveLive); ok && cached != nil {
				return cached
			}
		}
	}
	return row
}

// AddDailyLiveDuration 累加当日直播时长(下播且单场>30分钟时调用)
func AddDailyLiveDuration(roomId uint64, at time.Time, durationSec float64) {
	if roomId == 0 || durationSec <= 0 {
		return
	}
	date := entity.FormatDailyAnchorEffectiveLiveDate(at)
	row := GetDailyAnchorEffectiveLive(date, roomId)
	if row == nil {
		return
	}
	row.AddLiveDuration(durationSec)
	ForRoomGuild(roomId, func(guildId uint64) {
		AddDailyGuildLiveDuration(guildId, at, durationSec)
	})
}
