package controller

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/httpserver"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
	"xr-game-server/module/userinfo"
)

type avatarPartUploader func(context.Context, *multipart.Part, int64) (string, error)

// handleUploadAvatar 由控制器逐段解析 multipart，图片流直接交给上传模块。
func handleUploadAvatar(r *ghttp.Request) {
	storedName, err := parseUploadAvatarMultipart(r.Context(), r)
	if err != nil {
		r.SetError(err)
		return
	}
	res, err := userinfo.UpdateAvatarFromStoredFile(r.Context(), storedName)
	if err != nil {
		upload.DeleteUploadedFile(storedName)
		r.SetError(err)
		return
	}
	httpserver.SetHandlerResponseData(r, res)
}

func parseUploadAvatarMultipart(ctx context.Context, r *ghttp.Request) (string, error) {
	return parseUploadAvatarMultipartWithUploader(ctx, r, upload.UploadImagePartForApp)
}

func parseUploadAvatarMultipartWithUploader(ctx context.Context, r *ghttp.Request, uploader avatarPartUploader) (storedName string, err error) {
	if r == nil || r.Request == nil || uploader == nil {
		return "", errercode.CreateCode(errercode.InvalidParam)
	}
	reader, err := r.Request.MultipartReader()
	if err != nil {
		return "", errercode.CreateCode(errercode.InvalidParam)
	}
	defer func() {
		if err != nil && storedName != "" {
			upload.DeleteUploadedFile(storedName)
			storedName = ""
		}
	}()

	for {
		part, nextErr := reader.NextPart()
		if nextErr == io.EOF {
			break
		}
		if nextErr != nil {
			return storedName, mapAvatarMultipartErr(nextErr)
		}

		if part.FormName() != "file" {
			isUnexpectedFile := strings.TrimSpace(part.FileName()) != ""
			part.Close()
			if isUnexpectedFile {
				return storedName, errercode.CreateCode(errercode.InvalidParam)
			}
			continue
		}
		if storedName != "" || strings.TrimSpace(part.FileName()) == "" {
			part.Close()
			return storedName, errercode.CreateCode(errercode.InvalidParam)
		}
		storedName, err = uploader(ctx, part, int64(upload.GetAppImageMaxSize()))
		part.Close()
		if err != nil {
			return storedName, mapAvatarMultipartErr(err)
		}
	}
	if storedName == "" {
		return "", errercode.CreateCode(errercode.InvalidParam)
	}
	return storedName, nil
}

func mapAvatarMultipartErr(err error) error {
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
