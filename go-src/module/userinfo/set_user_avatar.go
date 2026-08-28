package userinfo

import (
	"context"
	"strings"

	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

// SetUserAvatar CMS 设置用户头像(avatar 为 uploadFile 返回的文件名,空字符串表示清除)
func SetUserAvatar(_ context.Context, req *accountdto.SetUserAvatarReq) (*accountdto.SetUserAvatarRes, error) {
	if req.AccountId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	user := userinfodao.GetUserInfoByUserId(req.AccountId)
	avatar := strings.TrimSpace(req.Avatar)
	user.SetAvatar(avatar)
	userinfodao.PublishUserInfo(user)
	return &accountdto.SetUserAvatarRes{
		Success: true,
		Avatar:  upload.ResolveAvatarUrlForUser(user.ID, user.Avatar),
	}, nil
}
