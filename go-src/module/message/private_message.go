package message

import (
	"context"
	"strings"
	"time"
	"xr-game-server/constants/cmd"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/push"
	"xr-game-server/dao/messagedao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/messagedto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/sensitiveword"
	"xr-game-server/module/upload"
)

// SendPrivateMessage App端用户之间发送私信
func SendPrivateMessage(ctx context.Context, req *messagedto.AppSendPrivateMessageReq) (*messagedto.AppSendPrivateMessageRes, error) {
	senderId := httpserver.GetAuthId(ctx)
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := sensitiveword.RequireTextCompliant(content); err != nil {
		return nil, err
	}
	if req.ReceiverId == senderId {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if userinfodao.GetUserInfoByUserId(req.ReceiverId) == nil {
		return nil, errercode.CreateCode(errercode.SysError)
	}

	msg := entity.NewUserMessage(entity.UserMessageTypePrivate, senderId, req.ReceiverId, "", content)
	messagedao.AddToCache(msg, req.ReceiverId)
	messagedao.AddToCache(msg, senderId)
	messagedao.InvalidateSenderReceiverMessageCache(req.ReceiverId, senderId)

	pushItem := buildPrivateMessagePushItem(msg)
	push.Data(req.ReceiverId, cmd.PrivateMessagePush, pushItem)
	push.Data(senderId, cmd.PrivateMessagePush, pushItem)

	unReadData := messagedao.GetUnReadByUserId(req.ReceiverId)
	unReadData.AddPrivateUnread(1)

	unReadDetail := messagedao.GetUnreadDetailByReceiverSender(req.ReceiverId, senderId)
	unReadDetail.AddUnread(1)
	messagedao.UpsertUnreadDetailToListCache(unReadDetail)

	return &messagedto.AppSendPrivateMessageRes{
		MessageId: msg.ID,
		Success:   true,
	}, nil
}

// ListPrivateMessageUnread App端分页查询私信未读明细
func ListPrivateMessageUnread(ctx context.Context, req *messagedto.AppPrivateMessageUnreadListReq) (*messagedto.AppPrivateMessageUnreadListRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	pageIndex := req.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}

	rows := messagedao.ListUnreadDetailByReceiverId(userId, pageIndex)
	list := make([]*messagedto.AppPrivateMessageUnreadDetailItem, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, toPrivateMessageUnreadDetailItem(row))
	}
	return &messagedto.AppPrivateMessageUnreadListRes{List: list}, nil
}

// ListPrivateMessageBySender App端按发送者查询私信内容(最多200条,超出返回空列表)
func ListPrivateMessageBySender(ctx context.Context, req *messagedto.AppPrivateMessageBySenderReq) (*messagedto.AppPrivateMessageBySenderRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if req.SenderId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	pageIndex := req.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 40
	}

	rows := messagedao.ListByReceiverAndSender(userId, req.SenderId, pageIndex, pageSize)
	list := make([]*messagedto.AppPrivateMessageItem, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, toPrivateMessageItem(row))
	}
	return &messagedto.AppPrivateMessageBySenderRes{List: list}, nil
}

// ClearPrivateMessageUnread App端清除指定玩家的私信未读
func ClearPrivateMessageUnread(ctx context.Context, req *messagedto.AppClearPrivateMessageUnreadReq) (*messagedto.AppClearPrivateMessageUnreadRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if req.SenderId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	unReadData := messagedao.GetUnReadByUserId(userId)
	clearedCount := uint64(0)

	unReadDetail := messagedao.GetUnreadDetailByReceiverSender(userId, req.SenderId)
	if unReadDetail != nil && unReadDetail.UnreadCount > 0 {
		clearedCount = unReadDetail.UnreadCount
		unReadDetail.ClearUnread()
		unReadData.SubPrivateUnread(clearedCount)
		messagedao.UpsertUnreadDetailToListCache(unReadDetail)
	}

	return &messagedto.AppClearPrivateMessageUnreadRes{
		Success:       true,
		ClearedCount:  clearedCount,
		PrivateUnread: unReadData.PrivateUnread,
	}, nil
}

func buildPrivateMessagePushItem(msg *entity.UserMessage) *messagedto.PrivateMessagePushItem {
	item := &messagedto.PrivateMessagePushItem{
		Id:         msg.ID,
		SenderId:   msg.SenderId,
		ReceiverId: msg.ReceiverId,
		Content:    msg.Content,
		SentAt:     formatMessageTime(msg.CreatedAt),
	}
	if sender := userinfodao.GetUserInfoByUserId(msg.SenderId); sender != nil {
		item.SenderName = sender.Nickname
		item.SenderAvatar = upload.GetUrlByName(sender.Avatar)
	}
	return item
}

func toPrivateMessageUnreadDetailItem(row *entity.UserMessageUnreadDetail) *messagedto.AppPrivateMessageUnreadDetailItem {
	item := &messagedto.AppPrivateMessageUnreadDetailItem{
		SenderId:    row.SenderId,
		ReceiverId:  row.ReceiverId,
		UnreadCount: row.UnreadCount,
		UpdatedAt:   formatMessageTime(row.UpdatedAt),
	}
	if sender := userinfodao.GetUserInfoByUserId(row.SenderId); sender != nil {
		item.SenderName = sender.Nickname
		item.SenderAvatar = upload.GetUrlByName(sender.Avatar)
	}
	return item
}

func toPrivateMessageItem(msg *entity.UserMessage) *messagedto.AppPrivateMessageItem {
	item := &messagedto.AppPrivateMessageItem{
		Id:         msg.ID,
		SenderId:   msg.SenderId,
		ReceiverId: msg.ReceiverId,
		Content:    msg.Content,
		CreatedAt:  formatMessageTime(msg.CreatedAt),
	}
	if sender := userinfodao.GetUserInfoByUserId(msg.SenderId); sender != nil {
		item.SenderName = sender.Nickname
		item.SenderAvatar = upload.GetUrlByName(sender.Avatar)
	}
	return item
}

func formatMessageTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
