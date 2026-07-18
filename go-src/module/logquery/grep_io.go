package logquery

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"xr-game-server/dto/logquerydto"
)

const (
	grepMaxOutputLines = 8000
	grepLineSeparator  = '\n'
)

var (
	grepAvailableOnce sync.Once
	grepAvailable     bool
	catAvailableOnce  sync.Once
	catAvailable      bool
)

func isGrepQueryEnabled() bool {
	if runtime.GOOS == "windows" {
		return false
	}
	grepAvailableOnce.Do(func() {
		_, err := exec.LookPath("grep")
		grepAvailable = err == nil
	})
	return grepAvailable
}

func isCatQueryEnabled() bool {
	if runtime.GOOS == "windows" {
		return false
	}
	catAvailableOnce.Do(func() {
		_, err := exec.LookPath("cat")
		catAvailable = err == nil
	})
	return catAvailable
}

func sanitizeGrepPattern(pattern string) (string, bool) {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		return "", false
	}
	if strings.ContainsAny(pattern, "\x00\r\n") {
		return "", false
	}
	return pattern, true
}

func buildDetailGrepPatterns(req *logquerydto.CMSQueryDetailLogsReq) (primary string, secondary []string) {
	secondary = make([]string, 0, 4)
	if traceId, ok := sanitizeGrepPattern(req.TraceId); ok {
		if strings.HasPrefix(traceId, "{") {
			return traceId, collectDetailSecondaryPatterns(req, "traceId")
		}
		return "{" + traceId + "}", collectDetailSecondaryPatterns(req, "traceId")
	}
	if reqId, ok := sanitizeGrepPattern(req.ReqId); ok {
		return reqId, collectDetailSecondaryPatterns(req, "reqId")
	}
	if url, ok := sanitizeGrepPattern(req.Url); ok {
		return url, collectDetailSecondaryPatterns(req, "url")
	}
	if keyword, ok := sanitizeGrepPattern(req.Keyword); ok {
		return keyword, collectDetailSecondaryPatterns(req, "keyword")
	}
	if authId, ok := sanitizeGrepPattern(req.AuthId); ok {
		return authId, collectDetailSecondaryPatterns(req, "authId")
	}
	return "", nil
}

func collectDetailSecondaryPatterns(req *logquerydto.CMSQueryDetailLogsReq, skipField string) []string {
	patterns := make([]string, 0, 4)
	if skipField != "reqId" {
		if reqId, ok := sanitizeGrepPattern(req.ReqId); ok {
			patterns = append(patterns, reqId)
		}
	}
	if skipField != "authId" {
		if authId, ok := sanitizeGrepPattern(req.AuthId); ok {
			patterns = append(patterns, authId)
		}
	}
	if skipField != "url" {
		if url, ok := sanitizeGrepPattern(req.Url); ok {
			patterns = append(patterns, url)
		}
	}
	if skipField != "keyword" {
		if keyword, ok := sanitizeGrepPattern(req.Keyword); ok {
			patterns = append(patterns, keyword)
		}
	}
	return uniqueNonEmptyStrings(patterns)
}

func buildAccessGrepPatterns(req *logquerydto.CMSQueryAccessLogsReq) (primary string, secondary []string) {
	if traceId, ok := sanitizeGrepPattern(req.TraceId); ok {
		if strings.HasPrefix(traceId, "{") {
			return traceId, collectAccessSecondaryPatterns(req, "traceId")
		}
		return "{" + traceId + "}", collectAccessSecondaryPatterns(req, "traceId")
	}
	if url, ok := sanitizeGrepPattern(req.Url); ok {
		return url, collectAccessSecondaryPatterns(req, "url")
	}
	if ip, ok := sanitizeGrepPattern(req.Ip); ok {
		return ip, collectAccessSecondaryPatterns(req, "ip")
	}
	return "", nil
}

func collectAccessSecondaryPatterns(req *logquerydto.CMSQueryAccessLogsReq, skipField string) []string {
	patterns := make([]string, 0, 2)
	if skipField != "url" {
		if url, ok := sanitizeGrepPattern(req.Url); ok {
			patterns = append(patterns, url)
		}
	}
	if skipField != "ip" {
		if ip, ok := sanitizeGrepPattern(req.Ip); ok {
			patterns = append(patterns, ip)
		}
	}
	return uniqueNonEmptyStrings(patterns)
}

func buildErrorGrepPatterns(req *logquerydto.CMSQueryErrorLogsReq) (primary string, secondary []string) {
	if traceId, ok := sanitizeGrepPattern(req.TraceId); ok {
		if strings.HasPrefix(traceId, "{") {
			return traceId, collectErrorSecondaryPatterns(req, "traceId")
		}
		return "{" + traceId + "}", collectErrorSecondaryPatterns(req, "traceId")
	}
	if url, ok := sanitizeGrepPattern(req.Url); ok {
		return url, collectErrorSecondaryPatterns(req, "url")
	}
	if ip, ok := sanitizeGrepPattern(req.Ip); ok {
		return ip, collectErrorSecondaryPatterns(req, "ip")
	}
	if keyword, ok := sanitizeGrepPattern(req.Keyword); ok {
		return keyword, collectErrorSecondaryPatterns(req, "keyword")
	}
	return "", nil
}

func collectErrorSecondaryPatterns(req *logquerydto.CMSQueryErrorLogsReq, skipField string) []string {
	patterns := make([]string, 0, 3)
	if skipField != "url" {
		if url, ok := sanitizeGrepPattern(req.Url); ok {
			patterns = append(patterns, url)
		}
	}
	if skipField != "ip" {
		if ip, ok := sanitizeGrepPattern(req.Ip); ok {
			patterns = append(patterns, ip)
		}
	}
	if skipField != "keyword" {
		if keyword, ok := sanitizeGrepPattern(req.Keyword); ok {
			patterns = append(patterns, keyword)
		}
	}
	return uniqueNonEmptyStrings(patterns)
}

func uniqueNonEmptyStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	ret := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		ret = append(ret, value)
	}
	return ret
}

func filterAllowedExistingLogFiles(filePaths []string, allowedDirs ...string) []string {
	ret := make([]string, 0, len(filePaths))
	for _, filePath := range filePaths {
		cleanPath := filepath.Clean(filePath)
		if !isPathUnderAllowedDirs(cleanPath, allowedDirs...) {
			continue
		}
		if stat, err := os.Stat(cleanPath); err != nil || stat.IsDir() {
			continue
		}
		ret = append(ret, cleanPath)
	}
	return ret
}

func isPathUnderAllowedDirs(filePath string, allowedDirs ...string) bool {
	filePath = filepath.Clean(filePath)
	for _, dir := range allowedDirs {
		dir = filepath.Clean(strings.TrimSpace(dir))
		if dir == "" {
			continue
		}
		if filePath == dir {
			return false
		}
		if strings.HasPrefix(filePath, dir+string(os.PathSeparator)) {
			return true
		}
	}
	return false
}

func grepLogLines(filePaths []string, primary string, secondary []string, maxLines int, useCat bool) ([]string, error) {
	if useCat && isCatQueryEnabled() && len(filePaths) > 0 {
		lines, err := grepWithCatPipe(filePaths, primary, secondary, maxLines)
		if err == nil {
			return lines, nil
		}
	}
	return grepFilesWithPatterns(primary, secondary, filePaths, maxLines)
}

func grepWithCatPipe(filePaths []string, primary string, extraPatterns []string, maxLines int) ([]string, error) {
	primary, ok := sanitizeGrepPattern(primary)
	if !ok || len(filePaths) == 0 {
		return nil, nil
	}
	if maxLines <= 0 {
		maxLines = grepMaxOutputLines
	}

	catCmd := exec.Command("cat", filePaths...)
	grepArgs := []string{"-i", "-F", "-a", "-m", itoa(maxLines), "--", primary}
	grepCmd := exec.Command("grep", grepArgs...)

	catStdout, err := catCmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	grepCmd.Stdin = catStdout

	var grepStderr bytes.Buffer
	grepCmd.Stderr = &grepStderr

	if err := catCmd.Start(); err != nil {
		return nil, err
	}

	out, err := grepCmd.Output()
	catErr := catCmd.Wait()
	if err != nil {
		if isGrepNoMatch(err) {
			return nil, catErr
		}
		return nil, err
	}
	if catErr != nil {
		return nil, catErr
	}

	lines := splitOutputLines(out, maxLines)
	for _, pattern := range extraPatterns {
		pattern, ok := sanitizeGrepPattern(pattern)
		if !ok {
			continue
		}
		lines = filterLinesContains(lines, pattern)
		if len(lines) == 0 {
			return nil, nil
		}
	}
	if len(lines) > maxLines {
		lines = lines[:maxLines]
	}
	return lines, nil
}

func grepFilesOr(firstPattern string, extraPatterns []string, filePaths []string, maxLines int) ([]string, error) {
	firstPattern, ok := sanitizeGrepPattern(firstPattern)
	if !ok || len(filePaths) == 0 {
		return nil, nil
	}
	if maxLines <= 0 {
		maxLines = grepMaxOutputLines
	}

	args := []string{"-h", "-i", "-F", "-a", "-m", itoa(maxLines), "--", firstPattern}
	args = append(args, filePaths...)

	out, err := exec.Command("grep", args...).Output()
	if err != nil {
		if isGrepNoMatch(err) {
			return nil, nil
		}
		return nil, err
	}

	lines := splitOutputLines(out, maxLines)
	for _, pattern := range extraPatterns {
		pattern, ok := sanitizeGrepPattern(pattern)
		if !ok {
			continue
		}
		lines = filterLinesContains(lines, pattern)
		if len(lines) == 0 {
			return nil, nil
		}
	}
	if len(lines) > maxLines {
		lines = lines[:maxLines]
	}
	return lines, nil
}

func grepFilesWithPatterns(primary string, secondary []string, filePaths []string, maxLines int) ([]string, error) {
	primary, ok := sanitizeGrepPattern(primary)
	if !ok || len(filePaths) == 0 {
		return nil, nil
	}
	return grepFilesOr(primary, uniqueNonEmptyStrings(secondary), filePaths, maxLines)
}

func isGrepNoMatch(err error) bool {
	exitErr, ok := err.(*exec.ExitError)
	return ok && exitErr.ExitCode() == 1
}

func splitOutputLines(out []byte, maxLines int) []string {
	if len(out) == 0 {
		return nil
	}
	lines := make([]string, 0, 128)
	start := 0
	for i, b := range out {
		if b != grepLineSeparator {
			continue
		}
		line := strings.TrimRight(string(out[start:i]), "\r")
		if line != "" {
			lines = append(lines, line)
			if len(lines) >= maxLines {
				return lines
			}
		}
		start = i + 1
	}
	if start < len(out) {
		line := strings.TrimRight(string(out[start:]), "\r")
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func filterLinesContains(lines []string, pattern string) []string {
	needle := strings.ToLower(strings.TrimSpace(pattern))
	if needle == "" {
		return lines
	}
	ret := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), needle) {
			ret = append(ret, line)
		}
	}
	return ret
}

func itoa(value int) string {
	if value <= 0 {
		return "0"
	}
	var buf [16]byte
	i := len(buf)
	for value > 0 {
		i--
		buf[i] = byte('0' + value%10)
		value /= 10
	}
	return string(buf[i:])
}
