package guild

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/guilddto"
	liveentity "xr-game-server/entity/live"
	"xr-game-server/errercode"
)

// QueryGuildDailyEffectiveLiveList CMS分页查询工会每日流水
func QueryGuildDailyEffectiveLiveList(_ context.Context, req *guilddto.GetGuildDailyEffectiveLiveListReq) (*httpserver.CMSQueryResp, error) {
	if req == nil || req.GuildId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	total, rows := liveroomdao.DailyGuildEffectiveLiveCMSList(&liveroomdao.DailyGuildEffectiveLiveCMSListFilter{
		GuildId:       req.GuildId,
		LiveDateStart: req.LiveDateStart,
		LiveDateEnd:   req.LiveDateEnd,
		Settled:       req.Settled,
		PageIndex:     req.PageIndex,
		PageSize:      req.PageSize,
	})
	list := make([]*guilddto.GuildDailyEffectiveLiveItem, 0, len(rows))
	for _, row := range rows {
		if item := toGuildDailyEffectiveLiveItem(row); item != nil {
			list = append(list, item)
		}
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}

func toGuildDailyEffectiveLiveItem(row *liveentity.DailyGuildEffectiveLive) *guilddto.GuildDailyEffectiveLiveItem {
	if row == nil || row.ID == "" {
		return nil
	}
	item := &guilddto.GuildDailyEffectiveLiveItem{
		ID:                        row.ID,
		GuildId:                   row.GuildId,
		LiveDate:                  row.LiveDate,
		LiveDuration:              row.LiveDuration,
		Settled:                   row.Settled,
		LiveRoomIncomeAmountsItem: toIncomeAmountsItem(&row.LiveRoomIncomeAmounts),
	}
	if !row.CreatedAt.IsZero() {
		createdAt := row.CreatedAt
		item.CreatedAt = &createdAt
	}
	if !row.UpdatedAt.IsZero() {
		updatedAt := row.UpdatedAt
		item.UpdatedAt = &updatedAt
	}
	return item
}
