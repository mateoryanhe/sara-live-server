package liveroom

import (
	"strconv"
	"time"

	"xr-game-server/constants/cmd"
	"xr-game-server/core/push"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/module/upload"
)

func broadcastAudienceJoin(roomId, userId uint64, onlineCount int) {
	payload := buildAudiencePresencePayload(roomId, userId, onlineCount)
	for _, o := range getOnline(roomId) {
		push.Data(o, cmd.LiveRoomAudienceJoin, payload)
	}
	pushAudiencePresenceToAnchor(roomId, cmd.LiveRoomAudienceJoin, payload)
}

func broadcastAudienceLeave(roomId, userId uint64, onlineCount int) {
	payload := &liveroomdto.AudienceLeavePushItem{
		RoomId:      strconv.FormatUint(roomId, 10),
		UserId:      strconv.FormatUint(userId, 10),
		OnlineCount: onlineCount,
		LeftAt:      time.Now().UnixMilli(),
	}
	if u := userinfodao.GetUserInfoByUserId(userId); u != nil {
		payload.Nickname = u.Nickname
		payload.Avatar = upload.ResolveAvatarUrlForUser(userId, u.Avatar)
		payload.VipLevel = u.VipLevel
	}
	for _, o := range getOnline(roomId) {
		push.Data(o, cmd.LiveRoomAudienceLeave, payload)
	}
	pushAudiencePresenceToAnchor(roomId, cmd.LiveRoomAudienceLeave, payload)
}

// pushAudiencePresenceToAnchor 单独通知主播(主播可能未在观众在线列表中)
func pushAudiencePresenceToAnchor(roomId uint64, pushCmd int, payload any) {
	if roomId == 0 {
		return
	}
	push.Data(roomId, pushCmd, payload)
}

func buildAudiencePresencePayload(roomId, userId uint64, onlineCount int) *liveroomdto.AudienceJoinPushItem {
	payload := &liveroomdto.AudienceJoinPushItem{
		RoomId:      strconv.FormatUint(roomId, 10),
		UserId:      strconv.FormatUint(userId, 10),
		OnlineCount: onlineCount,
		JoinedAt:    time.Now().UnixMilli(),
	}
	if u := userinfodao.GetUserInfoByUserId(userId); u != nil {
		payload.Nickname = u.Nickname
		payload.Avatar = upload.ResolveAvatarUrlForUser(userId, u.Avatar)
		payload.VipLevel = u.VipLevel
	}
	return payload
}
