package message

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/messagedao"
	"xr-game-server/dto/messagedto"
	"xr-game-server/entity/message"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

func addPersonalSystemMessage(userId uint64, messageTypeId uint32, params string) {
	if userId == 0 || messageTypeId == 0 {
		return
	}
	row := entity.NewUserPersonalSystemMessage(userId, messageTypeId, params)
	list := messagedao.GetUserPersonalSystemMessageListCache(userId)
	newList := make([]*entity.UserPersonalSystemMessage, 0, len(list)+1)
	newList = append(newList, row)
	for _, item := range list {
		if item != nil {
			newList = append(newList, item)
		}
	}
	if len(newList) > messagedao.UserPersonalSystemMessageListCacheMax {
		newList = newList[:messagedao.UserPersonalSystemMessageListCacheMax]
	}
	messagedao.FlushUserPersonalSystemMessageListCache(userId, newList)

	AddSystemMessageUnread(userId, entity.UserSystemMessageUnreadTypePersonal, 1)
	messagedao.GetUnReadByUserId(userId).AddSystemUnread(1)
}

// ListPersonalSystemMessage App端查询个人系统消息列表(只读缓存列表,不分页)
func ListPersonalSystemMessage(ctx context.Context, _ *messagedto.AppPersonalSystemMessageListReq) (*messagedto.AppPersonalSystemMessageListRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}

	rows := messagedao.GetUserPersonalSystemMessageListCache(userId)
	list := make([]*messagedto.AppPersonalSystemMessageItem, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, toAppPersonalSystemMessageItem(row))
	}
	return &messagedto.AppPersonalSystemMessageListRes{List: list}, nil
}

func toAppPersonalSystemMessageItem(row *entity.UserPersonalSystemMessage) *messagedto.AppPersonalSystemMessageItem {
	if row == nil {
		return nil
	}
	item := &messagedto.AppPersonalSystemMessageItem{
		Id:            row.ID,
		MessageTypeId: row.MessageTypeId,
		Params:        row.Params,
		CreatedAt:     formatMessageTime(row.CreatedAt),
	}
	if display := resolvePersonalSystemMessageDisplay(row.MessageTypeId); display != nil {
		item.IconEn = upload.GetUrlByName(display.IconEn)
		item.IconEs = upload.GetUrlByName(display.IconEs)
		item.IconPt = upload.GetUrlByName(display.IconPt)
		item.IconHi = upload.GetUrlByName(display.IconHi)
		item.IconId = upload.GetUrlByName(display.IconId)
		item.TitleEn = display.TitleEn
		item.TitleEs = display.TitleEs
		item.TitlePt = display.TitlePt
		item.TitleHi = display.TitleHi
		item.TitleId = display.TitleId
		item.ContentEn = display.ContentEn
		item.ContentEs = display.ContentEs
		item.ContentPt = display.ContentPt
		item.ContentHi = display.ContentHi
		item.ContentId = display.ContentId
	}
	return item
}

func resolvePersonalSystemMessageDisplay(messageTypeId uint32) *personalSystemMessageDisplay {
	switch messageTypeId {
	case entity.PersonalSystemMessageTypeWelcome:
		display := welcomePersonalMessageDisplay()
		return &display
	default:
		return nil
	}
}
