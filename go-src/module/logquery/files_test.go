package logquery

import (
	"strings"
	"testing"
)

func TestIsErrorLogFileName(t *testing.T) {
	cases := []struct {
		name string
		ok   bool
	}{
		{name: "error-20260717.log", ok: true},
		{name: "error-2026-07-17.log", ok: true},
		{name: "error.20260717.log", ok: true},
		{name: "access-2026-07-17.log", ok: false},
		{name: "2026-07-17.log", ok: false},
	}
	for _, tc := range cases {
		if isErrorLogFileName(tc.name) != tc.ok {
			t.Fatalf("isErrorLogFileName(%q) = %v, want %v", tc.name, !tc.ok, tc.ok)
		}
	}
}

func TestFilterErrorLogFiles(t *testing.T) {
	files := filterErrorLogFiles([]string{
		"/log/error-20260717.log",
		"/log/error-2026-07-17.log",
		"/log/access-2026-07-17.log",
		"/log/2026-07-17.log",
	})
	if len(files) != 2 {
		t.Fatalf("unexpected filtered files: %#v", files)
	}
}

func TestFilterDetailLogFiles(t *testing.T) {
	files := filterDetailLogFiles([]string{
		"/log/2026-07-17.log",
		"/log/error-20260717.log",
		"/log/access-2026-07-17.log",
	})
	if len(files) != 1 || files[0] != "/log/2026-07-17.log" {
		t.Fatalf("unexpected filtered detail files: %#v", files)
	}
}

func TestFilterAccessLogFiles(t *testing.T) {
	files := filterAccessLogFiles([]string{
		"/log/access-2026-07-17.log",
		"/log/access.20260717.log",
		"/log/error-20260717.log",
		"/log/2026-07-17.log",
	})
	if len(files) != 2 {
		t.Fatalf("unexpected filtered access files: %#v", files)
	}
}

func TestMatchTraceId(t *testing.T) {
	if !matchTraceId("abc123", "abc123") {
		t.Fatal("exact traceId should match")
	}
	if !matchTraceId("abc123def", "abc123") {
		t.Fatal("fuzzy traceId should match")
	}
}

func TestDetailEntryToErrorEntry(t *testing.T) {
	line := `2026-07-06T10:00:00.000 [ERRO] {abc123} response_middleware_util.go:95: ErrorLog source=Handler time=2026-07-06 10:00:00.000,reqId=1,authId=2,method=POST,url=/api/test err=panic stack=main.go:10`
	detail, ok := parseDetailLogLine(line)
	if !ok {
		t.Fatal("parse detail failed")
	}
	entry := detailEntryToErrorEntry(detail)
	if entry == nil {
		t.Fatal("expected error entry")
	}
	if entry.ErrorMessage != "ErrorLog/Handler" {
		t.Fatalf("unexpected error message: %s", entry.ErrorMessage)
	}
	if !strings.Contains(entry.Stack, "main.go:10") {
		t.Fatalf("unexpected stack: %s", entry.Stack)
	}
	if entry.Url != "/api/test" {
		t.Fatalf("unexpected url: %s", entry.Url)
	}
}
