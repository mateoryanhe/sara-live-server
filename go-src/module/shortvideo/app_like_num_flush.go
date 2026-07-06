package shortvideo

import (
	"context"
	"sort"
	"strconv"
	"sync/atomic"
	"time"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/xrtimer"
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
	var freeTime uint64
	if userId > 0 && row.IsPaid == entity.ShortVideoPaidYes {
		initFreeTime(userId, row.ID)
		watch := shortvideodao.GetOneShortVideoWatch(userId, row.ID)
		freeTime = watch.FreeTime
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
		AuthorNickname:   author.AuthorNickname,
		AuthorAvatar:     author.AuthorAvatar,
		LikeCount:        likeCount,
		ViewCount:        viewCount,
		WatchCount:       watchCount,
		Duration:         row.Duration,
		FreeWatchSeconds: row.FreeWatchSeconds,
		FreeTime:         freeTime,
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

func paginateAppShortVideoList(ctx context.Context, all []*shortvideodto.AppShortVideoItem, page, pageSize int) *shortvideodto.AppShortVideoListRes {
	page, pageSize = normalizeAppListPage(page, pageSize)
	total := len(all)
	start, end := appListPageRange(total, page, pageSize)
	userId := httpserver.GetAuthId(ctx)
	list := make([]*shortvideodto.AppShortVideoItem, 0, end-start)
	for _, item := range all[start:end] {
		if item == nil {
			continue
		}
		videoId, err := strconv.ParseUint(item.ID, 10, 64)
		if err != nil || videoId == 0 {
			continue
		}
		list = append(list, toAppShortVideoItem(shortvideodao.GetShortVideoById(videoId), shortvideodao.GetStatByVideoId(videoId), userId))
	}
	return &shortvideodto.AppShortVideoListRes{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     list,
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
