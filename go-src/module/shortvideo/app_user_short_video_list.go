package shortvideo

import (
	"context"
	"sort"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/livefollowdao"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/entity/shortvideo"
	"xr-game-server/errercode"
)

// GetAppUserShortVideoList App端查询指定用户短视频列表(仅已上架,按审核时间降序,走内存缓存)
func GetAppUserShortVideoList(ctx context.Context, req *shortvideodto.AppUserShortVideoListReq) (*shortvideodto.AppShortVideoListRes, error) {
	if req.UserId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	viewerId := httpserver.GetAuthId(ctx)
	if viewerId > 0 && livefollowdao.IsBlocked(viewerId, req.UserId) {
		page, pageSize := normalizeAppListPage(req.Page, req.PageSize)
		return &shortvideodto.AppShortVideoListRes{
			Total:    0,
			Page:     page,
			PageSize: pageSize,
			List:     emptyAppShortList,
		}, nil
	}

	rows := shortvideodao.GetAuthorOnShelfShortVideos(req.UserId)
	sort.Slice(rows, func(i, j int) bool {
		return compareShortVideoByUpdatedAt(rows[i], rows[j])
	})

	page, pageSize := normalizeAppListPage(req.Page, req.PageSize)
	total := len(rows)
	start, end := appListPageRange(total, page, pageSize)
	list := make([]*shortvideodto.AppShortVideoItem, 0, end-start)
	for _, row := range rows[start:end] {
		list = append(list, toAppShortVideoItem(row, shortvideodao.GetStatByVideoId(row.ID), viewerId))
	}
	return &shortvideodto.AppShortVideoListRes{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     list,
	}, nil
}

// compareShortVideoByUpdatedAt 按审核时间(updatedAt)降序,相同时间按 ID 降序
func compareShortVideoByUpdatedAt(a, b *entity.ShortVideo) bool {
	if a == nil || b == nil {
		return a != nil
	}
	if a.UpdatedAt.Equal(b.UpdatedAt) {
		return a.ID > b.ID
	}
	return a.UpdatedAt.After(b.UpdatedAt)
}
