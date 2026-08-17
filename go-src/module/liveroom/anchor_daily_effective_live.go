package liveroom

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/entity/live"
	"xr-game-server/errercode"
)

// QueryAnchorDailyEffectiveLiveList CMS分页查询主播每日直播时长
func QueryAnchorDailyEffectiveLiveList(_ context.Context, req *accountdto.GetAnchorDailyEffectiveLiveListReq) (*httpserver.CMSQueryResp, error) {
	if req.AnchorId == 0 || liveroomdao.ResolveRoom(req.AnchorId) == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	total, rows := liveroomdao.DailyAnchorEffectiveLiveCMSList(&liveroomdao.DailyAnchorEffectiveLiveCMSListFilter{
		RoomId:        req.AnchorId,
		LiveDateStart: req.LiveDateStart,
		LiveDateEnd:   req.LiveDateEnd,
		Settled:       req.Settled,
		PageIndex:     req.PageIndex,
		PageSize:      req.PageSize,
	})
	list := make([]*accountdto.AnchorDailyEffectiveLiveItem, 0, len(rows))
	for _, row := range rows {
		if item := toAnchorDailyEffectiveLiveItem(row); item != nil {
			list = append(list, item)
		}
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}

func toAnchorDailyEffectiveLiveItem(row *entity.DailyAnchorEffectiveLive) *accountdto.AnchorDailyEffectiveLiveItem {
	if row == nil || row.ID == "" {
		return nil
	}
	item := &accountdto.AnchorDailyEffectiveLiveItem{
		ID:           row.ID,
		RoomId:       row.RoomId,
		LiveDate:     row.LiveDate,
		LiveDuration: row.LiveDuration,
		Settled:      row.Settled,
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
