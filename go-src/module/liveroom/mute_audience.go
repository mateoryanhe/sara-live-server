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

// GetAudienceRestrictStatus 查询观众在指定直播间的禁言与被踢状态
func GetAudienceRestrictStatus(ctx context.Context, req *liveroomdto.GetAudienceRestrictStatusReq) (*liveroomdto.GetAudienceRestrictStatusRes, error) {
	authId := httpserver.GetAuthId(ctx)
	if authId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if authId != req.UserId && authId != req.RoomId {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}
	if liveroomdao.GetRoomById(req.RoomId) == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}

	onlineId := entity.BuildLiveRoomOnlineId(req.UserId, req.RoomId)
	online := liveroomdao.GetOnlineById(onlineId, req.UserId, req.RoomId)

	res := &liveroomdto.GetAudienceRestrictStatusRes{
		RoomId: strconv.FormatUint(req.RoomId, 10),
		UserId: strconv.FormatUint(req.UserId, 10),
	}
	if online == nil {
		return res, nil
	}

	res.Muted = online.Muted
	res.KickBanned = online.IsKickBanned()
	if online.KickTime != nil {
		kickTime := online.KickTime.Unix()
		expireAt := online.KickTime.Add(entity.LiveRoomKickBanDuration).Unix()
		res.KickTime = kickTime
		res.KickBanExpireAt = expireAt
		if res.KickBanned {
			remain := expireAt - time.Now().Unix()
			if remain > 0 {
				res.KickRemainSeconds = remain
			}
		}
	}
	return res, nil
}

// SetAudienceMute 主播对指定观众禁言或解禁
func SetAudienceMute(ctx context.Context, req *liveroomdto.SetAudienceMuteReq) (*liveroomdto.SetAudienceMuteRes, error) {
	anchorId := httpserver.GetAuthId(ctx)

	if req.UserId == anchorId {
		return nil, errercode.CreateCode(errercode.LiveRoomCannotMuteSelf)
	}
	if liveroomdao.GetRoomById(anchorId) == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}

	onlineId := entity.BuildLiveRoomOnlineId(req.UserId, anchorId)
	online := liveroomdao.GetOnlineById(onlineId, req.UserId, anchorId)
	if online == nil || online.Status != entity.LiveRoomOnlineStatusOnline {
		return nil, errercode.CreateCode(errercode.LiveRoomAudienceNotOnline)
	}
	if online.Muted == req.Muted {
		return &liveroomdto.SetAudienceMuteRes{Success: true}, nil
	}

	online.SetMuted(req.Muted)
	push.Data(req.UserId, cmd.LiveRoomAudienceMute, &liveroomdto.AudienceMutePushItem{
		RoomId: strconv.FormatUint(anchorId, 10),
		UserId: strconv.FormatUint(req.UserId, 10),
		Muted:  req.Muted,
	})

	return &liveroomdto.SetAudienceMuteRes{Success: true}, nil
}

// CancelAudienceMute 主播取消指定观众的禁言
func CancelAudienceMute(ctx context.Context, req *liveroomdto.CancelAudienceMuteReq) (*liveroomdto.CancelAudienceMuteRes, error) {
	anchorId := httpserver.GetAuthId(ctx)

	if req.UserId == anchorId {
		return nil, errercode.CreateCode(errercode.LiveRoomCannotMuteSelf)
	}
	if liveroomdao.GetRoomById(anchorId) == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}

	onlineId := entity.BuildLiveRoomOnlineId(req.UserId, anchorId)
	online := liveroomdao.GetOnlineById(onlineId, req.UserId, anchorId)
	if online == nil || !online.Muted {
		return &liveroomdto.CancelAudienceMuteRes{Success: true}, nil
	}

	online.SetMuted(false)
	push.Data(req.UserId, cmd.LiveRoomAudienceMute, &liveroomdto.AudienceMutePushItem{
		RoomId: strconv.FormatUint(anchorId, 10),
		UserId: strconv.FormatUint(req.UserId, 10),
		Muted:  false,
	})

	return &liveroomdto.CancelAudienceMuteRes{Success: true}, nil
}
