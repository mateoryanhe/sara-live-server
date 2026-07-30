package liveroom

import (
	"sort"
	"strconv"
	"strings"

	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/entity"
	"xr-game-server/module/upload"
)

func queryAnchorListFromMemory(req *accountdto.QueryAnchorListReq) (int, []*accountdto.AnchorListItem) {
	rooms := filterAnchorRooms(getRoomListCache(), req.Key)
	total := len(rooms)
	pageRooms := paginateAnchorRooms(rooms, req.PageIndex, req.PageSize)
	ret := make([]*accountdto.AnchorListItem, 0, len(pageRooms))
	for _, room := range pageRooms {
		if item := buildAnchorListItem(room); item != nil {
			ret = append(ret, item)
		}
	}
	return total, ret
}

func filterAnchorRooms(rooms []*entity.LiveRoom, key string) []*entity.LiveRoom {
	key = strings.TrimSpace(key)
	likeKey := strings.ToLower(key)
	filtered := make([]*entity.LiveRoom, 0, len(rooms))
	for _, room := range rooms {
		if room == nil || !isRegularAnchorRoom(room) {
			continue
		}
		if key != "" && !matchAnchorKey(room.ID, key, likeKey) {
			continue
		}
		filtered = append(filtered, room)
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].ID > filtered[j].ID })
	return filtered
}

func isRegularAnchorRoom(room *entity.LiveRoom) bool {
	user := userinfodao.GetUserInfoFromMemory(room.ID)
	return user != nil && (user.UserType == entity.UserTypeAnchor || user.UserType == entity.UserTypeSeniorAnchor)
}

func matchAnchorKey(id uint64, key, likeKey string) bool {
	if strings.Contains(strconv.FormatUint(id, 10), key) {
		return true
	}
	user := userinfodao.GetUserInfoFromMemory(id)
	if user == nil {
		return false
	}
	if strings.Contains(strings.ToLower(user.Nickname), likeKey) {
		return true
	}
	if strings.Contains(user.Phone, key) {
		return true
	}
	return strings.Contains(strings.ToLower(user.ShareCode), likeKey)
}

func paginateAnchorRooms(rooms []*entity.LiveRoom, pageIndex, pageSize int) []*entity.LiveRoom {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageIndex - 1) * pageSize
	if offset >= len(rooms) {
		return nil
	}
	end := offset + pageSize
	if end > len(rooms) {
		end = len(rooms)
	}
	return rooms[offset:end]
}

func buildAnchorListItem(room *entity.LiveRoom) *accountdto.AnchorListItem {
	if room == nil {
		return nil
	}
	userId := room.ID
	user := userinfodao.GetUserInfoFromMemory(userId)

	item := &accountdto.AnchorListItem{ID: userId}
	if user != nil {
		item.Nickname = user.Nickname
		item.Phone = user.Phone
		item.Avatar = user.Avatar
		item.GuildId = user.GuildId
		if !user.CreatedAt.IsZero() {
			createdAt := user.CreatedAt
			item.CreatedAt = &createdAt
		}
	}
	if accountCache, _, _ := accountdao.FindAccountInCacheByID(userId); accountCache != nil {
		item.IP = accountCache.IP
		if !accountCache.CreatedAt.IsZero() {
			registeredAt := accountCache.CreatedAt
			item.RegisteredAt = &registeredAt
		}
	}
	fillAnchorRoomFields(item, room)
	item.Avatar = upload.ResolveAvatarUrlForUser(userId, item.Avatar)
	return item
}

func fillAnchorRoomFields(item *accountdto.AnchorListItem, room *entity.LiveRoom) {
	if item == nil || room == nil {
		return
	}
	item.RoomTitle = room.Title
	item.RoomId = room.ID
	item.Category = room.Category
	item.PrivateInviteType = room.PrivateInviteType
	item.Ticket = room.Ticket
	item.Billing = room.Billing
	item.TotalIncome = room.TotalIncome
	item.TotalGiftIncome = room.TotalGiftIncome
	item.TotalPaidDanmakuIncome = room.TotalPaidDanmakuIncome
	item.TotalPrivateRoomTicketIncome = room.TotalPrivateRoomTicketIncome
	item.TotalPrivateRoomWatchIncome = room.TotalPrivateRoomWatchIncome
	item.TotalVideoCallIncome = room.TotalVideoCallIncome
	item.TotalVideoCallTicketIncome = room.TotalVideoCallTicketIncome
	item.TotalVideoCallBillingIncome = room.TotalVideoCallBillingIncome
	if room.LiveRecordId > 0 {
		item.LiveStatus = 1
	} else {
		item.LiveStatus = 0
	}
	item.BanApplyTime = room.BanApplyTime
	item.BanReason = room.BanReason
	item.Ban = IsRoomBanned(room)
}
