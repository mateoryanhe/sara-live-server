package shortvideo

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"strconv"
	"strings"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/snowflake"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

func PublishShortVideoApp(ctx context.Context, req *shortvideodto.AppPublishShortVideoReq) (*shortvideodto.AppPublishShortVideoRes, error) {
	authorId := httpserver.GetAuthId(ctx)
	if authorId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if err := validateShortVideoAuthorId(authorId); err != nil {
		return nil, err
	}
	if existing := shortvideodao.GetByTitle(req.Title); existing != nil {
		return nil, errercode.CreateCode(errercode.ShortVideoExist)
	}
	isPaid, diamondPerMinute := entity.ShortVideoPaidNo, float64(0)

	if err := validateShortVideoCategoryId(req.CategoryId); err != nil {
		return nil, err
	}
	videoName, err := uploadShortVideoFile(req.File)
	if err != nil {
		return nil, err
	}
	coverName, err := uploadShortVideoCoverFile(ctx, req.Cover)
	if err != nil {
		upload.DeleteUploadedFile(videoName)
		return nil, err
	}
	row := entity.NewShortVideo(
		snowflake.GetId(),
		req.Title,
		videoName,
		coverName,
		0,
		isPaid,
		diamondPerMinute,
		req.CategoryId,
		req.Source,
		authorId,
	)
	shortvideodao.AddShortVideoToCache(row)
	loadAppShortVideoListCache()
	res := &shortvideodto.AppPublishShortVideoRes{
		ID:    strconv.FormatUint(row.ID, 10),
		Video: upload.GetUrlByName(videoName),
	}
	if coverName != "" {
		res.Cover = upload.GetUrlByName(coverName)
	}
	return res, nil
}

func validateShortVideoAuthorId(authorId uint64) error {
	if authorId == 0 {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	if userinfodao.GetUserInfoByUserId(authorId) == nil {
		return errercode.CreateCode(errercode.SysError)
	}
	return nil
}

func uploadShortVideoFile(file *ghttp.UploadFile) (string, error) {
	if file == nil {
		return "", errercode.CreateCode(errercode.InvalidParam)
	}
	maxSize := getShortVideoMaxFileSize()
	if file.Size > int64(maxSize) {
		return "", errercode.CreateCode(errercode.ShortVideoFileTooLarge)
	}
	name, err := upload.UploadShortVideoFile(file, maxSize)
	if err != nil {
		if strings.Contains(err.Error(), "file too large") {
			return "", errercode.CreateCode(errercode.ShortVideoFileTooLarge)
		}
		return "", errercode.CreateCode(errercode.InvalidParam)
	}
	return name, nil
}

func uploadShortVideoCoverFile(ctx context.Context, file *ghttp.UploadFile) (string, error) {
	if file == nil {
		return "", nil
	}
	maxSize := int64(getShortVideoMaxCoverFileSize()) * 1024 * 1024
	if file.Size > maxSize {
		return "", errercode.CreateCode(errercode.ShortVideoFileTooLarge)
	}
	name, err := upload.UploadImageForApp(ctx, file)
	if err != nil {
		return "", err
	}
	return name, nil
}
