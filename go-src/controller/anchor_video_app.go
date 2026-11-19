package controller

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/anchorvideodto"
	"xr-game-server/dto/anchorvideopricetierdto"
	"xr-game-server/errercode"
	"xr-game-server/module/anchorvideo"
	"xr-game-server/module/anchorvideopricetier"
	"xr-game-server/module/livecfg"
	"xr-game-server/module/upload"
)

const (
	AnchorVideoAppUrl = "/anchorVideo"
)

type AnchorVideoAppController struct{}

type videoPartUploader func(ctx context.Context, part *multipart.Part, maxBytes int64) (string, error)

type uploadAnchorVideoForm struct {
	storedName  string
	price       float64
	description string
	cover       string
	fileType    string
	duration    uint32
	hasPrice    bool
}

func initAnchorVideoAppController() {
	httpserver.RegAPIHandler(AnchorVideoAppUrl, "/upload", handleUploadAnchorVideo)
	httpserver.RegAPI(AnchorVideoAppUrl, &AnchorVideoAppController{})
}

// HandleUploadAnchorVideo Sofie 等新前缀复用同一套上传解析
func HandleUploadAnchorVideo(r *ghttp.Request) {
	handleUploadAnchorVideo(r)
}

func handleUploadAnchorVideo(r *ghttp.Request) {
	form, err := parseUploadAnchorVideoMultipart(r.Context(), r)
	if err != nil {
		r.SetError(err)
		return
	}
	res, err := anchorvideo.AddVideoFromStoredFile(r.Context(), form.storedName, form.price, form.description, form.cover, form.fileType, form.duration)
	if err != nil {
		upload.DeleteUploadedFile(form.storedName)
		r.SetError(err)
		return
	}
	httpserver.SetHandlerResponseData(r, res)
}

func parseUploadAnchorVideoMultipart(ctx context.Context, r *ghttp.Request) (*uploadAnchorVideoForm, error) {
	return parseUploadAnchorVideoMultipartWithUploader(ctx, r, upload.StreamUploadShortVideoPartContext)
}

func parseUploadAnchorVideoMultipartWithUploader(ctx context.Context, r *ghttp.Request, uploader videoPartUploader) (form *uploadAnchorVideoForm, err error) {
	if r == nil || r.Request == nil || uploader == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	reader, err := r.Request.MultipartReader()
	if err != nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	form = &uploadAnchorVideoForm{}
	defer func() {
		if err != nil && form != nil && form.storedName != "" {
			upload.DeleteUploadedFile(form.storedName)
			form.storedName = ""
		}
	}()

	for {
		part, nextErr := reader.NextPart()
		if nextErr == io.EOF {
			break
		}
		if nextErr != nil {
			return form, mapAnchorVideoMultipartErr(nextErr)
		}
		name := part.FormName()
		switch name {
		case "file":
			if strings.TrimSpace(part.FileName()) == "" {
				part.Close()
				return form, errercode.CreateCode(errercode.InvalidParam)
			}
			if form.storedName != "" {
				part.Close()
				return form, errercode.CreateCode(errercode.InvalidParam) // 仅允许单个视频
			}
			storedName, upErr := uploader(ctx, part, livecfg.GetMaxAnchorVideoFileSizeBytes())
			part.Close()
			if upErr != nil {
				return form, mapAnchorVideoMultipartErr(upErr)
			}
			form.storedName = storedName
		case "price":
			raw, readErr := io.ReadAll(io.LimitReader(part, 64))
			part.Close()
			if readErr != nil {
				return form, errercode.CreateCode(errercode.InvalidParam)
			}
			price, parseErr := strconv.ParseFloat(strings.TrimSpace(string(raw)), 64)
			if parseErr != nil {
				return form, errercode.CreateCode(errercode.InvalidParam)
			}
			form.price = price
			form.hasPrice = true
		case "description":
			raw, readErr := io.ReadAll(io.LimitReader(part, 2048))
			part.Close()
			if readErr != nil {
				return form, errercode.CreateCode(errercode.InvalidParam)
			}
			form.description = string(raw)
		case "cover":
			raw, readErr := readAnchorVideoMultipartRaw(part, 4096)
			part.Close()
			if readErr != nil {
				return form, errercode.CreateCode(errercode.InvalidParam)
			}
			form.cover = raw
		case "fileType":
			raw, readErr := readAnchorVideoMultipartRaw(part, 256)
			part.Close()
			if readErr != nil {
				return form, errercode.CreateCode(errercode.InvalidParam)
			}
			form.fileType = raw
		case "duration":
			raw, readErr := readAnchorVideoMultipartRaw(part, 64)
			part.Close()
			if readErr != nil {
				return form, errercode.CreateCode(errercode.InvalidParam)
			}
			form.duration = parseAnchorVideoDurationRaw(raw)
		default:
			isUnexpectedFile := strings.TrimSpace(part.FileName()) != ""
			part.Close()
			if isUnexpectedFile {
				return form, errercode.CreateCode(errercode.InvalidParam)
			}
		}
	}
	if form.storedName == "" || !form.hasPrice {
		return form, errercode.CreateCode(errercode.InvalidParam)
	}
	return form, nil
}

func readAnchorVideoMultipartRaw(part *multipart.Part, max int64) (string, error) {
	if part == nil {
		return "", errercode.CreateCode(errercode.InvalidParam)
	}
	raw, err := io.ReadAll(io.LimitReader(part, max))
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

// parseAnchorVideoDurationRaw 不作业务校验:能解析为非负整数则入库,否则为 0.
func parseAnchorVideoDurationRaw(raw string) uint32 {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0
	}
	v, err := strconv.ParseUint(raw, 10, 32)
	if err != nil {
		return 0
	}
	return uint32(v)
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

func (c *AnchorVideoAppController) Update(ctx context.Context, req *anchorvideodto.UpdateVideoReq) (*anchorvideodto.UpdateVideoRes, error) {
	return anchorvideo.UpdateVideo(ctx, req)
}

func (c *AnchorVideoAppController) Delete(ctx context.Context, req *anchorvideodto.DeleteVideoReq) (*anchorvideodto.DeleteVideoRes, error) {
	return anchorvideo.DeleteVideo(ctx, req)
}

// Unlock 付费解锁主播相册视频(已购/本人不重复扣费)
func (c *AnchorVideoAppController) Unlock(ctx context.Context, req *anchorvideodto.UnlockVideoReq) (*anchorvideodto.UnlockVideoRes, error) {
	return anchorvideo.UnlockVideo(ctx, req)
}

// PriceTierList App 查询主播相册短视频价格挡位(仅内存缓存,不查库)
func (c *AnchorVideoAppController) PriceTierList(ctx context.Context, req *anchorvideopricetierdto.AppPriceTierListReq) (*anchorvideopricetierdto.AppPriceTierListRes, error) {
	return anchorvideopricetier.GetAppPriceTierList(ctx, req)
}
