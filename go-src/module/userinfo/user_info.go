package userinfo

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/userinfodto"
	"xr-game-server/module/upload"
)

// GetUserInfo 查询当前登录用户的基础信息
func GetUserInfo(ctx context.Context, req *userinfodto.GetUserInfoReq) (res *userinfodto.GetUserInfoRes, err error) {
	userId := httpserver.GetAuthId(ctx)
	data := userinfodao.GetUserInfoByUserId(userId)
	return &userinfodto.GetUserInfoRes{
		UserId:    data.ID,
		Nickname:  data.Nickname,
		Phone:     data.Phone,
		Avatar:    upload.GetUrlByName(data.Avatar),
		Remark:    data.Remark,
		Gold:      data.Gold,
		Diamond:   data.Diamond,
		ShareCode: data.ShareCode,
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

// UploadAvatar 上传头像,保存图片并写回 user_infos.avatar
func UploadAvatar(ctx context.Context, req *userinfodto.UploadAvatarReq) (res *userinfodto.UploadAvatarRes, err error) {
	userId := httpserver.GetAuthId(ctx)
	name, err := upload.UploadImage(req.File)
	if err != nil {
		return nil, err
	}
	data := userinfodao.GetUserInfoByUserId(userId)
	data.SetAvatar(name)
	return &userinfodto.UploadAvatarRes{
		Avatar: name,
	}, nil
}
