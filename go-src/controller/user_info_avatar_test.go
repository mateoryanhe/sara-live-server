package controller

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/net/ghttp"
)

func TestParseUploadAvatarMultipartStreamsFile(t *testing.T) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	filePart, err := writer.CreateFormFile("file", "avatar.jpg")
	if err != nil {
		t.Fatal(err)
	}
	avatarData := strings.Repeat("x", 64*1024)
	if _, err = io.WriteString(filePart, avatarData); err != nil {
		t.Fatal(err)
	}
	if err = writer.Close(); err != nil {
		t.Fatal(err)
	}

	httpReq := httptest.NewRequest("POST", "/userInfo/uploadAvatar", &body)
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	var gotData []byte
	storedName, err := parseUploadAvatarMultipartWithUploader(
		context.Background(),
		&ghttp.Request{Request: httpReq},
		func(_ context.Context, part *multipart.Part, _ int64) (string, error) {
			gotData, err = io.ReadAll(part)
			return "images/avatar.jpg", err
		},
	)
	if err != nil {
		t.Fatalf("parseUploadAvatarMultipartWithUploader() error = %v", err)
	}
	if storedName != "images/avatar.jpg" {
		t.Fatalf("storedName = %q", storedName)
	}
	if string(gotData) != avatarData {
		t.Fatalf("avatar data length = %d, want %d", len(gotData), len(avatarData))
	}
}

func TestParseUploadAvatarMultipartRequiresFile(t *testing.T) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("nickname", "test"); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	httpReq := httptest.NewRequest("POST", "/userInfo/uploadAvatar", &body)
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	_, err := parseUploadAvatarMultipartWithUploader(
		context.Background(),
		&ghttp.Request{Request: httpReq},
		func(context.Context, *multipart.Part, int64) (string, error) {
			t.Fatal("uploader should not be called")
			return "", nil
		},
	)
	if err == nil {
		t.Fatal("parseUploadAvatarMultipartWithUploader() expected missing file error")
	}
}
