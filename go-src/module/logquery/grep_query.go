package logquery

import (
	"strings"
	"time"

	"xr-game-server/dto/logquerydto"
)

func queryAccessLogsMatched(req *logquerydto.CMSQueryAccessLogsReq, rangeStart, rangeEnd time.Time) []*AccessLogEntry {
	if matched, ok := queryAccessLogsMatchedWithGrep(req, rangeStart, rangeEnd); ok {
		return matched
	}
	return queryAccessLogsMatchedWithScan(req, rangeStart, rangeEnd)
}

func queryAccessLogsMatchedWithGrep(req *logquerydto.CMSQueryAccessLogsReq, rangeStart, rangeEnd time.Time) ([]*AccessLogEntry, bool) {
	if !isGrepQueryEnabled() {
		return nil, false
	}
	primary, secondary := buildAccessGrepPatterns(req)
	if primary == "" {
		return nil, false
	}

	files := globLogFilesForRange(logGlobAccess, req.StartDate, req.EndDate)
	if len(files) == 0 {
		return nil, true
	}

	useCat := isFuzzyGrepQuery(req.TraceId)
	lines, err := grepLogLines(files, primary, secondary, grepMaxOutputLines, useCat)
	if err != nil {
		return nil, false
	}
	return parseAccessLogLinesForQuery(lines, req, rangeStart, rangeEnd), true
}

func queryAccessLogsMatchedWithScan(req *logquerydto.CMSQueryAccessLogsReq, rangeStart, rangeEnd time.Time) []*AccessLogEntry {
	var matched []*AccessLogEntry
	for _, filePath := range listAccessLogFilesForRange(req.StartDate, req.EndDate) {
		_ = scanAccessLogFile(filePath, func(entry *AccessLogEntry) bool {
			if !logTimeInRange(entry.Time, rangeStart, rangeEnd) {
				return true
			}
			if !matchAccessEntry(entry, req.TraceId, req.Url, req.Ip, req.StatusCode, req.MinHandlerMs, req.MaxHandlerMs) {
				return true
			}
			matched = append(matched, entry)
			return true
		})
	}
	sortAccessLogsByTimeDesc(matched)
	if len(matched) > maxMatchLines {
		matched = matched[:maxMatchLines]
	}
	return matched
}

func parseAccessLogLinesForQuery(lines []string, req *logquerydto.CMSQueryAccessLogsReq, rangeStart, rangeEnd time.Time) []*AccessLogEntry {
	matched := make([]*AccessLogEntry, 0, len(lines))
	for _, line := range lines {
		entry, ok := parseAccessLogLine(line)
		if !ok {
			continue
		}
		if !logTimeInRange(entry.Time, rangeStart, rangeEnd) {
			continue
		}
		if !matchAccessEntry(entry, req.TraceId, req.Url, req.Ip, req.StatusCode, req.MinHandlerMs, req.MaxHandlerMs) {
			continue
		}
		matched = append(matched, entry)
	}
	sortAccessLogsByTimeDesc(matched)
	if len(matched) > maxMatchLines {
		matched = matched[:maxMatchLines]
	}
	return matched
}

func queryErrorLogsMatched(req *logquerydto.CMSQueryErrorLogsReq, rangeStart, rangeEnd time.Time) []*ErrorLogEntry {
	if matched, ok := queryErrorLogsMatchedWithGrep(req, rangeStart, rangeEnd); ok {
		return matched
	}
	return queryErrorLogsMatchedWithScan(req, rangeStart, rangeEnd)
}

func queryErrorLogsMatchedWithGrep(req *logquerydto.CMSQueryErrorLogsReq, rangeStart, rangeEnd time.Time) ([]*ErrorLogEntry, bool) {
	if !isGrepQueryEnabled() {
		return nil, false
	}
	primary, secondary := buildErrorGrepPatterns(req)
	if primary == "" {
		return nil, false
	}

	files := globLogFilesForRange(logGlobError, req.StartDate, req.EndDate)
	if len(files) == 0 {
		return nil, true
	}

	useCat := isFuzzyGrepQuery(req.TraceId)
	lines, err := grepLogLines(files, primary, secondary, grepMaxOutputLines, useCat)
	if err != nil {
		return nil, false
	}
	return parseErrorLogLinesForQuery(lines, req, rangeStart, rangeEnd), true
}

func queryErrorLogsMatchedWithScan(req *logquerydto.CMSQueryErrorLogsReq, rangeStart, rangeEnd time.Time) []*ErrorLogEntry {
	var matched []*ErrorLogEntry
	scanErrorLogEntriesInRange(req.StartDate, req.EndDate, func(entry *ErrorLogEntry) bool {
		if !logTimeInRange(entry.Time, rangeStart, rangeEnd) {
			return true
		}
		if !matchErrorEntry(entry, req.TraceId, req.Url, req.Ip, req.Keyword, req.StatusCode) {
			return true
		}
		matched = append(matched, entry)
		return true
	})
	sortErrorLogsByTimeDesc(matched)
	if len(matched) > maxMatchLines {
		matched = matched[:maxMatchLines]
	}
	return matched
}

func parseErrorLogLinesForQuery(lines []string, req *logquerydto.CMSQueryErrorLogsReq, rangeStart, rangeEnd time.Time) []*ErrorLogEntry {
	matched := make([]*ErrorLogEntry, 0, len(lines))
	for _, line := range lines {
		entry := parseErrorLogLineFromGrep(line)
		if entry == nil {
			continue
		}
		if !logTimeInRange(entry.Time, rangeStart, rangeEnd) {
			continue
		}
		if !matchErrorEntry(entry, req.TraceId, req.Url, req.Ip, req.Keyword, req.StatusCode) {
			continue
		}
		matched = append(matched, entry)
	}
	sortErrorLogsByTimeDesc(matched)
	if len(matched) > maxMatchLines {
		matched = matched[:maxMatchLines]
	}
	return matched
}

func parseErrorLogLineFromGrep(line string) *ErrorLogEntry {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}
	if header, ok := parseErrorLogHeader(line); ok {
		finalizeErrorLogEntry(header, line)
		return header
	}
	return nil
}

func getTraceLogsWithGrep(req *logquerydto.CMSGetTraceLogsReq, rangeStart, rangeEnd time.Time) (*logquerydto.CMSGetTraceLogsRes, bool) {
	if !isGrepQueryEnabled() {
		return nil, false
	}
	traceId := strings.TrimSpace(req.TraceId)
	tracePattern := "{" + traceId + "}"
	if len(traceId) < 16 {
		tracePattern = traceId
	}
	tracePattern, ok := sanitizeGrepPattern(tracePattern)
	if !ok {
		return nil, false
	}

	detailFiles := globLogFilesForRange(logGlobDetail, req.StartDate, req.EndDate)
	accessFiles := globLogFilesForRange(logGlobAccess, req.StartDate, req.EndDate)
	errorFiles := globLogFilesForRange(logGlobError, req.StartDate, req.EndDate)

	useCat := isFuzzyGrepQuery(req.TraceId)
	detailLines, err := grepLogLines(detailFiles, tracePattern, nil, grepMaxOutputLines, useCat)
	if err != nil {
		return nil, false
	}
	accessLines, err := grepLogLines(accessFiles, tracePattern, nil, grepMaxOutputLines, useCat)
	if err != nil {
		return nil, false
	}
	errorLines, err := grepLogLines(errorFiles, tracePattern, nil, grepMaxOutputLines, useCat)
	if err != nil {
		return nil, false
	}

	detailLogs := make([]*DetailLogEntry, 0, len(detailLines))
	for _, line := range detailLines {
		entry, ok := parseDetailLogLine(line)
		if !ok || !logTimeInRange(entry.Time, rangeStart, rangeEnd) || !matchTraceId(entry.TraceId, req.TraceId) {
			continue
		}
		detailLogs = append(detailLogs, trimDetailLogEntryForQuery(entry))
	}

	accessLogs := make([]*AccessLogEntry, 0, len(accessLines))
	for _, line := range accessLines {
		entry, ok := parseAccessLogLine(line)
		if !ok || !logTimeInRange(entry.Time, rangeStart, rangeEnd) || !matchTraceId(entry.TraceId, req.TraceId) {
			continue
		}
		accessLogs = append(accessLogs, entry)
	}

	errorLogs := make([]*ErrorLogEntry, 0, len(errorLines))
	for _, line := range errorLines {
		entry := parseErrorLogLineFromGrep(line)
		if entry == nil || !logTimeInRange(entry.Time, rangeStart, rangeEnd) || !matchTraceId(entry.TraceId, req.TraceId) || isLogQueryRelatedErrorEntry(entry) {
			continue
		}
		errorLogs = append(errorLogs, entry)
	}

	sortDetailLogsByTimeAsc(detailLogs)
	sortAccessLogsByTimeAsc(accessLogs)
	sortErrorLogsByTimeAsc(errorLogs)

	return &logquerydto.CMSGetTraceLogsRes{
		TraceId:    req.TraceId,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		DetailLogs: detailLogs,
		AccessLogs: accessLogs,
		ErrorLogs:  errorLogs,
	}, true
}
