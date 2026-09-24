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
	GuildType                *uint8
	StartTime                int64
	EndTime                  int64
	TransferStartTime        int64
	TransferEndTime          int64
	HideTransferredBefore    int64
	PayoutOnly               bool
	Status                   *uint8
	OrderByReceivableUsdDesc bool
	IncludeDetail            bool
	IncludeTransfer          bool
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
	settlementAlias := "s"
	transferAlias := "t"
	m := g.Model(string(entity.TbGuildIncomeSettlementLog) + " " + settlementAlias).Ctx(ctx)
	guildIdCol := settlementAlias + "." + string(entity.GuildIncomeSettlementLogGuildId)
	if f.GuildType != nil {
		guildAlias := "g"
		m = m.InnerJoin(string(entity.TbLiveGuild)+" "+guildAlias,
			guildAlias+".id = "+guildIdCol).
			Where(guildAlias+"."+string(entity.LiveGuildGuildType)+" = ?", *f.GuildType)
	}
	includePayoutJoin := f.PayoutOnly || f.TransferStartTime > 0 || f.TransferEndTime > 0
	if includePayoutJoin {
		m = m.InnerJoin(string(entity.TbGuildIncomeSettlementTransfer)+" "+transferAlias,
			transferAlias+"."+string(entity.GuildIncomeSettlementTransferSettlementId)+" = "+settlementAlias+".id")
	}
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
		m = m.Where(settlementAlias+".created_at >= ?", time.Unix(f.StartTime, 0))
	}
	if f.EndTime > 0 {
		m = m.Where(settlementAlias+".created_at <= ?", time.Unix(f.EndTime, 0))
	}
	if f.TransferStartTime > 0 {
		m = m.Where("COALESCE("+transferAlias+"."+string(entity.GuildIncomeSettlementLogTransferAt)+", "+settlementAlias+".updated_at) >= ?", time.Unix(f.TransferStartTime, 0))
	}
	if f.TransferEndTime > 0 {
		m = m.Where("COALESCE("+transferAlias+"."+string(entity.GuildIncomeSettlementLogTransferAt)+", "+settlementAlias+".updated_at) <= ?", time.Unix(f.TransferEndTime, 0))
	}
	if f.HideTransferredBefore > 0 {
		m = m.Where("("+settlementAlias+"."+string(entity.GuildIncomeSettlementLogStatus)+" <> ? OR "+settlementAlias+".created_at >= ?)",
			entity.GuildIncomeSettlementStatusTransferred, time.Unix(f.HideTransferredBefore, 0))
	}
	if f.Status != nil {
		m = m.Where(settlementAlias+"."+string(entity.GuildIncomeSettlementLogStatus)+" = ?", *f.Status)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return 0, list
	}
	orderBy := settlementAlias + ".id desc"
	if includePayoutJoin {
		orderBy = "COALESCE(" + transferAlias + "." + string(entity.GuildIncomeSettlementLogTransferAt) + ", " + settlementAlias + ".updated_at) desc, " + settlementAlias + ".id desc"
	} else if f.OrderByReceivableUsdDesc {
		orderBy = settlementAlias + "." + string(entity.LiveRoomIncomeSettlementReceivableUsd) + " desc, " + settlementAlias + ".id desc"
	}
	_ = m.Clone().Fields(settlementAlias + ".*").Order(orderBy).
		Limit(f.PageSize).Offset((f.PageIndex - 1) * f.PageSize).
		Scan(&list)
	if f.IncludeDetail {
		hydrateGuildIncomeSettlementLogs(ctx, list)
	} else if f.IncludeTransfer {
		hydrateGuildIncomeSettlementTransfers(ctx, list)
	}
	list = mergeGuildIncomeSettlementLogsFromCache(list)
	if f.HideTransferredBefore > 0 {
		list = filterHistoricalTransferredGuildSettlements(list, time.Unix(f.HideTransferredBefore, 0))
	}
	return total, list
}

// filterHistoricalTransferredGuildSettlements 在缓存覆盖数据库状态后再次兜底过滤，
// 避免 syndb 尚未刷库时把已经转账成功的历史记录短暂显示出来。
func filterHistoricalTransferredGuildSettlements(rows []*entity.GuildIncomeSettlementLog, cutoff time.Time) []*entity.GuildIncomeSettlementLog {
	if len(rows) == 0 || cutoff.IsZero() {
		return rows
	}
	filtered := make([]*entity.GuildIncomeSettlementLog, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		if row.Status == entity.GuildIncomeSettlementStatusTransferred && row.CreatedAt.Before(cutoff) {
			continue
		}
		filtered = append(filtered, row)
	}
	return filtered
}
