package controller

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/net/ghttp"
)

func TestParseCreateRoomMultipart(t *testing.T) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for name, value := range map[string]string{
		"title":             "hi",
		"notice":            "notice",
		"category":          "2",
		"tagId":             "3",
		"gameCodes":         `["game-a","game-b"]`,
		"ticket":            "1000",
		"billing":           "25.5",
		"privateInviteType": "1",
	} {
		if err := writer.WriteField(name, value); err != nil {
			t.Fatal(err)
		}
	}
	coverPart, err := writer.CreateFormFile("cover", "cover.jpg")
	if err != nil {
		t.Fatal(err)
	}
	coverData := strings.Repeat("x", 64*1024)
	if _, err = io.WriteString(coverPart, coverData); err != nil {
		t.Fatal(err)
	}
	if err = writer.Close(); err != nil {
		t.Fatal(err)
	}

	httpReq := httptest.NewRequest("POST", "/liveRoom/create", &body)
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	var gotCover []byte
	req, coverName, err := parseCreateRoomMultipartWithUploader(
		context.Background(),
		&ghttp.Request{Request: httpReq},
		func(_ context.Context, part *multipart.Part, _ int64) (string, error) {
			gotCover, err = io.ReadAll(part)
			return "images/cover.jpg", err
		},
	)
	if err != nil {
		t.Fatalf("parseCreateRoomMultipartWithUploader() error = %v", err)
	}
	if coverName != "images/cover.jpg" {
		t.Fatalf("coverName = %q", coverName)
	}
	if string(gotCover) != coverData {
		t.Fatalf("cover data length = %d, want %d", len(gotCover), len(coverData))
	}
	if req.Title != "hi" || req.Notice != "notice" || !reflect.DeepEqual(req.GameCodes, []string{"game-a", "game-b"}) {
		t.Fatalf("unexpected parsed request: %+v", req)
	}
}

func TestBindCreateRoomMultipartFields(t *testing.T) {
	fields := map[string][]string{
		"title":             {"hi"},
		"notice":            {"notice"},
		"category":          {"2"},
		"tagId":             {"3"},
		"gameCodes":         {`["game-a","game-b"]`, "game-c,game-d"},
		"ticket":            {"1000"},
		"billing":           {"25.5"},
		"privateInviteType": {"1"},
	}

	req, err := bindCreateRoomMultipartFields(fields)
	if err != nil {
		t.Fatalf("bindCreateRoomMultipartFields() error = %v", err)
	}
	if req.Title != "hi" || req.Notice != "notice" || req.Category != 2 || req.TagId != 3 {
		t.Fatalf("unexpected basic fields: %+v", req)
	}
	if req.Ticket != 1000 || req.Billing != 25.5 || req.PrivateInviteType != 1 {
		t.Fatalf("unexpected pricing fields: %+v", req)
	}
	wantGameCodes := []string{"game-a", "game-b", "game-c", "game-d"}
	if !reflect.DeepEqual(req.GameCodes, wantGameCodes) {
		t.Fatalf("gameCodes = %#v, want %#v", req.GameCodes, wantGameCodes)
	}
}

func TestNormalizeCreateRoomFieldName(t *testing.T) {
	for _, name := range []string{"gameCodes[]", "gameCodes[0]", "gameCodes[12]"} {
		if got := normalizeCreateRoomFieldName(name); got != "gameCodes" {
			t.Fatalf("normalizeCreateRoomFieldName(%q) = %q", name, got)
		}
	}
}
