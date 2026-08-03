package shortvideo

import (
	"context"
	"sort"
	"strconv"
	"sync/atomic"
	"time"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/xrtimer"
	"xr-game-server/dao/livefollowdao"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/entity"
	"xr-game-server/module/upload"

	"github.com/gogf/gf/v2/os/gctx"
)

const (
	appListDefaultPageSize = 20
	appListMaxPageSize     = 800
	appListRefreshInterval = 15 * time.Minute
	appScrollDefaultCount  = 10
	appScrollMaxCount      = 50
)

var (
	appListCache      atomic.Value // []*shortvideodto.AppShortVideoItem
	emptyAppShortList = make([]*shortvideodto.AppShortVideoItem, 0)
)

// Init App端短视频列表缓存初始化,并每15分钟整体替换一次
func initAppListCache() {

	loadAppShortVideoListCache()
	xrtimer.AddSingleton(gctx.New(), appListRefreshInterval, func(ctx context.Context) {
		loadAppShortVideoListCache()
	})
}

// loadAppShortVideoListCache 构建新列表后整体替换,刷新过程中请求仍读取旧列表
func loadAppShortVideoListCache() {
	rows := shortvideodao.GetOnShelfShortVideos()

	list := make([]*shortvideodto.AppShortVideoItem, 0, len(rows))
	for _, row := range rows {
		stat := shortvideodao.GetStatByVideoId(row.ID)
		list = append(list, toAppShortVideoItem(row, stat, 0))
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].LikeCount > list[j].LikeCount
	})
	appListCache.Store(list)
}

// reloadAllAppShortVideoListCaches 立即刷新 App 端三个排序列表缓存
func reloadAllAppShortVideoListCaches() {
	loadAppShortVideoListCache()
	loadAppShortVideoPublishListCache()
	loadAppShortVideoViewListCache()
}

func getAppShortVideoListCache() []*shortvideodto.AppShortVideoItem {
	v := appListCache.Load()
	if v == nil {
		return emptyAppShortList
	}
	list, ok := v.([]*shortvideodto.AppShortVideoItem)
	if !ok || list == nil {
		return emptyAppShortList
	}
	return list
}

func toAppShortVideoItem(row *entity.ShortVideo, stat *entity.ShortVideoStat, userId uint64) *shortvideodto.AppShortVideoItem {
	if row == nil {
		return nil
	}
	var likeCount, viewCount, watchCount uint64
	if stat != nil {
		likeCount = stat.LikeCount
		viewCount = stat.ViewCount
		watchCount = stat.WatchCount
	}
	hasPaid := false
	if userId > 0 {
		watch := shortvideodao.GetOneShortVideoWatch(userId, row.ID)
		hasPaid = watch.PaidTime != nil
	}
	author := getShortVideoAuthorInfo(row.AuthorId)
	return &shortvideodto.AppShortVideoItem{
		ID:               strconv.FormatUint(row.ID, 10),
		Title:            row.Title,
		Video:            upload.GetUrlByName(row.Video),
		Cover:            upload.GetUrlByName(row.Cover),
		IsPaid:           row.IsPaid,
		PayDiamond:       row.PayDiamond,
		CategoryId:       row.CategoryId,
		Source:           row.Source,
		AuthorId:         author.AuthorId,
		AuthorType:       row.AuthorType,
		AuthorNickname:   author.AuthorNickname,
		AuthorAvatar:     author.AuthorAvatar,
		LikeCount:        likeCount,
		ViewCount:        viewCount,
		WatchCount:       watchCount,
		Duration:         row.Duration,
		FreeWatchSeconds: row.FreeWatchSeconds,
		HasPaid:          hasPaid,
	}
}

func getShortVideoAuthorInfo(authorId uint64) shortVideoAuthorInfo {
	ret := shortVideoAuthorInfo{AuthorId: strconv.FormatUint(authorId, 10)}
	if authorId == 0 {
		return ret
	}
	if u := userinfodao.GetUserInfoByUserId(authorId); u != nil {
		ret.AuthorNickname = u.Nickname
		ret.AuthorAvatar = upload.ResolveAvatarUrlForUser(authorId, u.Avatar)
	}
	return ret
}

type shortVideoAuthorInfo struct {
	AuthorId       string
	AuthorNickname string
	AuthorAvatar   string
}

// GetAppShortVideoList App端分页查询短视频列表(仅已上架,按点赞数排序,走内存缓存)
func GetAppShortVideoList(ctx context.Context, req *shortvideodto.AppShortVideoListReq) (*shortvideodto.AppShortVideoListRes, error) {
	return paginateAppShortVideoList(ctx, getAppShortVideoListCache(), req.Page, req.PageSize), nil
}

// GetAppShortVideoScrollList 以当前视频为锚点,向上/向下拉取n个短视频
func GetAppShortVideoScrollList(ctx context.Context, req *shortvideodto.AppShortVideoScrollReq) (*shortvideodto.AppShortVideoScrollRes, error) {
	all := getAppShortVideoListCacheBySort(req.SortType)
	userId := httpserver.GetAuthId(ctx)
	all = filterAppShortVideoListByBlocked(all, userId)

	count := normalizeAppScrollCount(req.Count)
	anchorIdx := findAppShortVideoIndex(all, req.VideoId)

	var slice []*shortvideodto.AppShortVideoItem
	hasMore := false
	if req.Direction == shortvideodto.AppShortVideoScrollNext {
		slice, hasMore = sliceAppShortVideoAfter(all, anchorIdx, count)
	} else {
		slice, hasMore = sliceAppShortVideoBefore(all, anchorIdx, count)
	}
	return &shortvideodto.AppShortVideoScrollRes{
		List:    buildAppShortVideoListForUser(ctx, slice),
		HasMore: hasMore,
	}, nil
}

func getAppShortVideoListCacheBySort(sortType int) []*shortvideodto.AppShortVideoItem {
	switch sortType {
	case shortvideodto.AppShortVideoSortView:
		return getAppShortVideoViewListCache()
	case shortvideodto.AppShortVideoSortPublish:
		return getAppShortVideoPublishListCache()
	default:
		return getAppShortVideoListCache()
	}
}

func findAppShortVideoIndex(list []*shortvideodto.AppShortVideoItem, videoId string) int {
	for i, item := range list {
		if item != nil && item.ID == videoId {
			return i
		}
	}
	return -1
}

// sliceAppShortVideoAfter 取锚点之后的n个(向下/下一个); 锚点不存在时从列表开头取
func sliceAppShortVideoAfter(list []*shortvideodto.AppShortVideoItem, anchorIdx, count int) ([]*shortvideodto.AppShortVideoItem, bool) {
	start := 0
	if anchorIdx >= 0 {
		start = anchorIdx + 1
	}
	if start >= len(list) {
		return emptyAppShortList, false
	}
	end := start + count
	hasMore := end < len(list)
	if end > len(list) {
		end = len(list)
	}
	return list[start:end], hasMore
}

// sliceAppShortVideoBefore 取锚点之前的n个(向上/上一个),按原序返回(不含锚点)
func sliceAppShortVideoBefore(list []*shortvideodto.AppShortVideoItem, anchorIdx, count int) ([]*shortvideodto.AppShortVideoItem, bool) {
	if anchorIdx <= 0 {
		return emptyAppShortList, false
	}
	start := anchorIdx - count
	if start < 0 {
		start = 0
	}
	hasMore := start > 0
	return list[start:anchorIdx], hasMore
}

func normalizeAppScrollCount(count int) int {
	if count <= 0 {
		return appScrollDefaultCount
	}
	if count > appScrollMaxCount {
		return appScrollMaxCount
	}
	return count
}

func buildAppShortVideoListForUser(ctx context.Context, items []*shortvideodto.AppShortVideoItem) []*shortvideodto.AppShortVideoItem {
	userId := httpserver.GetAuthId(ctx)
	list := make([]*shortvideodto.AppShortVideoItem, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		videoId, err := strconv.ParseUint(item.ID, 10, 64)
		if err != nil || videoId == 0 {
			continue
		}
		list = append(list, toAppShortVideoItem(shortvideodao.GetShortVideoById(videoId), shortvideodao.GetStatByVideoId(videoId), userId))
	}
	return list
}

func filterAppShortVideoListByBlocked(list []*shortvideodto.AppShortVideoItem, userId uint64) []*shortvideodto.AppShortVideoItem {
	if userId == 0 || len(list) == 0 {
		return list
	}
	filtered := make([]*shortvideodto.AppShortVideoItem, 0, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}
		if item.AuthorType == entity.ShortVideoAuthorTypeCMS {
			filtered = append(filtered, item)
			continue
		}
		authorId, err := strconv.ParseUint(item.AuthorId, 10, 64)
		if err != nil || authorId == 0 {
			filtered = append(filtered, item)
			continue
		}
		if livefollowdao.IsBlocked(userId, authorId) {
			continue
		}
		filtered = append(filtered, item)
	}
	return filtered
}

func paginateAppShortVideoList(ctx context.Context, all []*shortvideodto.AppShortVideoItem, page, pageSize int) *shortvideodto.AppShortVideoListRes {
	userId := httpserver.GetAuthId(ctx)
	all = filterAppShortVideoListByBlocked(all, userId)
	page, pageSize = normalizeAppListPage(page, pageSize)
	total := len(all)
	start, end := appListPageRange(total, page, pageSize)
	return &shortvideodto.AppShortVideoListRes{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     buildAppShortVideoListForUser(ctx, all[start:end]),
	}
}

func normalizeAppListPage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = appListDefaultPageSize
	}
	if pageSize > appListMaxPageSize {
		pageSize = appListMaxPageSize
	}
	return page, pageSize
}

func appListPageRange(total, page, pageSize int) (int, int) {
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
