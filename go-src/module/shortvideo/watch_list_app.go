package shortvideo

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/entity/shortvideo"
	"xr-game-server/errercode"
)

// GetAppShortVideoWatchList App端分页查询当前用户短视频观看记录
func GetAppShortVideoWatchList(ctx context.Context, req *shortvideodto.AppShortVideoWatchListReq) (*shortvideodto.AppShortVideoWatchListRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}

	page, pageSize := normalizeWatchListPage(req.Page, req.PageSize)
	filtered := filterAppShortVideoWatchList(shortvideodao.GetUserShortVideoWatchAll(userId), userId)
	rows := paginateShortVideoWatchList(filtered, page, pageSize)

	list := make([]*shortvideodto.AppShortVideoItem, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		video := shortvideodao.GetShortVideoById(row.VideoId)
		if video == nil {
			continue
		}
		list = append(list, toAppShortVideoItem(video, shortvideodao.GetStatByVideoId(row.VideoId), userId))
	}

	return &shortvideodto.AppShortVideoWatchListRes{
		Page:     page,
		PageSize: pageSize,
		List:     list,
	}, nil
}

func filterAppShortVideoWatchList(rows []*entity.ShortVideoWatch, userId uint64) []*entity.ShortVideoWatch {
	list := make([]*entity.ShortVideoWatch, 0, len(rows))
	for _, watch := range rows {
		if watch == nil || watch.VideoId == 0 {
			continue
		}
		if watch.ViewCounted != entity.ShortVideoWatchViewCountedYes {
			continue
		}
		video := shortvideodao.GetShortVideoById(watch.VideoId)
		if video == nil || video.Status != entity.ShortVideoStatusOnShelf {
			continue
		}
		if video.AuthorId == userId {
			continue
		}
		if video.IsPaid == entity.ShortVideoPaidYes && watch.PaidTime == nil {
			continue
		}
		list = append(list, watch)
	}
	return list
}

func paginateShortVideoWatchList(rows []*entity.ShortVideoWatch, page, pageSize int) []*entity.ShortVideoWatch {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = shortvideodao.WatchListCachePageSize
	}
	start := (page - 1) * pageSize
	if start >= len(rows) {
		return make([]*entity.ShortVideoWatch, 0)
	}
	end := start + pageSize
	if end > len(rows) {
		end = len(rows)
	}
	return rows[start:end]
}

func normalizeWatchListPage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = shortvideodao.WatchListCachePageSize
	}
	if pageSize > appListMaxPageSize {
		pageSize = appListMaxPageSize
	}
	return page, pageSize
}
