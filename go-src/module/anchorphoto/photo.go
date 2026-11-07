package anchorphoto

import (
	"context"
	"strings"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/anchorphotodao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/anchorphotodto"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

// GetPhotoList App 查询主播相册(读缓存)
func GetPhotoList(_ context.Context, req *anchorphotodto.PhotoListReq) (*anchorphotodto.PhotoListRes, error) {
	if req == nil || req.AnchorId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	rows := anchorphotodao.GetPhotoList(req.AnchorId)
	list := make([]*anchorphotodto.PhotoItem, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.ID == 0 || row.IsDeleted {
			continue
		}
		list = append(list, &anchorphotodto.PhotoItem{
			Id:    row.ID,
			Photo: upload.GetUrlByName(row.Photo),
		})
	}
	return &anchorphotodto.PhotoListRes{List: list}, nil
}

// AddPhotosFromStoredFiles 将已上传的对象路径批量写入主播相册;超上限整批失败
func AddPhotosFromStoredFiles(ctx context.Context, storedNames []string) (*anchorphotodto.UploadPhotoRes, error) {
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
	remain := anchorphotodao.MaxPhotos() - anchorphotodao.CountPhotos(userId)
	if remain <= 0 || len(names) > remain {
		return nil, errercode.CreateCode(errercode.AnchorPhotoLimit)
	}
	list := make([]*anchorphotodto.PhotoItem, 0, len(names))
	for _, name := range names {
		row := anchorphotodao.TryAddPhoto(userId, name)
		if row == nil {
			// 理论上已前置校验,仍失败则回滚本批已写入
			for _, item := range list {
				if item != nil && item.Id > 0 {
					anchorphotodao.SoftDeletePhoto(userId, item.Id)
				}
			}
			return nil, errercode.CreateCode(errercode.AnchorPhotoLimit)
		}
		list = append(list, &anchorphotodto.PhotoItem{
			Id:    row.ID,
			Photo: upload.GetUrlByName(row.Photo),
		})
	}
	return &anchorphotodto.UploadPhotoRes{List: list}, nil
}

// DeletePhoto 批量删除当前登录主播自己的相册照片
func DeletePhoto(ctx context.Context, req *anchorphotodto.DeletePhotoReq) (*anchorphotodto.DeletePhotoRes, error) {
	if req == nil || len(req.PhotoIds) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	rows := anchorphotodao.SoftDeletePhotos(userId, req.PhotoIds)
	if len(rows) == 0 {
		return nil, errercode.CreateCode(errercode.AnchorPhotoNonExist)
	}
	deletedIds := make([]uint64, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		upload.DeleteUploadedFile(row.Photo)
		deletedIds = append(deletedIds, row.ID)
	}
	return &anchorphotodto.DeletePhotoRes{
		Success:    true,
		DeletedIds: deletedIds,
	}, nil
}
