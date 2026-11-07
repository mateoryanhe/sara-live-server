package controller

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/anchorphotodto"
	"xr-game-server/errercode"
	"xr-game-server/module/anchorphoto"
	"xr-game-server/module/upload"
)

const AnchorPhotoAppUrl = "/anchorPhoto"

type AnchorPhotoAppController struct{}

func initAnchorPhotoAppController() {
	httpserver.RegAPIHandler(AnchorPhotoAppUrl, "/upload", handleUploadAnchorPhoto)
	httpserver.RegAPI(AnchorPhotoAppUrl, &AnchorPhotoAppController{})
}

// HandleUploadAnchorPhoto Sofie 等新前缀复用同一套上传解析
func HandleUploadAnchorPhoto(r *ghttp.Request) {
	handleUploadAnchorPhoto(r)
}

func handleUploadAnchorPhoto(r *ghttp.Request) {
	storedNames, err := parseUploadAnchorPhotoMultipart(r.Context(), r)
	if err != nil {
		r.SetError(err)
		return
	}
	res, err := anchorphoto.AddPhotosFromStoredFiles(r.Context(), storedNames)
	if err != nil {
		for _, name := range storedNames {
			upload.DeleteUploadedFile(name)
		}
		r.SetError(err)
		return
	}
	httpserver.SetHandlerResponseData(r, res)
}

func parseUploadAnchorPhotoMultipart(ctx context.Context, r *ghttp.Request) ([]string, error) {
	return parseUploadAnchorPhotoMultipartWithUploader(ctx, r, upload.UploadImagePartForApp)
}

func parseUploadAnchorPhotoMultipartWithUploader(ctx context.Context, r *ghttp.Request, uploader avatarPartUploader) (storedNames []string, err error) {
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
			return storedNames, mapAnchorPhotoMultipartErr(nextErr)
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
		storedName, upErr := uploader(ctx, part, int64(upload.GetAppImageMaxSize()))
		part.Close()
		if upErr != nil {
			return storedNames, mapAnchorPhotoMultipartErr(upErr)
		}
		storedNames = append(storedNames, storedName)
	}
	if len(storedNames) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	return storedNames, nil
}

func mapAnchorPhotoMultipartErr(err error) error {
	if err == nil {
		return nil
	}
	var bizErr *errercode.XError
	if errors.As(err, &bizErr) {
		return err
	}
	if upload.IsUploadImageFileTooLarge(err) {
		return errercode.CreateCode(errercode.AppImageFileTooLarge)
	}
	return errercode.CreateCode(errercode.InvalidParam)
}

func (c *AnchorPhotoAppController) List(ctx context.Context, req *anchorphotodto.PhotoListReq) (*anchorphotodto.PhotoListRes, error) {
	return anchorphoto.GetPhotoList(ctx, req)
}

func (c *AnchorPhotoAppController) Delete(ctx context.Context, req *anchorphotodto.DeletePhotoReq) (*anchorphotodto.DeletePhotoRes, error) {
	return anchorphoto.DeletePhoto(ctx, req)
}
