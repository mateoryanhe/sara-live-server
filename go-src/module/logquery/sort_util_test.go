package logquery

import "testing"

func TestCompareLogTimeDesc(t *testing.T) {
	if !compareLogTimeDesc("2026-07-06T10:00:01.000", "2026-07-06T10:00:00.000") {
		t.Fatal("newer time should sort first")
	}
	if compareLogTimeDesc("2026-07-06T10:00:00.000", "2026-07-06T10:00:01.000") {
		t.Fatal("older time should sort later")
	}
}

func TestSortDetailLogsByTimeDesc(t *testing.T) {
	logs := []*DetailLogEntry{
		{Time: "2026-07-06T10:00:00.000"},
		{Time: "2026-07-06T12:00:00.000"},
		{Time: "2026-07-06T11:00:00.000"},
	}
	sortDetailLogsByTimeDesc(logs)
	if logs[0].Time != "2026-07-06T12:00:00.000" {
		t.Fatalf("unexpected first: %s", logs[0].Time)
	}
}
