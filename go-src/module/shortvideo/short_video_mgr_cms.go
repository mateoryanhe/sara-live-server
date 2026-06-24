package shortvideo

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

// 后台对短视频的管理
func GetShortVideoList(_ context.Context, req *shortvideodto.ShortVideoListReq) (*httpserver.CMSQueryResp, error) {
	total, list := shortvideodao.GetShortVideoList(req)
	for _, row := range list {
		row.VideoName = row.Video
		row.CoverName = row.Cover
		row.Video = upload.GetUrlByName(row.VideoName)
		row.Cover = upload.GetUrlByName(row.CoverName)
	}
	return &httpserver.CMSQueryResp{Total: total, Data: list}, nil
}

func UpdateShortVideo(_ context.Context, req *shortvideodto.UpdateShortVideoReq) (*shortvideodto.UpdateShortVideoRes, error) {
	row := shortvideodao.GetShortVideoById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.ShortVideoNonExist)
	}
	if existing := shortvideodao.GetByTitle(req.Title); existing != nil && existing.ID != req.ID {
		return nil, errercode.CreateCode(errercode.ShortVideoExist)
	}
	isPaid, diamondPerMinute, err := normalizeShortVideoPaid(req.IsPaid, req.DiamondPerMinute)
	if err != nil {
		return nil, err
	}
	freeWatchSeconds := normalizeShortVideoFreeWatchSeconds(isPaid, req.FreeWatchSeconds)
	if err := validateShortVideoCategoryId(req.CategoryId); err != nil {
		return nil, err
	}
	row.SetTitle(req.Title)
	row.SetCover(req.Cover)
	row.SetSort(req.Sort)
	row.SetIsPaid(isPaid)
	row.SetDiamondPerMinute(diamondPerMinute)
	row.SetCategoryId(req.CategoryId)
	row.SetSource(req.Source)
	row.SetFreeWatchSeconds(freeWatchSeconds)
	shortvideodao.FlushShortVideo(row)
	loadAppShortVideoListCache()
	return &shortvideodto.UpdateShortVideoRes{Success: true}, nil
}

func DeleteShortVideo(_ context.Context, req *shortvideodto.DeleteShortVideoReq) (*shortvideodto.DeleteShortVideoRes, error) {
	row := shortvideodao.GetShortVideoById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.ShortVideoNonExist)
	}
	videoName := row.Video
	coverName := row.Cover
	if err := shortvideodao.Delete(req.ID); err != nil {
		return nil, err
	}
	_ = shortvideodao.DeleteByVideoId(req.ID)
	upload.DeleteUploadedFile(videoName)
	upload.DeleteUploadedFile(coverName)
	loadAppShortVideoListCache()
	return &shortvideodto.DeleteShortVideoRes{Success: true}, nil
}

func OnShelfShortVideo(_ context.Context, req *shortvideodto.OnShelfShortVideoReq) (*shortvideodto.OnShelfShortVideoRes, error) {
	row := shortvideodao.GetShortVideoById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.ShortVideoNonExist)
	}
	if row.Video == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if row.Status != entity.ShortVideoStatusOnShelf {
		shortvideodao.UpdateStatus(req.ID, entity.ShortVideoStatusOnShelf)
		loadAppShortVideoListCache()
	}
	return &shortvideodto.OnShelfShortVideoRes{Success: true, Status: entity.ShortVideoStatusOnShelf}, nil
}

func OffShelfShortVideo(_ context.Context, req *shortvideodto.OffShelfShortVideoReq) (*shortvideodto.OffShelfShortVideoRes, error) {
	row := shortvideodao.GetShortVideoById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.ShortVideoNonExist)
	}
	if row.Status != entity.ShortVideoStatusOffShelf {
		shortvideodao.UpdateStatus(req.ID, entity.ShortVideoStatusOffShelf)
		loadAppShortVideoListCache()
	}
	return &shortvideodto.OffShelfShortVideoRes{Success: true, Status: entity.ShortVideoStatusOffShelf}, nil
}

func normalizeShortVideoPaid(isPaid uint8, diamondPerMinute float64) (uint8, float64, error) {
	if isPaid != entity.ShortVideoPaidYes {
		return entity.ShortVideoPaidNo, 0, nil
	}
	if diamondPerMinute <= 0 {
		return 0, 0, errercode.CreateCode(errercode.InvalidParam)
	}
	return isPaid, diamondPerMinute, nil
}

func normalizeShortVideoFreeWatchSeconds(isPaid uint8, freeWatchSeconds uint32) uint32 {
	if isPaid != entity.ShortVideoPaidYes {
		return 0
	}
	return freeWatchSeconds
}

func validateShortVideoCategoryId(categoryId int) error {
	if categoryId <= 0 {
		return nil
	}
	if GetCategoryFromMemoryById(uint64(categoryId)) == nil {
		return errercode.CreateCode(errercode.ShortVideoCategoryNonExist)
	}
	return nil
}
