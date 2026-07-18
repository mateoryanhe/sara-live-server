package logquery

import (
	"bufio"
	"strings"
	"testing"
	"time"
)

func TestLineMatchesTraceIdFilterIgnoresEmbeddedTraceId(t *testing.T) {
	targetTraceId := "137dc2d8f1ffc218c7c68d668ff9b246"
	embeddedTraceId := "4436506d2fb4c218c410ce7cfde6e067"
	line := `2026-07-17T06:30:16.638 [INFO] {` + targetTraceId + `} log_util.go:227: respContent={"data":[{"traceId":"` + embeddedTraceId + `"}]}`

	if !lineMatchesTraceIdFilter(line, targetTraceId) {
		t.Fatal("expected target traceId to match line prefix")
	}
	if lineMatchesTraceIdFilter(line, embeddedTraceId) {
		t.Fatal("embedded traceId must not match line prefix filter")
	}
}

func TestQuickAcceptLineRejectsOutOfRangeTime(t *testing.T) {
	line := `2026-01-01T00:00:00.000 [INFO] {abc} log_util.go:1: test`
	start, _ := parseDate("2026-07-01")
	end := start.AddDate(0, 0, 1).Add(24*time.Hour - time.Millisecond)
	filter := detailQueryFilter{rangeStart: start, rangeEnd: end}
	if filter.quickAcceptLine(line) {
		t.Fatal("expected out-of-range line to be rejected")
	}
}

func TestPrepareDetailLogEntryForListClearsRawAndTruncatesMessage(t *testing.T) {
	entry := &DetailLogEntry{
		Message: strings.Repeat("m", maxDetailLogListMessageBytes+100),
		Raw:     "raw-content",
	}
	prepareDetailLogEntryForList(entry)
	if entry.Raw != "" {
		t.Fatal("expected raw to be cleared")
	}
	if !strings.Contains(entry.Message, detailLogPayloadTruncated) {
		t.Fatal("expected truncated message marker")
	}
	if len(entry.Message) != maxDetailLogListMessageBytes+len(detailLogPayloadTruncated) {
		t.Fatalf("unexpected message length: %d", len(entry.Message))
	}
}

func TestParseDetailLogLineOptSkipsHeadersWhenMissing(t *testing.T) {
	line := `2026-07-17T13:36:53.171 [INFO] {abc} log_util.go:65: 鉴权完成,reqId=1,authId=2,authMs=3ms,url=/api/test`
	entry, ok := parseDetailLogLineOpt(line)
	if !ok {
		t.Fatal("expected parse success")
	}
	if entry.ReqId != "1" || entry.AuthId != "2" {
		t.Fatalf("unexpected ids: reqId=%s authId=%s", entry.ReqId, entry.AuthId)
	}
}

func TestReadLogLineLimitedUsesSmallerCap(t *testing.T) {
	payload := strings.Repeat("x", 64*1024)
	line := "head," + payload + "\nnext\n"
	reader := bufio.NewReader(strings.NewReader(line))
	first, err := readLogLineLimited(reader, 1024)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if len(first) >= 64*1024 {
		t.Fatalf("expected limited line length, got %d", len(first))
	}
	if !strings.Contains(first, "...[line truncated]") {
		t.Fatal("expected truncated marker")
	}
}
