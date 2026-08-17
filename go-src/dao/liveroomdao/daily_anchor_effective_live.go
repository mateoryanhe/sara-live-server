package liveroomdao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
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

// DailyAnchorEffectiveLiveCMSListFilter CMS主播每日直播时长查询条件
type DailyAnchorEffectiveLiveCMSListFilter struct {
	RoomId        uint64
	LiveDateStart string
	LiveDateEnd   string
	Settled       int8 // -1全部,0未结算,1已结算
	PageIndex     int
	PageSize      int
}

// DailyAnchorEffectiveLiveCMSList CMS分页查询主播每日直播时长(按日期倒序)
func DailyAnchorEffectiveLiveCMSList(f *DailyAnchorEffectiveLiveCMSListFilter) (int, []*entity.DailyAnchorEffectiveLive) {
	list := make([]*entity.DailyAnchorEffectiveLive, 0)
	if f == nil || f.RoomId == 0 {
		return 0, list
	}
	if f.PageIndex <= 0 {
		f.PageIndex = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 20
	}
	ctx := gctx.New()
	m := g.Model(string(entity.TbDailyAnchorEffectiveLive)).Ctx(ctx).
		Where(string(entity.DailyAnchorEffectiveLiveRoomId)+" = ?", f.RoomId)
	if f.LiveDateStart != "" {
		m = m.Where(string(entity.DailyAnchorEffectiveLiveLiveDate)+" >= ?", f.LiveDateStart)
	}
	if f.LiveDateEnd != "" {
		m = m.Where(string(entity.DailyAnchorEffectiveLiveLiveDate)+" <= ?", f.LiveDateEnd)
	}
	if f.Settled == 0 {
		m = m.Where(string(entity.DailyAnchorEffectiveLiveSettled)+" = ?", false)
	} else if f.Settled == 1 {
		m = m.Where(string(entity.DailyAnchorEffectiveLiveSettled)+" = ?", true)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return 0, list
	}
	_ = m.Clone().
		Order(string(entity.DailyAnchorEffectiveLiveLiveDate) + " desc, id desc").
		Limit(f.PageSize).
		Offset((f.PageIndex - 1) * f.PageSize).
		Scan(&list)
	return total, list
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
