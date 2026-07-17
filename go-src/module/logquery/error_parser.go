package logquery

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	errorLogHeaderRe = regexp.MustCompile(`^(\S+)\s+\[(\w+)\]\s+\{([a-f0-9]+)\}\s+(\d+)\s+"(\w+)\s+\S+\s+\S+\s+(\S+)\s+HTTP/[\d.]+"\s+([\d.]+),\s+([^,]+),`)
	errorMetaRe      = regexp.MustCompile(`,\s*(-?\d+),\s*"([^"]*)"(?:,\s*(.*))?$`)
	errorSourceRe    = regexp.MustCompile(`(?:^|[,\s])source=(\S+)`)
	errorIpRe        = regexp.MustCompile(`(?:^|[,\s])ip=([^,\s]+)`)
)

const errorLogTag = "ErrorLog"

type ErrorLogEntry struct {
	Time         string  `json:"time"`
	Level        string  `json:"level"`
	TraceId      string  `json:"traceId"`
	StatusCode   int     `json:"statusCode"`
	Method       string  `json:"method"`
	Url          string  `json:"url"`
	HandlerMs    float64 `json:"handlerMs"`
	Ip           string  `json:"ip"`
	AuthId       string  `json:"authId"`
	ErrorCode    int     `json:"errorCode"`
	ErrorMessage string  `json:"errorMessage"`
	Detail       string  `json:"detail"`
	Stack        string  `json:"stack"`
	Raw          string  `json:"raw"`
}

func parseErrorLogHeader(line string) (*ErrorLogEntry, bool) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, false
	}
	m := errorLogHeaderRe.FindStringSubmatch(line)
	if len(m) < 9 {
		return nil, false
	}
	statusCode, _ := strconv.Atoi(m[4])
	handlerSec, _ := strconv.ParseFloat(m[7], 64)
	entry := &ErrorLogEntry{
		Time:       m[1],
		Level:      m[2],
		TraceId:    m[3],
		StatusCode: statusCode,
		Method:     m[5],
		Url:        m[6],
		HandlerMs:  handlerSec * 1000,
		Ip:         strings.TrimSpace(m[8]),
	}
	fillErrorLogMeta(entry, line)
	return entry, true
}

func fillErrorLogMeta(entry *ErrorLogEntry, line string) {
	if entry == nil {
		return
	}
	meta := errorMetaRe.FindStringSubmatch(line)
	if len(meta) < 3 {
		return
	}
	entry.ErrorCode, _ = strconv.Atoi(meta[1])
	entry.ErrorMessage = meta[2]
	if len(meta) >= 4 {
		entry.Detail = strings.TrimSpace(meta[3])
	}
}

func finalizeErrorLogEntry(entry *ErrorLogEntry, body string) {
	if entry == nil {
		return
	}
	entry.Raw = body
	if stackIdx := strings.Index(body, "\nStack:\n"); stackIdx >= 0 {
		entry.Stack = strings.TrimSpace(body[stackIdx+len("\nStack:\n"):])
		if entry.Detail == "" {
			entry.Detail = strings.TrimSpace(body[:stackIdx])
		}
		fillErrorLogAuthId(entry)
		return
	}
	if stackIdx := strings.Index(body, "Stack:\n"); stackIdx >= 0 {
		entry.Stack = strings.TrimSpace(body[stackIdx+len("Stack:\n"):])
	}
	fillErrorLogAuthId(entry)
}

func fillErrorLogAuthId(entry *ErrorLogEntry) {
	if entry == nil || entry.AuthId != "" {
		return
	}
	entry.AuthId = extractAuthIdFromMessage(entry.Detail)
	if entry.AuthId == "" {
		entry.AuthId = extractAuthIdFromMessage(entry.Raw)
	}
}

func matchErrorEntry(entry *ErrorLogEntry, traceId, url, ip, keyword string, statusCode int) bool {
	if entry == nil {
		return false
	}
	if traceId != "" && !matchTraceId(entry.TraceId, traceId) {
		return false
	}
	if url != "" && !fuzzyMatch(entry.Url, url) && !fuzzyMatch(entry.Raw, url) {
		return false
	}
	if ip != "" && !fuzzyMatch(entry.Ip, ip) {
		return false
	}
	if statusCode > 0 && entry.StatusCode != statusCode {
		return false
	}
	if keyword != "" && !fuzzyMatch(entry.Raw, keyword) &&
		!fuzzyMatch(entry.ErrorMessage, keyword) &&
		!fuzzyMatch(entry.Detail, keyword) &&
		!fuzzyMatch(entry.Stack, keyword) &&
		!fuzzyMatch(entry.Url, keyword) {
		return false
	}
	return true
}

func scanErrorLogFile(filePath string, fn func(entry *ErrorLogEntry) bool) error {
	file, err := openLogFile(filePath)
	if err != nil {
		return err
	}
	if file == nil {
		return nil
	}
	defer file.Close()

	scanner := newLogScanner(file)
	var current *ErrorLogEntry
	var body strings.Builder

	flush := func() bool {
		if current == nil {
			return true
		}
		finalizeErrorLogEntry(current, body.String())
		return fn(current)
	}

	for scanner.Scan() {
		line := scanner.Text()
		if header, ok := parseErrorLogHeader(line); ok {
			if !flush() {
				break
			}
			current = header
			body.Reset()
			body.WriteString(line)
			continue
		}
		if current == nil {
			continue
		}
		body.WriteByte('\n')
		body.WriteString(line)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	flush()
	return nil
}

func isDetailErrorEntry(entry *DetailLogEntry) bool {
	if entry == nil {
		return false
	}
	if strings.EqualFold(entry.Level, "ERRO") {
		return true
	}
	return strings.Contains(entry.Message, errorLogTag)
}

func detailEntryToErrorEntry(entry *DetailLogEntry) *ErrorLogEntry {
	if !isDetailErrorEntry(entry) {
		return nil
	}
	errEntry := &ErrorLogEntry{
		Time:    entry.Time,
		Level:   entry.Level,
		TraceId: entry.TraceId,
		Url:     entry.Url,
		Ip:      extractIpFromMessage(entry.Message),
		AuthId:  entry.AuthId,
		Detail:  entry.Message,
		Raw:     entry.Raw,
		Stack:   extractStackFromMessage(entry.Message),
	}
	if strings.Contains(entry.Message, errorLogTag) {
		errEntry.ErrorMessage = errorLogTag
		if sourceMatch := errorSourceRe.FindStringSubmatch(entry.Message); len(sourceMatch) > 1 {
			errEntry.ErrorMessage = errorLogTag + "/" + sourceMatch[1]
		}
	} else {
		errEntry.ErrorMessage = entry.Level
	}
	if entry.ElapsedMs != nil {
		errEntry.HandlerMs = *entry.ElapsedMs
	}
	return errEntry
}

func extractStackFromMessage(message string) string {
	const prefix = "stack="
	idx := strings.Index(message, prefix)
	if idx < 0 {
		return ""
	}
	return strings.TrimSpace(message[idx+len(prefix):])
}

func extractIpFromMessage(message string) string {
	if ipMatch := errorIpRe.FindStringSubmatch(message); len(ipMatch) > 1 {
		return ipMatch[1]
	}
	return ""
}

func scanDetailErrorLogFile(filePath string, fn func(entry *ErrorLogEntry) bool) error {
	return scanLogFile(filePath, func(line string) bool {
		detail, ok := parseDetailLogLine(line)
		if !ok {
			return true
		}
		errEntry := detailEntryToErrorEntry(detail)
		if errEntry == nil {
			return true
		}
		return fn(errEntry)
	})
}

func listDetailLogFilesForErrors() []string {
	paths := loadLogPaths()
	dirs := listDetailLogDirs()
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

func scanAllErrorLogEntries(fn func(entry *ErrorLogEntry) bool) {
	continueScan := true
	wrap := func(entry *ErrorLogEntry) bool {
		if !continueScan {
			return false
		}
		if !fn(entry) {
			continueScan = false
			return false
		}
		return true
	}
	for _, filePath := range listErrorLogFiles() {
		_ = scanErrorLogFile(filePath, wrap)
		if !continueScan {
			return
		}
	}
	for _, filePath := range listDetailLogFilesForErrors() {
		_ = scanDetailErrorLogFile(filePath, wrap)
		if !continueScan {
			return
		}
	}
}
