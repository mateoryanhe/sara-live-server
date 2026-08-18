package guild

import (
	"context"
	"strconv"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/guilddao"
	"xr-game-server/dto/guilddto"
	liveentity "xr-game-server/entity/live"
	"xr-game-server/errercode"
	"xr-game-server/module/incomesettlement"
)

// GetMyGuildAnchorIncomeSettlementLogList CMS分页查询当前用户名下工会的主播结算流水
func GetMyGuildAnchorIncomeSettlementLogList(ctx context.Context, req *guilddto.CMSMyGuildAnchorIncomeSettlementLogListReq) (*httpserver.CMSQueryResp, error) {
	cmsUserId, err := getCMSUserId(ctx)
	if err != nil {
		return nil, err
	}
	guildIds := collectOwnedGuildIds(guilddao.ListGuildsByLeaderId(cmsUserId))
	if len(guildIds) == 0 {
		return incomesettlement.GetAnchorCMSListByGuildIds(ctx, nil, req)
	}
	if filterGuildId := parseGuildIdFilter(req.GuildId); filterGuildId > 0 {
		if !containsGuildId(guildIds, filterGuildId) {
			return nil, errercode.CreateCode(errercode.NoPermission)
		}
		guildIds = []uint64{filterGuildId}
	}
	return incomesettlement.GetAnchorCMSListByGuildIds(ctx, guildIds, req)
}

func collectOwnedGuildIds(rows []*liveentity.LiveGuild) []uint64 {
	ids := make([]uint64, 0, len(rows))
	for _, row := range rows {
		if row != nil && row.ID > 0 {
			ids = append(ids, row.ID)
		}
	}
	return ids
}

func parseGuildIdFilter(val string) uint64 {
	if val == "" {
		return 0
	}
	id, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0
	}
	return id
}

func containsGuildId(ids []uint64, guildId uint64) bool {
	for _, id := range ids {
		if id == guildId {
			return true
		}
	}
	return false
}
