package liveroom

import (
	"context"
	"strconv"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/errercode"
	"xr-game-server/module/liverecord"
)

// GetLiveRecord 按直播间ID查询本场直播记录
func GetLiveRecord(ctx context.Context, req *liveroomdto.GetLiveRecordReq) (*liveroomdto.GetLiveRecordRes, error) {
	userId := httpserver.GetAuthId(ctx)

	room := liveroomdao.GetRoomById(req.RoomId)
	if room == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}
	if userId != room.ID && IsRoomOffShelf(room) {
		return nil, errercode.CreateCode(errercode.LiveRoomOffShelf)
	}
	if room.LiveRecordId == 0 {
		return nil, errercode.CreateCode(errercode.LiveRoomNotLive)
	}

	liveRecord := liveroomdao.GetLiveRecordById(room.LiveRecordId)
	if cached := liveroomdao.GetDataFromCache(room.LiveRecordId); cached != nil {
		liveRecord = cached
	}

	return &liveroomdto.GetLiveRecordRes{
		RoomId:       strconv.FormatUint(room.ID, 10),
		LiveRecordId: strconv.FormatUint(room.LiveRecordId, 10),
		LiveRecord:   liverecord.ToAppItem(liveRecord),
	}, nil
}
