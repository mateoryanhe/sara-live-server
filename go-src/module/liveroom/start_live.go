package liveroom

import (
	"context"
	"time"
	"xr-game-server/core/snowflake"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/liveroomdto"
)

// StartLive 开播
func StartLive(ctx context.Context, _ *liveroomdto.StartLiveReq) (*liveroomdto.StartLiveRes, error) {
	room, err := loadOwnRoom(ctx)
	if err != nil {
		return nil, err
	}
	if room.LiveRecordId > 0 {
		//重置

	}
	addTask(room.ID)
	liveRecordId := snowflake.GetId()

	room.SetLiveRecordId(liveRecordId)
	now := time.Now()
	room.SetHeartTime(&now)

	//记录开播
	liveRecord := liveroomdao.GetLiveRecordById(liveRecordId)
	liveRecord.SetStartTime(time.Now())
	liveRecord.SetAnchorId(room.ID)
	return &liveroomdto.StartLiveRes{}, nil
}
