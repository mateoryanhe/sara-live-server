package liveroom

import (
	"context"
	"strconv"
	"time"
	"xr-game-server/constants/cmd"
	"xr-game-server/core/snowflake"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/errercode"
)

// StartLive 开播
func StartLive(ctx context.Context, _ *liveroomdto.StartLiveReq) (*liveroomdto.StartLiveRes, error) {
	room, err := loadOwnRoom(ctx)
	if err != nil {
		return nil, err
	}
	if IsRoomBanned(room) {
		return nil, errercode.CreateCode(errercode.LiveRoomBanned)
	}
	if room.LiveRecordId > 0 {
		//重置

	}
	addTask(room.ID)
	//初始在线列表
	initRoomOnline(room.ID)

	liveRecordId := snowflake.GetId()

	room.SetLiveRecordId(liveRecordId)
	now := time.Now()
	room.SetHeartTime(&now)

	//记录开播
	liveRecord := liveroomdao.GetLiveRecordById(liveRecordId)
	liveRecord.SetStartTime(time.Now())
	liveRecord.SetAnchorId(room.ID)
	//
	flushOnlineLists(room.ID)
	pushAnchorStartLiveToAudience(room.ID, liveRecordId, liveRecord.StartTime.Unix())
	return &liveroomdto.StartLiveRes{}, nil
}

func pushAnchorStartLiveToAudience(roomId, liveRecordId uint64, startedAt int64) {
	if roomId == 0 || liveRecordId == 0 {
		return
	}
	payload := &liveroomdto.AnchorStartLivePushItem{
		RoomId:       strconv.FormatUint(roomId, 10),
		AnchorId:     strconv.FormatUint(roomId, 10),
		LiveRecordId: strconv.FormatUint(liveRecordId, 10),
		StartedAt:    startedAt,
	}
	PushToRoomAudience(roomId, cmd.LiveRoomStartLive, payload, roomId)
}
