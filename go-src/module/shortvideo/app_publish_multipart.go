package shortvideo

import (
	"context"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

type appPublishShortVideoInput struct {
	VideoName        string
	CoverName        string
	Title            string
	IsPaid           uint8
	PayDiamond       float64
	CategoryId       int
	Source           uint8
	Duration         uint32
	FreeWatchSeconds uint32
}

func parseAppPublishShortVideoMultipart(ctx context.Context, r *ghttp.Request) (*appPublishShortVideoInput, error) {
	if r == nil || r.Request == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	reader, err := r.Request.MultipartReader()
	if err != nil {
		return nil, mapAppPublishMultipartErr(err)
	}

	fields := make(map[string]string)
	ret := &appPublishShortVideoInput{}
	cleanup := func() {
		if ret.VideoName != "" {
			upload.DeleteUploadedFile(ret.VideoName)
		}
		if ret.CoverName != "" {
			upload.DeleteUploadedFile(ret.CoverName)
		}
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			cleanup()
			return nil, mapAppPublishMultipartErr(err)
		}

		formName := part.FormName()
		switch formName {
		case "file":
			if ret.VideoName != "" {
				part.Close()
				cleanup()
				return nil, errercode.CreateCode(errercode.InvalidParam)
			}
			ret.VideoName, err = upload.StreamUploadShortVideoPart(part, int64(getShortVideoMaxFileSize()))
		case "cover":
			if strings.TrimSpace(part.FileName()) == "" {
				err = nil
			} else if ret.CoverName != "" {
				err = errercode.CreateCode(errercode.InvalidParam)
			} else {
				maxCoverBytes := int64(getShortVideoMaxCoverFileSize()) * 1024 * 1024
				ret.CoverName, err = upload.StreamUploadImagePart(part, maxCoverBytes)
			}
		default:
			if part.FileName() == "" {
				value, readErr := readMultipartTextField(part)
				part.Close()
				if readErr != nil {
					cleanup()
					return nil, mapAppPublishMultipartErr(readErr)
				}
				fields[formName] = value
				continue
			}
			err = errercode.CreateCode(errercode.InvalidParam)
		}
		part.Close()
		if err != nil {
			cleanup()
			return nil, mapAppPublishMultipartErr(err)
		}
	}

	if ret.VideoName == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := fillAppPublishShortVideoInput(ret, fields); err != nil {
		cleanup()
		return nil, err
	}
	if ret.CoverName != "" {
		if err := upload.RequireAppImageCompliant(ctx, ret.CoverName); err != nil {
			cleanup()
			return nil, err
		}
	}
	return ret, nil
}

func readMultipartTextField(part io.Reader) (string, error) {
	body, err := io.ReadAll(io.LimitReader(part, 4096))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(body)), nil
}

func fillAppPublishShortVideoInput(ret *appPublishShortVideoInput, fields map[string]string) error {
	title := strings.TrimSpace(fields["title"])
	if title == "" || len(title) > 64 {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	ret.Title = title

	isPaid, err := parseUint8Field(fields["isPaid"], "isPaid")
	if err != nil {
		return err
	}
	if isPaid != 0 && isPaid != 1 {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	ret.IsPaid = isPaid

	ret.PayDiamond, err = parseFloat64Field(fields["payDiamond"])
	if err != nil {
		return err
	}

	categoryId, err := parseIntField(fields["categoryId"])
	if err != nil {
		return err
	}
	ret.CategoryId = categoryId

	source, err := parseUint8Field(fields["source"], "source")
	if err != nil {
		return err
	}
	if source != 1 && source != 2 && source != 3 {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	ret.Source = source

	duration, err := parseUint32Field(fields["duration"], "duration")
	if err != nil {
		return err
	}
	if duration == 0 {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	ret.Duration = duration

	freeWatchSeconds, err := parseUint32Field(fields["freeWatchSeconds"], "freeWatchSeconds")
	if err != nil {
		return err
	}
	ret.FreeWatchSeconds = freeWatchSeconds
	return nil
}

func parseUint8Field(raw, name string) (uint8, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		if name == "isPaid" || name == "freeWatchSeconds" {
			return 0, nil
		}
		return 0, errercode.CreateCode(errercode.InvalidParam)
	}
	v, err := strconv.ParseUint(raw, 10, 8)
	if err != nil {
		return 0, errercode.CreateCode(errercode.InvalidParam)
	}
	return uint8(v), nil
}

func parseUint32Field(raw, name string) (uint32, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		if name == "freeWatchSeconds" {
			return 0, nil
		}
		return 0, errercode.CreateCode(errercode.InvalidParam)
	}
	v, err := strconv.ParseUint(raw, 10, 32)
	if err != nil {
		return 0, errercode.CreateCode(errercode.InvalidParam)
	}
	return uint32(v), nil
}

func parseIntField(raw string) (int, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, nil
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, errercode.CreateCode(errercode.InvalidParam)
	}
	return v, nil
}

func parseFloat64Field(raw string) (float64, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, nil
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, errercode.CreateCode(errercode.InvalidParam)
	}
	return v, nil
}

func mapAppPublishMultipartErr(err error) error {
	if err == nil {
		return nil
	}
	var bizErr *errercode.XError
	if errors.As(err, &bizErr) {
		return err
	}
	if upload.IsUploadFileTooLarge(err) {
		return errercode.CreateCode(errercode.ShortVideoFileTooLarge)
	}
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	if strings.Contains(strings.ToLower(err.Error()), "unexpected eof") {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	return errercode.CreateCode(errercode.InvalidParam)
}
