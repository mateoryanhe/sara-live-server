package liveroom

import (
	"context"
	"errors"
	"strconv"
	"time"

	"xr-game-server/core/xrlog"
	"xr-game-server/core/xrtime"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/errercode"
	"xr-game-server/module/cmsvis"
)

type batchPlatformAnchorSettlementResult struct {
	settledCount int
	noDataCount  int
	failIds      []uint64
}

// BatchImmediateSettlePlatformAnchors 批量立即结算当前CMS用户有权查看的平台主播。
func BatchImmediateSettlePlatformAnchors(ctx context.Context, req *accountdto.BatchImmediateSettlePlatformAnchorsReq) (*accountdto.BatchImmediateSettlePlatformAnchorsRes, error) {
	if req == nil || len(req.AnchorIds) == 0 || len(req.AnchorIds) > accountdto.MaxBatchImmediateSettlePlatformAnchors {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	ids := uniquePlatformAnchorIds(req.AnchorIds)
	if len(ids) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	visibleIds, filterByVisibility := cmsvis.ResolvePlatformAnchorListVisibility(ctx)
	allowedIds, deniedIds := filterImmediateSettlementPlatformAnchorIds(ids, visibleIds, filterByVisibility)
	result, err := batchSettlePlatformAnchorsImmediately(ctx, allowedIds)
	if errors.Is(err, ErrSettlementRunning) {
		return nil, errercode.CreateCode(errercode.RequestTooFrequent)
	}
	if err != nil {
		return nil, err
	}
	if result == nil {
		result = &batchPlatformAnchorSettlementResult{}
	}
	failIds := append(deniedIds, result.failIds...)
	failAnchorIds := make([]string, 0, len(failIds))
	for _, id := range failIds {
		failAnchorIds = append(failAnchorIds, strconv.FormatUint(id, 10))
	}
	return &accountdto.BatchImmediateSettlePlatformAnchorsRes{
		SettledCount:  result.settledCount,
		NoDataCount:   result.noDataCount,
		FailCount:     len(failIds),
		FailAnchorIds: failAnchorIds,
	}, nil
}

func batchSettlePlatformAnchorsImmediately(ctx context.Context, anchorIds []uint64) (*batchPlatformAnchorSettlementResult, error) {
	if len(anchorIds) == 0 {
		return &batchPlatformAnchorSettlementResult{}, nil
	}
	if !settlementRunMu.TryLock() {
		return nil, ErrSettlementRunning
	}
	defer settlementRunMu.Unlock()

	result := &batchPlatformAnchorSettlementResult{failIds: make([]uint64, 0)}
	tierCfg := loadAnchorWeeklySettlementCfg()
	periodEnd := xrtime.NextWeekStart(time.Now())
	xrlog.DetailLog.Infof(ctx, "batch immediate platform anchor settlement start anchorIds=%v periodEnd=%s", anchorIds, periodEnd.Format(time.RFC3339))

	for _, anchorId := range anchorIds {
		room := liveroomdao.GetRoomById(anchorId)
		if room == nil || room.GuildId != 0 {
			result.failIds = append(result.failIds, anchorId)
			continue
		}
		switch settlePlatformAnchorTieredWithOutcome(room, tierCfg, periodEnd) {
		case tieredAnchorSettlementCreated:
			result.settledCount++
		case tieredAnchorSettlementNoData:
			result.noDataCount++
		default:
			result.failIds = append(result.failIds, anchorId)
		}
	}

	xrlog.DetailLog.Infof(ctx, "batch immediate platform anchor settlement done settled=%d noData=%d failed=%d failAnchorIds=%v",
		result.settledCount, result.noDataCount, len(result.failIds), result.failIds)
	return result, nil
}

func uniquePlatformAnchorIds(ids []uint64) []uint64 {
	ret := make([]uint64, 0, len(ids))
	seen := make(map[uint64]struct{}, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ret = append(ret, id)
	}
	return ret
}

func filterImmediateSettlementPlatformAnchorIds(ids, visibleIds []uint64, filter bool) (allowed, denied []uint64) {
	if !filter {
		return ids, nil
	}
	visible := make(map[uint64]struct{}, len(visibleIds))
	for _, id := range visibleIds {
		visible[id] = struct{}{}
	}
	allowed = make([]uint64, 0, len(ids))
	denied = make([]uint64, 0)
	for _, id := range ids {
		if _, ok := visible[id]; ok {
			allowed = append(allowed, id)
		} else {
			denied = append(denied, id)
		}
	}
	return allowed, denied
}
