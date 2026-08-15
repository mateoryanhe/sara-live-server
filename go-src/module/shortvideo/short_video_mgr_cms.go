package shortvideo

import (
	"context"
	"strconv"
	"strings"

	"xr-game-server/core/httpserver"
	"xr-game-server/core/snowflake"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/entity/shortvideo"
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

// CreateShortVideo CMS上传短视频(作者类型为 CMS)
func CreateShortVideo(ctx context.Context, req *shortvideodto.CreateShortVideoReq) (*shortvideodto.CreateShortVideoRes, error) {
	_ = ctx
	if existing := shortvideodao.GetByTitle(req.Title); existing != nil {
		return nil, errercode.CreateCode(errercode.ShortVideoExist)
	}
	isPaid, payDiamond, err := normalizeShortVideoPaid(req.IsPaid, req.PayDiamond)
	if err != nil {
		return nil, err
	}
	freeWatchSeconds := normalizeShortVideoFreeWatchSeconds(isPaid, req.FreeWatchSeconds)
	if err := validateShortVideoCategoryId(req.CategoryId); err != nil {
		return nil, err
	}
	if err := validateShortVideoDuration(req.Duration); err != nil {
		return nil, err
	}
	videoName := strings.TrimSpace(req.Video)
	if videoName == "" {
		if req.File == nil {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		var err error
		videoName, err = uploadShortVideoFile(req.File)
		if err != nil {
			return nil, err
		}
	}
	coverName := strings.TrimSpace(req.CoverName)
	if coverName == "" && req.Cover != nil {
		var err error
		coverName, err = uploadShortVideoCoverFile(ctx, req.Cover)
		if err != nil {
			upload.DeleteUploadedFile(videoName)
			return nil, err
		}
	}
	authorId, err := createCMSAuthorUser(resolveCMSAuthorNickname(req.AuthorNickname), coverName)
	if err != nil {
		upload.DeleteUploadedFile(videoName)
		upload.DeleteUploadedFile(coverName)
		return nil, errercode.CreateCode(errercode.SysError)
	}
	row := entity.NewShortVideo(
		snowflake.GetId(),
		req.Title,
		videoName,
		coverName,
		req.Sort,
		isPaid,
		payDiamond,
		req.CategoryId,
		req.Source,
		authorId,
		entity.ShortVideoAuthorTypeCMS,
		req.Duration,
		freeWatchSeconds,
	)
	shortvideodao.AddShortVideoToCache(row)
	reloadAllAppShortVideoListCaches()
	res := &shortvideodto.CreateShortVideoRes{
		ID:       strconv.FormatUint(row.ID, 10),
		Video:    upload.GetUrlByName(videoName),
		AuthorId: strconv.FormatUint(authorId, 10),
	}
	if coverName != "" {
		res.Cover = upload.GetUrlByName(coverName)
	}
	return res, nil
}

func UpdateShortVideo(_ context.Context, req *shortvideodto.UpdateShortVideoReq) (*shortvideodto.UpdateShortVideoRes, error) {
	// authorType 仅在创建(App/CMS上传)时写入,此处不可修改
	row := shortvideodao.GetShortVideoById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.ShortVideoNonExist)
	}
	if row.Status == entity.ShortVideoStatusOnShelf {
		return nil, errercode.CreateCode(errercode.ShortVideoOnShelfCannotUpdate)
	}
	if existing := shortvideodao.GetByTitle(req.Title); existing != nil && existing.ID != req.ID {
		return nil, errercode.CreateCode(errercode.ShortVideoExist)
	}
	isPaid, payDiamond, err := normalizeShortVideoPaid(req.IsPaid, req.PayDiamond)
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
	row.SetPayDiamond(payDiamond)
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
	if err := removeShortVideoRow(row); err != nil {
		return nil, err
	}
	return &shortvideodto.DeleteShortVideoRes{Success: true}, nil
}

func removeShortVideoRow(row *entity.ShortVideo) error {
	if row == nil {
		return errercode.CreateCode(errercode.ShortVideoNonExist)
	}
	videoName := row.Video
	coverName := row.Cover
	if err := shortvideodao.Delete(row.ID); err != nil {
		return err
	}
	if err := shortvideodao.DeleteWatchByVideoId(row.ID); err != nil {
		return err
	}
	//先刷新列表
	shortvideodao.RefreshAuthorStatListCacheOnVideoDelete(row.AuthorId, row.ID)
	//再移除单条缓存
	shortvideodao.RemoveStatCacheByVideoId(row.ID)
	upload.DeleteUploadedFile(videoName)
	upload.DeleteUploadedFile(coverName)
	reloadAllAppShortVideoListCaches()
	return nil
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
		shortvideodao.PrependStatToAuthorListCache(row.AuthorId, shortvideodao.GetStatByVideoId(req.ID))
		reloadAllAppShortVideoListCaches()
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
		reloadAllAppShortVideoListCaches()
	}
	return &shortvideodto.OffShelfShortVideoRes{Success: true, Status: entity.ShortVideoStatusOffShelf}, nil
}

func normalizeShortVideoPaid(isPaid uint8, payDiamond float64) (uint8, float64, error) {
	if isPaid != entity.ShortVideoPaidYes {
		return entity.ShortVideoPaidNo, 0, nil
	}
	if payDiamond <= 0 {
		return 0, 0, errercode.CreateCode(errercode.InvalidParam)
	}
	return isPaid, payDiamond, nil
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
