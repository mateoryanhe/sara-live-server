package shortvideo

import (
	"context"
	"strconv"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dao/shortvideolikedao"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

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
	row := &entity.ShortVideo{
		Title:       req.Title,
		Video:       req.Video,
		Cover:       req.Cover,
		Sort:        req.Sort,
		Status:      entity.ShortVideoStatusOffShelf,
		Description: req.Description,
	}
	if err := shortvideodao.Create(row); err != nil {
		return nil, err
	}
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
	row.Title = req.Title
	row.Video = req.Video
	row.Cover = req.Cover
	row.Sort = req.Sort
	row.Description = req.Description
	if err := shortvideodao.Update(row); err != nil {
		return nil, err
	}
	return &shortvideodto.UpdateShortVideoRes{Success: true}, nil
}

func DeleteShortVideo(_ context.Context, req *shortvideodto.DeleteShortVideoReq) (*shortvideodto.DeleteShortVideoRes, error) {
	if row := shortvideodao.GetShortVideoById(req.ID); row == nil {
		return nil, errercode.CreateCode(errercode.ShortVideoNonExist)
	}
	if err := shortvideodao.Delete(req.ID); err != nil {
		return nil, err
	}
	return &shortvideodto.DeleteShortVideoRes{Success: true}, nil
}

func OnShelfShortVideo(_ context.Context, req *shortvideodto.OnShelfShortVideoReq) (*shortvideodto.OnShelfShortVideoRes, error) {
	row := shortvideodao.GetShortVideoById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.ShortVideoNonExist)
	}
	if row.Status != entity.ShortVideoStatusOnShelf {
		if err := shortvideodao.UpdateStatus(req.ID, entity.ShortVideoStatusOnShelf); err != nil {
			return nil, err
		}
	}
	return &shortvideodto.OnShelfShortVideoRes{Success: true, Status: entity.ShortVideoStatusOnShelf}, nil
}

func OffShelfShortVideo(_ context.Context, req *shortvideodto.OffShelfShortVideoReq) (*shortvideodto.OffShelfShortVideoRes, error) {
	row := shortvideodao.GetShortVideoById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.ShortVideoNonExist)
	}
	if row.Status != entity.ShortVideoStatusOffShelf {
		if err := shortvideodao.UpdateStatus(req.ID, entity.ShortVideoStatusOffShelf); err != nil {
			return nil, err
		}
	}
	return &shortvideodto.OffShelfShortVideoRes{Success: true, Status: entity.ShortVideoStatusOffShelf}, nil
}

func LikeShortVideo(ctx context.Context, req *shortvideodto.LikeShortVideoReq) (*shortvideodto.LikeShortVideoRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	row := shortvideodao.GetShortVideoById(req.VideoId)
	if row == nil || row.Status != entity.ShortVideoStatusOnShelf {
		return nil, errercode.CreateCode(errercode.ShortVideoNonExist)
	}
	existing := shortvideolikedao.GetByUserVideo(userId, req.VideoId)
	if existing != nil && existing.Status == entity.ShortVideoLikeStatusLiked {
		return nil, errercode.CreateCode(errercode.ShortVideoAlreadyLiked)
	}
	if existing == nil {
		like := entity.NewShortVideoLike(userId, req.VideoId)
		shortvideolikedao.AddLikeToCache(like)
	} else {
		existing.SetStatus(entity.ShortVideoLikeStatusLiked)
		shortvideolikedao.AddLikeToCache(existing)
	}
	row.AddLikeCount(1)
	return &shortvideodto.LikeShortVideoRes{LikeCount: row.LikeCount}, nil
}
