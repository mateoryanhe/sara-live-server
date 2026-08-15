package shortvideo

import (
	"context"
	"strconv"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/entity/shortvideo"
	"xr-game-server/errercode"
)

// GetAppShortVideoStatList App端分页查询本人短视频统计数据(直接查库,按发布时间降序)
func GetAppShortVideoStatList(ctx context.Context, req *shortvideodto.AppShortVideoStatListReq) (*shortvideodto.AppShortVideoStatListRes, error) {
	authorId := httpserver.GetAuthId(ctx)
	if authorId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	page, pageSize := normalizeAppListPage(req.Page, req.PageSize)
	rows, err := shortvideodao.ListStatPageByAuthorId(authorId, page, pageSize)
	if err != nil {
		return nil, err
	}

	list := make([]*shortvideodto.AppShortVideoStatItem, 0, len(rows))
	for _, row := range rows {
		stat := row
		if cached := shortvideodao.GetStatFromCacheByVideoId(row.ID); cached != nil {
			stat = cached
		}
		if item := toAppShortVideoStatItem(stat); item != nil {
			list = append(list, item)
		}
	}
	return &shortvideodto.AppShortVideoStatListRes{
		Page:     page,
		PageSize: pageSize,
		List:     list,
	}, nil
}

func toAppShortVideoStatItem(row *entity.ShortVideoStat) *shortvideodto.AppShortVideoStatItem {
	if row == nil || row.ID == 0 {
		return nil
	}
	var publishedAt int64
	if !row.CreatedAt.IsZero() {
		publishedAt = row.CreatedAt.UnixMilli()
	}
	return &shortvideodto.AppShortVideoStatItem{
		VideoId:            strconv.FormatUint(row.ID, 10),
		Title:              row.Title,
		LikeCount:          row.LikeCount,
		ViewCount:          row.ViewCount,
		WatchCount:         row.WatchCount,
		TotalDiamondIncome: row.TotalDiamondIncome,
		PublishedAt:        publishedAt,
	}
}
