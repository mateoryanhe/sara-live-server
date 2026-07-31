package message

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/messagedao"
	"xr-game-server/dto/messagedto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

func addPersonalSystemMessage(
	userId uint64,
	icon, params string,
	titleEn, titleEs, titlePt, titleHi string,
	contentEn, contentEs, contentPt, contentHi string,
) {
	if userId == 0 {
		return
	}
	row := entity.NewUserPersonalSystemMessage(
		userId, icon, params,
		titleEn, titleEs, titlePt, titleHi,
		contentEn, contentEs, contentPt, contentHi,
	)
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
	return &messagedto.AppPersonalSystemMessageItem{
		Id:        row.ID,
		Icon:      row.Icon,
		TitleEn:   row.TitleEn,
		TitleEs:   row.TitleEs,
		TitlePt:   row.TitlePt,
		TitleHi:   row.TitleHi,
		ContentEn: row.ContentEn,
		ContentEs: row.ContentEs,
		ContentPt: row.ContentPt,
		ContentHi: row.ContentHi,
		Params:    row.Params,
		CreatedAt: formatMessageTime(row.CreatedAt),
	}
}
