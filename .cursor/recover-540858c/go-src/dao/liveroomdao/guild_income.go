package liveroomdao

import (
	"time"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/entity/live"
)

var (
	guildIncomeUnsettledCache = gmap.NewKVMap[uint64, *entity.GuildIncomeUnsettled](false)
	guildIncomeSettledCache   = gmap.NewKVMap[uint64, *entity.GuildIncomeSettled](false)
	guildIncomeTotalCache     = gmap.NewKVMap[uint64, *entity.GuildIncomeTotal](false)
)

// GetGuildIncomeUnsettled 工会未结算收益(仅内存,无则新建并写入缓存)
func GetGuildIncomeUnsettled(guildId uint64) *entity.GuildIncomeUnsettled {
	if guildId == 0 {
		return nil
	}
	if guildIncomeUnsettledCache.Contains(guildId) {
		return guildIncomeUnsettledCache.Get(guildId)
	}
	row := entity.NewGuildIncomeUnsettled(guildId)
	guildIncomeUnsettledCache.Set(guildId, row)
	return row
}

// GetGuildIncomeSettled 工会已结算收益(仅内存,无则新建并写入缓存)
func GetGuildIncomeSettled(guildId uint64) *entity.GuildIncomeSettled {
	if guildId == 0 {
		return nil
	}
	if guildIncomeSettledCache.Contains(guildId) {
		return guildIncomeSettledCache.Get(guildId)
	}
	row := entity.NewGuildIncomeSettled(guildId)
	guildIncomeSettledCache.Set(guildId, row)
	return row
}

// GetGuildIncomeTotal 工会生涯累计收益(仅内存,无则新建并写入缓存)
func GetGuildIncomeTotal(guildId uint64) *entity.GuildIncomeTotal {
	if guildId == 0 {
		return nil
	}
	if guildIncomeTotalCache.Contains(guildId) {
		return guildIncomeTotalCache.Get(guildId)
	}
	row := entity.NewGuildIncomeTotal(guildId)
	guildIncomeTotalCache.Set(guildId, row)
	return row
}

// ForRoomGuild 房间所属工会>0时回调(用于双写工会统计)
func ForRoomGuild(roomId uint64, fn func(guildId uint64)) {
	if roomId == 0 || fn == nil {
		return
	}
	guildId := GetAnchorGuildId(roomId)
	if guildId == 0 {
		return
	}
	fn(guildId)
}

// MirrorGuildGiftEarn 同步礼物收益到工会未结算+生涯+日表
func MirrorGuildGiftEarn(roomId uint64, amount float64) {
	at := time.Now()
	ForRoomGuild(roomId, func(guildId uint64) {
		if u := GetGuildIncomeUnsettled(guildId); u != nil {
			u.AddGiftEarn(amount)
		}
		if t := GetGuildIncomeTotal(guildId); t != nil {
			t.AddGiftEarn(amount)
		}
		MirrorDailyGuildGiftEarn(guildId, at, amount)
	})
}

// MirrorGuildPaidDanmakuEarn 同步付费弹幕收益到工会
func MirrorGuildPaidDanmakuEarn(roomId uint64, amount float64) {
	at := time.Now()
	ForRoomGuild(roomId, func(guildId uint64) {
		if u := GetGuildIncomeUnsettled(guildId); u != nil {
			u.AddPaidDanmakuEarn(amount)
		}
		if t := GetGuildIncomeTotal(guildId); t != nil {
			t.AddPaidDanmakuEarn(amount)
		}
		MirrorDailyGuildPaidDanmakuEarn(guildId, at, amount)
	})
}

// MirrorGuildPrivateRoomTicketEarn 同步私密房门票收益到工会
func MirrorGuildPrivateRoomTicketEarn(roomId uint64, amount float64) {
	at := time.Now()
	ForRoomGuild(roomId, func(guildId uint64) {
		if u := GetGuildIncomeUnsettled(guildId); u != nil {
			u.AddPrivateRoomTicketEarn(amount)
		}
		if t := GetGuildIncomeTotal(guildId); t != nil {
			t.AddPrivateRoomTicketEarn(amount)
		}
		MirrorDailyGuildPrivateRoomTicketEarn(guildId, at, amount)
	})
}

// MirrorGuildPrivateRoomWatchEarn 同步私密房观看收益到工会
func MirrorGuildPrivateRoomWatchEarn(roomId uint64, amount float64) {
	at := time.Now()
	ForRoomGuild(roomId, func(guildId uint64) {
		if u := GetGuildIncomeUnsettled(guildId); u != nil {
			u.AddPrivateRoomWatchEarn(amount)
		}
		if t := GetGuildIncomeTotal(guildId); t != nil {
			t.AddPrivateRoomWatchEarn(amount)
		}
		MirrorDailyGuildPrivateRoomWatchEarn(guildId, at, amount)
	})
}

// MirrorGuildShortVideoEarn 同步短视频付费观看收益到工会
func MirrorGuildShortVideoEarn(roomId uint64, amount float64) {
	at := time.Now()
	ForRoomGuild(roomId, func(guildId uint64) {
		if u := GetGuildIncomeUnsettled(guildId); u != nil {
			u.AddShortVideoEarn(amount)
		}
		if t := GetGuildIncomeTotal(guildId); t != nil {
			t.AddShortVideoEarn(amount)
		}
		MirrorDailyGuildShortVideoEarn(guildId, at, amount)
	})
}

// MirrorGuildLiveDuration 同步直播时长到工会
func MirrorGuildLiveDuration(roomId uint64, sec float64) {
	at := time.Now()
	ForRoomGuild(roomId, func(guildId uint64) {
		if u := GetGuildIncomeUnsettled(guildId); u != nil {
			u.AddTotalLiveDuration(sec)
		}
		if t := GetGuildIncomeTotal(guildId); t != nil {
			t.AddTotalLiveDuration(sec)
		}
		MirrorDailyGuildLiveDuration(guildId, at, sec)
	})
}

// MirrorGuildVideoCallIncomeDelta 同步通话收益增减到工会
func MirrorGuildVideoCallIncomeDelta(roomId uint64, amount float64, ticket, billing bool) {
	at := time.Now()
	ForRoomGuild(roomId, func(guildId uint64) {
		unsettled := GetGuildIncomeUnsettled(guildId)
		total := GetGuildIncomeTotal(guildId)
		if unsettled == nil || total == nil {
			return
		}
		entity.ApplyVideoCallIncomeDelta(entity.TbGuildIncomeUnsettled, unsettled.ID, &unsettled.LiveRoomIncomeAmounts, &unsettled.UpdatedAt, amount, ticket, billing)
		entity.ApplyVideoCallIncomeDelta(entity.TbGuildIncomeTotal, total.ID, &total.LiveRoomIncomeAmounts, &total.UpdatedAt, amount, ticket, billing)
		MirrorDailyGuildVideoCallIncomeDelta(guildId, at, amount, ticket, billing)
	})
}

// getGuildIncomeUnsettledForArchive 工会下架归档用:优先缓存,否则直查DB(不新建)
func getGuildIncomeUnsettledForArchive(guildId uint64) *entity.GuildIncomeUnsettled {
	if guildId == 0 {
		return nil
	}
	if guildIncomeUnsettledCache.Contains(guildId) {
		return guildIncomeUnsettledCache.Get(guildId)
	}
	var row entity.GuildIncomeUnsettled
	err := g.Model(string(entity.TbGuildIncomeUnsettled)).Unscoped().WherePri(guildId).Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

// ListGuildIncomeUnsettledTotalForCMS 批量查询工会未结算总收益(缓存优先,否则直查DB,不新建)
func ListGuildIncomeUnsettledTotalForCMS(guildIds []uint64) map[uint64]float64 {
	ret := make(map[uint64]float64)
	if len(guildIds) == 0 {
		return ret
	}
	seen := make(map[uint64]struct{}, len(guildIds))
	missing := make([]uint64, 0, len(guildIds))
	for _, guildId := range guildIds {
		if guildId == 0 {
			continue
		}
		if _, ok := seen[guildId]; ok {
			continue
		}
		seen[guildId] = struct{}{}
		if guildIncomeUnsettledCache.Contains(guildId) {
			if row := guildIncomeUnsettledCache.Get(guildId); row != nil {
				ret[guildId] = row.TotalIncome
			}
			continue
		}
		missing = append(missing, guildId)
	}
	if len(missing) == 0 {
		return ret
	}
	rows := make([]*entity.GuildIncomeUnsettled, 0, len(missing))
	_ = g.Model(string(entity.TbGuildIncomeUnsettled)).Unscoped().
		WhereIn(string(db.IdName), missing).Scan(&rows)
	for _, row := range rows {
		if row != nil && row.ID != 0 {
			ret[row.ID] = row.TotalIncome
		}
	}
	return ret
}

// GetGuildIncomeUnsettledForCMS 未结算收益(缓存优先,否则直查DB,不新建)
func GetGuildIncomeUnsettledForCMS(guildId uint64) *entity.GuildIncomeUnsettled {
	if guildId == 0 {
		return nil
	}
	if guildIncomeUnsettledCache.Contains(guildId) {
		return guildIncomeUnsettledCache.Get(guildId)
	}
	var row entity.GuildIncomeUnsettled
	err := g.Model(string(entity.TbGuildIncomeUnsettled)).Unscoped().WherePri(guildId).Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

// GetGuildIncomeSettledForCMS 已结算收益(缓存优先,否则直查DB,不新建)
func GetGuildIncomeSettledForCMS(guildId uint64) *entity.GuildIncomeSettled {
	if guildId == 0 {
		return nil
	}
	if guildIncomeSettledCache.Contains(guildId) {
		return guildIncomeSettledCache.Get(guildId)
	}
	var row entity.GuildIncomeSettled
	err := g.Model(string(entity.TbGuildIncomeSettled)).Unscoped().WherePri(guildId).Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

// GetGuildIncomeTotalForCMS 生涯累计收益(缓存优先,否则直查DB,不新建)
func GetGuildIncomeTotalForCMS(guildId uint64) *entity.GuildIncomeTotal {
	if guildId == 0 {
		return nil
	}
	if guildIncomeTotalCache.Contains(guildId) {
		return guildIncomeTotalCache.Get(guildId)
	}
	return GetGuildIncomeTotalFromDB(guildId)
}

// GetGuildIncomeTotalFromDB 直查数据库
func GetGuildIncomeTotalFromDB(guildId uint64) *entity.GuildIncomeTotal {
	if guildId == 0 {
		return nil
	}
	var row entity.GuildIncomeTotal
	err := g.Model(string(entity.TbGuildIncomeTotal)).Unscoped().WherePri(guildId).Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

// ListGuildIncomeUnsettledArchives 查询工会下架未结算归档记录
func ListGuildIncomeUnsettledArchives(guildId uint64, limit int) []*entity.GuildIncomeUnsettledArchive {
	if guildId == 0 {
		return nil
	}
	if limit <= 0 {
		limit = 50
	}
	rows := make([]*entity.GuildIncomeUnsettledArchive, 0)
	_ = g.Model(string(entity.TbGuildIncomeUnsettledArchive)).
		Where(string(entity.GuildIncomeUnsettledArchiveGuildId)+" = ?", guildId).
		Order("created_at desc").
		Limit(limit).
		Scan(&rows)
	return rows
}

// ArchiveAndClearGuildUnsettledIncome 工会下架时:新建一条归档记录并清零工会未结算表
func ArchiveAndClearGuildUnsettledIncome(guildId uint64) {
	unsettled := getGuildIncomeUnsettledForArchive(guildId)
	if unsettled == nil || unsettled.IsZero() {
		return
	}
	snap := unsettled.SnapshotAndClear()
	if snap.IsZero() {
		return
	}
	_ = entity.NewGuildIncomeUnsettledArchive(guildId, &snap)
}
