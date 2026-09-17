package userinfo

import (
	"context"
	"strconv"
	"strings"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/userinfodto"
	"xr-game-server/errercode"
)

// ResolveInviteCodeToUserId 邀请码即邀请者 userId，解析失败返回 0。
func ResolveInviteCodeToUserId(inviteCode string) uint64 {
	code := strings.TrimSpace(inviteCode)
	if code == "" {
		return 0
	}
	userId, err := strconv.ParseUint(code, 10, 64)
	if err != nil {
		return 0
	}
	return userId
}

// GetUserInfoByInviteCode 根据邀请码查询邀请人基础信息，应答与 GetUserInfo 保持一致。
func GetUserInfoByInviteCode(ctx context.Context, req *userinfodto.GetUserInfoByInviteCodeReq) (*userinfodto.GetUserInfoByInviteCodeRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	userId := ResolveInviteCodeToUserId(req.InviteCode)
	if userId == 0 || accountdao.GetAccountById(userId) == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	return GetUserInfo(ctx, &userinfodto.GetUserInfoReq{UserId: userId})
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
