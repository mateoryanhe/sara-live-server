package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/messagedto"
	"xr-game-server/module/message"
)

const MessageAppUrl = "/message"

type MessageAppController struct{}

func initMessageAppController() {
	httpserver.RegAPI(MessageAppUrl, &MessageAppController{})
}

// SendPrivateMessage App端发送私信
func (c *MessageAppController) SendPrivateMessage(ctx context.Context, req *messagedto.AppSendPrivateMessageReq) (*messagedto.AppSendPrivateMessageRes, error) {
	return message.SendPrivateMessage(ctx, req)
}

// PrivateMessageUnreadList App端查询私信未读明细列表
func (c *MessageAppController) PrivateMessageUnreadList(ctx context.Context, req *messagedto.AppPrivateMessageUnreadListReq) (*messagedto.AppPrivateMessageUnreadListRes, error) {
	return message.ListPrivateMessageUnread(ctx, req)
}

// ClearPrivateMessageUnread App端清除指定玩家私信未读
func (c *MessageAppController) ClearPrivateMessageUnread(ctx context.Context, req *messagedto.AppClearPrivateMessageUnreadReq) (*messagedto.AppClearPrivateMessageUnreadRes, error) {
	return message.ClearPrivateMessageUnread(ctx, req)
}

// SystemMessageList App端查询系统消息列表
func (c *MessageAppController) SystemMessageList(ctx context.Context, req *messagedto.AppSystemMessageListReq) (*messagedto.AppSystemMessageListRes, error) {
	return message.ListSystemMessage(ctx, req)
}

// ClearSystemMessageUnread App端清除系统消息未读(每次1条)
func (c *MessageAppController) ClearSystemMessageUnread(ctx context.Context, req *messagedto.AppClearSystemMessageUnreadReq) (*messagedto.AppClearSystemMessageUnreadRes, error) {
	return message.ClearSystemMessageUnread(ctx, req)
}
