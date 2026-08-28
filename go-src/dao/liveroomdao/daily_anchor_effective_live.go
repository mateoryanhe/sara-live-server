package liveroomdao

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity/live"
	userentity "xr-game-server/entity/user"
)

const recentUnsettledDailyEffectiveLiveLimit = 8

var dailyAnchorEffectiveLiveCacheMgr *cache.RowCache[*entity.DailyAnchorEffectiveLive]

func initDailyAnchorEffectiveLiveDao() {
	dailyAnchorEffectiveLiveCacheMgr = cache.NewRowCache[*entity.DailyAnchorEffectiveLive]()
}

// GetDailyAnchorEffectiveLive 按日期+直播间获取每日直播时长(缓存;无则新建缓冲对象)
func GetDailyAnchorEffectiveLive(date string, roomId uint64) *entity.DailyAnchorEffectiveLive {
	if date == "" || roomId == 0 || dailyAnchorEffectiveLiveCacheMgr == nil {
		return nil
	}
	id := entity.BuildDailyAnchorEffectiveLiveId(date, roomId)
	v := dailyAnchorEffectiveLiveCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.DailyAnchorEffectiveLive, error) {
		var row *entity.DailyAnchorEffectiveLive
		_ = g.Model(string(entity.TbDailyAnchorEffectiveLive)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewDailyAnchorEffectiveLive(date, roomId), nil
		}
		return row, nil
	})
	return v
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

// ListRecentUnsettledDailyEffectiveLives 直查DB:某主播最近N条未结算日表,命中缓存则用缓存数据
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
	return mergeDailyAnchorEffectiveLivesFromCache(rows)
}

func mergeDailyAnchorEffectiveLivesFromCache(rows []*entity.DailyAnchorEffectiveLive) []*entity.DailyAnchorEffectiveLive {
	if len(rows) == 0 {
		return rows
	}
	list := make([]*entity.DailyAnchorEffectiveLive, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.ID == "" {
			continue
		}
		list = append(list, resolveDailyAnchorEffectiveLiveTarget(row))
	}
	return list
}

// DailyAnchorEffectiveLiveCMSListFilter CMS主播每日流水查询条件
type DailyAnchorEffectiveLiveCMSListFilter struct {
	RoomId        uint64
	LiveDateStart string
	LiveDateEnd   string
	Settled       int8 // -1全部,0未结算,1已结算
	PageIndex     int
	PageSize      int
}

// DailyAnchorEffectiveLiveCMSList CMS分页查询主播每日流水(直查DB,命中缓存则用缓存数据)
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
	m = applyDailyAnchorEffectiveLiveCMSFilters(m, f.LiveDateStart, f.LiveDateEnd, f.Settled, "")
	total, err := m.Clone().Count()
	if err != nil {
		return 0, list
	}
	_ = m.Clone().
		Order(string(entity.DailyAnchorEffectiveLiveLiveDate) + " desc, id desc").
		Limit(f.PageSize).
		Offset((f.PageIndex - 1) * f.PageSize).
		Scan(&list)
	return total, mergeDailyAnchorEffectiveLivesFromCache(list)
}

// DailyAnchorEffectiveLiveCMSMultiListFilter CMS多主播每日流水查询条件
type DailyAnchorEffectiveLiveCMSMultiListFilter struct {
	RoomIds       []uint64
	LiveDateStart string
	LiveDateEnd   string
	Keyword       string
	Settled       int8
	PageIndex     int
	PageSize      int
}

// DailyAnchorEffectiveLiveCMSMultiList CMS分页查询多主播每日流水(直查DB,命中缓存则用缓存数据)
func DailyAnchorEffectiveLiveCMSMultiList(f *DailyAnchorEffectiveLiveCMSMultiListFilter) (int, []*entity.DailyAnchorEffectiveLive) {
	list := make([]*entity.DailyAnchorEffectiveLive, 0)
	if f == nil {
		return 0, list
	}
	if f.PageIndex <= 0 {
		f.PageIndex = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 20
	}
	ctx := gctx.New()
	keyword := strings.TrimSpace(f.Keyword)
	aliased := keyword != ""
	colPrefix := ""
	var m = g.Model(string(entity.TbDailyAnchorEffectiveLive)).Ctx(ctx)
	if aliased {
		like := "%" + keyword + "%"
		colPrefix = "d."
		m = g.Model(string(entity.TbDailyAnchorEffectiveLive)+" d").Ctx(ctx).
			LeftJoin(string(userentity.TbUserInfo)+" u", "u.id = d."+string(entity.DailyAnchorEffectiveLiveRoomId)).
			Where("(d.id LIKE ? OR CAST(d."+string(entity.DailyAnchorEffectiveLiveRoomId)+" AS CHAR) LIKE ? OR u."+string(userentity.UserInfoNickname)+" LIKE ?)", like, like, like)
	}
	if len(f.RoomIds) > 0 {
		m = m.Where(colPrefix+string(entity.DailyAnchorEffectiveLiveRoomId)+" IN (?)", f.RoomIds)
	}
	m = applyDailyAnchorEffectiveLiveCMSFilters(m, f.LiveDateStart, f.LiveDateEnd, f.Settled, colPrefix)
	total, err := m.Clone().Count()
	if err != nil {
		return 0, list
	}
	query := m.Clone().
		Order(colPrefix+string(entity.DailyAnchorEffectiveLiveLiveDate)+" desc, "+colPrefix+"id desc").
		Limit(f.PageSize).
		Offset((f.PageIndex - 1) * f.PageSize)
	if aliased {
		query = query.Fields("d.*")
	}
	_ = query.Scan(&list)
	return total, mergeDailyAnchorEffectiveLivesFromCache(list)
}

// ListWeeklyUnsettledIncomeByRoomIds 批量汇总本周未结算直播收益(按 room_id)
func ListWeeklyUnsettledIncomeByRoomIds(roomIds []uint64, liveDateStart, liveDateEnd string) map[uint64]float64 {
	ret := make(map[uint64]float64)
	if len(roomIds) == 0 || liveDateStart == "" || liveDateEnd == "" {
		return ret
	}
	type sumRow struct {
		RoomId uint64  `json:"room_id"`
		Total  float64 `json:"total"`
	}
	rows := make([]sumRow, 0, len(roomIds))
	ctx := gctx.New()
	_ = g.Model(string(entity.TbDailyAnchorEffectiveLive)).Ctx(ctx).
		Fields("room_id, SUM("+string(entity.LiveRoomIncomeTotalIncome)+") AS total").
		Where(string(entity.DailyAnchorEffectiveLiveRoomId)+" IN (?)", roomIds).
		Where(string(entity.DailyAnchorEffectiveLiveSettled)+" = ?", false).
		Where(string(entity.DailyAnchorEffectiveLiveLiveDate)+" >= ?", liveDateStart).
		Where(string(entity.DailyAnchorEffectiveLiveLiveDate)+" <= ?", liveDateEnd).
		Group(string(entity.DailyAnchorEffectiveLiveRoomId)).
		Scan(&rows)
	for _, row := range rows {
		if row.RoomId == 0 {
			continue
		}
		ret[row.RoomId] = row.Total
	}
	return ret
}

func applyDailyAnchorEffectiveLiveCMSFilters(m *gdb.Model, liveDateStart, liveDateEnd string, settled int8, colPrefix string) *gdb.Model {
	if liveDateStart != "" {
		m = m.Where(colPrefix+string(entity.DailyAnchorEffectiveLiveLiveDate)+" >= ?", liveDateStart)
	}
	if liveDateEnd != "" {
		m = m.Where(colPrefix+string(entity.DailyAnchorEffectiveLiveLiveDate)+" <= ?", liveDateEnd)
	}
	if settled == 0 {
		m = m.Where(colPrefix+string(entity.DailyAnchorEffectiveLiveSettled)+" = ?", false)
	} else if settled == 1 {
		m = m.Where(colPrefix+string(entity.DailyAnchorEffectiveLiveSettled)+" = ?", true)
	}
	return m
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
		if v, _ := dailyAnchorEffectiveLiveCacheMgr.GetRowCached(gctx.New(), row.ID); v != nil {
			return v
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
	PublishDailyAnchorEffectiveLive(row)
	ForRoomGuild(roomId, func(guildId uint64) {
		AddDailyGuildLiveDuration(guildId, at, durationSec)
	})
}
