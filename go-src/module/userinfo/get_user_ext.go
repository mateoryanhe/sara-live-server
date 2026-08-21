package userinfo

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/shortvideodao"
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
	if targetUserId == authUserId {
		ext = userinfodao.EnsureCancelCode(targetUserId)
	}
	authorStat := shortvideodao.GetAuthorStatByAuthorId(targetUserId)
	var shortVideoViewCount, shortVideoLikeCount uint64
	var shortVideoTotalIncome float64
	if authorStat != nil {
		shortVideoViewCount = authorStat.ViewCount
		shortVideoLikeCount = authorStat.LikeCount
		shortVideoTotalIncome = authorStat.TotalDiamondIncome
	}
	return &userinfodto.GetUserExtRes{
		UserId:                targetUserId,
		PrettyId:              ext.PrettyId,
		CanRank:               ext.CanRank,
		PackageName:           ext.PackageName,
		AppVersion:            ext.AppVersion,
		FollowCount:           ext.FollowCount,
		FollowerCount:         ext.FollowerCount,
		CancelCode:            ext.CancelCode,
		CancelCodeExpireAt:    ext.CancelCodeExpireAt,
		FirstRecharge:         ext.FirstRecharge,
		ShortVideoViewCount:   shortVideoViewCount,
		ShortVideoTotalIncome: shortVideoTotalIncome,
		ShortVideoLikeCount:   shortVideoLikeCount,
	}, nil
}
