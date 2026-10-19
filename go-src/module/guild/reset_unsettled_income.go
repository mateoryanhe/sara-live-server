package guild

import (
	"context"
	"strconv"

	"xr-game-server/dto/guilddto"
	"xr-game-server/errercode"
	"xr-game-server/module/liveroom"
)

// BatchResetGuildUnsettledIncome 批量清零工会及名下主播未结算收益(不写结算/归档日志)
func BatchResetGuildUnsettledIncome(ctx context.Context, req *guilddto.BatchResetGuildUnsettledIncomeReq) (*guilddto.BatchResetGuildUnsettledIncomeRes, error) {
	if req == nil || len(req.GuildIds) == 0 || len(req.GuildIds) > guilddto.MaxBatchImmediateSettleGuilds {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	ids := uniqueGuildIds(req.GuildIds)
	if len(ids) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	visibleIds, filterByVisibility := resolveGuildListVisibility(ctx)
	allowedIds, deniedIds := filterImmediateSettlementGuildIds(ids, visibleIds, filterByVisibility)
	result := liveroom.BatchResetGuildsUnsettledIncome(allowedIds)
	failIds := append(deniedIds, result.FailGuildIds...)
	failGuildIds := make([]string, 0, len(failIds))
	for _, id := range failIds {
		failGuildIds = append(failGuildIds, strconv.FormatUint(id, 10))
	}
	return &guilddto.BatchResetGuildUnsettledIncomeRes{
		ResetCount:   result.ResetCount,
		NoDataCount:  result.NoDataCount,
		FailCount:    len(failIds),
		FailGuildIds: failGuildIds,
	}, nil
}
