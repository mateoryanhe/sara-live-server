package logquery

import (
	"strings"
	"testing"
)

func TestTruncateResponseWriteRespContent(t *testing.T) {
	payload := strings.Repeat("x", 2000)
	raw := `2026-07-17T08:00:00.000 [INFO] {abc} log_util.go:247: 应答序列化,写入框架缓冲区,writeMs=4ms,url=/api/test,respContent=` + payload
	entry, ok := parseDetailLogLine(raw)
	if !ok {
		t.Fatal("parse detail failed")
	}
	trimDetailLogEntryForQuery(entry)
	if !strings.Contains(entry.Message, detailLogPayloadTruncated) {
		t.Fatalf("expected message truncated: %s", entry.Message)
	}
	idx := strings.Index(entry.Message, "respContent=")
	value := entry.Message[idx+len("respContent="):]
	if len(value) != maxDetailLogRespContentBytes+len(detailLogPayloadTruncated) {
		t.Fatalf("unexpected respContent length: %d", len(value))
	}
}

func TestTrimDetailLogEntrySkipsBodyContent(t *testing.T) {
	payload := strings.Repeat("y", 1500)
	raw := `2026-07-17T08:00:00.000 [INFO] {abc} log_util.go:198: 读取请求Body,url=/api/test,bodyContent=` + payload
	entry, ok := parseDetailLogLine(raw)
	if !ok {
		t.Fatal("parse detail failed")
	}
	before := entry.Message
	trimDetailLogEntryForQuery(entry)
	if entry.Message != before {
		t.Fatalf("bodyContent log should not be truncated")
	}
}
