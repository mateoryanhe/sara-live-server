package shortvideo

import (
	"context"
	"strconv"
	"sync/atomic"
	"time"
	"xr-game-server/core/xrtimer"
	"xr-game-server/dao/shortvideodao"
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
func Init() {
	reloadShortVideoCfgMemory()
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
		list = append(list, toAppShortVideoItem(row))
	}
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

func toAppShortVideoItem(row *entity.ShortVideo) *shortvideodto.AppShortVideoItem {
	return &shortvideodto.AppShortVideoItem{
		ID:          strconv.FormatUint(row.ID, 10),
		Title:       row.Title,
		Video:       upload.GetUrlByName(row.Video),
		Cover:       upload.GetUrlByName(row.Cover),
		Description: row.Description,
		IsPaid:      row.IsPaid,
		LikeCount:   row.LikeCount,
	}
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
