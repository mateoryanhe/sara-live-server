package message

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/cmd"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/push"
	"xr-game-server/core/xrpool"
	"xr-game-server/dao/livefollowdao"
	"xr-game-server/dao/messagedao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/messagedto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/aliyunmoderation"
	"xr-game-server/module/upload"
)

// SendPrivateMessage App端用户之间发送私信
func SendPrivateMessage(ctx context.Context, req *messagedto.AppSendPrivateMessageReq) (*messagedto.AppSendPrivateMessageRes, error) {
	senderId := httpserver.GetAuthId(ctx)
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := aliyunmoderation.RequireTextCompliant(aliyunmoderation.SceneChat, content); err != nil {
		return nil, err
	}
	if req.ReceiverId == senderId {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if userinfodao.GetUserInfoByUserId(req.ReceiverId) == nil {
		return nil, errercode.CreateCode(errercode.SysError)
	}
	if livefollowdao.IsBlocked(senderId, req.ReceiverId) || livefollowdao.IsBlocked(req.ReceiverId, senderId) {
		return nil, errercode.CreateCode(errercode.PrivateMessageUserBlocked)
	}

	msg := entity.NewUserMessage(entity.UserMessageTypePrivate, senderId, req.ReceiverId, "", content)

	//往发送者写入消息
	pushItem := buildPrivateMessagePushItem(msg)
	senderSession := entity.NewUserMessageSession(msg.ID, entity.BuildUserMessageSessionId(senderId, req.ReceiverId))
	push.Data(senderId, cmd.PrivateMessagePush, pushItem)

	//往接收者写入消息
	entity.NewUserMessageSession(msg.ID, entity.BuildUserMessageSessionId(req.ReceiverId, senderId))

	push.Data(req.ReceiverId, cmd.PrivateMessagePush, pushItem)
	//刷新未读消息缓存
	unReadData := messagedao.GetUnReadByUserId(req.ReceiverId)
	unReadData.AddPrivateUnread(1)

	unReadDetail := messagedao.GetUnreadDetailByReceiverSender(senderId, req.ReceiverId)
	unReadDetail.AddUnread(1)
	messagedao.MarkMutualPrivateChat(senderId, req.ReceiverId)
	messagedao.FlushUnreadDetailCache(unReadDetail)
	prependPrivateMessageUnreadListCache(req.ReceiverId, senderId, msg, senderSession.ID, unReadDetail.UnreadCount)

	return &messagedto.AppSendPrivateMessageRes{
		MessageId: msg.ID,
		SessionId: senderSession.ID,
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

	if pageIndex == 1 {
		list, ok := getPrivateMessageUnreadListCache(userId)
		if !ok {
			list = loadPrivateMessageUnreadListCache(userId)
		}
		return &messagedto.AppPrivateMessageUnreadListRes{
			List: pagePrivateMessageUnreadList(list, pageIndex),
		}, nil
	}

	offset := (pageIndex - 1) * privateMessageUnreadListPageSize
	rows := messagedao.ListPrivateMessageUnreadWithLastMessageFromDBLimit(userId, privateMessageUnreadListPageSize, offset)
	list := make([]*messagedto.AppPrivateMessageUnreadDetailItem, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		item := toPrivateMessageUnreadDetailItemFromRow(row)
		fillPrivateMessageUnreadSenderInfo(item)
		list = append(list, item)
	}
	return &messagedto.AppPrivateMessageUnreadListRes{List: list}, nil
}

// ListPrivateMessageBySender App端按目标用户查询私信内容
func ListPrivateMessageBySender(ctx context.Context, req *messagedto.AppPrivateMessageBySenderReq) (*messagedto.AppPrivateMessageBySenderRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if req.TargetId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 40
	}

	sessionId := entity.BuildUserMessageSessionId(userId, req.TargetId)
	rows, hasMore := messagedao.ListByReceiverAndSender(sessionId, req.LastCreatedAt, pageSize)
	list := make([]*messagedto.AppPrivateMessageItem, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.Message == nil {
			continue
		}
		list = append(list, toPrivateMessageItem(row.SessionId, row.Message))
	}
	return &messagedto.AppPrivateMessageBySenderRes{
		List:    list,
		HasMore: hasMore,
	}, nil
}

// ClearPrivateMessageUnread App端清除指定玩家的私信未读
func ClearPrivateMessageUnread(ctx context.Context, req *messagedto.AppClearPrivateMessageUnreadReq) (*messagedto.AppClearPrivateMessageUnreadRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if req.TargetId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	unReadData := messagedao.GetUnReadByUserId(userId)

	unReadDetail := messagedao.GetUnreadDetailByReceiverSender(req.TargetId, userId)

	clearedCount := unReadDetail.UnreadCount
	if req.ClearedCount > 0 {
		clearedCount = req.ClearedCount
	}

	if unReadDetail.UnreadCount > 0 {
		clearedCount = unReadDetail.UnreadCount
		unReadDetail.ClearUnread(clearedCount)
		unReadData.SubPrivateUnread(clearedCount)
		messagedao.FlushUnreadDetailCache(unReadDetail)
	}
	updatePrivateMessageUnreadListCacheUnread(userId, req.TargetId, 0)

	return &messagedto.AppClearPrivateMessageUnreadRes{
		Success:       true,
		ClearedCount:  clearedCount,
		PrivateUnread: unReadData.PrivateUnread,
	}, nil
}

// BatchDeletePrivateMessage App端一键删除与目标用户的全部私信(异步清除当前用户侧数据)
func BatchDeletePrivateMessage(ctx context.Context, req *messagedto.AppBatchDeletePrivateMessageReq) (*messagedto.AppBatchDeletePrivateMessageRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if req.TargetId == userId {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	unReadData := messagedao.GetUnReadByUserId(userId)
	unReadDetail := messagedao.GetUnreadDetailByReceiverSender(req.TargetId, userId)
	if unReadDetail != nil {
		if unReadDetail.UnreadCount > 0 {
			clearedCount := unReadDetail.UnreadCount
			if unReadData != nil {
				unReadData.SubPrivateUnread(clearedCount)
			}
			unReadDetail.ClearUnread(clearedCount)
		}
		if unReadDetail.MutualChat != entity.UserMessageUnreadMutualChatNo {
			unReadDetail.SetMutualChat(entity.UserMessageUnreadMutualChatNo)
		}
		messagedao.FlushUnreadDetailCache(unReadDetail)
	}
	removePrivateMessageUnreadListCacheSender(userId, req.TargetId)

	targetId := req.TargetId
	xrpool.AddWithRecover(ctx, func(poolCtx context.Context) {
		if err := messagedao.DeletePrivateConversationByTarget(poolCtx, userId, targetId); err != nil {
			g.Log().Errorf(poolCtx, "BatchDeletePrivateMessage userId=%d targetId=%d err=%v", userId, targetId, err)
		}
	})

	privateUnread := uint64(0)
	if unReadData != nil {
		privateUnread = unReadData.PrivateUnread
	}

	return &messagedto.AppBatchDeletePrivateMessageRes{
		Success:       true,
		PrivateUnread: privateUnread,
	}, nil
}

// ClearAllPrivateMessageUnread App端清除当前用户全部私信未读(异步直接更新数据库)
func ClearAllPrivateMessageUnread(ctx context.Context, req *messagedto.AppClearAllPrivateMessageUnreadReq) (*messagedto.AppClearAllPrivateMessageUnreadRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}

	unReadData := messagedao.GetUnReadByUserId(userId)
	clearedCount := uint64(0)
	if unReadData != nil {
		clearedCount = unReadData.PrivateUnread
	}

	if unReadData != nil && clearedCount > 0 {
		unReadData.SubPrivateUnread(clearedCount)
	}
	xrpool.AddWithRecover(ctx, func(poolCtx context.Context) {
		if err := messagedao.ClearAllPrivateUnreadInDB(poolCtx, userId); err != nil {
			g.Log().Errorf(poolCtx, "ClearAllPrivateUnreadInDB userId=%d err=%v", userId, err)
		}
	})
	messagedao.ClearAllPrivateUnreadCache(userId)
	clearAllPrivateMessageUnreadListCacheUnread(userId)

	return &messagedto.AppClearAllPrivateMessageUnreadRes{
		Success:       true,
		ClearedCount:  clearedCount,
		PrivateUnread: 0,
	}, nil
}

func buildPrivateMessagePushItem(msg *entity.UserMessage) *messagedto.PrivateMessagePushItem {
	item := &messagedto.PrivateMessagePushItem{
		Id:         msg.ID,
		SenderId:   msg.SenderId,
		ReceiverId: msg.ReceiverId,
		Content:    msg.Content,
		SentAt:     msg.CreatedAt.UnixMilli(),
	}
	if sender := userinfodao.GetUserInfoByUserId(msg.SenderId); sender != nil {
		item.SenderName = sender.Nickname
		item.SenderAvatar = upload.ResolveAvatarUrlForUser(msg.SenderId, sender.Avatar)
	}
	return item
}

func toPrivateMessageUnreadDetailItemFromRow(row *messagedao.PrivateMessageUnreadListRow) *messagedto.AppPrivateMessageUnreadDetailItem {
	item := &messagedto.AppPrivateMessageUnreadDetailItem{
		SenderId:    row.SenderId,
		UnreadCount: row.UnreadCount,
		UpdatedAt:   formatMessageTime(row.UpdatedAt),
	}
	if row.MessageId > 0 {
		msg := &entity.UserMessage{
			SenderId:   row.MessageSenderId,
			ReceiverId: row.MessageReceiverId,
			Content:    row.MessageContent,
		}
		msg.ID = row.MessageId
		msg.CreatedAt = row.MessageCreatedAt
		item.LastMessage = toPrivateMessageItem(row.SessionRowId, msg)
	}
	return item
}

func toPrivateMessageItem(sessionId uint64, msg *entity.UserMessage) *messagedto.AppPrivateMessageItem {
	item := &messagedto.AppPrivateMessageItem{
		Id:          msg.ID,
		SessionId:   sessionId,
		SenderId:    msg.SenderId,
		ReceiverId:  msg.ReceiverId,
		Content:     msg.Content,
		CreatedAt:   formatMessageTime(msg.CreatedAt),
		CreatedAtMs: msg.CreatedAt.UnixMilli(),
	}
	if sender := userinfodao.GetUserInfoByUserId(msg.SenderId); sender != nil {
		item.SenderName = sender.Nickname
		item.SenderAvatar = upload.ResolveAvatarUrlForUser(msg.SenderId, sender.Avatar)
	}
	return item
}

func formatMessageTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006/01/02 15:04:05")
}
