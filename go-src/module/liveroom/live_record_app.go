package liveroom

import (
	"context"
	"strconv"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/liverecorddto"
	"xr-game-server/entity/live"
)

const appLiveRecordListPageSize = liveroomdao.AppLiveRecordListCachePageSize

// ToAppItem 将直播记录实体转为 App 端条目
func ToAppItem(v *entity.LiveRecord) *liverecorddto.AppLiveRecordItem {
	if v == nil {
		return nil
	}
	item := &liverecorddto.AppLiveRecordItem{
		Id:                           strconv.FormatUint(v.ID, 10),
		StartTime:                    v.StartTime.UnixMilli(),
		TotalAudience:                v.TotalAudience,
		TotalLiveDuration:            v.TotalLiveDuration,
		TotalIncome:                  v.TotalIncome,
		TotalGiftIncome:              v.TotalGiftIncome,
		TotalPaidDanmakuIncome:       v.TotalPaidDanmakuIncome,
		TotalPrivateRoomIncome:       v.TotalPrivateRoomIncome,
		TotalPrivateRoomTicketIncome: v.TotalPrivateRoomTicketIncome,
		TotalPrivateRoomWatchIncome:  v.TotalPrivateRoomWatchIncome,
		TotalVideoCallIncome:         v.TotalVideoCallIncome,
		TotalVideoCallTicketIncome:   v.TotalVideoCallTicketIncome,
		TotalVideoCallBillingIncome:  v.TotalVideoCallBillingIncome,
		TotalGameBet:                 v.TotalGameBet,
		TotalGiftSender:              v.TotalGiftSender,
		TotalNewFollower:             v.TotalNewFollower,
	}
	if v.EndTime != nil && !v.EndTime.IsZero() {
		item.EndTime = v.EndTime.UnixMilli()
	}
	return item
}

// GetAppList App端分页查询当前主播直播数据
func GetAppList(ctx context.Context, req *liverecorddto.AppLiveRecordListReq) (*liverecorddto.AppLiveRecordListRes, error) {
	anchorId := httpserver.GetAuthId(ctx)
	page, _ := normalizeAppPage(req.Page, req.PageSize)

	rows := liveroomdao.GetAppLiveRecordsByAnchor(anchorId, page, appLiveRecordListPageSize)
	list := make([]*liverecorddto.AppLiveRecordItem, 0, len(rows))
	for _, row := range rows {
		record := row
		if cached := liveroomdao.GetDataFromCache(row.ID); cached != nil {
			record = cached
		}
		list = append(list, ToAppItem(record))
	}

	return &liverecorddto.AppLiveRecordListRes{
		Page:     page,
		PageSize: appLiveRecordListPageSize,
		List:     list,
	}, nil
}

func normalizeAppPage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = appLiveRecordListPageSize
	}
	if pageSize > appLiveRecordListPageSize {
		pageSize = appLiveRecordListPageSize
	}
	return page, pageSize
}
