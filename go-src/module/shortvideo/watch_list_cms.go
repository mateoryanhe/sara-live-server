package shortvideo

import (
	"context"
	"strconv"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/shortvideodto"
)

func GetShortVideoWatchList(_ context.Context, req *shortvideodto.ShortVideoWatchListReq) (*httpserver.CMSQueryResp, error) {
	total, rows := shortvideodao.GetList(&shortvideodao.WatchListFilter{
		UserId:    parseUint64Filter(req.UserId),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		PageIndex: req.PageIndex,
		PageSize:  req.PageSize,
	})

	list := make([]*shortvideodto.ShortVideoWatchListItem, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		item := &shortvideodto.ShortVideoWatchListItem{
			ID:            row.ID,
			UserId:        strconv.FormatUint(row.UserId, 10),
			VideoId:       strconv.FormatUint(row.VideoId, 10),
			BilledSeconds: row.BilledSeconds,
			CreatedAt:     formatShortVideoCfgTime(row.CreatedAt),
			UpdatedAt:     formatShortVideoCfgTime(row.UpdatedAt),
		}
		if user := userinfodao.GetUserInfoByUserId(row.UserId); user != nil {
			item.Nickname = user.Nickname
		}
		if video := shortvideodao.GetShortVideoById(row.VideoId); video != nil {
			item.VideoTitle = video.Title
		}
		list = append(list, item)
	}
	return &httpserver.CMSQueryResp{Total: total, Data: list}, nil
}

func parseUint64Filter(val string) uint64 {
	if val == "" {
		return 0
	}
	n, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0
	}
	return n
}
