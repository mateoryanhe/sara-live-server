package message

import (
	"context"
	"strings"
	"xr-game-server/constants/cmd"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/push"
	"xr-game-server/dao/messagedao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/messagedto"
	"xr-game-server/entity"
	"xr-game-server/gameevent"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func onSystemMessage(val any) {
	data, ok := val.(*gameevent.SystemMessageEventData)
	if !ok || data == nil || data.ReceiverId == 0 {
		g.Log().Errorf(gctx.New(), "SystemMessageEvent payload type error: %T", val)
		return
	}
	title := strings.TrimSpace(data.Title)
	content := strings.TrimSpace(data.Content)
	if content == "" {
		return
	}
	if userinfodao.GetUserInfoByUserId(data.ReceiverId) == nil {
		return
	}

	msg := entity.NewUserMessage(entity.UserMessageTypeSystem, 0, data.ReceiverId, title, content)
	messagedao.AddSystemToCache(msg)

	pushItem := buildSystemMessagePushItem(msg)
	push.Data(data.ReceiverId, cmd.SystemMessagePush, pushItem)

	//加入未读数
	unReadData := messagedao.GetUnReadByUserId(data.ReceiverId)
	unReadData.AddSystemUnread(1)
}

// ListSystemMessage App端分页查询当前用户的系统消息
func ListSystemMessage(ctx context.Context, req *messagedto.AppSystemMessageListReq) (*messagedto.AppSystemMessageListRes, error) {
	userId := httpserver.GetAuthId(ctx)
	pageIndex := req.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}

	rows := messagedao.ListSystemByReceiverId(userId, pageIndex)
	list := make([]*messagedto.AppSystemMessageItem, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, toSystemMessageItem(row))
	}
	return &messagedto.AppSystemMessageListRes{List: list}, nil
}

func buildSystemMessagePushItem(msg *entity.UserMessage) *messagedto.SystemMessagePushItem {
	return &messagedto.SystemMessagePushItem{
		Id:         msg.ID,
		ReceiverId: msg.ReceiverId,
		Title:      msg.Title,
		Content:    msg.Content,
		SentAt:     formatMessageTime(msg.CreatedAt),
	}
}

func toSystemMessageItem(msg *entity.UserMessage) *messagedto.AppSystemMessageItem {
	return &messagedto.AppSystemMessageItem{
		Id:         msg.ID,
		ReceiverId: msg.ReceiverId,
		Title:      msg.Title,
		Content:    msg.Content,
		CreatedAt:  formatMessageTime(msg.CreatedAt),
	}
}
