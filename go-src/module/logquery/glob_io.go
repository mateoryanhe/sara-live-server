package logquery

import (
	"path/filepath"
	"strings"
)

type logGlobKind int

const (
	logGlobDetail logGlobKind = iota
	logGlobAccess
	logGlobError
)

func dirsForLogGlobKind(kind logGlobKind, paths logPaths) []string {
	switch kind {
	case logGlobAccess:
		return uniqueNonEmptyDirs(paths.AccessLogDir, paths.DetailLogDir, paths.ErrorLogDir)
	case logGlobError:
		return uniqueNonEmptyDirs(paths.ErrorLogDir, paths.DetailLogDir, paths.AccessLogDir)
	default:
		return uniqueNonEmptyDirs(paths.DetailLogDir, paths.AccessLogDir, paths.ErrorLogDir)
	}
}

func filterLogFilesByKind(files []string, kind logGlobKind) []string {
	switch kind {
	case logGlobAccess:
		return filterAccessLogFiles(files)
	case logGlobError:
		return filterErrorLogFiles(files)
	default:
		return filterDetailLogFiles(files)
	}
}

func globLogFilesForRange(kind logGlobKind, startDate, endDate string) []string {
	paths := loadLogPaths()
	dirs := dirsForLogGlobKind(kind, paths)

	fileSet := make(map[string]struct{})
	for _, dir := range dirs {
		dir = normalizeDir(dir)
		if dir == "" {
			continue
		}
		matches, err := filepath.Glob(filepath.Join(dir, "*.log"))
		if err != nil || len(matches) == 0 {
			continue
		}
		matches = filterLogFilesByKind(matches, kind)
		matches = filterLogFilesByDateRange(matches, startDate, endDate)
		for _, filePath := range matches {
			fileSet[filePath] = struct{}{}
		}
	}

	files := make([]string, 0, len(fileSet))
	for filePath := range fileSet {
		files = append(files, filePath)
	}
	sortLogFilesDesc(files)
	return filterAllowedExistingLogFiles(files, dirs...)
}

func isFuzzyGrepQuery(traceId string) bool {
	traceId = strings.TrimSpace(traceId)
	if traceId == "" {
		return true
	}
	if strings.HasPrefix(traceId, "{") && strings.HasSuffix(traceId, "}") {
		traceId = strings.Trim(traceId, "{}")
	}
	// 完整 traceId 走 grep 直读文件；部分匹配走 cat + *.log 目录扫描
	return len(traceId) < 32
}
