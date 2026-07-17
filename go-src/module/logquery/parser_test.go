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

func TestParseDetailLogLineWithHeaders(t *testing.T) {
	line := `2026-07-17T13:36:53.170 [INFO] {a1c6e4392fd2d251trace001} log_util.go:24: 收到前端请求,enterTime=2026-07-17 05:36:53.170,从队列进入到中间件时间间隔Ms=0ms,method=POST,url=/liveRoom/roomList,ip=113.109.204.150,headers={"Accept":["application/json"],"Authorization":["2076611433029701632.1ghdc9jjdqcdjzuxiqh8ogq100qvd5z7"],"Reqid":["1784266609644"],"User-Agent":["Dart/3.11 (dart:io)"]}`
	entry, ok := parseDetailLogLine(line)
	if !ok {
		t.Fatal("expected parse success")
	}
	if entry.ReqId != "1784266609644" {
		t.Fatalf("unexpected reqId: %s", entry.ReqId)
	}
	if entry.AuthId != "2076611433029701632" {
		t.Fatalf("unexpected authId: %s", entry.AuthId)
	}
	if entry.Url != "/liveRoom/roomList" {
		t.Fatalf("unexpected url: %s", entry.Url)
	}
	if !matchDetailEntry(entry, "", "1784266609644", "2076611433029701632", "", "") {
		t.Fatal("expected reqId/authId filter to match")
	}
	if entry.ElapsedMs == nil || *entry.ElapsedMs != 0 {
		t.Fatalf("unexpected elapsedMs: %#v", entry.ElapsedMs)
	}
}

func TestParseDetailLogLineElapsedMs(t *testing.T) {
	cases := []struct {
		line string
		want float64
	}{
		{
			line: `2026-07-17T13:36:53.171 [INFO] {abc} log_util.go:65: 鉴权完成,reqId=1,authId=2,authMs=3ms,url=/api/test`,
			want: 3,
		},
		{
			line: `2026-07-17T13:36:53.172 [INFO] {abc} log_util.go:213: Handler执行完成,reqId=1,authId=2,handlerMs=85ms,url=/api/test`,
			want: 85,
		},
		{
			line: `2026-07-17T13:36:53.175 [INFO] {abc} log_util.go:266: 应答写入到系统缓冲区,输出完成,reqId=1,authId=2,afterOutputMs=5ms,gzip=false,totalMs=120ms,url=/api/test`,
			want: 120,
		},
	}
	for _, tc := range cases {
		entry, ok := parseDetailLogLine(tc.line)
		if !ok {
			t.Fatalf("parse failed: %s", tc.line)
		}
		if entry.ElapsedMs == nil || *entry.ElapsedMs != tc.want {
			t.Fatalf("unexpected elapsedMs %#v, want %v, line=%s", entry.ElapsedMs, tc.want, tc.line)
		}
	}
}

func TestMatchTraceIdIgnoresEmbeddedTraceIdInRaw(t *testing.T) {
	targetTraceId := "137dc2d8f1ffc218c7c68d668ff9b246"
	embeddedTraceId := "4436506d2fb4c218c410ce7cfde6e067"
	line := `2026-07-17T06:30:16.638 [INFO] {` + targetTraceId + `} log_util.go:227: 应答序列化,writeMs=4ms,url=/logQuery/queryDetailLogs,respContent={"data":[{"traceId":"` + embeddedTraceId + `"}]}`
	entry, ok := parseDetailLogLine(line)
	if !ok {
		t.Fatal("expected parse success")
	}
	if !matchTraceId(entry.TraceId, targetTraceId) {
		t.Fatal("expected target traceId to match")
	}
	if matchTraceId(entry.TraceId, embeddedTraceId) {
		t.Fatal("embedded traceId must not match parsed traceId field")
	}
	if !fuzzyMatch(entry.Raw, embeddedTraceId) {
		t.Fatal("embedded traceId exists in raw content")
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
