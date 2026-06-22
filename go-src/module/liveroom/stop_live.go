package liveroom

import (
	"context"
	"strconv"
	"time"
	"xr-game-server/constants/cmd"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/push"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/liveroomdto"
)

// StopLive 下播
func StopLive(ctx context.Context, _ *liveroomdto.StopLiveReq) (*liveroomdto.StopLiveRes, error) {

	userId := httpserver.GetAuthId(ctx)
	stopLive(userId)
	return &liveroomdto.StopLiveRes{}, nil
}

func stopLive(anchorId uint64) {

	//移除在线列表
	clearRoomAudienceCaches(anchorId)

	room := liveroomdao.GetRoomById(anchorId)
	if room == nil {
		return
	}
	taskMap.Remove(anchorId)
	liveRecordId := room.LiveRecordId
	broadcastAnchorStopLive(anchorId, liveRecordId)
	room.SetLiveRecordId(0)
	room.SetHeartTime(nil)
	if liveRecordId == 0 {
		return
	}
	liveRecord := liveroomdao.GetLiveRecordById(liveRecordId)
	if liveRecord == nil {
		return
	}
	now := time.Now()
	liveRecord.SetEndTime(&now)
}

func broadcastAnchorStopLive(roomId, liveRecordId uint64) {
	if liveRecordId == 0 {
		return
	}
	payload := &liveroomdto.AnchorStopLivePushItem{
		RoomId:       strconv.FormatUint(roomId, 10),
		AnchorId:     strconv.FormatUint(roomId, 10),
		LiveRecordId: strconv.FormatUint(liveRecordId, 10),
		StoppedAt:    time.Now().Unix(),
	}
	for _, o := range getOnline(roomId) {
		if o == roomId {
			continue
		}
		push.Data(o, cmd.LiveRoomStopLive, payload)
	}
}
