package game

import (
	"context"
	"strconv"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/gamebetdao"
	"xr-game-server/dto/gamebetdto"
	"xr-game-server/entity"
)

const appGameBetListPageSize = gamebetdao.AppGameBetListCachePageSize

// GetAppBetList App 端分页查询当前用户游戏下注记录
func GetAppBetList(ctx context.Context, req *gamebetdto.AppGameBetListReq) (*gamebetdto.AppGameBetListRes, error) {
	userId := httpserver.GetAuthId(ctx)
	page := 1
	if req != nil && req.Page > 0 {
		page = req.Page
	}

	rows := gamebetdao.GetAppGameBetLogsByUser(userId, page, appGameBetListPageSize)
	list := make([]*gamebetdto.AppGameBetListItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, toAppGameBetListItem(row))
	}
	return &gamebetdto.AppGameBetListRes{
		Page:     page,
		PageSize: appGameBetListPageSize,
		List:     list,
	}, nil
}

func toAppGameBetListItem(row *entity.GameBetLog) *gamebetdto.AppGameBetListItem {
	if row == nil {
		return nil
	}
	return &gamebetdto.AppGameBetListItem{
		Id:           strconv.FormatUint(row.ID, 10),
		GameCode:     row.GameCode,
		NameEn:       row.NameEn,
		Cover:        BuildGameCoverUrl(row.Cover),
		Amount:       row.Amount,
		PlatformType: row.PlatformType,
		OrderId:      row.OrderId,
		LiveRoomId:   formatOptionalUint64(row.LiveRoomId),
		LiveRecordId: formatOptionalUint64(row.LiveRecordId),
		CreatedAt:    row.CreatedAt.UnixMilli(),
	}
}

func formatOptionalUint64(v uint64) string {
	if v == 0 {
		return ""
	}
	return strconv.FormatUint(v, 10)
}
