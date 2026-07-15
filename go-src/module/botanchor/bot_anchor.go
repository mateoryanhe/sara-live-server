package botanchor

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/util/guid"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/botanchordto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/auth"
	"xr-game-server/module/liveroom"
	"xr-game-server/module/userinfo"
)

// QueryBotAnchorList CMS分页查询机器人主播列表(基于内存ID列表)
func QueryBotAnchorList(_ context.Context, req *botanchordto.QueryBotAnchorListReq) (*httpserver.CMSQueryResp, error) {
	total, data := queryBotAnchorListFromMemory(req)
	return httpserver.NewCMSQueryResp(total, data), nil
}

// CreateBotAnchor CMS创建机器人主播(写入 accounts / user_infos / live_rooms)
func CreateBotAnchor(_ context.Context, req *botanchordto.CreateBotAnchorReq) (*botanchordto.CreateBotAnchorRes, error) {
	nickname := strings.TrimSpace(req.Nickname)
	if nickname == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	category := req.Category
	if category == 0 {
		category = entity.LiveRoomCategoryHot
	}
	if err := validateRoomTagId(req.TagId); err != nil {
		return nil, err
	}

	openId := fmt.Sprintf("bot_%s", guid.S())
	account := entity.NewAccount(openId, auth.BotAnchorChannel)

	user := userinfodao.GetUserInfoByUserId(account.ID)
	user.SetNickname(nickname)
	if req.Avatar != "" {
		user.SetAvatar(strings.TrimSpace(req.Avatar))
	}
	user.SetUserType(entity.UserTypeBotAnchor)
	user.SetBotAnchorStatus(entity.BotAnchorStatusEnabled)
	user.SetHasLiveRoom(true)
	if req.GuildId > 0 {
		user.SetGuildId(req.GuildId)
	}

	room := liveroom.EnsureAnchorRoom(account.ID, user.GuildId)
	room.SetTitle(strings.TrimSpace(req.RoomTitle))
	room.SetCategory(category)
	room.SetPrivateInviteType(entity.DefaultPrivateInviteType(category))
	room.SetTagId(req.TagId)
	if req.CloudPlayerVideo != "" {
		room.SetCloudPlayerVideo(strings.TrimSpace(req.CloudPlayerVideo))
	}
	room.SetPushStream(req.PushStream)
	userinfodao.GetUserCumulativeStatByUserId(account.ID)
	addBotAnchorId(account.ID)
	userinfo.AddIdToCache(account.ID)

	return &botanchordto.CreateBotAnchorRes{ID: account.ID}, nil
}

// UpdateBotAnchor CMS更新机器人主播资料
func UpdateBotAnchor(_ context.Context, req *botanchordto.UpdateBotAnchorReq) (*botanchordto.UpdateBotAnchorRes, error) {
	user, err := getBotAnchorUser(req.ID)
	if err != nil {
		return nil, err
	}
	nickname := strings.TrimSpace(req.Nickname)
	if nickname == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := validateRoomTagId(req.TagId); err != nil {
		return nil, err
	}
	user.SetNickname(nickname)
	if req.Avatar != nil {
		user.SetAvatar(strings.TrimSpace(*req.Avatar))
	}

	category := req.Category
	if category == 0 {
		category = entity.LiveRoomCategoryHot
	}
	room := liveroom.EnsureAnchorRoom(req.ID, user.GuildId)
	room.SetTitle(strings.TrimSpace(req.RoomTitle))
	room.SetCategory(category)
	room.SetPrivateInviteType(entity.DefaultPrivateInviteType(category))
	room.SetTagId(req.TagId)
	if req.CloudPlayerVideo != nil {
		room.SetCloudPlayerVideo(strings.TrimSpace(*req.CloudPlayerVideo))
	}
	room.SetPushStream(req.PushStream)

	userinfo.AddIdToCache(req.ID)
	return &botanchordto.UpdateBotAnchorRes{Success: true}, nil
}

// SetBotAnchorStatus CMS启用/停用机器人主播
func SetBotAnchorStatus(_ context.Context, req *botanchordto.SetBotAnchorStatusReq) (*botanchordto.SetBotAnchorStatusRes, error) {
	user, err := getBotAnchorUser(req.ID)
	if err != nil {
		return nil, err
	}
	if req.Status != entity.BotAnchorStatusDisabled && req.Status != entity.BotAnchorStatusEnabled {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	user.SetBotAnchorStatus(req.Status)
	userinfo.AddIdToCache(req.ID)
	return &botanchordto.SetBotAnchorStatusRes{Success: true}, nil
}

func getBotAnchorUser(userId uint64) (*entity.UserInfo, error) {
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if !isBotAnchorId(userId) {
		return nil, errercode.CreateCode(errercode.BotAnchorNonExist)
	}
	return userinfodao.GetUserInfoByUserId(userId), nil
}

func validateRoomTagId(tagId uint64) error {
	if tagId == 0 {
		return nil
	}
	if liveroomdao.GetRoomTagById(tagId) == nil {
		return errercode.CreateCode(errercode.LiveRoomTagNonExist)
	}
	return nil
}
