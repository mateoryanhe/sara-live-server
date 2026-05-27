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
	//清除房间状态
	taskMap.Remove(anchorId)
	room.SetLiveRecordId(0)
	room.SetHeartTime(nil)
	//开始计算直播时长
	liveRecord := liveroomdao.GetLiveRecordById(room.LiveRecordId)
	now := time.Now()
	liveRecord.SetEndTime(&now)
}
