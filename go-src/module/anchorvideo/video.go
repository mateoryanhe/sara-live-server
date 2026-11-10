package anchorvideo

import (
	"context"
	"strings"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/anchorvideodao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/anchorvideodto"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

// GetVideoList App 查询主播视频列表(读缓存)
func GetVideoList(_ context.Context, req *anchorvideodto.VideoListReq) (*anchorvideodto.VideoListRes, error) {
	if req == nil || req.AnchorId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	rows := anchorvideodao.GetVideoList(req.AnchorId)
	list := make([]*anchorvideodto.VideoItem, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.ID == 0 || row.IsDeleted {
			continue
		}
		list = append(list, &anchorvideodto.VideoItem{
			Id:    row.ID,
			Video: upload.GetUrlByName(row.Video),
		})
	}
	return &anchorvideodto.VideoListRes{List: list}, nil
}

// AddVideosFromStoredFiles 将已上传的对象路径批量写入主播视频库;超上限整批失败
func AddVideosFromStoredFiles(ctx context.Context, storedNames []string) (*anchorvideodto.UploadVideoRes, error) {
	names := make([]string, 0, len(storedNames))
	for _, name := range storedNames {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		names = append(names, name)
	}
	if len(names) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	info := userinfodao.GetUserInfoByUserId(userId)
	if info == nil || !info.IsAnchor() {
		return nil, errercode.CreateCode(errercode.LiveRoomNotAnchor)
	}
	remain := anchorvideodao.MaxVideos() - anchorvideodao.CountVideos(userId)
	if remain <= 0 || len(names) > remain {
		return nil, errercode.CreateCode(errercode.AnchorVideoLimit)
	}
	list := make([]*anchorvideodto.VideoItem, 0, len(names))
	for _, name := range names {
		row := anchorvideodao.TryAddVideo(userId, name)
		if row == nil {
			for _, item := range list {
				if item != nil && item.Id > 0 {
					anchorvideodao.SoftDeleteVideos(userId, []uint64{item.Id})
				}
			}
			return nil, errercode.CreateCode(errercode.AnchorVideoLimit)
		}
		list = append(list, &anchorvideodto.VideoItem{
			Id:    row.ID,
			Video: upload.GetUrlByName(row.Video),
		})
	}
	return &anchorvideodto.UploadVideoRes{List: list}, nil
}

// DeleteVideo 批量删除当前登录主播自己的相册视频
func DeleteVideo(ctx context.Context, req *anchorvideodto.DeleteVideoReq) (*anchorvideodto.DeleteVideoRes, error) {
	if req == nil || len(req.VideoIds) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	rows := anchorvideodao.SoftDeleteVideos(userId, req.VideoIds)
	if len(rows) == 0 {
		return nil, errercode.CreateCode(errercode.AnchorVideoNonExist)
	}
	deletedIds := make([]uint64, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		upload.DeleteUploadedFile(row.Video)
		deletedIds = append(deletedIds, row.ID)
	}
	return &anchorvideodto.DeleteVideoRes{
		Success:    true,
		DeletedIds: deletedIds,
	}, nil
}
