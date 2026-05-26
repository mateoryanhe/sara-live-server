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
	room.SetLiveRecordId(0)
	room.ClearFailNum()
	//开始计算直播时长
	liveRecord := liveroomdao.GetLiveRecordById(room.LiveRecordId)
	now := time.Now()
	if liveRecord.EndTime == nil {
		liveRecord.SetEndTime(&now)
	}
	liveRecord.SetTotalLiveDuration(liveRecord.EndTime.Sub(liveRecord.StartTime).Seconds())
}
