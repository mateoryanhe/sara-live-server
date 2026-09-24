package guild

import (
	"context"
	"errors"
	"strconv"

	"xr-game-server/dto/guilddto"
	"xr-game-server/errercode"
	"xr-game-server/module/liveroom"
)

// BatchImmediateSettleGuilds 批量立即结算当前CMS用户有权查看的工会。
func BatchImmediateSettleGuilds(ctx context.Context, req *guilddto.BatchImmediateSettleGuildsReq) (*guilddto.BatchImmediateSettleGuildsRes, error) {
	if req == nil || len(req.GuildIds) == 0 || len(req.GuildIds) > guilddto.MaxBatchImmediateSettleGuilds {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	ids := uniqueGuildIds(req.GuildIds)
	if len(ids) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	visibleIds, filterByVisibility := resolveGuildListVisibility(ctx)
	allowedIds, deniedIds := filterImmediateSettlementGuildIds(ids, visibleIds, filterByVisibility)
	result, err := liveroom.BatchSettleGuildsImmediately(ctx, allowedIds)
	if errors.Is(err, liveroom.ErrSettlementRunning) {
		return nil, errercode.CreateCode(errercode.RequestTooFrequent)
	}
	if err != nil {
		return nil, err
	}
	if result == nil {
		result = &liveroom.BatchGuildSettlementResult{}
	}
	failIds := append(deniedIds, result.FailGuildIds...)
	failGuildIds := make([]string, 0, len(failIds))
	for _, id := range failIds {
		failGuildIds = append(failGuildIds, strconv.FormatUint(id, 10))
	}
	return &guilddto.BatchImmediateSettleGuildsRes{
		SettledCount: result.SettledCount,
		NoDataCount:  result.NoDataCount,
		FailCount:    len(failIds),
		FailGuildIds: failGuildIds,
	}, nil
}

func uniqueGuildIds(ids []uint64) []uint64 {
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

func filterImmediateSettlementGuildIds(ids, visibleIds []uint64, filter bool) (allowed, denied []uint64) {
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
