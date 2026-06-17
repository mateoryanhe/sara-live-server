package liverecord

import (
	"context"
	"strconv"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/liverecorddto"
	"xr-game-server/entity"
)

const (
	appDefaultPageSize = 20
	appMaxPageSize     = 100
)

func toAppItem(v *entity.LiveRecord) *liverecorddto.AppLiveRecordItem {
	if v == nil {
		return nil
	}
	item := &liverecorddto.AppLiveRecordItem{
		Id:                     strconv.FormatUint(v.ID, 10),
		StartTime:              v.StartTime.Format("2006/01/02 15:04:05"),
		TotalAudience:          v.TotalAudience,
		TotalLiveDuration:      v.TotalLiveDuration,
		TotalIncome:            v.TotalIncome,
		TotalGiftIncome:        v.TotalGiftIncome,
		TotalPaidDanmakuIncome: v.TotalPaidDanmakuIncome,
		TotalPrivateRoomIncome: v.TotalPrivateRoomIncome,
		TotalGameBet:           v.TotalGameBet,
		TotalGiftSender:        v.TotalGiftSender,
		TotalNewFollower:       v.TotalNewFollower,
	}
	if v.EndTime != nil && !v.EndTime.IsZero() {
		item.EndTime = v.EndTime.Format("2006/01/02 15:04:05")
	}
	return item
}

// GetAppList App端分页查询当前主播直播数据
func GetAppList(ctx context.Context, req *liverecorddto.AppLiveRecordListReq) (*liverecorddto.AppLiveRecordListRes, error) {
	anchorId := httpserver.GetAuthId(ctx)
	page, pageSize := normalizeAppPage(req.Page, req.PageSize)

	total, rows := liveroomdao.LiveRecordCMSList(&liveroomdao.LiveRecordCMSListFilter{
		AnchorId:  anchorId,
		PageIndex: page,
		PageSize:  pageSize,
	})

	list := make([]*liverecorddto.AppLiveRecordItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, toAppItem(row))
	}

	return &liverecorddto.AppLiveRecordListRes{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     list,
	}, nil
}

func normalizeAppPage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = appDefaultPageSize
	}
	if pageSize > appMaxPageSize {
		pageSize = appMaxPageSize
	}
	return page, pageSize
}
