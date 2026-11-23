package controller

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

type appImagePartUploader func(context.Context, *multipart.Part, int64) (string, error)

// parseAppImageMultipartFile 解析 multipart 中单张图片 file 字段并上传(App 审核+大小限制).
func parseAppImageMultipartFile(ctx context.Context, r *ghttp.Request) (string, error) {
	return parseAppImageMultipartFileWithUploader(ctx, r, upload.UploadImagePartForApp)
}

func parseAppImageMultipartFileWithUploader(ctx context.Context, r *ghttp.Request, uploader appImagePartUploader) (storedName string, err error) {
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
			return storedName, mapAppImageMultipartErr(nextErr)
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
			return storedName, mapAppImageMultipartErr(err)
		}
	}
	if storedName == "" {
		return "", errercode.CreateCode(errercode.InvalidParam)
	}
	return storedName, nil
}

func mapAppImageMultipartErr(err error) error {
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