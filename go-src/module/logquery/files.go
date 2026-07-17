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
	return listDetailLogFilesInDir(paths.DetailLogDir, paths.DetailLogPattern)
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
	addFiles(listLogFilesByPattern(dir, pattern))

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
	return listLogFilesByPattern(paths.AccessLogDir, paths.AccessLogPattern)
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
