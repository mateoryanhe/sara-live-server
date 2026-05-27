package liveroom

import (
	"context"
	"strings"
	"time"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

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
