package logquery

import (
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var logFileDateRes = []*regexp.Regexp{
	regexp.MustCompile(`(\d{4}-\d{2}-\d{2})`),
	regexp.MustCompile(`(\d{4})(\d{2})(\d{2})`),
}

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
	dirs := uniqueNonEmptyDirs(paths.DetailLogDir, paths.AccessLogDir, paths.ErrorLogDir)
	fileSet := make(map[string]struct{})
	for _, dir := range dirs {
		for _, filePath := range listDetailLogFilesInDir(dir, paths.DetailLogPattern) {
			fileSet[filePath] = struct{}{}
		}
	}
	files := make([]string, 0, len(fileSet))
	for filePath := range fileSet {
		files = append(files, filePath)
	}
	sort.Strings(files)
	return files
}

func listDetailLogFilesInDir(dir, pattern string) []string {
	dir = normalizeDir(dir)
	if dir == "" {
		return nil
	}
	fileSet := make(map[string]struct{})
	addFiles := func(list []string) {
		for _, filePath := range list {
			if filePath == "" {
				continue
			}
			fileSet[filePath] = struct{}{}
		}
	}
	addFiles(filterDetailLogFiles(listDirLogFiles(dir)))
	addFiles(filterDetailLogFiles(listLogFilesByPattern(dir, pattern)))

	files := make([]string, 0, len(fileSet))
	for filePath := range fileSet {
		files = append(files, filePath)
	}
	sort.Strings(files)
	return files
}

func listDetailLogDirs() []string {
	paths := loadLogPaths()
	return uniqueNonEmptyDirs(paths.DetailLogDir, paths.ErrorLogDir)
}

func listErrorLogFiles() []string {
	paths := loadLogPaths()
	dirs := uniqueNonEmptyDirs(paths.ErrorLogDir, paths.DetailLogDir)
	fileSet := make(map[string]struct{})
	addFiles := func(list []string) {
		for _, filePath := range list {
			if filePath == "" {
				continue
			}
			fileSet[filePath] = struct{}{}
		}
	}
	for _, dir := range dirs {
		addFiles(filterErrorLogFiles(listDirLogFiles(dir)))
		addFiles(listLogFilesByPattern(dir, paths.ErrorLogPattern))
		addFiles(listLogFilesByPattern(dir, "error-*.log"))
		addFiles(listLogFilesByPattern(dir, "error.*.log"))
	}

	files := make([]string, 0, len(fileSet))
	for filePath := range fileSet {
		files = append(files, filePath)
	}
	sort.Strings(files)
	return files
}

func listAccessLogFiles() []string {
	paths := loadLogPaths()
	dirs := uniqueNonEmptyDirs(paths.AccessLogDir, paths.DetailLogDir, paths.ErrorLogDir)
	fileSet := make(map[string]struct{})
	addFiles := func(list []string) {
		for _, filePath := range list {
			if filePath == "" {
				continue
			}
			fileSet[filePath] = struct{}{}
		}
	}
	for _, dir := range dirs {
		addFiles(filterAccessLogFiles(listDirLogFiles(dir)))
		addFiles(listLogFilesByPattern(dir, paths.AccessLogPattern))
		addFiles(listLogFilesByPattern(dir, "access-*.log"))
		addFiles(listLogFilesByPattern(dir, "access.*.log"))
	}

	files := make([]string, 0, len(fileSet))
	for filePath := range fileSet {
		files = append(files, filePath)
	}
	sort.Strings(files)
	return files
}

func filterAccessLogFiles(files []string) []string {
	ret := make([]string, 0, len(files))
	for _, filePath := range files {
		if isAccessLogFileName(filepath.Base(filePath)) {
			ret = append(ret, filePath)
		}
	}
	return ret
}

func filterErrorLogFiles(files []string) []string {
	ret := make([]string, 0, len(files))
	for _, filePath := range files {
		if isErrorLogFileName(filepath.Base(filePath)) {
			ret = append(ret, filePath)
		}
	}
	return ret
}

func filterDetailLogFiles(files []string) []string {
	ret := make([]string, 0, len(files))
	for _, filePath := range files {
		name := filepath.Base(filePath)
		if isDetailLogFileName(name) {
			ret = append(ret, filePath)
		}
	}
	return ret
}

func isErrorLogFileName(name string) bool {
	name = strings.ToLower(strings.TrimSpace(name))
	return strings.HasPrefix(name, "error") && strings.HasSuffix(name, ".log")
}

func isAccessLogFileName(name string) bool {
	name = strings.ToLower(strings.TrimSpace(name))
	return strings.HasPrefix(name, "access") && strings.HasSuffix(name, ".log")
}

func isDetailLogFileName(name string) bool {
	name = strings.ToLower(strings.TrimSpace(name))
	if !strings.HasSuffix(name, ".log") {
		return false
	}
	return !isErrorLogFileName(name) && !isAccessLogFileName(name)
}

func uniqueNonEmptyDirs(dirs ...string) []string {
	seen := make(map[string]struct{})
	ret := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		dir = normalizeDir(dir)
		if dir == "" {
			continue
		}
		if _, ok := seen[dir]; ok {
			continue
		}
		seen[dir] = struct{}{}
		ret = append(ret, dir)
	}
	return ret
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

func extractLogFileDate(filePath string) (string, bool) {
	name := filepath.Base(filePath)
	if match := logFileDateRes[0].FindStringSubmatch(name); len(match) > 1 {
		return match[1], true
	}
	if match := logFileDateRes[1].FindStringSubmatch(name); len(match) > 3 {
		return match[1] + "-" + match[2] + "-" + match[3], true
	}
	return "", false
}

func filterLogFilesByDateRange(files []string, startDate, endDate string) []string {
	allowedDates := listDates(startDate, endDate)
	if len(allowedDates) == 0 {
		return files
	}
	allowed := make(map[string]struct{}, len(allowedDates))
	for _, dateStr := range allowedDates {
		allowed[dateStr] = struct{}{}
	}

	ret := make([]string, 0, len(files))
	for _, filePath := range files {
		fileDate, ok := extractLogFileDate(filePath)
		if !ok {
			ret = append(ret, filePath)
			continue
		}
		if _, inRange := allowed[fileDate]; inRange {
			ret = append(ret, filePath)
		}
	}
	return ret
}

func sortLogFilesDesc(files []string) {
	sort.Slice(files, func(i, j int) bool {
		return files[i] > files[j]
	})
}

func collectLogFilesForRange(startDate, endDate string, datedFiles func(string) string, fallbackFiles func() []string) []string {
	fileSet := make(map[string]struct{})
	addFile := func(filePath string) {
		if filePath == "" {
			return
		}
		fileSet[filePath] = struct{}{}
	}
	for _, dateStr := range listDates(startDate, endDate) {
		addFile(datedFiles(dateStr))
	}
	for _, filePath := range filterLogFilesByDateRange(fallbackFiles(), startDate, endDate) {
		addFile(filePath)
	}
	files := make([]string, 0, len(fileSet))
	for filePath := range fileSet {
		files = append(files, filePath)
	}
	sortLogFilesDesc(files)
	return files
}

func listDetailLogFilesForRange(startDate, endDate string) []string {
	files := collectLogFilesForRange(startDate, endDate, resolveDetailLogFile, listDetailLogFiles)
	ret := make([]string, 0, len(files))
	for _, filePath := range files {
		if isDetailLogFileName(filepath.Base(filePath)) {
			ret = append(ret, filePath)
		}
	}
	sortLogFilesDesc(ret)
	return ret
}

func listAccessLogFilesForRange(startDate, endDate string) []string {
	return collectLogFilesForRange(startDate, endDate, resolveAccessLogFile, listAccessLogFiles)
}

func listErrorLogFilesForRange(startDate, endDate string) []string {
	return collectLogFilesForRange(startDate, endDate, resolveErrorLogFile, listErrorLogFiles)
}
