package liveroom

import (
	"context"
	"strconv"
	"strings"
	"time"
	"xr-game-server/constants/cmd"
	"xr-game-server/core/push"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
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

// IsRoomBanned 直播间是否处于封禁有效期内
func IsRoomBanned(room *entity.LiveRoom) bool {
	if room == nil || !room.Ban {
		return false
	}
	if room.BanApplyTime == nil {
		return false
	}
	return room.BanApplyTime.After(time.Now())
}

// BanAnchor CMS封禁主播:写入 live_rooms 封禁状态并向 App 推送
func BanAnchor(_ context.Context, req *accountdto.BanAnchorReq) (*accountdto.BanRes, error) {
	user := userinfodao.GetUserInfoByUserId(req.AccountId)
	if user == nil || !user.IsAnchor {
		return nil, errercode.CreateCode(errercode.LiveRoomNotAnchor)
	}
	if req.BanApplyTime == nil || !req.BanApplyTime.After(time.Now()) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	banReason := strings.TrimSpace(req.BanReason)
	if banReason == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	room := EnsureAnchorRoom(req.AccountId, user.GuildId)

	room.SetBan(true)
	room.SetBanApplyTime(req.BanApplyTime)
	room.SetBanReason(banReason)
	NotifyAnchorBanned(req.AccountId, req.BanApplyTime, banReason)
	return &accountdto.BanRes{}, nil
}

// UnBanAnchor CMS解封主播直播间
func UnBanAnchor(_ context.Context, req *accountdto.UnBanAnchorReq) (bool, error) {
	room := liveroomdao.GetRoomByAnchor(req.AccountId)
	if room == nil {
		return false, errercode.CreateCode(errercode.LiveRoomNotExist)
	}
	room.SetBan(false)
	room.SetBanApplyTime(nil)
	room.SetBanReason("")
	return true, nil
}
