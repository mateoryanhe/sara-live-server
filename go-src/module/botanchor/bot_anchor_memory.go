package botanchor

import (
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/dao/botanchordao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/botanchordto"
	"xr-game-server/entity"
	"xr-game-server/module/liveroom"
	"xr-game-server/module/upload"
)

var (
	botAnchorIdMu         sync.RWMutex
	botAnchorIds          []uint64
	botAnchorIdSet        map[uint64]struct{}
	enabledBotAnchorIdSet map[uint64]struct{}
)

func initBotAnchorMemory() {
	reloadBotAnchorMemory()
}

func reloadBotAnchorMemory() {
	ids := botanchordao.LoadAllBotAnchorIds()
	set := make(map[uint64]struct{}, len(ids))
	enabledSet := make(map[uint64]struct{}, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		set[id] = struct{}{}
		user := userinfodao.GetUserInfoByUserId(id)
		if user != nil && user.BotAnchorStatus == entity.BotAnchorStatusEnabled {
			enabledSet[id] = struct{}{}
			preloadBotAnchorCache(id)
			continue
		}
		liveroomdao.RemoveRoomFromCache(id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] > ids[j] })

	botAnchorIdMu.Lock()
	botAnchorIds = ids
	botAnchorIdSet = set
	enabledBotAnchorIdSet = enabledSet
	botAnchorIdMu.Unlock()

	liveroom.RefreshRoomListCache(gctx.New())
}

func preloadBotAnchorCache(userId uint64) {
	userinfodao.GetUserInfoByUserId(userId)
	liveroomdao.GetRoomByAnchor(userId)
}

func addBotAnchorId(userId uint64) {
	if userId == 0 {
		return
	}
	botAnchorIdMu.Lock()
	defer botAnchorIdMu.Unlock()
	if _, ok := botAnchorIdSet[userId]; ok {
		return
	}
	botAnchorIds = append(botAnchorIds, userId)
	botAnchorIdSet[userId] = struct{}{}
	sort.Slice(botAnchorIds, func(i, j int) bool { return botAnchorIds[i] > botAnchorIds[j] })
}

func addEnabledBotAnchorId(userId uint64) {
	if userId == 0 {
		return
	}
	botAnchorIdMu.Lock()
	defer botAnchorIdMu.Unlock()
	if enabledBotAnchorIdSet == nil {
		enabledBotAnchorIdSet = make(map[uint64]struct{})
	}
	enabledBotAnchorIdSet[userId] = struct{}{}
}

func removeEnabledBotAnchorId(userId uint64) {
	if userId == 0 {
		return
	}
	botAnchorIdMu.Lock()
	defer botAnchorIdMu.Unlock()
	delete(enabledBotAnchorIdSet, userId)
}

func enableBotAnchorRoomCache(anchorId, guildId uint64) {
	addEnabledBotAnchorId(anchorId)
	room := liveroom.EnsureAnchorRoom(anchorId, guildId)
	liveroomdao.FlushRoomCache(room)
	preloadBotAnchorCache(anchorId)
}

func disableBotAnchorRoomCache(anchorId uint64) {
	removeEnabledBotAnchorId(anchorId)
	liveroomdao.RemoveRoomFromCache(anchorId)
}

func isBotAnchorId(userId uint64) bool {
	botAnchorIdMu.RLock()
	defer botAnchorIdMu.RUnlock()
	_, ok := botAnchorIdSet[userId]
	return ok
}

func getBotAnchorIds() []uint64 {
	botAnchorIdMu.RLock()
	defer botAnchorIdMu.RUnlock()
	ret := make([]uint64, len(botAnchorIds))
	copy(ret, botAnchorIds)
	return ret
}

func filterBotAnchorIds(ids []uint64, key string) []uint64 {
	key = strings.TrimSpace(key)
	if key == "" {
		return ids
	}
	likeKey := strings.ToLower(key)
	ret := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if matchBotAnchorKey(id, key, likeKey) {
			ret = append(ret, id)
		}
	}
	return ret
}

func matchBotAnchorKey(id uint64, key, likeKey string) bool {
	if strconv.FormatUint(id, 10) == key {
		return true
	}
	user := userinfodao.GetUserInfoByUserId(id)
	if user == nil {
		return false
	}
	nickname := strings.ToLower(user.Nickname)
	return strings.Contains(nickname, likeKey)
}

func paginateBotAnchorIds(ids []uint64, pageIndex, pageSize int) []uint64 {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageIndex - 1) * pageSize
	if offset >= len(ids) {
		return nil
	}
	end := offset + pageSize
	if end > len(ids) {
		end = len(ids)
	}
	return ids[offset:end]
}

func buildBotAnchorListItem(userId uint64) *botanchordto.BotAnchorListItem {
	user := userinfodao.GetUserInfoByUserId(userId)
	room := liveroomdao.GetRoomByAnchor(userId)
	if room == nil {
		room = liveroom.EnsureAnchorRoom(userId, 0)
	}

	item := &botanchordto.BotAnchorListItem{
		ID:     userId,
		RoomId: userId,
	}
	if user != nil {
		item.Nickname = user.Nickname
		item.Avatar = user.Avatar
		item.GuildId = user.GuildId
		item.BotAnchorStatus = user.BotAnchorStatus
		if !user.CreatedAt.IsZero() {
			createdAt := user.CreatedAt
			item.CreatedAt = &createdAt
		}
		if !user.UpdatedAt.IsZero() {
			updatedAt := user.UpdatedAt
			item.UpdatedAt = &updatedAt
		}
	}
	if room != nil {
		item.RoomTitle = room.Title
		item.Category = room.Category
		item.TagId = room.TagId
		item.CloudPlayerVideoFile = room.CloudPlayerVideo
		item.PushStream = room.PushStream
		if room.LiveRecordId > 0 {
			item.LiveStatus = 1
		}
		if tag := liveroom.GetRoomTagFromMemoryById(room.TagId); tag != nil {
			item.TagName = tag.Name
		}
	}
	item.Avatar = upload.ResolveAvatarUrlForUser(userId, item.Avatar)
	item.CloudPlayerVideo = upload.ResolveCloudPlayerVideoUrl(item.CloudPlayerVideoFile)
	return item
}

func queryBotAnchorListFromMemory(req *botanchordto.QueryBotAnchorListReq) (int, []*botanchordto.BotAnchorListItem) {
	ids := filterBotAnchorIds(getBotAnchorIds(), req.Key)
	total := len(ids)
	pageIds := paginateBotAnchorIds(ids, req.PageIndex, req.PageSize)
	ret := make([]*botanchordto.BotAnchorListItem, 0, len(pageIds))
	for _, id := range pageIds {
		if item := buildBotAnchorListItem(id); item != nil {
			ret = append(ret, item)
		}
	}
	return total, ret
}
