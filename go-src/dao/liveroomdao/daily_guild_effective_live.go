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

const recentUnsettledDailyGuildEffectiveLiveLimit = 8

var dailyGuildEffectiveLiveCacheMgr *cache.RowCache[*entity.DailyGuildEffectiveLive]

func initDailyGuildEffectiveLiveDao() {
	dailyGuildEffectiveLiveCacheMgr = cache.NewRowCache[*entity.DailyGuildEffectiveLive]()
}

func GetDailyGuildEffectiveLive(date string, guildId uint64) *entity.DailyGuildEffectiveLive {
	if date == "" || guildId == 0 || dailyGuildEffectiveLiveCacheMgr == nil {
		return nil
	}
	id := entity.BuildDailyGuildEffectiveLiveId(date, guildId)
	v := dailyGuildEffectiveLiveCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.DailyGuildEffectiveLive, error) {
		var row *entity.DailyGuildEffectiveLive
		_ = g.Model(string(entity.TbDailyGuildEffectiveLive)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewDailyGuildEffectiveLive(date, guildId), nil
		}
		return row, nil
	})
	return v
}

// ListRecentUnsettledDailyGuildEffectiveLives 直查DB:某工会最近N条未结算日表,命中缓存则用缓存数据
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
	return mergeDailyGuildEffectiveLivesFromCache(rows)
}

func mergeDailyGuildEffectiveLivesFromCache(rows []*entity.DailyGuildEffectiveLive) []*entity.DailyGuildEffectiveLive {
	if len(rows) == 0 {
		return rows
	}
	list := make([]*entity.DailyGuildEffectiveLive, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.ID == "" {
			continue
		}
		list = append(list, resolveDailyGuildEffectiveLiveTarget(row))
	}
	return list
}

// DailyGuildEffectiveLiveCMSListFilter CMS工会每日流水查询条件
type DailyGuildEffectiveLiveCMSListFilter struct {
	GuildId       uint64
	LiveDateStart string
	LiveDateEnd   string
	Settled       int8 // -1全部,0未结算,1已结算
	PageIndex     int
	PageSize      int
}

// DailyGuildEffectiveLiveCMSList CMS分页查询工会每日流水(直查DB,命中缓存则用缓存数据)
func DailyGuildEffectiveLiveCMSList(f *DailyGuildEffectiveLiveCMSListFilter) (int, []*entity.DailyGuildEffectiveLive) {
	list := make([]*entity.DailyGuildEffectiveLive, 0)
	if f == nil || f.GuildId == 0 {
		return 0, list
	}
	if f.PageIndex <= 0 {
		f.PageIndex = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 20
	}
	ctx := gctx.New()
	m := g.Model(string(entity.TbDailyGuildEffectiveLive)).Ctx(ctx).
		Where(string(entity.DailyGuildEffectiveLiveGuildId)+" = ?", f.GuildId)
	if f.LiveDateStart != "" {
		m = m.Where(string(entity.DailyGuildEffectiveLiveLiveDate)+" >= ?", f.LiveDateStart)
	}
	if f.LiveDateEnd != "" {
		m = m.Where(string(entity.DailyGuildEffectiveLiveLiveDate)+" <= ?", f.LiveDateEnd)
	}
	if f.Settled == 0 {
		m = m.Where(string(entity.DailyGuildEffectiveLiveSettled)+" = ?", false)
	} else if f.Settled == 1 {
		m = m.Where(string(entity.DailyGuildEffectiveLiveSettled)+" = ?", true)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return 0, list
	}
	_ = m.Clone().
		Order(string(entity.DailyGuildEffectiveLiveLiveDate) + " desc, id desc").
		Limit(f.PageSize).
		Offset((f.PageIndex - 1) * f.PageSize).
		Scan(&list)
	return total, mergeDailyGuildEffectiveLivesFromCache(list)
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

// MarkDailyGuildEffectiveLivesSettled 将工会日表记录标记为已结算(优先写缓存对象)
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
		if v, _ := dailyGuildEffectiveLiveCacheMgr.GetRowCached(gctx.New(), row.ID); v != nil {
			return v
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
	PublishDailyGuildEffectiveLive(row)
}
