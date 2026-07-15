package liveroom

import (
	"context"
	"github.com/gogf/gf/v2/os/gctx"
	"strconv"
	"time"
	"xr-game-server/constants/cmd"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/push"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/liverecord"
)

// StopLive 下播
func StopLive(ctx context.Context, _ *liveroomdto.StopLiveReq) (*liveroomdto.StopLiveRes, error) {

	userId := httpserver.GetAuthId(ctx)
	liveRecord := stopLive(userId)
	return &liveroomdto.StopLiveRes{
		LiveRecord: liverecord.ToAppItem(liveRecord),
	}, nil
}

// StopLiveForBotAnchor CMS机器人主播下播
func StopLiveForBotAnchor(ctx context.Context, anchorId uint64) error {
	_ = ctx
	room := liveroomdao.GetRoomByAnchor(anchorId)
	if room == nil || room.LiveRecordId == 0 {
		return errercode.CreateCode(errercode.LiveRoomNotLive)
	}
	stopLive(anchorId)
	liveroomdao.FlushRoomCache(room)
	return nil
}

func stopLive(anchorId uint64) *entity.LiveRecord {

	//移除在线列表
	clearRoomAudienceCaches(anchorId)

	room := liveroomdao.GetRoomById(anchorId)
	if room == nil {
		return nil
	}
	//清除直播间
	taskMap.Remove(anchorId)
	liveRecordId := room.LiveRecordId
	broadcastAnchorStopLive(anchorId, liveRecordId)
	room.SetLiveRecordId(0)
	room.SetHeartTime(nil)
	flushRoomList(gctx.New())

	//清除直播记录
	if liveRecordId == 0 {
		return nil
	}
	liveRecord := liveroomdao.GetLiveRecordById(liveRecordId)
	if liveRecord == nil {
		return nil
	}
	now := time.Now()
	liveRecord.SetEndTime(&now)
	return liveRecord
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
	push.Broadcast(cmd.LiveRoomStopLive, payload)
}
