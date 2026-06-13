package liveroom

import (
	"context"
	"sort"
	"strconv"
	"sync/atomic"
	"time"
	"xr-game-server/constants/userstatus"
	"xr-game-server/core/xrtimer"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/entity"
	"xr-game-server/module/upload"

	"github.com/gogf/gf/v2/os/gctx"
)

const (
	roomListDefaultPageSize = 20
	roomListMaxPageSize     = 100
	roomListRefreshInterval = 10 * time.Minute
)

var (
	roomListCache      atomic.Value // []*liveroomdto.LiveRoomListItem
	emptyRoomListCache = make([]*liveroomdto.LiveRoomListItem, 0)
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

	list := make([]*liveroomdto.LiveRoomListItem, 0, len(filtered))
	for _, room := range filtered {
		list = append(list, toLiveRoomListItem(room))
	}
	roomListCache.Store(list)
}

func liveRoomHeartTimeUnix(room *entity.LiveRoom) int64 {
	if room == nil || room.HeartTime == nil {
		return 0
	}
	return room.HeartTime.Unix()
}

func toLiveRoomListItem(room *entity.LiveRoom) *liveroomdto.LiveRoomListItem {
	status := userstatus.LiveRoomStatusClosed
	if room.LiveRecordId > 0 {
		status = userstatus.LiveRoomStatusLive
	}
	item := &liveroomdto.LiveRoomListItem{
		RoomId:   strconv.FormatUint(room.ID, 10),
		GuildId:  strconv.FormatUint(room.GuildId, 10),
		Title:    room.Title,
		Cover:    upload.GetUrlByName(room.Cover),
		Notice:   room.Notice,
		Status:   status,
		Category: room.Category,
		Ticket:   room.Ticket,
		Billing:  room.Billing,
		CreateAt: room.CreatedAt.Unix(),
	}
	if u := userinfodao.GetUserInfoByUserId(room.ID); u != nil {
		item.AnchorNickname = u.Nickname
		item.AnchorAvatar = upload.ResolveAvatarUrl(u.Avatar)
	}
	return item
}

func getRoomListCache() []*liveroomdto.LiveRoomListItem {
	v := roomListCache.Load()
	if v == nil {
		return nil
	}
	list, ok := v.([]*liveroomdto.LiveRoomListItem)
	if !ok || len(list) == 0 {
		return nil
	}
	return list
}

func filterRoomListByStatus(all []*liveroomdto.LiveRoomListItem, statusFilter int) []*liveroomdto.LiveRoomListItem {
	if statusFilter == 0 {
		return all
	}
	filtered := make([]*liveroomdto.LiveRoomListItem, 0, len(all))
	for _, item := range all {
		if item == nil {
			continue
		}
		switch statusFilter {
		case 1:
			if item.Status == userstatus.LiveRoomStatusLive {
				filtered = append(filtered, item)
			}
		case 2:
			if item.Status == userstatus.LiveRoomStatusClosed {
				filtered = append(filtered, item)
			}
		default:
			filtered = append(filtered, item)
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
func GetRoomList(_ context.Context, req *liveroomdto.GetLiveRoomListReq) (*liveroomdto.GetLiveRoomListRes, error) {
	page, pageSize := normalizeRoomListPage(req.Page, req.PageSize)

	all := getRoomListCache()
	if all == nil {
		return &liveroomdto.GetLiveRoomListRes{
			Total:    0,
			Page:     page,
			PageSize: pageSize,
			List:     emptyRoomListCache,
		}, nil
	}

	all = filterRoomListByStatus(all, req.StatusFilter)
	total := len(all)
	start, end := roomListPageRange(total, page, pageSize)

	return &liveroomdto.GetLiveRoomListRes{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     all[start:end],
	}, nil
}
