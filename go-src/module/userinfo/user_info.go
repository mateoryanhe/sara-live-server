package userinfo

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/userinfodto"
)

// GetUserInfo 查询当前登录用户的基础信息
func GetUserInfo(ctx context.Context, req *userinfodto.GetUserInfoReq) (res *userinfodto.GetUserInfoRes, err error) {
	userId := httpserver.GetAuthId(ctx)
	data := userinfodao.GetUserInfoByUserId(userId)
	return &userinfodto.GetUserInfoRes{
		UserId:   data.ID,
		Nickname: data.Nickname,
		Phone:    data.Phone,
		Avatar:   data.Avatar,
		Remark:   data.Remark,
	}, nil
}

// UpdateNickname 修改昵称
func UpdateNickname(ctx context.Context, req *userinfodto.UpdateNicknameReq) (res *userinfodto.UpdateNicknameRes, err error) {
	userId := httpserver.GetAuthId(ctx)
	data := userinfodao.GetUserInfoByUserId(userId)
	data.SetNickname(req.Nickname)
	return &userinfodto.UpdateNicknameRes{
		Nickname: data.Nickname,
	}, nil
}
