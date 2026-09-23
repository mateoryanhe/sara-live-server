package liveroomdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity/live"
)

var guildIncomeSettlementLogCacheMgr *cache.RowCache[*entity.GuildIncomeSettlementLog]

func initGuildIncomeSettlementLogDao() {
	guildIncomeSettlementLogCacheMgr = cache.NewRowCache[*entity.GuildIncomeSettlementLog]()
}

func loadGuildIncomeSettlementLogFromDB(ctx context.Context, id uint64) (*entity.GuildIncomeSettlementLog, error) {
	if id == 0 {
		return nil, nil
	}
	var row entity.GuildIncomeSettlementLog
	err := g.Model(string(entity.TbGuildIncomeSettlementLog)).Ctx(ctx).WherePri(id).Scan(&row)
	if err != nil || row.ID == 0 {
		return nil, err
	}
	hydrateGuildIncomeSettlementLogs(ctx, []*entity.GuildIncomeSettlementLog{&row})
	return &row, nil
}

// hydrateGuildIncomeSettlementLogs 批量装配结算明细和代付过程，避免列表逐条查询。
func hydrateGuildIncomeSettlementLogs(ctx context.Context, rows []*entity.GuildIncomeSettlementLog) {
	hydrateGuildIncomeSettlementDetails(ctx, rows)
	hydrateGuildIncomeSettlementTransfers(ctx, rows)
}

func guildIncomeSettlementLogRowIndex(rows []*entity.GuildIncomeSettlementLog) ([]uint64, map[uint64]*entity.GuildIncomeSettlementLog) {
	if len(rows) == 0 {
		return nil, nil
	}
	ids := make([]uint64, 0, len(rows))
	rowMap := make(map[uint64]*entity.GuildIncomeSettlementLog, len(rows))
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		ids = append(ids, row.ID)
		rowMap[row.ID] = row
	}
	return ids, rowMap
}

func hydrateGuildIncomeSettlementDetails(ctx context.Context, rows []*entity.GuildIncomeSettlementLog) {
	ids, rowMap := guildIncomeSettlementLogRowIndex(rows)
	if len(ids) == 0 {
		return
	}
	details := make([]*entity.GuildIncomeSettlementDetail, 0, len(ids))
	_ = g.Model(string(entity.TbGuildIncomeSettlementDetail)).Ctx(ctx).
		WhereIn(string(entity.GuildIncomeSettlementDetailSettlementId), ids).Scan(&details)
	for _, detail := range details {
		if detail != nil {
			if row := rowMap[detail.SettlementId]; row != nil {
				row.ApplyDetail(detail)
			}
		}
	}
}

func hydrateGuildIncomeSettlementTransfers(ctx context.Context, rows []*entity.GuildIncomeSettlementLog) {
	ids, rowMap := guildIncomeSettlementLogRowIndex(rows)
	if len(ids) == 0 {
		return
	}
	transfers := make([]*entity.GuildIncomeSettlementTransfer, 0, len(ids))
	_ = g.Model(string(entity.TbGuildIncomeSettlementTransfer)).Ctx(ctx).
		WhereIn(string(entity.GuildIncomeSettlementTransferSettlementId), ids).Scan(&transfers)
	for _, transfer := range transfers {
		if transfer != nil {
			if row := rowMap[transfer.SettlementId]; row != nil {
				row.ApplyTransfer(transfer)
			}
		}
	}
}

// GetGuildIncomeSettlementLogById 按主键查询工会结算流水，优先使用进程缓存。
func GetGuildIncomeSettlementLogById(id uint64) *entity.GuildIncomeSettlementLog {
	if id == 0 {
		return nil
	}
	if guildIncomeSettlementLogCacheMgr == nil {
		row, _ := loadGuildIncomeSettlementLogFromDB(gctx.New(), id)
		return row
	}
	return guildIncomeSettlementLogCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.GuildIncomeSettlementLog, error) {
		return loadGuildIncomeSettlementLogFromDB(ctx, id)
	})
}

// GetGuildIncomeSettlementLogFromCacheById 只读取工会结算流水缓存，未命中时不查询数据库。
func GetGuildIncomeSettlementLogFromCacheById(id uint64) *entity.GuildIncomeSettlementLog {
	if id == 0 || guildIncomeSettlementLogCacheMgr == nil {
		return nil
	}
	row, _ := guildIncomeSettlementLogCacheMgr.GetRowCached(gctx.New(), id)
	return row
}

// PublishGuildIncomeSettlementLog 在结算流水变更后刷新缓存。
func PublishGuildIncomeSettlementLog(row *entity.GuildIncomeSettlementLog) {
	if row == nil || row.ID == 0 || guildIncomeSettlementLogCacheMgr == nil {
		return
	}
	guildIncomeSettlementLogCacheMgr.PublishRow(gctx.New(), row.ID, row)
}

// mergeGuildIncomeSettlementLogsFromCache 让数据库列表优先展示缓存中的最新状态。
func mergeGuildIncomeSettlementLogsFromCache(rows []*entity.GuildIncomeSettlementLog) []*entity.GuildIncomeSettlementLog {
	for i, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		if cached := GetGuildIncomeSettlementLogFromCacheById(row.ID); cached != nil {
			rows[i] = cached
		}
	}
	return rows
}
