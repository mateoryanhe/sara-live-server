package liveroomdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	live "xr-game-server/entity/live"
)

var anchorIncomeSettlementLogCacheMgr *cache.RowCache[*live.AnchorIncomeSettlementLog]

func initAnchorIncomeSettlementLogDao() {
	anchorIncomeSettlementLogCacheMgr = cache.NewRowCache[*live.AnchorIncomeSettlementLog]()
}

func loadAnchorIncomeSettlementLogFromDB(ctx context.Context, id uint64) (*live.AnchorIncomeSettlementLog, error) {
	if id == 0 {
		return nil, nil
	}
	var row live.AnchorIncomeSettlementLog
	err := g.Model(string(live.TbAnchorIncomeSettlementLog)).Ctx(ctx).WherePri(id).Scan(&row)
	if err != nil || row.ID == 0 {
		return nil, err
	}
	return &row, nil
}

func GetAnchorIncomeSettlementLogById(id uint64) *live.AnchorIncomeSettlementLog {
	if id == 0 {
		return nil
	}
	if anchorIncomeSettlementLogCacheMgr == nil {
		row, _ := loadAnchorIncomeSettlementLogFromDB(gctx.New(), id)
		return row
	}
	return anchorIncomeSettlementLogCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*live.AnchorIncomeSettlementLog, error) {
		return loadAnchorIncomeSettlementLogFromDB(ctx, id)
	})
}

func GetAnchorIncomeSettlementLogFromCacheById(id uint64) *live.AnchorIncomeSettlementLog {
	if id == 0 || anchorIncomeSettlementLogCacheMgr == nil {
		return nil
	}
	row, _ := anchorIncomeSettlementLogCacheMgr.GetRowCached(gctx.New(), id)
	return row
}

func PublishAnchorIncomeSettlementLog(row *live.AnchorIncomeSettlementLog) {
	if row == nil || row.ID == 0 || anchorIncomeSettlementLogCacheMgr == nil {
		return
	}
	anchorIncomeSettlementLogCacheMgr.PublishRow(gctx.New(), row.ID, row)
}

func mergeAnchorIncomeSettlementLogsFromCache(rows []*live.AnchorIncomeSettlementLog) []*live.AnchorIncomeSettlementLog {
	for i, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		if cached := GetAnchorIncomeSettlementLogFromCacheById(row.ID); cached != nil {
			rows[i] = cached
		}
	}
	return rows
}
