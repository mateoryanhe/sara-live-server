package logquery

import "testing"

func TestIsFuzzyGrepQuery(t *testing.T) {
	if !isFuzzyGrepQuery("") {
		t.Fatal("empty traceId should use fuzzy cat path")
	}
	if isFuzzyGrepQuery("137dc2d8f1ffc218c7c68d668ff9b246") {
		t.Fatal("full traceId should not use cat path")
	}
	if !isFuzzyGrepQuery("137dc2d8") {
		t.Fatal("partial traceId should use cat path")
	}
}

func TestFilterLogFilesByKind(t *testing.T) {
	files := []string{
		"/log/2026-07-17.log",
		"/log/access-2026-07-17.log",
		"/log/error-20260717.log",
	}
	detail := filterLogFilesByKind(files, logGlobDetail)
	if len(detail) != 1 || detail[0] != files[0] {
		t.Fatalf("unexpected detail files: %#v", detail)
	}
	access := filterLogFilesByKind(files, logGlobAccess)
	if len(access) != 1 || access[0] != files[1] {
		t.Fatalf("unexpected access files: %#v", access)
	}
}
