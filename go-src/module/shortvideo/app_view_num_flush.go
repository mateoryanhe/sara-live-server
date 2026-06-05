package shortvideo

import (
	"context"
	"sort"
	"sync/atomic"
	"xr-game-server/core/xrtimer"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dto/shortvideodto"

	"github.com/gogf/gf/v2/os/gctx"
)

var (
	appViewListCache      atomic.Value // []*shortvideodto.AppShortVideoItem
	emptyAppShortViewList = make([]*shortvideodto.AppShortVideoItem, 0)
)

func initAppViewListCache() {
	loadAppShortVideoViewListCache()
	xrtimer.AddSingleton(gctx.New(), appListRefreshInterval, func(ctx context.Context) {
		loadAppShortVideoViewListCache()
	})
}

func loadAppShortVideoViewListCache() {
	rows := shortvideodao.GetOnShelfShortVideos()

	list := make([]*shortvideodto.AppShortVideoItem, 0, len(rows))
	for _, row := range rows {
		stat := shortvideodao.GetStatByVideoId(row.ID)
		list = append(list, toAppShortVideoItem(row, stat))
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].ViewCount > list[j].ViewCount
	})
	appViewListCache.Store(list)
}

func getAppShortVideoViewListCache() []*shortvideodto.AppShortVideoItem {
	v := appViewListCache.Load()
	if v == nil {
		return emptyAppShortViewList
	}
	list, ok := v.([]*shortvideodto.AppShortVideoItem)
	if !ok || list == nil {
		return emptyAppShortViewList
	}
	return list
}

// GetAppShortVideoViewList App端分页查询短视频列表(仅已上架,按观看人数排序,走内存缓存)
func GetAppShortVideoViewList(_ context.Context, req *shortvideodto.AppShortVideoViewListReq) (*shortvideodto.AppShortVideoListRes, error) {
	page, pageSize := normalizeAppListPage(req.Page, req.PageSize)
	all := getAppShortVideoViewListCache()
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
