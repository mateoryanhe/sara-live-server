package shortvideo

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/errercode"
)

// GetAppShortVideoWatchList App端分页查询当前用户短视频观看记录
func GetAppShortVideoWatchList(ctx context.Context, req *shortvideodto.AppShortVideoWatchListReq) (*shortvideodto.AppShortVideoListRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}

	page, pageSize := normalizeAppListPage(req.Page, req.PageSize)
	rows := shortvideodao.GetShortVideoWatch(userId, page, pageSize)

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
	list = filterAppShortVideoListByBlocked(list, userId)

	return &shortvideodto.AppShortVideoListRes{
		Page:     page,
		PageSize: pageSize,
		List:     list,
	}, nil
}
