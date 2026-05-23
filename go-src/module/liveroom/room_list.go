package liveroom

import (
	"context"
	"strconv"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/module/upload"
)

// GetRoomList App 分页查询直播间列表(直连数据库)
func GetRoomList(_ context.Context, req *liveroomdto.GetLiveRoomListReq) (*liveroomdto.GetLiveRoomListRes, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	total, rows := liveroomdao.ListRooms(page, pageSize, req.StatusFilter)
	list := make([]*liveroomdto.LiveRoomListItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, &liveroomdto.LiveRoomListItem{
			RoomId:         strconv.FormatUint(row.ID, 10),
			GuildId:        strconv.FormatUint(row.GuildId, 10),
			Title:          row.Title,
			Cover:          upload.GetUrlByName(row.Cover),
			Notice:         row.Notice,
			Status:         row.Status,
			CreateAt:       row.CreatedAt,
			AnchorNickname: row.Nickname,
			AnchorAvatar:   upload.GetUrlByName(row.Avatar),
		})
	}

	return &liveroomdto.GetLiveRoomListRes{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     list,
	}, nil
}
