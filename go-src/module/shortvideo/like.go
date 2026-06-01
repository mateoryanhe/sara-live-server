package shortvideo

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

func LikeShortVideo(ctx context.Context, req *shortvideodto.LikeShortVideoReq) (*shortvideodto.LikeShortVideoRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	row := shortvideodao.GetShortVideoById(req.VideoId)
	if row == nil || row.Status != entity.ShortVideoStatusOnShelf {
		return nil, errercode.CreateCode(errercode.ShortVideoNonExist)
	}
	existing := shortvideodao.GetShortVideoLikeByUserVideo(userId, req.VideoId)
	if existing != nil && existing.Status == entity.ShortVideoLikeStatusLiked {
		return nil, errercode.CreateCode(errercode.ShortVideoAlreadyLiked)
	}
	if existing == nil {
		like := entity.NewShortVideoLike(userId, req.VideoId)
		shortvideodao.AddLikeToCache(like)
	} else {
		existing.SetStatus(entity.ShortVideoLikeStatusLiked)
		shortvideodao.AddLikeToCache(existing)
	}
	stat := shortvideodao.GetStatByVideoId(req.VideoId)
	if stat == nil {
		return nil, errercode.CreateCode(errercode.ShortVideoNonExist)
	}
	stat.AddLikeCount(1)
	return &shortvideodto.LikeShortVideoRes{LikeCount: stat.LikeCount}, nil
}
