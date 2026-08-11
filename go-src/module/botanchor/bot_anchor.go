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
	room.SetIsTest(req.IsTest)
	userinfodao.GetUserCumulativeStatByUserId(account.ID)
	addBotAnchorId(account.ID)
	addEnabledBotAnchorId(account.ID)
	liveroomdao.FlushRoomCache(room)

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
	room.SetIsTest(req.IsTest)

	return &botanchordto.UpdateBotAnchorRes{Success: true}, nil
}

// SetBotAnchorStatus CMS启用/停用机器人主播
func SetBotAnchorStatus(ctx context.Context, req *botanchordto.SetBotAnchorStatusReq) (*botanchordto.SetBotAnchorStatusRes, error) {
	user, err := getBotAnchorUser(req.ID)
	if err != nil {
		return nil, err
	}
	if req.Status != entity.BotAnchorStatusDisabled && req.Status != entity.BotAnchorStatusEnabled {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	user.SetBotAnchorStatus(req.Status)
	switch req.Status {
	case entity.BotAnchorStatusEnabled:
		enableBotAnchorRoomCache(req.ID, user.GuildId)
	case entity.BotAnchorStatusDisabled:
		disableBotAnchorRoomCache(req.ID)
	}
	liveroom.RefreshRoomListCache(ctx)
	return &botanchordto.SetBotAnchorStatusRes{Success: true}, nil
}

// StartBotAnchorLive CMS机器人主播开播(按是否推流决定是否调用声网云播放)
func StartBotAnchorLive(ctx context.Context, req *botanchordto.StartBotAnchorLiveReq) (*botanchordto.StartBotAnchorLiveRes, error) {
	user, err := getBotAnchorUser(req.ID)
	if err != nil {
		return nil, err
	}
	if user.BotAnchorStatus != entity.BotAnchorStatusEnabled {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := liveroom.StartLiveForBotAnchor(ctx, req.ID, user.GuildId); err != nil {
		return nil, err
	}
	return &botanchordto.StartBotAnchorLiveRes{Success: true}, nil
}

// StopBotAnchorLive CMS机器人主播下播
func StopBotAnchorLive(ctx context.Context, req *botanchordto.StopBotAnchorLiveReq) (*botanchordto.StopBotAnchorLiveRes, error) {
	if _, err := getBotAnchorUser(req.ID); err != nil {
		return nil, err
	}
	if err := liveroom.StopLiveForBotAnchor(ctx, req.ID); err != nil {
		return nil, err
	}
	return &botanchordto.StopBotAnchorLiveRes{Success: true}, nil
}

// BatchStartBotAnchorLive CMS批量机器人主播开播
func BatchStartBotAnchorLive(ctx context.Context, req *botanchordto.BatchStartBotAnchorLiveReq) (*botanchordto.BatchBotAnchorLiveRes, error) {
	return batchBotAnchorLive(ctx, req.IDs, func(id uint64) error {
		_, err := StartBotAnchorLive(ctx, &botanchordto.StartBotAnchorLiveReq{ID: id})
		return err
	})
}

// BatchStopBotAnchorLive CMS批量机器人主播下播
func BatchStopBotAnchorLive(ctx context.Context, req *botanchordto.BatchStopBotAnchorLiveReq) (*botanchordto.BatchBotAnchorLiveRes, error) {
	return batchBotAnchorLive(ctx, req.IDs, func(id uint64) error {
		_, err := StopBotAnchorLive(ctx, &botanchordto.StopBotAnchorLiveReq{ID: id})
		return err
	})
}

func batchBotAnchorLive(ctx context.Context, ids []uint64, action func(uint64) error) (*botanchordto.BatchBotAnchorLiveRes, error) {
	res := &botanchordto.BatchBotAnchorLiveRes{}
	seen := make(map[uint64]struct{}, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		if err := action(id); err != nil {
			res.FailCount++
			res.FailIds = append(res.FailIds, id)
			continue
		}
		res.SuccessCount++
	}
	if res.SuccessCount > 0 {
		liveroom.RefreshRoomListCache(ctx)
	}
	return res, nil
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
