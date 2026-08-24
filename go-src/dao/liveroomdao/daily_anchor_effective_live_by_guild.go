package liveroomdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity/live"
)

// DailyAnchorEffectiveLiveCMSListByGuildFilter CMS按工会查询名下主播每日流水
type DailyAnchorEffectiveLiveCMSListByGuildFilter struct {
	GuildId       uint64
	RoomId        uint64
	LiveDateStart string
	LiveDateEnd   string
	Settled       int8 // -1全部,0未结算,1已结算
	PageIndex     int
	PageSize      int
}

// DailyAnchorEffectiveLiveCMSListByGuild CMS分页查询工会名下主播每日流水(直查DB,命中缓存则用缓存数据)
func DailyAnchorEffectiveLiveCMSListByGuild(f *DailyAnchorEffectiveLiveCMSListByGuildFilter) (int, []*entity.DailyAnchorEffectiveLive) {
	list := make([]*entity.DailyAnchorEffectiveLive, 0)
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
	dailyTable := string(entity.TbDailyAnchorEffectiveLive)
	roomTable := string(entity.TbLiveRoom)
	roomIdCol := string(entity.DailyAnchorEffectiveLiveRoomId)
	guildIdCol := string(entity.LiveRoomGuildId)
	m := g.Model(dailyTable+" d").Ctx(ctx).
		InnerJoin(roomTable+" r", "r.id = d."+roomIdCol).
		Where("r."+guildIdCol+" = ?", f.GuildId)
	if f.RoomId > 0 {
		m = m.Where("d."+roomIdCol+" = ?", f.RoomId)
	}
	if f.LiveDateStart != "" {
		m = m.Where("d."+string(entity.DailyAnchorEffectiveLiveLiveDate)+" >= ?", f.LiveDateStart)
	}
	if f.LiveDateEnd != "" {
		m = m.Where("d."+string(entity.DailyAnchorEffectiveLiveLiveDate)+" <= ?", f.LiveDateEnd)
	}
	if f.Settled == 0 {
		m = m.Where("d."+string(entity.DailyAnchorEffectiveLiveSettled)+" = ?", false)
	} else if f.Settled == 1 {
		m = m.Where("d."+string(entity.DailyAnchorEffectiveLiveSettled)+" = ?", true)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return 0, list
	}
	_ = m.Clone().Fields("d.*").
		Order("d." + string(entity.DailyAnchorEffectiveLiveLiveDate) + " desc, d.id desc").
		Limit(f.PageSize).
		Offset((f.PageIndex - 1) * f.PageSize).
		Scan(&list)
	return total, mergeDailyAnchorEffectiveLivesFromCache(list)
}

// DailyAnchorEffectiveLiveCMSListByGuildIdsFilter CMS按多个工会查询名下主播每日流水
type DailyAnchorEffectiveLiveCMSListByGuildIdsFilter struct {
	GuildIds      []uint64
	RoomId        uint64
	LiveDateStart string
	LiveDateEnd   string
	Settled       int8
	PageIndex     int
	PageSize      int
}

// DailyAnchorEffectiveLiveCMSListByGuildIds CMS分页查询多个工会名下主播每日流水
func DailyAnchorEffectiveLiveCMSListByGuildIds(f *DailyAnchorEffectiveLiveCMSListByGuildIdsFilter) (int, []*entity.DailyAnchorEffectiveLive) {
	list := make([]*entity.DailyAnchorEffectiveLive, 0)
	if f == nil || len(f.GuildIds) == 0 {
		return 0, list
	}
	if f.PageIndex <= 0 {
		f.PageIndex = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 20
	}
	ctx := gctx.New()
	dailyTable := string(entity.TbDailyAnchorEffectiveLive)
	roomTable := string(entity.TbLiveRoom)
	roomIdCol := string(entity.DailyAnchorEffectiveLiveRoomId)
	guildIdCol := string(entity.LiveRoomGuildId)
	m := g.Model(dailyTable+" d").Ctx(ctx).
		InnerJoin(roomTable+" r", "r.id = d."+roomIdCol).
		WhereIn("r."+guildIdCol, f.GuildIds)
	if f.RoomId > 0 {
		m = m.Where("d."+roomIdCol+" = ?", f.RoomId)
	}
	if f.LiveDateStart != "" {
		m = m.Where("d."+string(entity.DailyAnchorEffectiveLiveLiveDate)+" >= ?", f.LiveDateStart)
	}
	if f.LiveDateEnd != "" {
		m = m.Where("d."+string(entity.DailyAnchorEffectiveLiveLiveDate)+" <= ?", f.LiveDateEnd)
	}
	if f.Settled == 0 {
		m = m.Where("d."+string(entity.DailyAnchorEffectiveLiveSettled)+" = ?", false)
	} else if f.Settled == 1 {
		m = m.Where("d."+string(entity.DailyAnchorEffectiveLiveSettled)+" = ?", true)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return 0, list
	}
	_ = m.Clone().Fields("d.*").
		Order("d." + string(entity.DailyAnchorEffectiveLiveLiveDate) + " desc, d.id desc").
		Limit(f.PageSize).
		Offset((f.PageIndex - 1) * f.PageSize).
		Scan(&list)
	return total, mergeDailyAnchorEffectiveLivesFromCache(list)
}
