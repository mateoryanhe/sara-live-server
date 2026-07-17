package logquery

import "sort"

func compareLogTimeDesc(a, b string) bool {
	ta, oka := parseLogTime(a)
	tb, okb := parseLogTime(b)
	if oka && okb {
		return ta.After(tb)
	}
	if oka {
		return true
	}
	if okb {
		return false
	}
	return a > b
}

func compareLogTimeAsc(a, b string) bool {
	return compareLogTimeDesc(b, a)
}

func sortDetailLogsByTimeDesc(entries []*DetailLogEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return compareLogTimeDesc(entries[i].Time, entries[j].Time)
	})
}

func sortAccessLogsByTimeDesc(entries []*AccessLogEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return compareLogTimeDesc(entries[i].Time, entries[j].Time)
	})
}

func sortErrorLogsByTimeDesc(entries []*ErrorLogEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return compareLogTimeDesc(entries[i].Time, entries[j].Time)
	})
}

func sortDetailLogsByTimeAsc(entries []*DetailLogEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return compareLogTimeAsc(entries[i].Time, entries[j].Time)
	})
}

func sortAccessLogsByTimeAsc(entries []*AccessLogEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return compareLogTimeAsc(entries[i].Time, entries[j].Time)
	})
}

func sortErrorLogsByTimeAsc(entries []*ErrorLogEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return compareLogTimeAsc(entries[i].Time, entries[j].Time)
	})
}
