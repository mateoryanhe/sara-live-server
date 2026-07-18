package logquery

import (
	"strings"
	"time"

	"xr-game-server/dto/logquerydto"
)

const (
	maxDetailMatchLines                = 2000
	maxDetailQueryLineBytes            = 32 * 1024
	maxDetailQueryLineBytesWithKeyword = 256 * 1024
	maxDetailLogListMessageBytes       = 512
)

type detailQueryFilter struct {
	traceId    string
	reqId      string
	authId     string
	url        string
	keyword    string
	rangeStart time.Time
	rangeEnd   time.Time
}

func buildDetailQueryFilter(req *logquerydto.CMSQueryDetailLogsReq, rangeStart, rangeEnd time.Time) detailQueryFilter {
	return detailQueryFilter{
		traceId:    strings.TrimSpace(req.TraceId),
		reqId:      strings.TrimSpace(req.ReqId),
		authId:     strings.TrimSpace(req.AuthId),
		url:        strings.TrimSpace(req.Url),
		keyword:    strings.TrimSpace(req.Keyword),
		rangeStart: rangeStart,
		rangeEnd:   rangeEnd,
	}
}

func maxDetailQueryLineBytesForReq(req *logquerydto.CMSQueryDetailLogsReq) int {
	if strings.TrimSpace(req.Keyword) != "" {
		return maxDetailQueryLineBytesWithKeyword
	}
	return maxDetailQueryLineBytes
}

func (f detailQueryFilter) quickAcceptLine(line string) bool {
	if line == "" {
		return false
	}
	if logTime, ok := quickLogTimeFromLine(line); ok {
		if logTime.Before(f.rangeStart) || logTime.After(f.rangeEnd) {
			return false
		}
	}
	if f.traceId != "" && !lineMatchesTraceIdFilter(line, f.traceId) {
		return false
	}
	if f.reqId != "" && !fuzzyMatch(line, f.reqId) {
		return false
	}
	if f.authId != "" && !detailLineMightMatchAuthId(line, f.authId) {
		return false
	}
	if f.url != "" && !fuzzyMatch(line, f.url) {
		return false
	}
	if f.keyword != "" && !fuzzyMatch(line, f.keyword) {
		return false
	}
	return true
}

func quickLogTimeFromLine(line string) (time.Time, bool) {
	space := strings.IndexByte(line, ' ')
	if space <= 0 {
		return time.Time{}, false
	}
	return parseLogTime(line[:space])
}

func isLikelyDetailLogLine(line string) bool {
	if len(line) < 24 {
		return false
	}
	space := strings.IndexByte(line, ' ')
	if space <= 0 || space+1 >= len(line) || line[space+1] != '[' {
		return false
	}
	return strings.IndexByte(line[space:], '{') > 0
}

func lineMatchesTraceIdFilter(line, traceId string) bool {
	traceId = strings.TrimSpace(traceId)
	if traceId == "" {
		return true
	}
	start := strings.IndexByte(line, '{')
	if start < 0 {
		return false
	}
	end := strings.IndexByte(line[start+1:], '}')
	if end < 0 {
		return false
	}
	return matchTraceId(line[start+1:start+1+end], traceId)
}

func detailLineMightMatchAuthId(line, authId string) bool {
	if authId == "" {
		return true
	}
	lowerLine := strings.ToLower(line)
	needle := strings.ToLower(authId)
	if strings.Contains(lowerLine, needle) {
		return true
	}
	if strings.Contains(lowerLine, "authid="+needle) {
		return true
	}
	if strings.Contains(lowerLine, "userid="+needle) {
		return true
	}
	return strings.Contains(lowerLine, "玩家="+needle)
}

func parseDetailLogLineOpt(line string) (*DetailLogEntry, bool) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, false
	}
	m := detailLogRe.FindStringSubmatch(line)
	if len(m) < 5 {
		return nil, false
	}
	message := strings.TrimSpace(m[4])
	entry := &DetailLogEntry{
		Time:    m[1],
		Level:   m[2],
		TraceId: m[3],
		Message: message,
	}
	if reqMatch := reqIdRe.FindStringSubmatch(message); len(reqMatch) > 1 {
		entry.ReqId = reqMatch[1]
	}
	entry.AuthId = extractAuthIdFromMessage(message)
	if (entry.ReqId == "" || entry.AuthId == "") && strings.Contains(message, "headers=") {
		headerReqId, headerAuthId := extractIdsFromLogHeaders(message)
		if entry.ReqId == "" {
			entry.ReqId = headerReqId
		}
		if entry.AuthId == "" {
			entry.AuthId = headerAuthId
		}
	}
	if urlMatch := urlRe.FindStringSubmatch(message); len(urlMatch) > 1 {
		entry.Url = urlMatch[1]
	}
	if elapsedMs, ok := extractElapsedMsFromMessage(message); ok {
		entry.ElapsedMs = &elapsedMs
	}
	return entry, true
}

func prepareDetailLogEntryForList(entry *DetailLogEntry) *DetailLogEntry {
	if entry == nil {
		return nil
	}
	entry.Raw = ""
	if len(entry.Message) > maxDetailLogListMessageBytes {
		entry.Message = entry.Message[:maxDetailLogListMessageBytes] + detailLogPayloadTruncated
	}
	return entry
}

func queryDetailLogsMatched(req *logquerydto.CMSQueryDetailLogsReq, rangeStart, rangeEnd time.Time) []*DetailLogEntry {
	if matched, ok := queryDetailLogsMatchedWithGrep(req, rangeStart, rangeEnd); ok {
		return matched
	}
	return queryDetailLogsMatchedWithScan(req, rangeStart, rangeEnd)
}

func queryDetailLogsMatchedWithGrep(req *logquerydto.CMSQueryDetailLogsReq, rangeStart, rangeEnd time.Time) ([]*DetailLogEntry, bool) {
	if !isGrepQueryEnabled() {
		return nil, false
	}
	primary, secondary := buildDetailGrepPatterns(req)
	if primary == "" {
		return nil, false
	}

	files := globLogFilesForRange(logGlobDetail, req.StartDate, req.EndDate)
	if len(files) == 0 {
		return nil, true
	}

	useCat := isFuzzyGrepQuery(req.TraceId)
	lines, err := grepLogLines(files, primary, secondary, grepMaxOutputLines, useCat)
	if err != nil {
		return nil, false
	}
	return parseDetailLogLinesForQuery(lines, req, rangeStart, rangeEnd), true
}

func queryDetailLogsMatchedWithScan(req *logquerydto.CMSQueryDetailLogsReq, rangeStart, rangeEnd time.Time) []*DetailLogEntry {
	filter := buildDetailQueryFilter(req, rangeStart, rangeEnd)
	maxLineBytes := maxDetailQueryLineBytesForReq(req)

	matched := make([]*DetailLogEntry, 0, 128)
	for _, filePath := range listDetailLogFilesForRange(req.StartDate, req.EndDate) {
		_ = scanLogFileWithMaxLineBytes(filePath, maxLineBytes, func(line string) bool {
			if !filter.quickAcceptLine(line) {
				return true
			}
			if !isLikelyDetailLogLine(line) {
				return true
			}
			entry, ok := parseDetailLogLineOpt(line)
			if !ok {
				return true
			}
			if !logTimeInRange(entry.Time, rangeStart, rangeEnd) {
				return true
			}
			if !matchDetailEntryWithRaw(entry, line, req.TraceId, req.ReqId, req.AuthId, req.Url, req.Keyword) {
				return true
			}
			entry.Message = truncateLogFieldValue(entry.Message, "respContent=", maxDetailLogRespContentBytes)
			matched = append(matched, prepareDetailLogEntryForList(entry))
			return true
		})
	}
	if len(matched) > maxDetailMatchLines {
		sortDetailLogsByTimeDesc(matched)
		matched = matched[:maxDetailMatchLines]
	}
	return matched
}

func parseDetailLogLinesForQuery(lines []string, req *logquerydto.CMSQueryDetailLogsReq, rangeStart, rangeEnd time.Time) []*DetailLogEntry {
	matched := make([]*DetailLogEntry, 0, len(lines))
	for _, line := range lines {
		if !isLikelyDetailLogLine(line) {
			continue
		}
		entry, ok := parseDetailLogLineOpt(line)
		if !ok {
			continue
		}
		if !logTimeInRange(entry.Time, rangeStart, rangeEnd) {
			continue
		}
		if !matchDetailEntryWithRaw(entry, line, req.TraceId, req.ReqId, req.AuthId, req.Url, req.Keyword) {
			continue
		}
		entry.Message = truncateLogFieldValue(entry.Message, "respContent=", maxDetailLogRespContentBytes)
		matched = append(matched, prepareDetailLogEntryForList(entry))
	}
	if len(matched) > maxDetailMatchLines {
		matched = matched[:maxDetailMatchLines]
	}
	sortDetailLogsByTimeDesc(matched)
	return matched
}

func matchDetailEntryWithRaw(entry *DetailLogEntry, raw, traceId, reqId, authId, url, keyword string) bool {
	if entry == nil {
		return false
	}
	if traceId != "" && !matchTraceId(entry.TraceId, traceId) {
		return false
	}
	if reqId != "" && !fuzzyMatch(entry.ReqId, reqId) && !fuzzyMatch(entry.Message, reqId) {
		return false
	}
	if authId != "" && !messageMatchesAuthId(entry, authId) && !fuzzyMatch(raw, authId) {
		return false
	}
	if url != "" && !fuzzyMatch(entry.Url, url) && !fuzzyMatch(entry.Message, url) {
		return false
	}
	if keyword != "" && !fuzzyMatch(raw, keyword) {
		return false
	}
	return true
}
