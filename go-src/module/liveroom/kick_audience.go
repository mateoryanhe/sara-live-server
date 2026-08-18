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
	"xr-game-server/entity/live"
	"xr-game-server/errercode"
)

// kickAudience 踢出指定观众并推送通知(刷新封禁时间)
func kickAudience(anchorId, userId uint64) {
	pushKickAudience(anchorId, userId, true)
}

// notifyKickBannedAudience 封禁期内尝试进房:推送踢出通知并移出在线列表,但不刷新封禁时间
func notifyKickBannedAudience(anchorId, userId uint64) {
	pushKickAudience(anchorId, userId, false)
}

func pushKickAudience(anchorId, userId uint64, refreshKickTime bool) {
	onlineId := entity.BuildLiveRoomOnlineId(userId, anchorId)
	online := liveroomdao.GetOnlineById(onlineId, userId, anchorId)
	if online == nil {
		return
	}

	var kickAt time.Time
	var banSeconds int64
	if refreshKickTime {
		kickAt = time.Now()
		online.SetKickTime(&kickAt)
		banSeconds = int64(entity.LiveRoomKickBanDuration / time.Second)
	} else if online.KickTime == nil || !online.IsKickBanned() {
		return
	} else {
		kickAt = *online.KickTime
		remain := int64(time.Until(kickAt.Add(entity.LiveRoomKickBanDuration)).Seconds())
		if remain < 0 {
			remain = 0
		}
		banSeconds = remain
	}

	if online.Status == entity.LiveRoomOnlineStatusOnline {
		exitRoom(userId, anchorId)
	}
	push.Data(userId, cmd.LiveRoomAudienceKick, &liveroomdto.AudienceKickPushItem{
		RoomId:     strconv.FormatUint(anchorId, 10),
		UserId:     strconv.FormatUint(userId, 10),
		KickTime:   kickAt.UnixMilli(),
		BanSeconds: banSeconds,
	})
}

// KickAudience 主播踢出指定观众,默认30分钟内不可再次进入
func KickAudience(ctx context.Context, req *liveroomdto.KickAudienceReq) (*liveroomdto.KickAudienceRes, error) {
	anchorId := httpserver.GetAuthId(ctx)

	if req.UserId == anchorId {
		return nil, errercode.CreateCode(errercode.LiveRoomCannotKickSelf)
	}
	if liveroomdao.GetRoomById(anchorId) == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}

	onlineId := entity.BuildLiveRoomOnlineId(req.UserId, anchorId)
	online := liveroomdao.GetOnlineById(onlineId, req.UserId, anchorId)
	if online == nil || online.Status != entity.LiveRoomOnlineStatusOnline {
		return nil, errercode.CreateCode(errercode.LiveRoomAudienceNotOnline)
	}

	kickAudience(anchorId, req.UserId)

	return &liveroomdto.KickAudienceRes{Success: true}, nil
}

// CancelKickBan 主播取消指定观众的进入限制
func CancelKickBan(ctx context.Context, req *liveroomdto.CancelKickBanReq) (*liveroomdto.CancelKickBanRes, error) {
	anchorId := httpserver.GetAuthId(ctx)

	if liveroomdao.GetRoomById(anchorId) == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}

	onlineId := entity.BuildLiveRoomOnlineId(req.UserId, anchorId)
	online := liveroomdao.GetOnlineById(onlineId, req.UserId, anchorId)
	if online == nil || online.KickTime == nil {
		return &liveroomdto.CancelKickBanRes{Success: true}, nil
	}

	online.SetKickTime(nil)
	push.Data(req.UserId, cmd.LiveRoomAudienceKickCancel, &liveroomdto.AudienceKickCancelPushItem{
		RoomId: strconv.FormatUint(anchorId, 10),
		UserId: strconv.FormatUint(req.UserId, 10),
	})

	return &liveroomdto.CancelKickBanRes{Success: true}, nil
}
