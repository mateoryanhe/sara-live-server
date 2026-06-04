package userinfo

import (
	"context"
	"strings"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/userinfodto"
	"xr-game-server/errercode"
	"xr-game-server/module/sensitiveword"
	"xr-game-server/module/upload"
)

const (
	Def_Url = "https://img.yonogames.com/headimg/man/147.png"
)

// GetUserInfo 查询当前登录用户的基础信息
func GetUserInfo(ctx context.Context, req *userinfodto.GetUserInfoReq) (res *userinfodto.GetUserInfoRes, err error) {
	userId := httpserver.GetAuthId(ctx)
	data := userinfodao.GetUserInfoByUserId(userId)
	ret := &userinfodto.GetUserInfoRes{
		UserId:      data.ID,
		Nickname:    data.Nickname,
		Phone:       data.Phone,
		Avatar:      upload.GetUrlByName(data.Avatar),
		Remark:      data.Remark,
		Gold:        data.Gold,
		Diamond:     data.Diamond,
		ShareCode:   data.ShareCode,
		VipLevel:    data.VipLevel,
		IsAnchor:    data.IsAnchor,
		HasLiveRoom: data.HasLiveRoom,
		Gender:      data.Gender,
		Birthday:    formatBirthday(data.Birthday),
	}
	if data.Avatar == "" {
		ret.Avatar = Def_Url
	}
	return ret, nil
}

// UpdateNickname 修改昵称
func UpdateNickname(ctx context.Context, req *userinfodto.UpdateNicknameReq) (res *userinfodto.UpdateNicknameRes, err error) {
	nickname := strings.TrimSpace(req.Nickname)
	if nickname == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := sensitiveword.RequireTextCompliant(nickname); err != nil {
		return nil, err
	}

	userId := httpserver.GetAuthId(ctx)
	data := userinfodao.GetUserInfoByUserId(userId)
	data.SetNickname(nickname)
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
