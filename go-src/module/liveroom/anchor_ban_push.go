package liveroom

import (
	"strconv"
	"time"
	"xr-game-server/constants/cmd"
	"xr-game-server/core/push"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/liveroomdto"
)

// NotifyAnchorBanned 向主播及其直播间在线观众推送封禁通知,并强制下播、踢下线主播
func NotifyAnchorBanned(anchorId uint64, banApplyTime *time.Time, banReason string) {
	if banApplyTime == nil {
		return
	}

	roomId := anchorId
	payload := &liveroomdto.AnchorBanPushItem{
		RoomId:       strconv.FormatUint(roomId, 10),
		AnchorId:     strconv.FormatUint(anchorId, 10),
		BanApplyTime: banApplyTime.Unix(),
		BannedAt:     time.Now().Unix(),
		BanReason:    banReason,
	}

	pushed := make(map[uint64]struct{})
	pushTo := func(userId uint64) {
		if userId == 0 {
			return
		}
		if _, ok := pushed[userId]; ok {
			return
		}
		push.Data(userId, cmd.LiveRoomAnchorBan, payload)
		pushed[userId] = struct{}{}
	}

	pushTo(anchorId)
	for _, o := range getOnline(roomId) {
		pushTo(o)
	}

	room := liveroomdao.GetRoomByAnchor(anchorId)
	if room != nil && room.LiveRecordId > 0 {
		stopLive(anchorId)
	}
	push.Kick(anchorId)
}
