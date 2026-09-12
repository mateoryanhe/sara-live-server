package userinfo

import (
	"context"
	"strings"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/userinfodto"
	"xr-game-server/errercode"
)

// ResolveInviteCodeToUserId 邀请码 → 邀请者 userId。
// 当前仅支持 shareCode；后续新格式在此扩展(如带前缀码)，勿散落各调用点。
func ResolveInviteCodeToUserId(inviteCode string) uint64 {
	code := strings.TrimSpace(inviteCode)
	if code == "" {
		return 0
	}
	// 现有格式: user_infos.share_code
	if id := userinfodao.GetUserIdByShareCode(code); id != 0 {
		return id
	}
	// TODO: 新格式邀请码解析
	return 0
}

// ReportInviter App端上报邀请码(仅首次生效,写入 user_exts)
func ReportInviter(ctx context.Context, req *userinfodto.ReportInviterReq) (*userinfodto.ReportInviterRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	inviterId := ResolveInviteCodeToUserId(req.InviteCode)
	if inviterId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if inviterId == userId {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if accountdao.GetAccountById(inviterId) == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	ext := userinfodao.GetUserExtByUserId(userId)
	if ext.InviterId == 0 {
		ext.SetInviterId(inviterId)
	}
	// 与历史 user_infos.inviter_id 保持一致(仅空时写入)
	info := userinfodao.GetUserInfoByUserId(userId)
	if info.InviterId == 0 {
		info.SetInviterId(ext.InviterId)
		userinfodao.PublishUserInfo(info)
	}

	return &userinfodto.ReportInviterRes{
		Success:   true,
		InviterId: ext.InviterId,
	}, nil
}
