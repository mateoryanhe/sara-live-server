package logquery

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func runCommand(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.Output()
}

func runCommandAllowNoMatch(name string, args ...string) ([]byte, error) {
	out, err := runCommand(name, args...)
	if err == nil {
		return out, nil
	}
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
		return nil, nil
	}
	return out, err
}

func writeFile(path string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0o644)
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	_, err := runCommand("cp", "-f", src, dst)
	return err
}

func removeFile(path string) {
	_ = os.Remove(path)
}

func countLines(path string) (int, error) {
	out, err := runCommand("wc", "-l", path)
	if err != nil {
		return 0, err
	}
	fields := strings.Fields(string(out))
	if len(fields) == 0 {
		return 0, nil
	}
	return strconv.Atoi(fields[0])
}

func sanitizeShellPattern(pattern string) (string, bool) {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		return "", false
	}
	if strings.ContainsAny(pattern, "\x00\r\n") {
		return "", false
	}
	return pattern, true
}

func concatFilesToFile(files []string, maxLines int, outPath string) error {
	if len(files) == 0 {
		return writeFile(outPath, nil)
	}
	args := append([]string{"--"}, files...)
	out, err := runCommand("cat", args...)
	if err != nil {
		return err
	}
	lines := splitLines(out)
	if len(lines) > maxLines {
		lines = lines[:maxLines]
	}
	if len(lines) == 0 {
		return writeFile(outPath, nil)
	}
	return writeFile(outPath, []byte(strings.Join(lines, "\n")+"\n"))
}

func grepFilesToFile(patterns []string, files []string, maxLines int, outPath string) error {
	if len(files) == 0 {
		return writeFile(outPath, nil)
	}
	patterns = uniqueNonEmpty(patterns)
	if len(patterns) == 0 {
		return concatFilesToFile(files, maxLines, outPath)
	}
	first, ok := sanitizeShellPattern(patterns[0])
	if !ok {
		return writeFile(outPath, nil)
	}

	args := []string{"-h", "-i", "-F", "-a", "-m", strconv.Itoa(maxLines), "--", first}
	args = append(args, files...)
	out, err := runCommandAllowNoMatch("grep", args...)
	if err != nil {
		return err
	}
	lines := splitLines(out)
	for _, pattern := range patterns[1:] {
		pattern, ok := sanitizeShellPattern(pattern)
		if !ok {
			continue
		}
		lines = filterContains(lines, pattern)
	}
	if len(lines) > maxLines {
		lines = lines[:maxLines]
	}
	return writeFile(outPath, []byte(strings.Join(lines, "\n")+"\n"))
}

func paginateFile(srcPath, dstPath string, pageIndex, pageSize int) error {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 50
	}
	start := (pageIndex-1)*pageSize + 1
	tailOut, err := runCommand("tail", "-n", "+"+strconv.Itoa(start), srcPath)
	if err != nil {
		return writeFile(dstPath, nil)
	}
	lines := splitLines(tailOut)
	if len(lines) > pageSize {
		lines = lines[:pageSize]
	}
	if len(lines) == 0 {
		return writeFile(dstPath, nil)
	}
	return writeFile(dstPath, []byte(strings.Join(lines, "\n")+"\n"))
}

func splitLines(data []byte) []string {
	if len(data) == 0 {
		return nil
	}
	raw := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	ret := make([]string, 0, len(raw))
	for _, line := range raw {
		line = strings.TrimRight(line, "\r")
		if line != "" {
			ret = append(ret, line)
		}
	}
	return ret
}

func filterContains(lines []string, pattern string) []string {
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

func uniqueNonEmpty(values []string) []string {
	seen := map[string]struct{}{}
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

func buildTracePatterns(traceId string) []string {
	traceId = strings.TrimSpace(traceId)
	if traceId == "" {
		return nil
	}
	if strings.HasPrefix(traceId, "{") {
		return []string{traceId}
	}
	return []string{"{" + traceId + "}", traceId}
}

func buildDetailPatterns(traceId, reqId, authId, url, keyword string) []string {
	return buildCommonPatterns(traceId, reqId, authId, url, keyword)
}

func buildAccessPatterns(traceId, url, ip string, statusCode int) []string {
	patterns := buildCommonPatterns(traceId, "", "", url, "")
	if ip != "" {
		patterns = append(patterns, ip)
	}
	if statusCode > 0 {
		patterns = append(patterns, strconv.Itoa(statusCode))
	}
	return uniqueNonEmpty(patterns)
}

func buildErrorPatterns(traceId, url, ip, keyword string, statusCode int) []string {
	patterns := buildCommonPatterns(traceId, "", "", url, keyword)
	if ip != "" {
		patterns = append(patterns, ip)
	}
	if statusCode > 0 {
		patterns = append(patterns, strconv.Itoa(statusCode))
	}
	return uniqueNonEmpty(patterns)
}

func buildCommonPatterns(traceId, reqId, authId, url, keyword string) []string {
	var patterns []string
	if traceId = strings.TrimSpace(traceId); traceId != "" {
		if strings.HasPrefix(traceId, "{") {
			patterns = append(patterns, traceId)
		} else {
			patterns = append(patterns, "{"+traceId+"}", traceId)
		}
	}
	if reqId = strings.TrimSpace(reqId); reqId != "" {
		patterns = append(patterns, reqId)
	}
	if authId = strings.TrimSpace(authId); authId != "" {
		patterns = append(patterns, authId)
	}
	if url = strings.TrimSpace(url); url != "" {
		patterns = append(patterns, url)
	}
	if keyword = strings.TrimSpace(keyword); keyword != "" {
		patterns = append(patterns, keyword)
	}
	return uniqueNonEmpty(patterns)
}
