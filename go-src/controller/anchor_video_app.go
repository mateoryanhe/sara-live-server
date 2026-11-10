package controller

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/anchorvideodto"
	"xr-game-server/errercode"
	"xr-game-server/module/anchorvideo"
	"xr-game-server/module/upload"
)

const (
	AnchorVideoAppUrl          = "/anchorVideo"
	defaultAnchorVideoMaxBytes = 100 * 1024 * 1024 // 100MB
)

type AnchorVideoAppController struct{}

type videoPartUploader func(ctx context.Context, part *multipart.Part, maxBytes int64) (string, error)

func initAnchorVideoAppController() {
	httpserver.RegAPIHandler(AnchorVideoAppUrl, "/upload", handleUploadAnchorVideo)
	httpserver.RegAPI(AnchorVideoAppUrl, &AnchorVideoAppController{})
}

// HandleUploadAnchorVideo Sofie 等新前缀复用同一套上传解析
func HandleUploadAnchorVideo(r *ghttp.Request) {
	handleUploadAnchorVideo(r)
}

func handleUploadAnchorVideo(r *ghttp.Request) {
	storedNames, err := parseUploadAnchorVideoMultipart(r.Context(), r)
	if err != nil {
		r.SetError(err)
		return
	}
	res, err := anchorvideo.AddVideosFromStoredFiles(r.Context(), storedNames)
	if err != nil {
		for _, name := range storedNames {
			upload.DeleteUploadedFile(name)
		}
		r.SetError(err)
		return
	}
	httpserver.SetHandlerResponseData(r, res)
}

func parseUploadAnchorVideoMultipart(ctx context.Context, r *ghttp.Request) ([]string, error) {
	return parseUploadAnchorVideoMultipartWithUploader(ctx, r, upload.StreamUploadShortVideoPartContext)
}

func parseUploadAnchorVideoMultipartWithUploader(ctx context.Context, r *ghttp.Request, uploader videoPartUploader) (storedNames []string, err error) {
	if r == nil || r.Request == nil || uploader == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	reader, err := r.Request.MultipartReader()
	if err != nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	defer func() {
		if err != nil {
			for _, name := range storedNames {
				upload.DeleteUploadedFile(name)
			}
			storedNames = nil
		}
	}()

	for {
		part, nextErr := reader.NextPart()
		if nextErr == io.EOF {
			break
		}
		if nextErr != nil {
			return storedNames, mapAnchorVideoMultipartErr(nextErr)
		}
		if part.FormName() != "file" {
			isUnexpectedFile := strings.TrimSpace(part.FileName()) != ""
			part.Close()
			if isUnexpectedFile {
				return storedNames, errercode.CreateCode(errercode.InvalidParam)
			}
			continue
		}
		if strings.TrimSpace(part.FileName()) == "" {
			part.Close()
			return storedNames, errercode.CreateCode(errercode.InvalidParam)
		}
		storedName, upErr := uploader(ctx, part, defaultAnchorVideoMaxBytes)
		part.Close()
		if upErr != nil {
			return storedNames, mapAnchorVideoMultipartErr(upErr)
		}
		storedNames = append(storedNames, storedName)
	}
	if len(storedNames) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	return storedNames, nil
}

func mapAnchorVideoMultipartErr(err error) error {
	if err == nil {
		return nil
	}
	var bizErr *errercode.XError
	if errors.As(err, &bizErr) {
		return err
	}
	if upload.IsUploadVideoFileTooLarge(err) {
		return errercode.CreateCode(errercode.ShortVideoFileTooLarge)
	}
	return errercode.CreateCode(errercode.InvalidParam)
}

func (c *AnchorVideoAppController) List(ctx context.Context, req *anchorvideodto.VideoListReq) (*anchorvideodto.VideoListRes, error) {
	return anchorvideo.GetVideoList(ctx, req)
}

func (c *AnchorVideoAppController) Delete(ctx context.Context, req *anchorvideodto.DeleteVideoReq) (*anchorvideodto.DeleteVideoRes, error) {
	return anchorvideo.DeleteVideo(ctx, req)
}
