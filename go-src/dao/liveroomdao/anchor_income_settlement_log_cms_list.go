package liveroomdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity/live"
)

// AnchorIncomeSettlementLogCMSListFilter CMS主播结算流水查询条件
type AnchorIncomeSettlementLogCMSListFilter struct {
	RoomId                   uint64
	RoomIds                  []uint64
	StartTime                int64
	EndTime                  int64
	TransferStartTime        int64
	TransferEndTime          int64
	PayoutOnly               bool
	HideTransferredBefore    int64
	Status                   *uint8
	DirectPayout             *bool
	OrderByReceivableUsdDesc bool
	PageIndex                int
	PageSize                 int
}

// AnchorIncomeSettlementLogCMSList CMS分页查询主播结算流水(按ID倒序)
func AnchorIncomeSettlementLogCMSList(f *AnchorIncomeSettlementLogCMSListFilter) (int, []*entity.AnchorIncomeSettlementLog) {
	list := make([]*entity.AnchorIncomeSettlementLog, 0)
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
	m := g.Model(string(entity.TbAnchorIncomeSettlementLog)).Ctx(ctx)
	roomIdCol := string(entity.AnchorIncomeSettlementLogRoomId)
	if len(f.RoomIds) > 0 {
		m = m.Where(roomIdCol+" IN (?)", f.RoomIds)
	} else if f.RoomId > 0 {
		m = m.Where(roomIdCol+" = ?", f.RoomId)
	}
	if f.StartTime > 0 {
		m = m.Where("created_at >= ?", time.Unix(f.StartTime, 0))
	}
	if f.EndTime > 0 {
		m = m.Where("created_at <= ?", time.Unix(f.EndTime, 0))
	}
	if f.PayoutOnly {
		m = m.Where("("+string(entity.AnchorIncomeSettlementLogTransferOrderId)+" <> '' OR "+
			string(entity.AnchorIncomeSettlementLogTransferFailMsg)+" <> '' OR "+
			string(entity.AnchorIncomeSettlementLogStatus)+" IN (?, ?))",
			entity.AnchorIncomeSettlementStatusTransferred, entity.AnchorIncomeSettlementStatusTransferring)
	}
	payoutTimeExpr := "COALESCE(" + string(entity.AnchorIncomeSettlementLogTransferAt) + ", updated_at)"
	if f.TransferStartTime > 0 {
		m = m.Where(payoutTimeExpr+" >= ?", time.Unix(f.TransferStartTime, 0))
	}
	if f.TransferEndTime > 0 {
		m = m.Where(payoutTimeExpr+" <= ?", time.Unix(f.TransferEndTime, 0))
	}
	if f.HideTransferredBefore > 0 {
		m = m.Where("("+string(entity.AnchorIncomeSettlementLogStatus)+" <> ? OR created_at >= ?)",
			entity.AnchorIncomeSettlementStatusTransferred, time.Unix(f.HideTransferredBefore, 0))
	}
	if f.Status != nil {
		m = m.Where(string(entity.AnchorIncomeSettlementLogStatus)+" = ?", *f.Status)
	}
	if f.DirectPayout != nil {
		m = m.Where(string(entity.AnchorIncomeSettlementLogDirectPayout)+" = ?", *f.DirectPayout)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return 0, list
	}
	orderBy := "id desc"
	if f.PayoutOnly || f.TransferStartTime > 0 || f.TransferEndTime > 0 {
		orderBy = payoutTimeExpr + " desc, id desc"
	} else if f.OrderByReceivableUsdDesc {
		orderBy = string(entity.LiveRoomIncomeSettlementReceivableUsd) + " desc, id desc"
	}
	_ = m.Clone().Order(orderBy).
		Limit(f.PageSize).Offset((f.PageIndex - 1) * f.PageSize).
		Scan(&list)
	list = mergeAnchorIncomeSettlementLogsFromCache(list)
	if f.HideTransferredBefore > 0 {
		list = filterHistoricalTransferredAnchorSettlements(list, time.Unix(f.HideTransferredBefore, 0))
	}
	return total, list
}

// filterHistoricalTransferredAnchorSettlements 在缓存覆盖数据库状态后再次兜底过滤，
// 避免 syndb 尚未刷库时把已经转账成功的历史记录短暂显示出来。
func filterHistoricalTransferredAnchorSettlements(rows []*entity.AnchorIncomeSettlementLog, cutoff time.Time) []*entity.AnchorIncomeSettlementLog {
	if len(rows) == 0 || cutoff.IsZero() {
		return rows
	}
	filtered := make([]*entity.AnchorIncomeSettlementLog, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		if row.Status == entity.AnchorIncomeSettlementStatusTransferred && row.CreatedAt.Before(cutoff) {
			continue
		}
		filtered = append(filtered, row)
	}
	return filtered
}

// AnchorIncomeSettlementLogCMSListByGuildIdsFilter CMS按工会ID列表查询主播结算流水
type AnchorIncomeSettlementLogCMSListByGuildIdsFilter struct {
	GuildIds  []uint64
	RoomId    uint64
	StartTime int64
	EndTime   int64
	PageIndex int
	PageSize  int
}

// AnchorIncomeSettlementLogCMSListByGuildIds CMS分页查询指定工会下主播结算流水(关联 live_rooms.guild_id)
func AnchorIncomeSettlementLogCMSListByGuildIds(f *AnchorIncomeSettlementLogCMSListByGuildIdsFilter) (int, []*entity.AnchorIncomeSettlementLog) {
	list := make([]*entity.AnchorIncomeSettlementLog, 0)
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
	logTable := string(entity.TbAnchorIncomeSettlementLog)
	roomTable := string(entity.TbLiveRoom)
	roomIdCol := string(entity.AnchorIncomeSettlementLogRoomId)
	guildIdCol := string(entity.LiveRoomGuildId)
	m := g.Model(logTable+" a").Ctx(ctx).
		InnerJoin(roomTable+" r", "r.id = a."+roomIdCol).
		Where("r."+guildIdCol+" IN (?)", f.GuildIds)
	if f.RoomId > 0 {
		m = m.Where("a."+roomIdCol+" = ?", f.RoomId)
	}
	if f.StartTime > 0 {
		m = m.Where("a.created_at >= ?", time.Unix(f.StartTime, 0))
	}
	if f.EndTime > 0 {
		m = m.Where("a.created_at <= ?", time.Unix(f.EndTime, 0))
	}
	total, err := m.Clone().Count()
	if err != nil {
		return 0, list
	}
	_ = m.Clone().Fields("a.*").Order("a.id desc").
		Limit(f.PageSize).Offset((f.PageIndex - 1) * f.PageSize).
		Scan(&list)
	return total, mergeAnchorIncomeSettlementLogsFromCache(list)
}
