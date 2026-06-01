package shortvideo

import (
	"context"
	"strconv"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/snowflake"
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

func CreateShortVideo(_ context.Context, req *shortvideodto.CreateShortVideoReq) (*shortvideodto.CreateShortVideoRes, error) {
	if existing := shortvideodao.GetByTitle(req.Title); existing != nil {
		return nil, errercode.CreateCode(errercode.ShortVideoExist)
	}
	isPaid, diamondPerSecond, err := normalizeShortVideoPaid(req.IsPaid, req.DiamondPerSecond)
	if err != nil {
		return nil, err
	}
	row := &entity.ShortVideo{
		Title:            req.Title,
		Video:            req.Video,
		Cover:            req.Cover,
		Sort:             req.Sort,
		Status:           entity.ShortVideoStatusOffShelf,
		IsPaid:           isPaid,
		DiamondPerSecond: diamondPerSecond,
		Description:      req.Description,
	}
	row.ID = snowflake.GetId()
	if err := shortvideodao.Create(row); err != nil {
		return nil, err
	}
	loadAppShortVideoListCache()

	return &shortvideodto.CreateShortVideoRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func UpdateShortVideo(_ context.Context, req *shortvideodto.UpdateShortVideoReq) (*shortvideodto.UpdateShortVideoRes, error) {
	row := shortvideodao.GetShortVideoById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.ShortVideoNonExist)
	}
	if existing := shortvideodao.GetByTitle(req.Title); existing != nil && existing.ID != req.ID {
		return nil, errercode.CreateCode(errercode.ShortVideoExist)
	}
	isPaid, diamondPerSecond, err := normalizeShortVideoPaid(req.IsPaid, req.DiamondPerSecond)
	if err != nil {
		return nil, err
	}
	row.Title = req.Title
	row.Video = req.Video
	row.Cover = req.Cover
	row.Sort = req.Sort
	row.IsPaid = isPaid
	row.DiamondPerSecond = diamondPerSecond
	row.Description = req.Description
	if err := shortvideodao.Update(row); err != nil {
		return nil, err
	}
	loadAppShortVideoListCache()
	return &shortvideodto.UpdateShortVideoRes{Success: true}, nil
}

func DeleteShortVideo(_ context.Context, req *shortvideodto.DeleteShortVideoReq) (*shortvideodto.DeleteShortVideoRes, error) {
	if row := shortvideodao.GetShortVideoById(req.ID); row == nil {
		return nil, errercode.CreateCode(errercode.ShortVideoNonExist)
	}
	if err := shortvideodao.Delete(req.ID); err != nil {
		return nil, err
	}
	_ = shortvideodao.DeleteByVideoId(req.ID)
	loadAppShortVideoListCache()
	return &shortvideodto.DeleteShortVideoRes{Success: true}, nil
}

func OnShelfShortVideo(_ context.Context, req *shortvideodto.OnShelfShortVideoReq) (*shortvideodto.OnShelfShortVideoRes, error) {
	row := shortvideodao.GetShortVideoById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.ShortVideoNonExist)
	}
	if row.Status != entity.ShortVideoStatusOnShelf {
		shortvideodao.UpdateStatus(req.ID, entity.ShortVideoStatusOnShelf)
		row.Status = entity.ShortVideoStatusOnShelf
		loadAppShortVideoListCache()
		shortvideodao.FlushShortVideo(row)
		return nil, nil
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
		row.Status = entity.ShortVideoStatusOffShelf
		loadAppShortVideoListCache()
		shortvideodao.FlushShortVideo(row)
		return nil, nil
	}
	return &shortvideodto.OffShelfShortVideoRes{Success: true, Status: entity.ShortVideoStatusOffShelf}, nil
}

func normalizeShortVideoPaid(isPaid uint8, diamondPerSecond uint64) (uint8, uint64, error) {
	if isPaid != entity.ShortVideoPaidYes {
		return entity.ShortVideoPaidNo, 0, nil
	}
	if diamondPerSecond == 0 {
		return 0, 0, errercode.CreateCode(errercode.InvalidParam)
	}
	return isPaid, diamondPerSecond, nil
}
