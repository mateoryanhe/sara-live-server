package userinfo

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/userinfodto"
	"xr-game-server/errercode"
)

// ReportInviter App端上报邀请者用户ID(仅首次生效,写入 user_exts)
func ReportInviter(ctx context.Context, req *userinfodto.ReportInviterReq) (*userinfodto.ReportInviterRes, error) {
	if req == nil || req.InviterId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if req.InviterId == userId {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if accountdao.GetAccountById(req.InviterId) == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	ext := userinfodao.GetUserExtByUserId(userId)
	if ext.InviterId == 0 {
		ext.SetInviterId(req.InviterId)
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
