package liveroom

import (
	"context"
	"strconv"

	"xr-game-server/constants/cmd"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/push"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

// SetAudienceMute 主播对指定观众禁言或解禁
func SetAudienceMute(ctx context.Context, req *liveroomdto.SetAudienceMuteReq) (*liveroomdto.SetAudienceMuteRes, error) {
	anchorId := httpserver.GetAuthId(ctx)
	if anchorId != req.RoomId {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}
	if req.UserId == anchorId {
		return nil, errercode.CreateCode(errercode.LiveRoomCannotMuteSelf)
	}
	if liveroomdao.GetRoomById(req.RoomId) == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}

	onlineId := entity.BuildLiveRoomOnlineId(req.UserId, req.RoomId)
	online := liveroomdao.GetOnlineById(onlineId, req.UserId, req.RoomId)
	if online == nil || online.Status != entity.LiveRoomOnlineStatusOnline {
		return nil, errercode.CreateCode(errercode.LiveRoomAudienceNotOnline)
	}
	if online.Muted == req.Muted {
		return &liveroomdto.SetAudienceMuteRes{Success: true}, nil
	}

	online.SetMuted(req.Muted)
	push.Data(req.UserId, cmd.LiveRoomAudienceMute, &liveroomdto.AudienceMutePushItem{
		RoomId: strconv.FormatUint(req.RoomId, 10),
		UserId: strconv.FormatUint(req.UserId, 10),
		Muted:  req.Muted,
	})

	return &liveroomdto.SetAudienceMuteRes{Success: true}, nil
}
