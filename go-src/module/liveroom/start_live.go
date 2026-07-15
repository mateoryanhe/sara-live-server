package liveroom

import (
	"context"
	"strconv"
	"time"
	"xr-game-server/constants/cmd"
	"xr-game-server/core/push"
	"xr-game-server/core/snowflake"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

// StartLive 开播
func StartLive(ctx context.Context, _ *liveroomdto.StartLiveReq) (*liveroomdto.StartLiveRes, error) {
	room, err := loadOwnRoom(ctx)
	if err != nil {
		return nil, err
	}
	_, err = startAnchorLive(ctx, room, time.Now())
	if err != nil {
		return nil, err
	}
	return &liveroomdto.StartLiveRes{}, nil
}

// StartLiveForBotAnchor CMS机器人主播开播(不调声网,心跳写入10年后防止被心跳任务下播)
func StartLiveForBotAnchor(ctx context.Context, anchorId, guildId uint64) error {
	room := liveroomdao.GetRoomByAnchor(anchorId)
	if room == nil {
		room = EnsureAnchorRoom(anchorId, guildId)
	}
	heartTime := time.Now().AddDate(botAnchorLiveHeartYears, 0, 0)
	_, err := startAnchorLive(ctx, room, heartTime)
	return err
}

const botAnchorLiveHeartYears = 10

func startAnchorLive(ctx context.Context, room *entity.LiveRoom, heartTime time.Time) (uint64, error) {
	if room == nil {
		return 0, errercode.CreateCode(errercode.LiveRoomNotExist)
	}
	if IsRoomBanned(room) {
		return 0, errercode.CreateCode(errercode.LiveRoomBanned)
	}
	addTask(room.ID)
	initRoomOnline(room.ID)

	liveRecordId := snowflake.GetId()
	room.SetLiveRecordId(liveRecordId)
	room.SetHeartTime(&heartTime)

	liveRecord := liveroomdao.GetLiveRecordById(liveRecordId)
	liveRecord.SetStartTime(time.Now())
	liveRecord.SetAnchorId(room.ID)
	liveroomdao.FlushRoomCache(room)
	flushRoomList(ctx)
	flushOnlineLists(room.ID)
	broadcastAnchorStartLive(room.ID, liveRecordId, liveRecord.StartTime.Unix())
	return liveRecordId, nil
}

func broadcastAnchorStartLive(roomId, liveRecordId uint64, startedAt int64) {
	if roomId == 0 || liveRecordId == 0 {
		return
	}
	payload := &liveroomdto.AnchorStartLivePushItem{
		RoomId:       strconv.FormatUint(roomId, 10),
		AnchorId:     strconv.FormatUint(roomId, 10),
		LiveRecordId: strconv.FormatUint(liveRecordId, 10),
		StartedAt:    startedAt,
	}
	push.Broadcast(cmd.LiveRoomStartLive, payload)
}
