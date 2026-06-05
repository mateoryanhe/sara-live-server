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
	appPublishListCache      atomic.Value // []*shortvideodto.AppShortVideoItem
	emptyAppShortPublishList = make([]*shortvideodto.AppShortVideoItem, 0)
)

func initAppPublishListCache() {
	loadAppShortVideoPublishListCache()
	xrtimer.AddSingleton(gctx.New(), appListRefreshInterval, func(ctx context.Context) {
		loadAppShortVideoPublishListCache()
	})
}

func loadAppShortVideoPublishListCache() {
	rows := shortvideodao.GetOnShelfShortVideos()
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].CreatedAt.After(rows[j].CreatedAt)
	})

	list := make([]*shortvideodto.AppShortVideoItem, 0, len(rows))
	for _, row := range rows {
		stat := shortvideodao.GetStatByVideoId(row.ID)
		list = append(list, toAppShortVideoItem(row, stat))
	}
	appPublishListCache.Store(list)
}

func getAppShortVideoPublishListCache() []*shortvideodto.AppShortVideoItem {
	v := appPublishListCache.Load()
	if v == nil {
		return emptyAppShortPublishList
	}
	list, ok := v.([]*shortvideodto.AppShortVideoItem)
	if !ok || list == nil {
		return emptyAppShortPublishList
	}
	return list
}

// GetAppShortVideoPublishList App端分页查询短视频列表(仅已上架,按发布时间降序,走内存缓存)
func GetAppShortVideoPublishList(_ context.Context, req *shortvideodto.AppShortVideoPublishListReq) (*shortvideodto.AppShortVideoListRes, error) {
	page, pageSize := normalizeAppListPage(req.Page, req.PageSize)
	all := getAppShortVideoPublishListCache()
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
