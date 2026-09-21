package userinfo

import (
	"context"
	"strings"
	"time"
	"xr-game-server/core/event"
	"xr-game-server/gameevent"

	"xr-game-server/constants/country"
	"xr-game-server/constants/followstatus"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/livefollowdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dao/userloginlocationdao"
	"xr-game-server/dto/userinfodto"
	"xr-game-server/errercode"
	"xr-game-server/module/aliyunmoderation"
	"xr-game-server/module/anchorrank"
	"xr-game-server/module/countryflagdeploy"
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
	userExt := userinfodao.GetUserExtByUserId(targetUserId)
	prettyId := userExt.PrettyId
	if prettyId == 0 {
		prettyId = data.ID
	}
	ret := &userinfodto.GetUserInfoRes{
		UserId:        data.ID,
		Nickname:      data.Nickname,
		Phone:         data.Phone,
		Avatar:        upload.ResolveAvatarUrlForUser(data.ID, data.Avatar),
		Remark:        data.Remark,
		Gold:          data.Gold,
		Diamond:       data.Diamond,
		ShareCode:     data.ShareCode,
		VipLevel:      data.VipLevel,
		IsAnchor:      data.IsAnchor(),
		UserType:      data.UserType,
		PrettyId:      prettyId,
		Gender:        data.Gender,
		Birthday:      formatBirthday(data.Birthday),
		FollowCount:   int(userExt.FollowCount),
		FollowerCount: int(userExt.FollowerCount),
		FollowStatus:  resolveFollowStatus(authUserId, targetUserId),
		TotalIncome:   float64(anchorrank.GetUserLast30DayRevenue(targetUserId)),
		Age:           calcAge(data.Birthday),
		FlagIcon:      ResolveUserFlagIcon(targetUserId),
	}
	if req.UserId == 0 {
		now := time.Now()
		data.SetLastLoginTime(&now)
		userinfodao.PublishUserInfo(data)
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

// ResolveUserFlagIcon 优先注册国,空则登录国;拼当前国旗版本完整 URL。
func ResolveUserFlagIcon(userId uint64) string {
	return ResolveUserFlagIcons([]uint64{userId})[userId]
}

// ResolveUserFlagIcons 批量解析用户国旗完整 URL，列表接口只查询一次登录地域表。
func ResolveUserFlagIcons(userIds []uint64) map[uint64]string {
	result := make(map[uint64]string, len(userIds))
	locations := userloginlocationdao.GetByUserIds(userIds)
	version := countryflagdeploy.CurrentVersion()
	for userId, location := range locations {
		code := strings.TrimSpace(location.RegisterCountry)
		if code == "" {
			code = strings.TrimSpace(location.LoginCountry)
		}
		if !country.Exists(code) {
			continue
		}
		rel := country.RelPath(code, version)
		if rel == "" {
			continue
		}
		result[userId] = upload.GetUrlByName(rel)
	}
	return result
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
	userinfodao.PublishUserInfo(data)
	return &userinfodto.UpdateNicknameRes{
		Nickname: data.Nickname,
	}, nil
}

// UploadAvatar 上传头像,保存图片并写回 user_infos.avatar
func UploadAvatar(ctx context.Context, req *userinfodto.UploadAvatarReq) (res *userinfodto.UploadAvatarRes, err error) {
	name, err := upload.UploadImageForApp(ctx, req.File)
	if err != nil {
		return nil, err
	}
	res, err = UpdateAvatarFromStoredFile(ctx, name)
	if err != nil {
		upload.DeleteUploadedFile(name)
	}
	return res, err
}

// UpdateAvatarFromStoredFile 将上传模块返回的对象路径写入用户资料。
func UpdateAvatarFromStoredFile(ctx context.Context, name string) (res *userinfodto.UploadAvatarRes, err error) {
	if strings.TrimSpace(name) == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	userId := httpserver.GetAuthId(ctx)
	data := userinfodao.GetUserInfoByUserId(userId)
	data.SetAvatar(name)
	userinfodao.PublishUserInfo(data)
	return &userinfodto.UploadAvatarRes{
		Avatar: upload.ResolveAvatarUrlForUser(data.ID, data.Avatar),
	}, nil
}
