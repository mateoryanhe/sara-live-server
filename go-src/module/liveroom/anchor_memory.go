package liveroom

import (
	"sort"
	"strconv"
	"strings"

	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/guilddao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	liveentity "xr-game-server/entity/live"
	userentity "xr-game-server/entity/user"
	"xr-game-server/module/upload"
)

func queryAnchorListFromMemory(req *accountdto.QueryAnchorListReq) (int, []*accountdto.AnchorListItem) {
	rooms := filterAnchorRooms(getRoomListCache(), req.Key, req.GuildId, req.PlatformOnly, req.GuildOnly, req.LiveStatus)
	total := len(rooms)
	pageRooms := paginateAnchorRooms(rooms, req.PageIndex, req.PageSize)
	ret := make([]*accountdto.AnchorListItem, 0, len(pageRooms))
	for _, room := range pageRooms {
		if item := buildAnchorListItem(room); item != nil {
			ret = append(ret, item)
		}
	}
	fillAnchorListGuildNames(ret)
	return total, ret
}

func fillAnchorListGuildNames(items []*accountdto.AnchorListItem) {
	if len(items) == 0 {
		return
	}
	guildIds := make([]uint64, 0, len(items))
	for _, item := range items {
		if item != nil && item.GuildId > 0 {
			guildIds = append(guildIds, item.GuildId)
		}
	}
	guildNameMap := guilddao.GetNameMapByIds(guildIds)
	if guildNameMap == nil {
		return
	}
	for _, item := range items {
		if item != nil && item.GuildId > 0 {
			item.GuildName = guildNameMap[item.GuildId]
		}
	}
}

func collectRoomGuildIds(rooms []*liveentity.LiveRoom) []uint64 {
	guildIdSet := make(map[uint64]struct{})
	for _, room := range rooms {
		if room != nil && room.GuildId > 0 {
			guildIdSet[room.GuildId] = struct{}{}
		}
	}
	if len(guildIdSet) == 0 {
		return nil
	}
	guildIds := make([]uint64, 0, len(guildIdSet))
	for guildId := range guildIdSet {
		guildIds = append(guildIds, guildId)
	}
	return guildIds
}

func filterAnchorRooms(rooms []*liveentity.LiveRoom, key string, guildId uint64, platformOnly, guildOnly bool, liveStatus *uint8) []*liveentity.LiveRoom {
	key = strings.TrimSpace(key)
	likeKey := strings.ToLower(key)
	guildNameMap := guilddao.GetNameMapByIds(collectRoomGuildIds(rooms))
	filtered := make([]*liveentity.LiveRoom, 0, len(rooms))
	for _, room := range rooms {
		if room == nil || !isRegularAnchorRoom(room) {
			continue
		}
		if platformOnly {
			if room.GuildId != 0 {
				continue
			}
		} else if guildOnly {
			if room.GuildId == 0 {
				continue
			}
		} else if guildId > 0 {
			if room.GuildId != guildId {
				continue
			}
		}
		if liveStatus != nil && roomLiveStatus(room) != *liveStatus {
			continue
		}
		if key != "" && !matchAnchorKey(room.ID, room.GuildId, key, likeKey, guildNameMap) {
			continue
		}
		filtered = append(filtered, room)
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].ID > filtered[j].ID })
	return filtered
}

func filterAnchorRoomsByGuildIds(rooms []*liveentity.LiveRoom, guildIds []uint64, key string, liveStatus *uint8) []*liveentity.LiveRoom {
	if len(guildIds) == 0 {
		return nil
	}
	guildIdSet := make(map[uint64]struct{}, len(guildIds))
	for _, guildId := range guildIds {
		if guildId > 0 {
			guildIdSet[guildId] = struct{}{}
		}
	}
	if len(guildIdSet) == 0 {
		return nil
	}
	key = strings.TrimSpace(key)
	likeKey := strings.ToLower(key)
	guildNameMap := guilddao.GetNameMapByIds(collectRoomGuildIds(rooms))
	filtered := make([]*liveentity.LiveRoom, 0, len(rooms))
	for _, room := range rooms {
		if room == nil || !isRegularAnchorRoom(room) {
			continue
		}
		if _, ok := guildIdSet[room.GuildId]; !ok {
			continue
		}
		if liveStatus != nil && roomLiveStatus(room) != *liveStatus {
			continue
		}
		if key != "" && !matchAnchorKey(room.ID, room.GuildId, key, likeKey, guildNameMap) {
			continue
		}
		filtered = append(filtered, room)
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].ID > filtered[j].ID })
	return filtered
}

func queryAnchorListByGuildIdsFromMemory(guildIds []uint64, req *accountdto.QueryAnchorListReq) (int, []*accountdto.AnchorListItem) {
	rooms := filterAnchorRoomsByGuildIds(getRoomListCache(), guildIds, req.Key, req.LiveStatus)
	total := len(rooms)
	pageRooms := paginateAnchorRooms(rooms, req.PageIndex, req.PageSize)
	ret := make([]*accountdto.AnchorListItem, 0, len(pageRooms))
	for _, room := range pageRooms {
		if item := buildAnchorListItem(room); item != nil {
			ret = append(ret, item)
		}
	}
	fillAnchorListGuildNames(ret)
	return total, ret
}

func roomLiveStatus(room *liveentity.LiveRoom) uint8 {
	if room != nil && room.LiveRecordId > 0 {
		return 1
	}
	return 0
}

func isRegularAnchorRoom(room *liveentity.LiveRoom) bool {
	user := userinfodao.GetUserInfoFromMemory(room.ID)
	return user != nil && (user.UserType == userentity.UserTypeAnchor || user.UserType == userentity.UserTypeSeniorAnchor)
}

func matchAnchorKey(id, guildId uint64, key, likeKey string, guildNameMap map[uint64]string) bool {
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
	if strings.Contains(strings.ToLower(user.ShareCode), likeKey) {
		return true
	}
	if guildId > 0 && guildNameMap != nil {
		return strings.Contains(strings.ToLower(guildNameMap[guildId]), likeKey)
	}
	return false
}

func paginateAnchorRooms(rooms []*liveentity.LiveRoom, pageIndex, pageSize int) []*liveentity.LiveRoom {
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

func buildAnchorListItem(room *liveentity.LiveRoom) *accountdto.AnchorListItem {
	if room == nil {
		return nil
	}
	userId := room.ID
	user := userinfodao.GetUserInfoFromMemory(userId)

	item := &accountdto.AnchorListItem{ID: userId, GuildId: room.GuildId}
	if user != nil {
		item.Nickname = user.Nickname
		item.Phone = user.Phone
		item.Avatar = user.Avatar
		item.UserType = user.UserType
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

func fillAnchorRoomFields(item *accountdto.AnchorListItem, room *liveentity.LiveRoom) {
	if item == nil || room == nil {
		return
	}
	item.RoomTitle = room.Title
	item.RoomId = room.ID
	if cfg := liveroomdao.GetLiveRoomCfgFromCache(room.ID); cfg != nil {
		item.Category = cfg.Category
		item.PrivateInviteType = liveentity.NormalizePrivateInviteType(cfg.PrivateInviteType, cfg.Category)
		item.Ticket = cfg.Ticket
		item.Billing = cfg.Billing
	}
	// 主播列表收益只读永久缓存(生涯累计),不查库
	if income := liveroomdao.GetLiveRoomIncomeTotalFromCache(room.ID); income != nil {
		item.TotalIncome = income.TotalIncome
		item.TotalGiftIncome = income.TotalGiftIncome
		item.TotalPaidDanmakuIncome = income.TotalPaidDanmakuIncome
		item.TotalPrivateRoomTicketIncome = income.TotalPrivateRoomTicketIncome
		item.TotalPrivateRoomWatchIncome = income.TotalPrivateRoomWatchIncome
		item.TotalVideoCallIncome = income.TotalVideoCallIncome
		item.TotalVideoCallTicketIncome = income.TotalVideoCallTicketIncome
		item.TotalVideoCallBillingIncome = income.TotalVideoCallBillingIncome
	}
	item.LiveStatus = roomLiveStatus(room)
	item.BanApplyTime = room.BanApplyTime
	item.BanReason = room.BanReason
	item.Ban = IsRoomBanned(room)
	item.Status = room.Status
}
