package liveroom

import (
	"context"
	"strconv"
	"strings"
	"time"
	"xr-game-server/constants/cmd"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/push"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/aliyunmoderation"
	"xr-game-server/module/upload"
)

// SendChat 直播间文字消息
//  1. 校验房间存在、消息内容非空
//  2. 构造推送载荷,向房间内全体在线用户(含发送者自身)推送 cmd.LiveRoomChat
//
// 注:消息只走推送通道,不做持久化(与送礼一致),业务侧需要历史记录可后续再加。
func SendChat(ctx context.Context, req *liveroomdto.SendChatReq) (*liveroomdto.SendChatRes, error) {
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, errercode.CreateCode(errercode.SysError)
	}

	senderId := httpserver.GetAuthId(ctx)

	if liveroomdao.GetRoomById(req.RoomId) == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}

	if err := aliyunmoderation.RequireTextCompliant(aliyunmoderation.SceneChat, content); err != nil {
		return nil, err
	}

	// 主播本人不受禁言限制;观众需校验在线记录上的禁言标记
	if senderId != req.RoomId {
		onlineId := entity.BuildLiveRoomOnlineId(senderId, req.RoomId)
		if online := liveroomdao.GetOnlineById(onlineId, senderId, req.RoomId); online != nil && online.Muted {
			return nil, errercode.CreateCode(errercode.LiveRoomChatMuted)
		}
	}

	pushRoomChat(req.RoomId, senderId, content)
	return &liveroomdto.SendChatRes{Success: true}, nil
}

// SendPrivateRoomChat 私密房文字消息(不执行禁言逻辑,仅发送者与目标用户可见)
func SendPrivateRoomChat(ctx context.Context, req *liveroomdto.SendPrivateRoomChatReq) (*liveroomdto.SendPrivateRoomChatRes, error) {
	content := strings.TrimSpace(req.Content)

	if content == "" {
		return nil, errercode.CreateCode(errercode.SysError)
	}

	senderId := httpserver.GetAuthId(ctx)
	targetId := req.TargetId
	if targetId == 0 || targetId == senderId {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if userinfodao.GetUserInfoByUserId(targetId) == nil {
		return nil, errercode.CreateCode(errercode.SysError)
	}

	if err := aliyunmoderation.RequireTextCompliant(aliyunmoderation.SceneChat, content); err != nil {
		return nil, err
	}

	pushPrivateRoomChat(senderId, targetId, content)
	return &liveroomdto.SendPrivateRoomChatRes{Success: true}, nil
}

func resolvePrivateRoomAnchorId(userId, targetId uint64) uint64 {
	if room := liveroomdao.GetRoomById(userId); room != nil && room.Category == entity.LiveRoomCategoryPrivate {
		return userId
	}
	if room := liveroomdao.GetRoomById(targetId); room != nil && room.Category == entity.LiveRoomCategoryPrivate {
		return targetId
	}
	return 0
}

func pushPrivateRoomChat(senderId, targetId uint64, content string) {
	sender := userinfodao.GetUserInfoByUserId(senderId)
	payload := &liveroomdto.PrivateRoomChatPushItem{
		TargetId: strconv.FormatUint(targetId, 10),
		SenderId: strconv.FormatUint(senderId, 10),
		Content:  content,
		SentAt:   time.Now().Unix(),
	}
	if sender != nil {
		payload.SenderName = sender.Nickname
		payload.SenderAvatar = upload.ResolveAvatarUrlForUser(sender.ID, sender.Avatar)
		payload.VipLevel = sender.VipLevel
	}

	push.Data(senderId, cmd.LiveRoomPrivateChat, payload)
	push.Data(targetId, cmd.LiveRoomPrivateChat, payload)
}

func pushRoomChat(roomId, senderId uint64, content string) {
	sender := userinfodao.GetUserInfoByUserId(senderId)
	payload := &liveroomdto.ChatPushItem{
		RoomId:   strconv.FormatUint(roomId, 10),
		SenderId: strconv.FormatUint(senderId, 10),
		Content:  content,
		SentAt:   time.Now().Unix(),
	}
	if sender != nil {
		payload.SenderName = sender.Nickname
		payload.SenderAvatar = upload.ResolveAvatarUrlForUser(sender.ID, sender.Avatar)
		payload.VipLevel = sender.VipLevel
	}

	for _, o := range getOnline(roomId) {
		push.Data(o, cmd.LiveRoomChat, payload)
	}
	push.Data(roomId, cmd.LiveRoomChat, payload)
}
