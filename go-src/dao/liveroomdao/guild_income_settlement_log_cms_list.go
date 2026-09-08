package liveroomdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity/live"
)

// GuildIncomeSettlementLogCMSListFilter CMS工会结算流水查询条件
type GuildIncomeSettlementLogCMSListFilter struct {
	GuildId                  uint64
	GuildIds                 []uint64
	FilterByGuild            bool // true 时仅返回 GuildIds 对应工会(空则无数据)
	StartTime                int64
	EndTime                  int64
	Status                   *uint8
	OrderByReceivableUsdDesc bool
	PageIndex                int
	PageSize                 int
}

// GuildIncomeSettlementLogCMSList CMS分页查询工会结算流水
func GuildIncomeSettlementLogCMSList(f *GuildIncomeSettlementLogCMSListFilter) (int, []*entity.GuildIncomeSettlementLog) {
	list := make([]*entity.GuildIncomeSettlementLog, 0)
	if f == nil {
		return 0, list
	}
	if f.FilterByGuild && len(f.GuildIds) == 0 {
		return 0, list
	}
	if f.PageIndex <= 0 {
		f.PageIndex = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 20
	}
	ctx := gctx.New()
	m := g.Model(string(entity.TbGuildIncomeSettlementLog)).Ctx(ctx)
	guildIdCol := string(entity.GuildIncomeSettlementLogGuildId)
	if f.GuildId > 0 {
		if f.FilterByGuild {
			allowed := false
			for _, id := range f.GuildIds {
				if id == f.GuildId {
					allowed = true
					break
				}
			}
			if !allowed {
				return 0, list
			}
		}
		m = m.Where(guildIdCol+" = ?", f.GuildId)
	} else if f.FilterByGuild {
		m = m.WhereIn(guildIdCol, f.GuildIds)
	}
	if f.StartTime > 0 {
		m = m.Where("created_at >= ?", time.Unix(f.StartTime, 0))
	}
	if f.EndTime > 0 {
		m = m.Where("created_at <= ?", time.Unix(f.EndTime, 0))
	}
	if f.Status != nil {
		m = m.Where(string(entity.GuildIncomeSettlementLogStatus)+" = ?", *f.Status)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return 0, list
	}
	orderBy := "id desc"
	if f.OrderByReceivableUsdDesc {
		orderBy = string(entity.LiveRoomIncomeSettlementReceivableUsd) + " desc, id desc"
	}
	_ = m.Clone().Order(orderBy).
		Limit(f.PageSize).Offset((f.PageIndex - 1) * f.PageSize).
		Scan(&list)
	return total, list
}

// GetGuildIncomeSettlementLogById 按ID查询工会结算流水
func GetGuildIncomeSettlementLogById(id uint64) *entity.GuildIncomeSettlementLog {
	if id == 0 {
		return nil
	}
	var row entity.GuildIncomeSettlementLog
	err := g.Model(string(entity.TbGuildIncomeSettlementLog)).WherePri(id).Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}
