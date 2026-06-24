package shortvideo

import (
	"context"
	"sort"
	"strconv"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

// GetAppShortVideoWatchList App端分页查询当前用户短视频观看记录
func GetAppShortVideoWatchList(ctx context.Context, req *shortvideodto.AppShortVideoWatchListReq) (*shortvideodto.AppShortVideoWatchListRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}

	page, pageSize := normalizeAppListPage(req.Page, req.PageSize)
	all := shortvideodao.GetShortVideoWatch(userId)
	sort.Slice(all, func(i, j int) bool {
		return all[i].UpdatedAt.After(all[j].UpdatedAt)
	})
	total := len(all)
	start, end := appListPageRange(total, page, pageSize)

	list := make([]*shortvideodto.AppShortVideoWatchListItem, 0, end-start)
	for _, row := range all[start:end] {
		if row == nil {
			continue
		}
		list = append(list, toAppShortVideoWatchListItem(row))
	}

	return &shortvideodto.AppShortVideoWatchListRes{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     list,
	}, nil
}

func toAppShortVideoWatchListItem(row *entity.ShortVideoWatch) *shortvideodto.AppShortVideoWatchListItem {
	item := &shortvideodto.AppShortVideoWatchListItem{
		ID:            row.ID,
		VideoId:       strconv.FormatUint(row.VideoId, 10),
		BilledSeconds: row.BilledSeconds,
		WatchSeconds:  row.WatchSeconds,
		UpdatedAt:     formatShortVideoCfgTime(row.UpdatedAt),
	}
	if video := shortvideodao.GetShortVideoById(row.VideoId); video != nil {
		item.Title = video.Title
		item.Cover = upload.GetUrlByName(video.Cover)
		author := getShortVideoAuthorInfo(video.AuthorId)
		item.AuthorId = author.AuthorId
		item.AuthorNickname = author.AuthorNickname
		item.AuthorAvatar = author.AuthorAvatar
	}
	return item
}
