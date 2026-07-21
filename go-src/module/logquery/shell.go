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

func runShellScript(script string) error {
	_, err := exec.Command("bash", "-lc", script).Output()
	return err
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", `'\'\'\'`) + "'"
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

func shellJoinPaths(paths []string) string {
	if len(paths) == 0 {
		return ""
	}
	quoted := make([]string, len(paths))
	for i, path := range paths {
		quoted[i] = shellQuote(path)
	}
	return strings.Join(quoted, " ")
}

func buildSearchPipeline(patterns []string, files []string) string {
	patterns = uniqueNonEmpty(patterns)
	fileArgs := shellJoinPaths(files)
	if len(patterns) == 0 {
		return "cat -- " + fileArgs
	}
	first, ok := sanitizeShellPattern(patterns[0])
	if !ok {
		return "cat -- " + fileArgs
	}
	var builder strings.Builder
	builder.WriteString("grep -h -i -F -a -- ")
	builder.WriteString(shellQuote(first))
	builder.WriteString(" ")
	builder.WriteString(fileArgs)
	for _, pattern := range patterns[1:] {
		pattern, ok := sanitizeShellPattern(pattern)
		if !ok {
			continue
		}
		builder.WriteString(" | grep -F -i -- ")
		builder.WriteString(shellQuote(pattern))
	}
	return builder.String()
}

func writePipelineToFile(pipeline string, maxLines int, outPath string, reverse bool) error {
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	var script strings.Builder
	script.WriteString("( ")
	script.WriteString(pipeline)
	script.WriteString(" || true )")
	if maxLines > 0 {
		script.WriteString(" | tail -n ")
		script.WriteString(strconv.Itoa(maxLines))
	}
	if reverse {
		script.WriteString(" | tac")
	}
	script.WriteString(" > ")
	script.WriteString(shellQuote(outPath))
	return runShellScript(script.String())
}

func concatFilesToFile(files []string, maxLines int, outPath string) error {
	if len(files) == 0 {
		return writeFile(outPath, nil)
	}
	return writePipelineToFile(buildSearchPipeline(nil, files), maxLines, outPath, false)
}

func grepFilesToFile(patterns []string, files []string, maxLines int, outPath string) error {
	if len(files) == 0 {
		return writeFile(outPath, nil)
	}
	return writePipelineToFile(buildSearchPipeline(patterns, files), maxLines, outPath, false)
}

func grepFilesToReversedFile(patterns []string, files []string, maxLines int, outPath string) error {
	if len(files) == 0 {
		return writeFile(outPath, nil)
	}
	return writePipelineToFile(buildSearchPipeline(patterns, files), maxLines, outPath, true)
}

func buildLogContinuationAwkScript(patterns []string) string {
	patterns = uniqueNonEmpty(patterns)
	var builder strings.Builder
	builder.WriteString(`function is_log_header(line) {
  return (line ~ /^[0-9]{4}-[0-9]{2}-[0-9]{2}T/ && line ~ / \[[A-Z]+\] \{/)
}
function matches(line) {
`)
	for _, pattern := range patterns {
		builder.WriteString("  if (index(line, ")
		builder.WriteString(awkQuote(pattern))
		builder.WriteString(") == 0) return 0\n")
	}
	builder.WriteString(`  return 1
}
BEGIN { capturing = 0 }
{
  if (is_log_header($0)) {
    if (matches($0)) {
      print $0
      capturing = 1
    } else {
      capturing = 0
    }
    next
  }
  if (capturing) {
    print $0
  }
}`)
	return builder.String()
}

func awkQuote(value string) string {
	return `"` + strings.ReplaceAll(strings.ReplaceAll(value, `\`, `\\`), `"`, `\"`) + `"`
}

func buildLogContinuationPipeline(patterns []string, files []string) string {
	fileArgs := shellJoinPaths(files)
	return "cat -- " + fileArgs + " | awk " + shellQuote(buildLogContinuationAwkScript(patterns))
}

func grepFilesWithLogContinuationToFile(patterns []string, files []string, maxLines int, outPath string, reverse bool) error {
	patterns = uniqueNonEmpty(patterns)
	if len(files) == 0 {
		return writeFile(outPath, nil)
	}
	if len(patterns) == 0 {
		return writePipelineToFile("cat -- "+shellJoinPaths(files), maxLines, outPath, reverse)
	}
	return writePipelineToFile(buildLogContinuationPipeline(patterns, files), maxLines, outPath, reverse)
}

func paginateFile(srcPath, dstPath string, pageIndex, pageSize int) error {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 50
	}
	if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
		return err
	}
	start := (pageIndex-1)*pageSize + 1
	script := "tail -n +" + strconv.Itoa(start) + " " + shellQuote(srcPath) +
		" | head -n " + strconv.Itoa(pageSize) +
		" > " + shellQuote(dstPath) + " || true"
	return runShellScript(script)
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
