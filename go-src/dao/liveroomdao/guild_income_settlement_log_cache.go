package liveroomdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity/live"
)

var (
	guildIncomeSettlementLogCacheMgr             *cache.RowCache[*entity.GuildIncomeSettlementLog]
	coinMerchantGuildIncomeSettlementLogCacheMgr *cache.RowCache[*entity.GuildIncomeSettlementLog]
)

func initGuildIncomeSettlementLogDao() {
	guildIncomeSettlementLogCacheMgr = cache.NewRowCache[*entity.GuildIncomeSettlementLog]()
	coinMerchantGuildIncomeSettlementLogCacheMgr = cache.NewRowCache[*entity.GuildIncomeSettlementLog]()
}

func guildIncomeSettlementLogCacheMgrFor(storage entity.GuildIncomeSettlementStorage) *cache.RowCache[*entity.GuildIncomeSettlementLog] {
	if storage == entity.GuildIncomeSettlementStorageCoinMerchant {
		return coinMerchantGuildIncomeSettlementLogCacheMgr
	}
	return guildIncomeSettlementLogCacheMgr
}

func loadGuildIncomeSettlementLogFromDB(ctx context.Context, storage entity.GuildIncomeSettlementStorage, id uint64) (*entity.GuildIncomeSettlementLog, error) {
	if id == 0 {
		return nil, nil
	}
	logTb, _, _ := storage.Tables()
	tb := logTb
	var row entity.GuildIncomeSettlementLog
	err := g.Model(string(tb)).Ctx(ctx).WherePri(id).Scan(&row)
	if err != nil || row.ID == 0 {
		return nil, err
	}
	row.Storage = storage
	hydrateGuildIncomeSettlementLogs(ctx, []*entity.GuildIncomeSettlementLog{&row})
	return &row, nil
}

// hydrateGuildIncomeSettlementLogs 批量装配结算明细和代付过程，避免列表逐条查询。
func hydrateGuildIncomeSettlementLogs(ctx context.Context, rows []*entity.GuildIncomeSettlementLog) {
	hydrateGuildIncomeSettlementDetails(ctx, rows)
	hydrateGuildIncomeSettlementTransfers(ctx, rows)
}

func guildIncomeSettlementLogRowIndex(rows []*entity.GuildIncomeSettlementLog) (map[entity.GuildIncomeSettlementStorage][]uint64, map[uint64]*entity.GuildIncomeSettlementLog) {
	if len(rows) == 0 {
		return nil, nil
	}
	idsByStorage := make(map[entity.GuildIncomeSettlementStorage][]uint64)
	rowMap := make(map[uint64]*entity.GuildIncomeSettlementLog, len(rows))
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		storage := row.Storage
		idsByStorage[storage] = append(idsByStorage[storage], row.ID)
		rowMap[row.ID] = row
	}
	return idsByStorage, rowMap
}

func hydrateGuildIncomeSettlementDetails(ctx context.Context, rows []*entity.GuildIncomeSettlementLog) {
	idsByStorage, rowMap := guildIncomeSettlementLogRowIndex(rows)
	if len(rowMap) == 0 {
		return
	}
	for storage, ids := range idsByStorage {
		if len(ids) == 0 {
			continue
		}
		_, detailTb, _ := storage.Tables()
		tb := detailTb
		details := make([]*entity.GuildIncomeSettlementDetail, 0, len(ids))
		_ = g.Model(string(tb)).Ctx(ctx).
			WhereIn(string(entity.GuildIncomeSettlementDetailSettlementId), ids).Scan(&details)
		for _, detail := range details {
			if detail != nil {
				if row := rowMap[detail.SettlementId]; row != nil {
					row.ApplyDetail(detail)
				}
			}
		}
	}
}

func hydrateGuildIncomeSettlementTransfers(ctx context.Context, rows []*entity.GuildIncomeSettlementLog) {
	idsByStorage, rowMap := guildIncomeSettlementLogRowIndex(rows)
	if len(rowMap) == 0 {
		return
	}
	for storage, ids := range idsByStorage {
		if len(ids) == 0 {
			continue
		}
		_, _, transferTb := storage.Tables()
		tb := transferTb
		transfers := make([]*entity.GuildIncomeSettlementTransfer, 0, len(ids))
		_ = g.Model(string(tb)).Ctx(ctx).
			WhereIn(string(entity.GuildIncomeSettlementTransferSettlementId), ids).Scan(&transfers)
		for _, transfer := range transfers {
			if transfer != nil {
				if row := rowMap[transfer.SettlementId]; row != nil {
					row.ApplyTransfer(transfer)
				}
			}
		}
	}
}

func getGuildIncomeSettlementLogByStorage(id uint64, storage entity.GuildIncomeSettlementStorage) *entity.GuildIncomeSettlementLog {
	if id == 0 {
		return nil
	}
	mgr := guildIncomeSettlementLogCacheMgrFor(storage)
	if mgr == nil {
		row, _ := loadGuildIncomeSettlementLogFromDB(gctx.New(), storage, id)
		return row
	}
	return mgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.GuildIncomeSettlementLog, error) {
		return loadGuildIncomeSettlementLogFromDB(ctx, storage, id)
	})
}

// GetGuildIncomeSettlementLogById 按主键查询普通工会结算流水，优先使用进程缓存。
func GetGuildIncomeSettlementLogById(id uint64) *entity.GuildIncomeSettlementLog {
	return getGuildIncomeSettlementLogByStorage(id, entity.GuildIncomeSettlementStorageNormal)
}

// GetCoinMerchantGuildIncomeSettlementLogById 按主键查询币商工会结算流水，优先使用进程缓存。
func GetCoinMerchantGuildIncomeSettlementLogById(id uint64) *entity.GuildIncomeSettlementLog {
	return getGuildIncomeSettlementLogByStorage(id, entity.GuildIncomeSettlementStorageCoinMerchant)
}

// ResolveGuildIncomeSettlementLogById 按主键解析结算流水（先普通表，后币商表）。
func ResolveGuildIncomeSettlementLogById(id uint64) *entity.GuildIncomeSettlementLog {
	if row := GetGuildIncomeSettlementLogById(id); row != nil {
		return row
	}
	return GetCoinMerchantGuildIncomeSettlementLogById(id)
}

func getGuildIncomeSettlementLogFromCacheByStorage(id uint64, storage entity.GuildIncomeSettlementStorage) *entity.GuildIncomeSettlementLog {
	if id == 0 {
		return nil
	}
	mgr := guildIncomeSettlementLogCacheMgrFor(storage)
	if mgr == nil {
		return nil
	}
	row, _ := mgr.GetRowCached(gctx.New(), id)
	return row
}

// GetGuildIncomeSettlementLogFromCacheById 只读取普通工会结算流水缓存，未命中时不查询数据库。
func GetGuildIncomeSettlementLogFromCacheById(id uint64) *entity.GuildIncomeSettlementLog {
	return getGuildIncomeSettlementLogFromCacheByStorage(id, entity.GuildIncomeSettlementStorageNormal)
}

// PublishGuildIncomeSettlementLog 在结算流水变更后刷新缓存。
func PublishGuildIncomeSettlementLog(row *entity.GuildIncomeSettlementLog) {
	if row == nil || row.ID == 0 {
		return
	}
	mgr := guildIncomeSettlementLogCacheMgrFor(row.Storage)
	if mgr == nil {
		return
	}
	mgr.PublishRow(gctx.New(), row.ID, row)
}

// mergeGuildIncomeSettlementLogsFromCache 让数据库列表优先展示缓存中的最新状态。
func mergeGuildIncomeSettlementLogsFromCache(rows []*entity.GuildIncomeSettlementLog) []*entity.GuildIncomeSettlementLog {
	for i, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		if cached := getGuildIncomeSettlementLogFromCacheByStorage(row.ID, row.Storage); cached != nil {
			rows[i] = cached
		}
	}
	return rows
}
