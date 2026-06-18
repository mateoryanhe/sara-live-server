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
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

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

	now := time.Now()
	online.SetKickTime(&now)
	exitRoom(req.UserId, anchorId)

	push.Data(req.UserId, cmd.LiveRoomAudienceKick, &liveroomdto.AudienceKickPushItem{
		RoomId:     strconv.FormatUint(anchorId, 10),
		UserId:     strconv.FormatUint(req.UserId, 10),
		KickTime:   now.Unix(),
		BanSeconds: int64(entity.LiveRoomKickBanDuration / time.Second),
	})

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
