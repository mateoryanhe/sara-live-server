package message

import (
	"context"
	"sort"
	"time"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/messagedao"
	"xr-game-server/dto/messagedto"
	"xr-game-server/entity/message"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

// ListUserActivityMessage App端查询活动消息列表
func ListUserActivityMessage(ctx context.Context, _ *messagedto.AppActivityMessageListReq) (*messagedto.AppActivityMessageListRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}

	syncNewActivityMessages(userId)

	detailMap := make(map[uint64]*entity.ActivityMessage)
	for _, msg := range messagedao.GetAllActivityMessagesCached() {
		if msg != nil && msg.ID != 0 && msg.Status == entity.ActivityMessageStatusPublished {
			detailMap[msg.ID] = msg
		}
	}

	userRows := messagedao.GetUserActivityMessageListCacheA(userId)
	list := make([]*messagedto.AppActivityMessageItem, 0, len(userRows))
	for _, row := range userRows {
		if row == nil {
			continue
		}
		detail := detailMap[row.ActivityMessageId]
		if detail == nil {
			continue
		}
		list = append(list, toAppActivityMessageItem(detail))
	}
	return &messagedto.AppActivityMessageListRes{List: list}, nil
}

func toAppActivityMessageItem(msg *entity.ActivityMessage) *messagedto.AppActivityMessageItem {
	if msg == nil {
		return nil
	}
	publishedAt := ""
	if msg.PublishedAt != nil {
		publishedAt = formatMessageTime(*msg.PublishedAt)
	}
	return &messagedto.AppActivityMessageItem{
		Id:          msg.ID,
		IconEn:      upload.GetUrlByName(msg.IconEn),
		IconEs:      upload.GetUrlByName(msg.IconEs),
		IconPt:      upload.GetUrlByName(msg.IconPt),
		IconHi:      upload.GetUrlByName(msg.IconHi),
		IconId:      upload.GetUrlByName(msg.IconId),
		BgEn:        upload.GetUrlByName(msg.BgEn),
		BgEs:        upload.GetUrlByName(msg.BgEs),
		BgPt:        upload.GetUrlByName(msg.BgPt),
		BgHi:        upload.GetUrlByName(msg.BgHi),
		BgId:        upload.GetUrlByName(msg.BgId),
		TitleEn:     msg.TitleEn,
		TitleEs:     msg.TitleEs,
		TitlePt:     msg.TitlePt,
		TitleHi:     msg.TitleHi,
		TitleId:     msg.TitleId,
		ContentEn:   msg.ContentEn,
		ContentEs:   msg.ContentEs,
		ContentPt:   msg.ContentPt,
		ContentHi:   msg.ContentHi,
		ContentId:   msg.ContentId,
		UrlEn:       msg.UrlEn,
		UrlEs:       msg.UrlEs,
		UrlPt:       msg.UrlPt,
		UrlHi:       msg.UrlHi,
		UrlId:       msg.UrlId,
		PublishedAt: publishedAt,
	}
}

func syncNewActivityMessages(userId uint64) {
	published := GetPublishedActivityMessages()
	if len(published) == 0 {
		return
	}

	userList := messagedao.GetUserActivityMessageListCacheA(userId)
	owned := make(map[uint64]struct{}, len(userList))
	for _, row := range userList {
		if row != nil {
			owned[row.ActivityMessageId] = struct{}{}
		}
	}

	newRows := make([]*entity.UserActivityMessage, 0)
	for _, msg := range published {
		if msg == nil || msg.ID == 0 {
			continue
		}
		if _, ok := owned[msg.ID]; ok {
			continue
		}
		newRows = append(newRows, entity.NewUserActivityMessage(userId, msg.ID, msg.PublishedAt))
	}
	if len(newRows) == 0 {
		return
	}

	messagedao.FlushUserActivityMessageListCacheA(userId, mergeUserActivityMessageList(newRows, userList))

	delta := uint64(len(newRows))
	AddSystemMessageUnread(userId, entity.UserSystemMessageUnreadTypeActivity, delta)
	unread := messagedao.GetUnReadByUserId(userId)
	unread.AddSystemUnread(delta)
	messagedao.PublishMessageUnread(unread)
}

func mergeUserActivityMessageList(newRows, existing []*entity.UserActivityMessage) []*entity.UserActivityMessage {
	seen := make(map[uint64]struct{}, len(newRows)+len(existing))
	merged := make([]*entity.UserActivityMessage, 0, len(newRows)+len(existing))
	appendUnique := func(row *entity.UserActivityMessage) {
		if row == nil || row.ActivityMessageId == 0 {
			return
		}
		if _, ok := seen[row.ActivityMessageId]; ok {
			return
		}
		seen[row.ActivityMessageId] = struct{}{}
		merged = append(merged, row)
	}
	for _, row := range newRows {
		appendUnique(row)
	}
	for _, row := range existing {
		appendUnique(row)
	}
	sort.Slice(merged, func(i, j int) bool {
		return userActivityMessagePublishedAt(merged[i]).After(userActivityMessagePublishedAt(merged[j]))
	})
	if len(merged) > messagedao.UserActivityMessageListCacheASize {
		merged = merged[:messagedao.UserActivityMessageListCacheASize]
	}
	return merged
}

func prependPublishedActivityMessageToCachedUsers(msg *entity.ActivityMessage) {
	if msg == nil || msg.ID == 0 || msg.PublishedAt == nil {
		return
	}
	for _, userId := range messagedao.GetCachedUserActivityMessageUserIds() {
		newRow := entity.NewUserActivityMessage(userId, msg.ID, msg.PublishedAt)
		if !messagedao.PrependUserActivityMessageToListCacheA(userId, newRow) {
			continue
		}
		AddSystemMessageUnread(userId, entity.UserSystemMessageUnreadTypeActivity, 1)
		unread := messagedao.GetUnReadByUserId(userId)
		unread.AddSystemUnread(1)
		messagedao.PublishMessageUnread(unread)
	}
}

func removeUnpublishedActivityMessageFromCachedUsers(activityMessageId uint64) {
	if activityMessageId == 0 {
		return
	}
	for _, userId := range messagedao.GetCachedUserActivityMessageUserIds() {
		messagedao.RemoveUserActivityMessageFromListCacheA(userId, activityMessageId)
	}
}

func userActivityMessagePublishedAt(row *entity.UserActivityMessage) time.Time {
	if row == nil || row.PublishedAt == nil {
		return time.Time{}
	}
	return *row.PublishedAt
}
