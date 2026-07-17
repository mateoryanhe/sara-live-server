package logquery

import (
	"path/filepath"
	"sort"
	"strings"
)

func listDirLogFiles(dir string) []string {
	dir = normalizeDir(dir)
	if dir == "" {
		return nil
	}
	matches, err := filepath.Glob(filepath.Join(dir, "*.log"))
	if err != nil || len(matches) == 0 {
		return nil
	}
	sort.Strings(matches)
	return matches
}

func listDetailLogFiles() []string {
	paths := loadLogPaths()
	files := listDirLogFiles(paths.DetailLogDir)
	if len(files) > 0 {
		return files
	}
	return listLogFilesByPattern(paths.DetailLogDir, paths.DetailLogPattern)
}

func listAccessLogFiles() []string {
	paths := loadLogPaths()
	files := listDirLogFiles(paths.AccessLogDir)
	if len(files) > 0 {
		return files
	}
	return listLogFilesByPattern(paths.AccessLogDir, paths.AccessLogPattern)
}

func listLogFilesByPattern(dir, pattern string) []string {
	dir = normalizeDir(dir)
	if dir == "" {
		return nil
	}
	globPattern := patternToGlob(pattern)
	matches, err := filepath.Glob(filepath.Join(dir, globPattern))
	if err != nil || len(matches) == 0 {
		return nil
	}
	sort.Strings(matches)
	return matches
}

func patternToGlob(pattern string) string {
	s := pattern
	for _, placeholder := range []string{"{Y-m-d}", "{Ymd}", "{Y-m}", "{Y}", "{m}", "{d}"} {
		s = strings.ReplaceAll(s, placeholder, "*")
	}
	return s
}
