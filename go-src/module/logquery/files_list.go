package logquery

import (
	"path/filepath"
	"sort"
	"strings"
)

func listLogFilesByPrefix(logDir, prefix, startDate, endDate string) []string {
	if logDir == "" || prefix == "" {
		return nil
	}
	glob := filepath.Join(logDir, prefix+"*.log")
	matches, err := filepath.Glob(glob)
	if err != nil || len(matches) == 0 {
		return nil
	}
	allowed := map[string]struct{}{}
	for _, dateStr := range listDates(startDate, endDate) {
		allowed[dateStr] = struct{}{}
	}
	ret := make([]string, 0, len(matches))
	for _, filePath := range matches {
		base := filepath.Base(filePath)
		if !strings.HasPrefix(base, prefix) {
			continue
		}
		if fileDate, ok := extractDateFromFileName(base); ok {
			if _, inRange := allowed[fileDate]; !inRange {
				continue
			}
		}
		ret = append(ret, filePath)
	}
	sortLogFilesByDate(ret)
	return ret
}

func sortLogFilesByDate(files []string) {
	sort.Slice(files, func(i, j int) bool {
		di, iok := extractDateFromFileName(filepath.Base(files[i]))
		dj, jok := extractDateFromFileName(filepath.Base(files[j]))
		if iok && jok {
			return di < dj
		}
		return files[i] < files[j]
	})
}

func extractDateFromFileName(name string) (string, bool) {
	// access-2026-07-30.log → suffix "2026-07-30" (need len >= 14)
	if len(name) >= 14 {
		candidate := name[len(name)-14 : len(name)-4]
		if len(candidate) == 10 && candidate[4] == '-' && candidate[7] == '-' {
			return candidate, true
		}
	}
	for i := 0; i+10 <= len(name); i++ {
		part := name[i : i+10]
		if part[4] == '-' && part[7] == '-' {
			return part, true
		}
	}
	if len(name) >= 8 {
		for i := 0; i+8 <= len(name); i++ {
			part := name[i : i+8]
			if isDigits(part) {
				return part[0:4] + "-" + part[4:6] + "-" + part[6:8], true
			}
		}
	}
	return "", false
}

func isDigits(s string) bool {
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}
