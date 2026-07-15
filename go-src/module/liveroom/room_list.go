package liveroom

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"xr-game-server/constants/userstatus"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/xrtimer"
	"xr-game-server/dao/livefollowdao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/agoradto"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/entity"
	"xr-game-server/module/agora"
	"xr-game-server/module/upload"

	"github.com/gogf/gf/v2/os/gctx"
)

const (
	roomListDefaultPageSize = 20
	roomListMaxPageSize     = 100
	roomListRefreshInterval = 10 * time.Minute

	nearbyLiveRoomDefaultCount = 1
	nearbyLiveRoomMaxCount     = 20
)

var (
	roomListCache      atomic.Value // []*entity.LiveRoom
	emptyRoomListCache = make([]*entity.LiveRoom, 0)
)

func initRoomList() {
	roomListCache.Store(emptyRoomListCache)
	ctx := gctx.New()
	flushRoomList(ctx)
	xrtimer.AddSingleton(ctx, roomListRefreshInterval, flushRoomList)
}

func flushRoomList(ctx context.Context) {
	_ = ctx
	allData := liveroomdao.GetAllLiveRoom()
	filtered := make([]*entity.LiveRoom, 0, len(allData))
	for _, room := range allData {
		if room == nil || IsRoomBanned(room) || isDisabledBotAnchorRoom(room) {
			continue
		}
		filtered = append(filtered, room)
	}

	onlineCounts := make(map[uint64]int, len(filtered))
	for _, room := range filtered {
		onlineCounts[room.ID] = countAudienceInRoom(room.ID)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return compareLiveRoomsForList(filtered[i], filtered[j], onlineCounts)
	})

	roomListCache.Store(filtered)
}

// RefreshRoomListCache 主动刷新直播间列表缓存
func RefreshRoomListCache(ctx context.Context) {
	flushRoomList(ctx)
}

func isDisabledBotAnchorRoom(room *entity.LiveRoom) bool {
	if room == nil {
		return false
	}
	user := userinfodao.GetUserInfoByUserId(room.ID)
	return user != nil && user.IsBotAnchor() && user.BotAnchorStatus != entity.BotAnchorStatusEnabled
}

func resolveBotAnchorRoomInfo(anchorId uint64, cloudPlayerVideo string) (isBotAnchor bool, cloudPlayerVideoUrl string) {
	if anchorId > 0 {
		if user := userinfodao.GetUserInfoByUserId(anchorId); user != nil {
			isBotAnchor = user.IsBotAnchor()
		}
	}
	cloudPlayerVideoUrl = upload.ResolveCloudPlayerVideoUrl(cloudPlayerVideo)
	return
}

func compareLiveRoomsForList(a, b *entity.LiveRoom, onlineCounts map[uint64]int) bool {
	if a == nil {
		return false
	}
	if b == nil {
		return true
	}
	liveA := a.LiveRecordId > 0
	liveB := b.LiveRecordId > 0
	if liveA != liveB {
		return liveA
	}
	onlineA := onlineCounts[a.ID]
	onlineB := onlineCounts[b.ID]
	if onlineA != onlineB {
		return onlineA > onlineB
	}
	if liveRoomHeartTimeUnix(a) != liveRoomHeartTimeUnix(b) {
		return liveRoomHeartTimeUnix(a) > liveRoomHeartTimeUnix(b)
	}
	return a.ID > b.ID
}

func liveRoomHeartTimeUnix(room *entity.LiveRoom) int64 {
	if room == nil || room.HeartTime == nil {
		return 0
	}
	return room.HeartTime.Unix()
}

func toLiveRoomListItem(room *entity.LiveRoom, userId uint64) *liveroomdto.LiveRoomListItem {
	status := userstatus.LiveRoomStatusClosed
	if room.LiveRecordId > 0 {
		status = userstatus.LiveRoomStatusLive
	}
	item := &liveroomdto.LiveRoomListItem{
		RoomId:        strconv.FormatUint(room.ID, 10),
		GuildId:       strconv.FormatUint(room.GuildId, 10),
		Title:         room.Title,
		Cover:         room.Cover,
		Notice:        room.Notice,
		Status:        status,
		Category:      room.Category,
		Ticket:        room.Ticket,
		Billing:       room.Billing,
		AllowCallIcon: allowShowCallIcon(room, userId),
		CreateAt:      room.CreatedAt.Unix(),
	}
	if room.TagId > 0 {
		item.TagId = strconv.FormatUint(room.TagId, 10)
		item.TagName = getRoomTagName(room.TagId)
	}
	if item.Cover != "" {
		item.Cover = upload.GetUrlByName(room.Cover)
	}

	if u := userinfodao.GetUserInfoByUserId(room.ID); u != nil {
		item.AnchorNickname = u.Nickname
		item.AnchorAvatar = upload.ResolveAvatarUrlForUser(room.ID, u.Avatar)
	}
	item.IsBotAnchor, item.CloudPlayerVideo = resolveBotAnchorRoomInfo(room.ID, room.CloudPlayerVideo)
	item.OnlineCount = countAudienceInRoom(room.ID)

	if userId > 0 {
		channelName := strconv.FormatUint(room.ID, 10)
		role := agoradto.RTCRoleSubscriber
		if userId == room.ID {
			role = agoradto.RTCRolePublisher
		}
		if token, expireAt, err := agora.ResolveChannelToken(userId, channelName, role); err == nil {
			item.AgoraToken = token
			item.AgoraTokenExpireAt = expireAt
		}
	}
	return item
}

func buildLiveRoomListItems(rooms []*entity.LiveRoom, userId uint64) []*liveroomdto.LiveRoomListItem {
	list := make([]*liveroomdto.LiveRoomListItem, 0, len(rooms))
	for _, room := range rooms {
		if room == nil {
			continue
		}
		list = append(list, toLiveRoomListItem(room, userId))
	}
	return list
}

func getRoomListCache() []*entity.LiveRoom {
	v := roomListCache.Load()
	if v == nil {
		return nil
	}
	list, ok := v.([]*entity.LiveRoom)
	if !ok || len(list) == 0 {
		return nil
	}
	return list
}

func filterRoomsByStatus(rooms []*entity.LiveRoom, statusFilter int) []*entity.LiveRoom {
	if statusFilter == 0 {
		return rooms
	}
	filtered := make([]*entity.LiveRoom, 0, len(rooms))
	for _, room := range rooms {
		if room == nil {
			continue
		}
		status := userstatus.LiveRoomStatusClosed
		if room.LiveRecordId > 0 {
			status = userstatus.LiveRoomStatusLive
		}
		switch statusFilter {
		case 1:
			if status == userstatus.LiveRoomStatusLive {
				filtered = append(filtered, room)
			}
		case 2:
			if status == userstatus.LiveRoomStatusClosed {
				filtered = append(filtered, room)
			}
		default:
			filtered = append(filtered, room)
		}
	}
	return filtered
}

func filterRoomsByQuery(rooms []*entity.LiveRoom, tagId uint64, title, notice string) []*entity.LiveRoom {
	titleKey := strings.ToLower(strings.TrimSpace(title))
	noticeKey := strings.ToLower(strings.TrimSpace(notice))
	if tagId == 0 && titleKey == "" && noticeKey == "" {
		return rooms
	}
	filtered := make([]*entity.LiveRoom, 0, len(rooms))
	for _, room := range rooms {
		if room == nil {
			continue
		}
		if tagId > 0 && room.TagId != tagId {
			continue
		}
		if titleKey != "" && !strings.Contains(strings.ToLower(room.Title), titleKey) {
			continue
		}
		if noticeKey != "" && !strings.Contains(strings.ToLower(room.Notice), noticeKey) {
			continue
		}
		filtered = append(filtered, room)
	}
	return filtered
}

func normalizeRoomListPage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = roomListDefaultPageSize
	}
	if pageSize > roomListMaxPageSize {
		pageSize = roomListMaxPageSize
	}
	return page, pageSize
}

func roomListPageRange(total, page, pageSize int) (int, int) {
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	return start, end
}

// GetRoomList App 分页查询直播间列表(走内存缓存)
func GetRoomList(ctx context.Context, req *liveroomdto.GetLiveRoomListReq) (*liveroomdto.GetLiveRoomListRes, error) {
	userId := httpserver.GetAuthId(ctx)
	page, pageSize := normalizeRoomListPage(req.Page, req.PageSize)

	cached := getRoomListCache()
	if cached == nil {
		return &liveroomdto.GetLiveRoomListRes{
			Total:    0,
			Page:     page,
			PageSize: pageSize,
			List:     make([]*liveroomdto.LiveRoomListItem, 0),
		}, nil
	}

	filtered := filterRoomsByStatus(cached, req.StatusFilter)
	filtered = filterRoomsByQuery(filtered, req.TagId, req.Title, req.Notice)
	total := len(filtered)
	start, end := roomListPageRange(total, page, pageSize)

	return &liveroomdto.GetLiveRoomListRes{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     buildLiveRoomListItems(filtered[start:end], userId),
	}, nil
}

// GetFollowedRoomList App 分页查询当前用户关注的直播间(按关注时间倒序)
func GetFollowedRoomList(ctx context.Context, req *liveroomdto.GetFollowedLiveRoomListReq) (*liveroomdto.GetLiveRoomListRes, error) {
	userId := httpserver.GetAuthId(ctx)
	page, pageSize := normalizeRoomListPage(req.Page, req.PageSize)

	followings := livefollowdao.GetFollowingsByUser(userId)
	rooms := make([]*entity.LiveRoom, 0, len(followings))
	for _, f := range followings {
		if f == nil {
			continue
		}
		room := liveroomdao.GetRoomById(f.AnchorId)
		if room == nil || IsRoomBanned(room) {
			continue
		}
		rooms = append(rooms, room)
	}

	filtered := filterRoomsByStatus(rooms, req.StatusFilter)
	total := len(filtered)
	start, end := roomListPageRange(total, page, pageSize)

	list := buildLiveRoomListItems(filtered[start:end], userId)
	if len(list) == 0 {
		list = make([]*liveroomdto.LiveRoomListItem, 0)
	}

	return &liveroomdto.GetLiveRoomListRes{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     list,
	}, nil
}

func normalizeNearbyLiveRoomCount(count int) int {
	if count <= 0 {
		return nearbyLiveRoomDefaultCount
	}
	if count > nearbyLiveRoomMaxCount {
		return nearbyLiveRoomMaxCount
	}
	return count
}

func findLiveRoomIndex(rooms []*entity.LiveRoom, roomId uint64) int {
	for i, room := range rooms {
		if room != nil && room.ID == roomId {
			return i
		}
	}
	return -1
}

func collectNearbyLiveRooms(rooms []*entity.LiveRoom, currentIdx, direction, count int) []*entity.LiveRoom {
	if currentIdx < 0 || len(rooms) == 0 || count <= 0 {
		return make([]*entity.LiveRoom, 0)
	}
	result := make([]*entity.LiveRoom, 0, count)
	switch direction {
	case liveroomdto.NearbyLiveRoomDirectionDown:
		for i := currentIdx + 1; i < len(rooms) && len(result) < count; i++ {
			result = append(result, rooms[i])
		}
	case liveroomdto.NearbyLiveRoomDirectionUp:
		for i := currentIdx - 1; i >= 0 && len(result) < count; i-- {
			result = append(result, rooms[i])
		}
	}
	return result
}

func filterHotLiveRooms(rooms []*entity.LiveRoom) []*entity.LiveRoom {
	hotTagId := getRoomTagIdByName("Hot")
	filtered := make([]*entity.LiveRoom, 0, len(rooms))
	for _, room := range rooms {
		if room == nil || room.LiveRecordId == 0 {
			continue
		}
		if hotTagId > 0 {
			if room.TagId != hotTagId {
				continue
			}
		} else if room.Category != entity.LiveRoomCategoryHot {
			continue
		}
		filtered = append(filtered, room)
	}
	return filtered
}

func buildHotLiveRoomListItems(rooms []*entity.LiveRoom, userId uint64, startRank int) []*liveroomdto.HotLiveRoomListItem {
	list := make([]*liveroomdto.HotLiveRoomListItem, 0, len(rooms))
	rank := startRank
	for _, room := range rooms {
		if room == nil {
			continue
		}
		list = append(list, &liveroomdto.HotLiveRoomListItem{
			LiveRoomListItem: *toLiveRoomListItem(room, userId),
			Rank:             rank,
		})
		rank++
	}
	return list
}

// GetHotLiveRoomList App 分页查询 Hot 分类直播中房间列表(走内存缓存排序,含排名)
func GetHotLiveRoomList(ctx context.Context, req *liveroomdto.GetHotLiveRoomListReq) (*liveroomdto.GetHotLiveRoomListRes, error) {
	userId := httpserver.GetAuthId(ctx)
	pageIndex, pageSize := normalizeRoomListPage(req.PageIndex, req.PageSize)

	cached := getRoomListCache()
	if cached == nil {
		return &liveroomdto.GetHotLiveRoomListRes{
			Total:     0,
			PageIndex: pageIndex,
			PageSize:  pageSize,
			List:      make([]*liveroomdto.HotLiveRoomListItem, 0),
		}, nil
	}

	filtered := filterHotLiveRooms(cached)
	total := len(filtered)
	start, end := roomListPageRange(total, pageIndex, pageSize)
	list := buildHotLiveRoomListItems(filtered[start:end], userId, start+1)
	if len(list) == 0 {
		list = make([]*liveroomdto.HotLiveRoomListItem, 0)
	}

	return &liveroomdto.GetHotLiveRoomListRes{
		Total:     total,
		PageIndex: pageIndex,
		PageSize:  pageSize,
		List:      list,
	}, nil
}

// GetNearbyLiveRoomList App 以当前直播间为锚点,获取列表中相邻的直播中直播间(走内存缓存排序)
func GetNearbyLiveRoomList(ctx context.Context, req *liveroomdto.GetNearbyLiveRoomListReq) (*liveroomdto.GetNearbyLiveRoomListRes, error) {
	userId := httpserver.GetAuthId(ctx)
	count := normalizeNearbyLiveRoomCount(req.Count)

	cached := getRoomListCache()
	if cached == nil {
		return &liveroomdto.GetNearbyLiveRoomListRes{
			List: make([]*liveroomdto.LiveRoomListItem, 0),
		}, nil
	}

	liveRooms := filterRoomsByStatus(cached, int(userstatus.LiveRoomStatusLive))
	currentIdx := findLiveRoomIndex(liveRooms, req.RoomId)
	nearbyRooms := collectNearbyLiveRooms(liveRooms, currentIdx, req.Direction, count)

	list := buildLiveRoomListItems(nearbyRooms, userId)
	if len(list) == 0 {
		list = make([]*liveroomdto.LiveRoomListItem, 0)
	}
	return &liveroomdto.GetNearbyLiveRoomListRes{List: list}, nil
}
