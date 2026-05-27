package liveroom

import (
	"context"
	"time"
	"xr-game-server/core/httpserver"
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
	room := liveroomdao.GetRoomById(anchorId)
	if room == nil {
		return
	}
	taskMap.Remove(anchorId)
	liveRecordId := room.LiveRecordId
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
