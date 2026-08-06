package message

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/messagedao"
	"xr-game-server/dto/messagedto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

func systemMessageUnreadRow(userId uint64, msgType uint8) *entity.UserSystemMessageUnread {
	for _, row := range messagedao.GetSystemMessageUnreadListCache(userId) {
		if row != nil && row.Type == msgType {
			return row
		}
	}
	return entity.NewUserSystemMessageUnread(userId, msgType)
}

func saveSystemMessageUnreadList(userId uint64, row *entity.UserSystemMessageUnread) {
	list := messagedao.GetSystemMessageUnreadListCache(userId)
	newList := make([]*entity.UserSystemMessageUnread, 0, len(list)+1)
	if row.UnreadCount > 0 {
		newList = append(newList, row)
	}
	for _, item := range list {
		if item != nil && item.Type != row.Type && item.UnreadCount > 0 {
			newList = append(newList, item)
		}
	}
	if len(newList) > messagedao.SystemMessageUnreadListCacheMax {
		newList = newList[:messagedao.SystemMessageUnreadListCacheMax]
	}
	messagedao.FlushSystemMessageUnreadListCache(userId, newList)
}

func markSystemMessageRead(userId uint64, msgType uint8, readCount uint64) uint64 {
	row := systemMessageUnreadRow(userId, msgType)
	if row.UnreadCount == 0 {
		return 0
	}
	cleared := readCount
	if readCount == 0 || readCount > row.UnreadCount {
		cleared = row.UnreadCount
	}
	row.SubUnread(cleared)
	saveSystemMessageUnreadList(userId, row)
	return cleared
}

func clearAllSystemMessageUnread(userId uint64) {
	if userId == 0 {
		return
	}
	unReadData := messagedao.GetUnReadByUserId(userId)
	list := messagedao.GetSystemMessageUnreadListCache(userId)
	for _, row := range list {
		if row == nil || row.UnreadCount == 0 {
			continue
		}
		row.SubUnread(row.UnreadCount)
	}
	messagedao.FlushSystemMessageUnreadListCache(userId, make([]*entity.UserSystemMessageUnread, 0))
	if unReadData != nil && unReadData.SystemUnread > 0 {
		unReadData.SubSystemUnread(unReadData.SystemUnread)
	}
}

// ListSystemMessageUnread App端查询系统消息未读列表
func ListSystemMessageUnread(ctx context.Context, _ *messagedto.AppSystemMessageUnreadListReq) (*messagedto.AppSystemMessageUnreadListRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}

	rows := messagedao.GetSystemMessageUnreadListCache(userId)
	list := make([]*messagedto.AppSystemMessageUnreadListItem, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.UnreadCount == 0 {
			continue
		}
		list = append(list, &messagedto.AppSystemMessageUnreadListItem{
			MsgType:     row.Type,
			UnreadCount: row.UnreadCount,
			UpdatedAt:   formatMessageTime(row.UpdatedAt),
			UpdatedAtMs: row.UpdatedAt.UnixMilli(),
		})
	}
	return &messagedto.AppSystemMessageUnreadListRes{List: list}, nil
}

// ClearSystemMessageUnread App端上报系统消息已读(按类型;已读数=0时该类型未读清零)
func ClearSystemMessageUnread(ctx context.Context, req *messagedto.AppClearSystemMessageUnreadReq) (*messagedto.AppClearSystemMessageUnreadRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}

	cleared := markSystemMessageRead(userId, req.MsgType, req.ReadCount)
	unReadData := messagedao.GetUnReadByUserId(userId)
	if cleared > 0 {
		unReadData.SubSystemUnread(cleared)
	}

	row := systemMessageUnreadRow(userId, req.MsgType)
	return &messagedto.AppClearSystemMessageUnreadRes{
		Success:      true,
		UnreadCount:  row.UnreadCount,
		SystemUnread: unReadData.SystemUnread,
	}, nil
}

// AddSystemMessageUnread 增加指定类型系统消息未读数
func AddSystemMessageUnread(userId uint64, msgType uint8, delta uint64) {
	if userId == 0 || msgType == 0 || delta == 0 {
		return
	}
	row := systemMessageUnreadRow(userId, msgType)
	row.AddUnread(delta)
	saveSystemMessageUnreadList(userId, row)
}
