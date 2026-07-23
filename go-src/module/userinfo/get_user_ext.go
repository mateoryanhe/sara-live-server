package userinfo

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/userinfodto"
)

// GetUserExt 查询用户扩展信息(不传 userId 时查当前登录用户)
func GetUserExt(ctx context.Context, req *userinfodto.GetUserExtReq) (*userinfodto.GetUserExtRes, error) {
	authUserId := httpserver.GetAuthId(ctx)
	targetUserId := authUserId
	if req.UserId > 0 {
		targetUserId = req.UserId
	}
	ext := userinfodao.GetUserExtByUserId(targetUserId)
	return &userinfodto.GetUserExtRes{
		UserId:        targetUserId,
		CanRank:       ext.CanRank,
		PackageName:   ext.PackageName,
		AppVersion:    ext.AppVersion,
		FollowCount:   ext.FollowCount,
		FollowerCount: ext.FollowerCount,
		CancelCode:    ext.CancelCode,
	}, nil
}
