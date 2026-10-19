package liveroom

import (
	"context"
	"strconv"

	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/errercode"
	"xr-game-server/module/cmsvis"
)

type batchResetPlatformAnchorUnsettledResult struct {
	resetCount int
	noDataCount int
	failIds    []uint64
}

// BatchResetPlatformAnchorUnsettledIncome 批量清零平台主播未结算收益(不生成结算单/归档/流水日志)
func BatchResetPlatformAnchorUnsettledIncome(ctx context.Context, req *accountdto.BatchResetPlatformAnchorUnsettledIncomeReq) (*accountdto.BatchResetPlatformAnchorUnsettledIncomeRes, error) {
	if req == nil || len(req.AnchorIds) == 0 || len(req.AnchorIds) > accountdto.MaxBatchImmediateSettlePlatformAnchors {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	ids := uniquePlatformAnchorIds(req.AnchorIds)
	if len(ids) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	visibleIds, filterByVisibility := cmsvis.ResolvePlatformAnchorListVisibility(ctx)
	allowedIds, deniedIds := filterImmediateSettlementPlatformAnchorIds(ids, visibleIds, filterByVisibility)
	result := batchResetPlatformAnchorsUnsettledIncome(allowedIds)
	failIds := append(deniedIds, result.failIds...)
	failAnchorIds := make([]string, 0, len(failIds))
	for _, id := range failIds {
		failAnchorIds = append(failAnchorIds, strconv.FormatUint(id, 10))
	}
	return &accountdto.BatchResetPlatformAnchorUnsettledIncomeRes{
		ResetCount:    result.resetCount,
		NoDataCount:   result.noDataCount,
		FailCount:     len(failIds),
		FailAnchorIds: failAnchorIds,
	}, nil
}

func batchResetPlatformAnchorsUnsettledIncome(anchorIds []uint64) *batchResetPlatformAnchorUnsettledResult {
	result := &batchResetPlatformAnchorUnsettledResult{failIds: make([]uint64, 0)}
	for _, anchorId := range anchorIds {
		room := liveroomdao.GetRoomById(anchorId)
		if room == nil || room.GuildId != 0 {
			result.failIds = append(result.failIds, anchorId)
			continue
		}
		switch liveroomdao.ClearLiveRoomIncomeUnsettledForCMS(anchorId) {
		case liveroomdao.ClearUnsettledIncomeCleared:
			result.resetCount++
		case liveroomdao.ClearUnsettledIncomeAlreadyZero:
			result.noDataCount++
		default:
			result.failIds = append(result.failIds, anchorId)
		}
	}
	return result
}
