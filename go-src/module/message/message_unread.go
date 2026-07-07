package message

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/messagedao"
	"xr-game-server/dto/messagedto"
	"xr-game-server/errercode"
)

// GetMessageUnreadCount App端查询系统消息与私信未读数
func GetMessageUnreadCount(ctx context.Context, _ *messagedto.AppMessageUnreadCountReq) (*messagedto.AppMessageUnreadCountRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}

	unReadData := messagedao.GetUnReadByUserId(userId)
	return &messagedto.AppMessageUnreadCountRes{
		SystemUnread:  unReadData.SystemUnread,
		PrivateUnread: unReadData.PrivateUnread,
	}, nil
}
