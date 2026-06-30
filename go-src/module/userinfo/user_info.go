package userinfo

import (
	"context"
	"strings"
	"xr-game-server/core/event"
	"xr-game-server/gameevent"

	"xr-game-server/constants/followstatus"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/livefollowdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/userinfodto"
	"xr-game-server/errercode"
	"xr-game-server/module/aliyunmoderation"
	"xr-game-server/module/anchorrank"
	"xr-game-server/module/upload"
)

// GetUserInfo 查询用户基础信息(不传 userId 时查当前登录用户)
func GetUserInfo(ctx context.Context, req *userinfodto.GetUserInfoReq) (res *userinfodto.GetUserInfoRes, err error) {
	authUserId := httpserver.GetAuthId(ctx)
	targetUserId := authUserId
	if req.UserId > 0 {
		targetUserId = req.UserId
	}
	data := userinfodao.GetUserInfoByUserId(targetUserId)
	ret := &userinfodto.GetUserInfoRes{
		UserId:        data.ID,
		Nickname:      data.Nickname,
		Phone:         data.Phone,
		Avatar:        upload.ResolveAvatarUrl(data.Avatar),
		Remark:        data.Remark,
		Gold:          data.Gold,
		Diamond:       data.Diamond,
		ShareCode:     data.ShareCode,
		VipLevel:      data.VipLevel,
		IsAnchor:      data.IsAnchor,
		HasLiveRoom:   data.HasLiveRoom,
		Gender:        data.Gender,
		Birthday:      formatBirthday(data.Birthday),
		FollowCount:   len(livefollowdao.GetFollowingsByUser(targetUserId)),
		FollowerCount: len(livefollowdao.GetFollowersByAnchor(targetUserId)),
		FollowStatus:  resolveFollowStatus(authUserId, targetUserId),
		TotalIncome:   float64(anchorrank.GetUserLast30DayRevenue(targetUserId)),
		Age:           calcAge(data.Birthday),
	}
	if req.UserId == authUserId {
		event.Pub(gameevent.LoginEvent, authUserId)
	}
	return ret, nil
}

func resolveFollowStatus(viewerId, targetUserId uint64) uint8 {
	if viewerId == 0 || targetUserId == 0 || viewerId == targetUserId {
		return followstatus.NotFollowing
	}
	iFollow := livefollowdao.IsFollowing(viewerId, targetUserId)
	theyFollow := livefollowdao.IsFollowing(targetUserId, viewerId)
	if iFollow && theyFollow {
		return followstatus.Mutual
	}
	if iFollow {
		return followstatus.Following
	}
	return followstatus.NotFollowing
}

// UpdateNickname 修改昵称
func UpdateNickname(ctx context.Context, req *userinfodto.UpdateNicknameReq) (res *userinfodto.UpdateNicknameRes, err error) {
	nickname := strings.TrimSpace(req.Nickname)
	if nickname == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := aliyunmoderation.RequireTextCompliant(aliyunmoderation.SceneNickname, nickname); err != nil {
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
	name, err := upload.UploadImageForApp(ctx, req.File)
	if err != nil {
		return nil, err
	}
	data := userinfodao.GetUserInfoByUserId(userId)
	data.SetAvatar(name)
	return &userinfodto.UploadAvatarRes{
		Avatar: name,
	}, nil
}
