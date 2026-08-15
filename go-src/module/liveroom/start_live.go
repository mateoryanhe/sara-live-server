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
	"xr-game-server/entity/live"
	"xr-game-server/errercode"
	"xr-game-server/module/agora"
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

// StartLiveForBotAnchor CMS机器人主播开播(按是否推流决定是否调用声网云播放)
func StartLiveForBotAnchor(ctx context.Context, anchorId, guildId uint64) error {
	room := liveroomdao.GetRoomByAnchor(anchorId)
	if room == nil {
		room = EnsureAnchorRoom(anchorId, guildId)
	}
	heartTime := time.Now().AddDate(botAnchorLiveHeartYears, 0, 0)
	_, err := startAnchorLive(ctx, room, heartTime)
	if err != nil {
		return err
	}
	if !room.PushStream {
		return nil
	}
	playerId, tokenExpireAt, err := agora.StartBotAnchorCloudPlayer(ctx, anchorId, room.CloudPlayerVideo)
	if err != nil {
		stopLive(anchorId)
		return err
	}
	room.SetCloudPlayerId(playerId)
	tokenExpireTime := time.Unix(tokenExpireAt, 0)
	room.SetCloudPlayerTokenExpireAt(&tokenExpireTime)
	liveroomdao.FlushRoomCache(room)
	agora.ScheduleCloudPlayerTokenRefresh(anchorId, playerId, tokenExpireTime)
	return nil
}

const botAnchorLiveHeartYears = 10

func startAnchorLive(ctx context.Context, room *entity.LiveRoom, heartTime time.Time) (uint64, error) {
	if room == nil {
		return 0, errercode.CreateCode(errercode.LiveRoomNotExist)
	}
	if IsRoomBanned(room) {
		return 0, errercode.CreateCode(errercode.LiveRoomBanned)
	}
	if IsRoomOffShelf(room) {
		return 0, errercode.CreateCode(errercode.LiveRoomOffShelf)
	}
	addTask(room.ID)
	initRoomAudienceCaches(room.ID)

	liveRecordId := snowflake.GetId()
	room.SetLiveRecordId(liveRecordId)
	room.SetHeartTime(&heartTime)

	liveRecord := liveroomdao.GetLiveRecordById(liveRecordId)
	liveRecord.SetStartTime(time.Now())
	liveRecord.SetAnchorId(room.ID)
	liveroomdao.PrependLiveRecordToAppListCache(room.ID, liveRecord)
	liveroomdao.FlushRoomCache(room)
	flushRoomList(ctx)
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
