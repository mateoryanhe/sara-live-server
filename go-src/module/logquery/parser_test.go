package logquery

import (
	"strings"
	"testing"
)

func TestParseDetailLogLine(t *testing.T) {
	line := `2026-07-06T00:00:13.764 [INFO] {a2ebcfb5468abf18bff439322487c3ed} log_util.go:94: time=2026-07-06 00:00:13.764,收到前端请求,reqId=1783160418858,url=/liveRoom/reportLiveStartStatus,ip=113.109.204.237,authId=2070436154057953280,请求数据={"RoomId":"2070436154057953280"}`
	entry, ok := parseDetailLogLine(line)
	if !ok {
		t.Fatal("expected parse success")
	}
	if entry.TraceId != "a2ebcfb5468abf18bff439322487c3ed" {
		t.Fatalf("unexpected traceId: %s", entry.TraceId)
	}
	if entry.ReqId != "1783160418858" {
		t.Fatalf("unexpected reqId: %s", entry.ReqId)
	}
	if entry.AuthId != "2070436154057953280" {
		t.Fatalf("unexpected authId: %s", entry.AuthId)
	}
	if entry.Url != "/liveRoom/reportLiveStartStatus" {
		t.Fatalf("unexpected url: %s", entry.Url)
	}
}

func TestParseAccessLogLine(t *testing.T) {
	line := `2026-07-06T00:00:13.765 {a2ebcfb5468abf18bff439322487c3ed} 200 "POST https www.bigtktool.shop /liveRoom/reportLiveStartStatus HTTP/1.1" 0.123, 113.109.204.237, "", "Dart/3.11 (dart:io)"`
	entry, ok := parseAccessLogLine(line)
	if !ok {
		t.Fatal("expected parse success")
	}
	if entry.StatusCode != 200 {
		t.Fatalf("unexpected status: %d", entry.StatusCode)
	}
	if entry.Url != "/liveRoom/reportLiveStartStatus" {
		t.Fatalf("unexpected url: %s", entry.Url)
	}
	if entry.HandlerMs != 123 {
		t.Fatalf("unexpected handlerMs: %f", entry.HandlerMs)
	}
	if entry.Ip != "113.109.204.237" {
		t.Fatalf("unexpected ip: %s", entry.Ip)
	}
}

func TestParseErrorLogBlock(t *testing.T) {
	line1 := `2026-07-06T00:00:13.765 [ERRO] {a2ebcfb5468abf18bff439322487c3ed} 500 "POST https www.bigtktool.shop /liveRoom/reportLiveStartStatus HTTP/1.1" 0.123, 113.109.204.237, "", "Dart/3.11", -1, "INTERNAL PANIC", "exception recovered: custom error"`
	line2 := `Stack:`
	line3 := `main.go:10`
	var entry *ErrorLogEntry
	body := strings.Builder{}
	for i, line := range []string{line1, line2, line3} {
		if header, ok := parseErrorLogHeader(line); ok {
			entry = header
			body.Reset()
			body.WriteString(line)
			continue
		}
		if entry != nil && i > 0 {
			body.WriteByte('\n')
			body.WriteString(line)
		}
	}
	if entry == nil {
		t.Fatal("expected parse success")
	}
	finalizeErrorLogEntry(entry, body.String())
	if entry.StatusCode != 500 {
		t.Fatalf("unexpected status: %d", entry.StatusCode)
	}
	if entry.ErrorMessage != "INTERNAL PANIC" {
		t.Fatalf("unexpected error message: %s", entry.ErrorMessage)
	}
	if !strings.Contains(entry.Stack, "main.go:10") {
		t.Fatalf("unexpected stack: %s", entry.Stack)
	}
}

func TestParsePushDetailLogLines(t *testing.T) {
	cases := []struct {
		line   string
		authId string
	}{
		{
			line:   `2026-07-06T10:00:00.000 [INFO] {abc123} push_mgr.go:34: authId=1001,发送数据cmd=11,data={"gold":100}`,
			authId: "1001",
		},
		{
			line:   `2026-07-06T10:00:01.000 [INFO] {abc124} push_mgr.go:51: userid=1001,发送数据cmd=2`,
			authId: "1001",
		},
		{
			line:   `2026-07-06T10:00:02.000 [INFO] {abc125} push_mgr.go:123: ip=127.0.0.1:1234,玩家=1001,上线了`,
			authId: "1001",
		},
	}
	for _, tc := range cases {
		entry, ok := parseDetailLogLine(tc.line)
		if !ok {
			t.Fatalf("parse failed: %s", tc.line)
		}
		if entry.AuthId != tc.authId {
			t.Fatalf("unexpected authId %s, want %s, line=%s", entry.AuthId, tc.authId, tc.line)
		}
		if !messageMatchesAuthId(entry, tc.authId) {
			t.Fatalf("authId filter mismatch: %s", tc.line)
		}
	}
}
