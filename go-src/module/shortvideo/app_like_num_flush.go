package shortvideo

import (
	"context"
	"sort"
	"strconv"
	"sync/atomic"
	"time"
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
		list = append(list, toAppShortVideoItem(row, stat))
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

func toAppShortVideoItem(row *entity.ShortVideo, stat *entity.ShortVideoStat) *shortvideodto.AppShortVideoItem {
	var likeCount, viewCount uint64
	if stat != nil {
		likeCount = stat.LikeCount
		viewCount = stat.ViewCount
	}
	author := getShortVideoAuthorInfo(row.AuthorId)
	item := &shortvideodto.AppShortVideoItem{
		ID:               strconv.FormatUint(row.ID, 10),
		Title:            row.Title,
		Video:            upload.GetUrlByName(row.Video),
		Cover:            upload.GetUrlByName(row.Cover),
		IsPaid:           row.IsPaid,
		DiamondPerMinute: row.DiamondPerMinute,
		CategoryId:       row.CategoryId,
		Source:           row.Source,
		AuthorId:         author.AuthorId,
		AuthorNickname:   author.AuthorNickname,
		AuthorAvatar:     author.AuthorAvatar,
		LikeCount:        likeCount,
		ViewCount:        viewCount,
		Duration:         row.Duration,
		FreeWatchSeconds: row.FreeWatchSeconds,
	}
	return item
}

func getShortVideoAuthorInfo(authorId uint64) shortVideoAuthorInfo {
	ret := shortVideoAuthorInfo{AuthorId: strconv.FormatUint(authorId, 10)}
	if authorId == 0 {
		return ret
	}
	if u := userinfodao.GetUserInfoByUserId(authorId); u != nil {
		ret.AuthorNickname = u.Nickname
		ret.AuthorAvatar = upload.ResolveAvatarUrl(u.Avatar)
	}
	return ret
}

type shortVideoAuthorInfo struct {
	AuthorId       string
	AuthorNickname string
	AuthorAvatar   string
}

// GetAppShortVideoList App端分页查询短视频列表(仅已上架,按点赞数排序,走内存缓存)
func GetAppShortVideoList(_ context.Context, req *shortvideodto.AppShortVideoListReq) (*shortvideodto.AppShortVideoListRes, error) {
	page, pageSize := normalizeAppListPage(req.Page, req.PageSize)
	all := getAppShortVideoListCache()
	total := len(all)
	start, end := appListPageRange(total, page, pageSize)
	pageData := all[start:end]

	return &shortvideodto.AppShortVideoListRes{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     pageData,
	}, nil
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
