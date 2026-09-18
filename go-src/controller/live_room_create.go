package controller

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/xrjson"
	"xr-game-server/core/xrlog"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

const createRoomTextFieldMaxBytes = 64 * 1024

type createRoomImageUploader func(context.Context, *multipart.Part, int64) (string, error)

// parseCreateRoomMultipart 在控制器层逐段解析 HTTP multipart，封面流直接交给上传模块。
func parseCreateRoomMultipart(ctx context.Context, r *ghttp.Request) (*liveroomdto.CreateLiveRoomReq, string, error) {
	return parseCreateRoomMultipartWithUploader(ctx, r, upload.UploadImagePartForApp)
}

func parseCreateRoomMultipartWithUploader(ctx context.Context, r *ghttp.Request, uploader createRoomImageUploader) (req *liveroomdto.CreateLiveRoomReq, coverName string, err error) {
	if r == nil || r.Request == nil || uploader == nil {
		return nil, "", errercode.CreateCode(errercode.InvalidParam)
	}
	reader, err := r.Request.MultipartReader()
	if err != nil {
		return nil, "", mapCreateRoomMultipartErr(err)
	}

	rawFields := make(map[string][]string)
	fields := make(map[string][]string)
	coverFilename := ""
	coverContentType := ""
	defer func() {
		if err != nil && coverName != "" {
			upload.DeleteUploadedFile(coverName)
			coverName = ""
		}
	}()
	defer func() {
		xrlog.DetailLog.Infof(ctx,
			"liveRoom/create raw multipart,reqId=%v,anchorId=%d,contentType=%q,contentLength=%d,fields=%s,coverFilename=%q,coverContentType=%q,coverStoredName=%q,err=%v",
			r.GetHeader(httpserver.ReqId, ""),
			httpserver.GetAuthId(ctx),
			r.Header.Get("Content-Type"),
			r.ContentLength,
			string(xrjson.MustMarshal(rawFields)),
			coverFilename,
			coverContentType,
			coverName,
			err,
		)
	}()

	for {
		part, nextErr := reader.NextPart()
		if nextErr == io.EOF {
			break
		}
		if nextErr != nil {
			return nil, coverName, mapCreateRoomMultipartErr(nextErr)
		}

		formName := part.FormName()
		filename := strings.TrimSpace(part.FileName())
		if filename != "" {
			if formName != "cover" || coverName != "" {
				part.Close()
				return nil, coverName, errercode.CreateCode(errercode.InvalidParam)
			}
			coverFilename = filename
			coverContentType = part.Header.Get("Content-Type")
			coverName, err = uploader(ctx, part, int64(upload.GetAppImageMaxSize()))
			part.Close()
			if err != nil {
				return nil, coverName, mapCreateRoomMultipartErr(err)
			}
			continue
		}

		value, readErr := readCreateRoomTextField(part)
		part.Close()
		if readErr != nil {
			return nil, coverName, mapCreateRoomMultipartErr(readErr)
		}
		rawFields[formName] = append(rawFields[formName], value)
		fieldName := normalizeCreateRoomFieldName(formName)
		fields[fieldName] = append(fields[fieldName], value)
	}

	req, err = bindCreateRoomMultipartFields(fields)
	if err != nil {
		return nil, coverName, err
	}
	return req, coverName, nil
}

func readCreateRoomTextField(part io.Reader) (string, error) {
	body, err := io.ReadAll(io.LimitReader(part, createRoomTextFieldMaxBytes+1))
	if err != nil {
		return "", err
	}
	if len(body) > createRoomTextFieldMaxBytes {
		return "", errercode.CreateCode(errercode.InvalidParam)
	}
	return string(body), nil
}

func normalizeCreateRoomFieldName(name string) string {
	if name == "gameCodes[]" || strings.HasPrefix(name, "gameCodes[") && strings.HasSuffix(name, "]") {
		return "gameCodes"
	}
	return name
}

func bindCreateRoomMultipartFields(fields map[string][]string) (*liveroomdto.CreateLiveRoomReq, error) {
	req := &liveroomdto.CreateLiveRoomReq{
		Title:  createRoomLastField(fields, "title"),
		Notice: createRoomLastField(fields, "notice"),
	}
	var err error
	if req.Category, err = parseCreateRoomUint8(createRoomLastField(fields, "category")); err != nil {
		return nil, err
	}
	if req.TagId, err = parseCreateRoomUint64(createRoomLastField(fields, "tagId")); err != nil {
		return nil, err
	}
	if req.Ticket, err = parseCreateRoomFloat(createRoomLastField(fields, "ticket")); err != nil {
		return nil, err
	}
	if req.Billing, err = parseCreateRoomFloat(createRoomLastField(fields, "billing")); err != nil {
		return nil, err
	}
	if req.PrivateInviteType, err = parseCreateRoomUint8(createRoomLastField(fields, "privateInviteType")); err != nil {
		return nil, err
	}
	if req.PrivateInviteType != 0 && req.PrivateInviteType != 1 && req.PrivateInviteType != 3 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if req.GameCodes, err = parseCreateRoomGameCodes(fields["gameCodes"]); err != nil {
		return nil, err
	}
	return req, nil
}

func createRoomLastField(fields map[string][]string, name string) string {
	values := fields[name]
	if len(values) == 0 {
		return ""
	}
	return values[len(values)-1]
}

func parseCreateRoomUint8(raw string) (uint8, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, nil
	}
	value, err := strconv.ParseUint(raw, 10, 8)
	if err != nil {
		return 0, errercode.CreateCode(errercode.InvalidParam)
	}
	return uint8(value), nil
}

func parseCreateRoomUint64(raw string) (uint64, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, nil
	}
	value, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, errercode.CreateCode(errercode.InvalidParam)
	}
	return value, nil
}

func parseCreateRoomFloat(raw string) (float64, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, nil
	}
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, errercode.CreateCode(errercode.InvalidParam)
	}
	return value, nil
}

func parseCreateRoomGameCodes(rawValues []string) ([]string, error) {
	result := make([]string, 0, len(rawValues))
	for _, raw := range rawValues {
		value := strings.TrimSpace(raw)
		if value == "" || value == "[]" {
			continue
		}
		if strings.HasPrefix(value, "[") {
			var values []string
			if err := json.Unmarshal([]byte(value), &values); err != nil {
				return nil, errercode.CreateCode(errercode.InvalidParam)
			}
			for _, item := range values {
				if item = strings.TrimSpace(item); item != "" {
					result = append(result, item)
				}
			}
			continue
		}
		if strings.Contains(value, ",") {
			for _, item := range strings.Split(value, ",") {
				if item = strings.TrimSpace(item); item != "" {
					result = append(result, item)
				}
			}
			continue
		}
		result = append(result, value)
	}
	return result, nil
}

func mapCreateRoomMultipartErr(err error) error {
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
