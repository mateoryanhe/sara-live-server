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
	liveentity "xr-game-server/entity/live"
	userentity "xr-game-server/entity/user"
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
	roomListCache      atomic.Value // []*liveentity.LiveRoom
	emptyRoomListCache = make([]*liveentity.LiveRoom, 0)
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
	rooms := make([]*liveentity.LiveRoom, 0, len(allData))
	for _, room := range allData {
		if room == nil {
			continue
		}
		userinfodao.GetUserInfoByUserId(room.ID)
		rooms = append(rooms, room)
	}

	onlineCounts := make(map[uint64]int, len(rooms))
	for _, room := range rooms {
		onlineCounts[room.ID] = countAudienceInRoom(room.ID)
	}

	sort.Slice(rooms, func(i, j int) bool {
		return compareLiveRoomsForList(rooms[i], rooms[j], onlineCounts)
	})

	roomListCache.Store(rooms)
}

// RefreshRoomListCache 主动刷新直播间列表缓存
func RefreshRoomListCache(ctx context.Context) {
	flushRoomList(ctx)
}

func isDisabledBotAnchorRoom(room *liveentity.LiveRoom) bool {
	if room == nil {
		return false
	}
	user := userinfodao.GetUserInfoByUserId(room.ID)
	return user != nil && user.IsBotAnchor() && user.BotAnchorStatus != userentity.BotAnchorStatusEnabled
}

// filterRoomsForApp App 端列表查询时过滤不可见直播间
func filterRoomsForApp(rooms []*liveentity.LiveRoom) []*liveentity.LiveRoom {
	if len(rooms) == 0 {
		return rooms
	}
	filtered := make([]*liveentity.LiveRoom, 0, len(rooms))
	for _, room := range rooms {
		if room == nil || IsRoomBanned(room) || IsRoomOffShelf(room) || isDisabledBotAnchorRoom(room) {
			continue
		}
		filtered = append(filtered, room)
	}
	return filtered
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

func compareLiveRoomsForList(a, b *liveentity.LiveRoom, onlineCounts map[uint64]int) bool {
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

func liveRoomHeartTimeUnix(room *liveentity.LiveRoom) int64 {
	if room == nil || room.HeartTime == nil {
		return 0
	}
	return room.HeartTime.Unix()
}

func toLiveRoomListItem(room *liveentity.LiveRoom, userId uint64) *liveroomdto.LiveRoomListItem {
	status := userstatus.LiveRoomStatusClosed
	if room.LiveRecordId > 0 {
		status = userstatus.LiveRoomStatusLive
	}
	cfg := liveroomdao.GetLiveRoomCfg(room.ID)
	item := &liveroomdto.LiveRoomListItem{
		RoomId:   strconv.FormatUint(room.ID, 10),
		GuildId:  strconv.FormatUint(room.GuildId, 10),
		Title:    room.Title,
		Cover:    room.Cover,
		Notice:   room.Notice,
		Status:   status,
		CreateAt: room.CreatedAt.Unix(),
	}
	if cfg != nil {
		item.Category = cfg.Category
		item.Ticket = cfg.Ticket
		item.Billing = cfg.Billing
		item.AllowCallIcon = allowShowCallIcon(room, cfg, userId)
		if cfg.TagId > 0 {
			item.TagId = strconv.FormatUint(cfg.TagId, 10)
			item.TagName = getRoomTagName(cfg.TagId)
		}
		item.IsBotAnchor, item.CloudPlayerVideo = resolveBotAnchorRoomInfo(room.ID, cfg.CloudPlayerVideo)
		item.IsTest = cfg.IsTest
	}
	if item.Cover != "" {
		item.Cover = upload.GetUrlByName(room.Cover)
	}

	if u := userinfodao.GetUserInfoByUserId(room.ID); u != nil {
		item.AnchorNickname = u.Nickname
		item.AnchorAvatar = upload.ResolveAvatarUrlForUser(room.ID, u.Avatar)
		item.UserType = u.UserType
	}
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

func buildLiveRoomListItems(rooms []*liveentity.LiveRoom, userId uint64) []*liveroomdto.LiveRoomListItem {
	list := make([]*liveroomdto.LiveRoomListItem, 0, len(rooms))
	for _, room := range rooms {
		if room == nil {
			continue
		}
		list = append(list, toLiveRoomListItem(room, userId))
	}
	return list
}

func getRoomListCache() []*liveentity.LiveRoom {
	v := roomListCache.Load()
	if v == nil {
		return nil
	}
	list, ok := v.([]*liveentity.LiveRoom)
	if !ok || len(list) == 0 {
		return nil
	}
	return list
}

func filterRoomsByStatus(rooms []*liveentity.LiveRoom, statusFilter int) []*liveentity.LiveRoom {
	if statusFilter == 0 {
		return rooms
	}
	filtered := make([]*liveentity.LiveRoom, 0, len(rooms))
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

func filterRoomsByQuery(rooms []*liveentity.LiveRoom, tagId uint64, title, notice string) []*liveentity.LiveRoom {
	titleKey := strings.ToLower(strings.TrimSpace(title))
	noticeKey := strings.ToLower(strings.TrimSpace(notice))
	if tagId == 0 && titleKey == "" && noticeKey == "" {
		return rooms
	}
	filtered := make([]*liveentity.LiveRoom, 0, len(rooms))
	for _, room := range rooms {
		if room == nil {
			continue
		}
		if tagId > 0 {
			cfg := liveroomdao.GetLiveRoomCfgFromCache(room.ID)
			if cfg == nil || cfg.TagId != tagId {
				continue
			}
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

func buildFollowingAnchorSet(userId uint64) map[uint64]struct{} {
	set := make(map[uint64]struct{})
	if userId == 0 {
		return set
	}
	followingTotal := userinfodao.GetFollowCount(userId)
	followings := livefollowdao.GetFollowingsByUser(userId, 1, followingTotal)
	for _, f := range followings {
		if f != nil && f.AnchorId > 0 {
			set[f.AnchorId] = struct{}{}
		}
	}
	return set
}

func filterRoomsByFollowing(rooms []*liveentity.LiveRoom, userId uint64) []*liveentity.LiveRoom {
	if userId == 0 {
		return make([]*liveentity.LiveRoom, 0)
	}
	following := buildFollowingAnchorSet(userId)
	if len(following) == 0 {
		return make([]*liveentity.LiveRoom, 0)
	}
	filtered := make([]*liveentity.LiveRoom, 0, len(following))
	for _, room := range rooms {
		if room == nil {
			continue
		}
		if _, ok := following[room.ID]; ok {
			filtered = append(filtered, room)
		}
	}
	return filtered
}

func filterRoomsByBlocked(rooms []*liveentity.LiveRoom, userId uint64) []*liveentity.LiveRoom {
	if userId == 0 || len(rooms) == 0 {
		return rooms
	}
	filtered := make([]*liveentity.LiveRoom, 0, len(rooms))
	for _, room := range rooms {
		if room == nil {
			continue
		}
		if livefollowdao.IsBlocked(userId, room.ID) {
			continue
		}
		filtered = append(filtered, room)
	}
	return filtered
}

// viewerCanSeeSeniorAnchorRoom App 列表中高级主播直播间对已登录用户可见(含 VIP=0,主播本人始终可见)
func viewerCanSeeSeniorAnchorRoom(viewerUserId uint64, room *liveentity.LiveRoom) bool {
	if room == nil {
		return false
	}
	if viewerUserId == room.ID {
		return true
	}
	anchor := userinfodao.GetUserInfoByUserId(room.ID)
	if anchor == nil || anchor.UserType != userentity.UserTypeSeniorAnchor {
		return true
	}
	if viewerUserId == 0 {
		return false
	}
	return userinfodao.GetUserInfoByUserId(viewerUserId) != nil
}

func filterRoomsBySeniorAnchor(rooms []*liveentity.LiveRoom, viewerUserId uint64) []*liveentity.LiveRoom {
	if len(rooms) == 0 {
		return rooms
	}
	filtered := make([]*liveentity.LiveRoom, 0, len(rooms))
	for _, room := range rooms {
		if room == nil {
			continue
		}
		if viewerCanSeeSeniorAnchorRoom(viewerUserId, room) {
			filtered = append(filtered, room)
		}
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

	cached := filterRoomsForApp(getRoomListCache())
	if cached == nil {
		return &liveroomdto.GetLiveRoomListRes{
			Total:    0,
			Page:     page,
			PageSize: pageSize,
			List:     make([]*liveroomdto.LiveRoomListItem, 0),
		}, nil
	}

	filtered := filterRoomsByStatus(cached, req.StatusFilter)
	tagId := req.TagId
	switch resolveSpecialRoomTagFilterMode(tagId) {
	case specialRoomTagNameAll:
		tagId = 0
	case specialRoomTagNameMy:
		filtered = filterRoomsByFollowing(filtered, userId)
		tagId = 0
	}
	filtered = filterRoomsByQuery(filtered, tagId, req.Title, req.Notice)
	filtered = filterRoomsBySeniorAnchor(filtered, userId)
	filtered = filterRoomsByBlocked(filtered, userId)
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

	followingTotal := userinfodao.GetFollowCount(userId)
	followings := livefollowdao.GetFollowingsByUser(userId, 1, followingTotal)
	rooms := make([]*liveentity.LiveRoom, 0, len(followings))
	for _, f := range followings {
		if f == nil {
			continue
		}
		room := liveroomdao.GetRoomById(f.AnchorId)
		if room == nil || IsRoomBanned(room) || IsRoomOffShelf(room) || isDisabledBotAnchorRoom(room) {
			continue
		}
		rooms = append(rooms, room)
	}

	filtered := filterRoomsByStatus(rooms, req.StatusFilter)
	filtered = filterRoomsBySeniorAnchor(filtered, userId)
	filtered = filterRoomsByBlocked(filtered, userId)
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

func findLiveRoomIndex(rooms []*liveentity.LiveRoom, roomId uint64) int {
	for i, room := range rooms {
		if room != nil && room.ID == roomId {
			return i
		}
	}
	return -1
}

func collectNearbyLiveRooms(rooms []*liveentity.LiveRoom, currentIdx, direction, count int) []*liveentity.LiveRoom {
	if currentIdx < 0 || len(rooms) == 0 || count <= 0 {
		return make([]*liveentity.LiveRoom, 0)
	}
	result := make([]*liveentity.LiveRoom, 0, count)
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

func buildHotLiveRoomListItems(rooms []*liveentity.LiveRoom, userId uint64, startRank int) []*liveroomdto.HotLiveRoomListItem {
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

// GetHotLiveRoomList App 分页查询 Hot 列表(暂拉取全部直播间,走内存缓存排序,含排名)
func GetHotLiveRoomList(ctx context.Context, req *liveroomdto.GetHotLiveRoomListReq) (*liveroomdto.GetHotLiveRoomListRes, error) {
	userId := httpserver.GetAuthId(ctx)
	pageIndex, pageSize := normalizeRoomListPage(req.PageIndex, req.PageSize)

	cached := filterRoomsForApp(getRoomListCache())
	if cached == nil {
		return &liveroomdto.GetHotLiveRoomListRes{
			Total:     0,
			PageIndex: pageIndex,
			PageSize:  pageSize,
			List:      make([]*liveroomdto.HotLiveRoomListItem, 0),
		}, nil
	}

	filtered := filterRoomsBySeniorAnchor(cached, userId)
	filtered = filterRoomsByBlocked(filtered, userId)
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

	cached := filterRoomsForApp(getRoomListCache())
	if cached == nil {
		return &liveroomdto.GetNearbyLiveRoomListRes{
			List: make([]*liveroomdto.LiveRoomListItem, 0),
		}, nil
	}

	liveRooms := filterRoomsByStatus(cached, int(userstatus.LiveRoomStatusLive))
	liveRooms = filterRoomsBySeniorAnchor(liveRooms, userId)
	liveRooms = filterRoomsByBlocked(liveRooms, userId)
	currentIdx := findLiveRoomIndex(liveRooms, req.RoomId)
	nearbyRooms := collectNearbyLiveRooms(liveRooms, currentIdx, req.Direction, count)

	list := buildLiveRoomListItems(nearbyRooms, userId)
	if len(list) == 0 {
		list = make([]*liveroomdto.LiveRoomListItem, 0)
	}
	return &liveroomdto.GetNearbyLiveRoomListRes{List: list}, nil
}
