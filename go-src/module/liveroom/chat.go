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
	"xr-game-server/errercode"
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

	sender := userinfodao.GetUserInfoByUserId(senderId)
	payload := &liveroomdto.ChatPushItem{
		RoomId:   strconv.FormatUint(req.RoomId, 10),
		SenderId: strconv.FormatUint(senderId, 10),
		Content:  content,
		SentAt:   time.Now().Unix(),
	}
	if sender != nil {
		payload.SenderName = sender.Nickname
		payload.SenderAvatar = upload.GetUrlByName(sender.Avatar)
	}

	for _, o := range liveroomdao.GetOnlinesByRoom(req.RoomId) {
		push.Data(o.UserId, cmd.LiveRoomChat, payload)
	}

	return &liveroomdto.SendChatRes{Success: true}, nil
}
