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
		if room == nil || IsRoomBanned(room) {
			continue
		}
		filtered = append(filtered, room)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return liveRoomHeartTimeUnix(filtered[i]) > liveRoomHeartTimeUnix(filtered[j])
	})

	roomListCache.Store(filtered)
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
		RoomId:   strconv.FormatUint(room.ID, 10),
		GuildId:  strconv.FormatUint(room.GuildId, 10),
		Title:    room.Title,
		Cover:    room.Cover,
		Notice:   room.Notice,
		Status:   status,
		Category: room.Category,
		Ticket:   room.Ticket,
		Billing:  room.Billing,
		CreateAt: room.CreatedAt.Unix(),
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
		item.AnchorAvatar = upload.ResolveAvatarUrl(u.Avatar)
	}

	if userId > 0 {
		if token, expireAt, err := agora.BuildLiveRoomToken(userId, room.ID); err == nil {
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
