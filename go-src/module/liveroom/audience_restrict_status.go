package liveroom

import (
	"context"
	"strconv"
	"time"

	"xr-game-server/core/httpserver"
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
