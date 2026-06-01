package shortvideo

import (
	"context"
	"strings"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

// UploadShortVideoCMS CMS后台上传短视频,按数据库配置校验文件大小
func UploadShortVideoCMS(_ context.Context, req *shortvideodto.UploadShortVideoCMSReq) (*shortvideodto.UploadShortVideoCMSRes, error) {
	if req.File == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	maxSize := getShortVideoMaxFileSize()
	if req.File.Size > int64(maxSize) {
		return nil, errercode.CreateCode(errercode.ShortVideoFileTooLarge)
	}

	name, err := upload.UploadShortVideoFile(req.File, maxSize)
	if err != nil {
		if strings.Contains(err.Error(), "file too large") {
			return nil, errercode.CreateCode(errercode.ShortVideoFileTooLarge)
		}
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	return &shortvideodto.UploadShortVideoCMSRes{
		FileName: name,
		Url:      upload.GetUrlByName(name),
	}, nil
}
