package logquery

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	detailLogRe  = regexp.MustCompile(`^(\S+)\s+\[(\w+)\]\s+\{([a-f0-9]+)\}\s+(.+)$`)
	accessLogRe  = regexp.MustCompile(`^(\S+)\s+\{([a-f0-9]+)\}\s+(\d+)\s+"(\w+)\s+\S+\s+\S+\s+(\S+)\s+HTTP/[\d.]+"\s+([\d.]+),\s+([^,]+),`)
	reqIdRe      = regexp.MustCompile(`reqId=(\d+)`)
	authIdRe     = regexp.MustCompile(`authId=(\d+)`)
	userIdRe     = regexp.MustCompile(`userid=(\d+)`)
	playerRe     = regexp.MustCompile(`玩家=(\d+)`)
	urlRe        = regexp.MustCompile(`url=([^,\s]+)`)
	headersRe    = regexp.MustCompile(`headers=(\{.*\})$`)
	elapsedMsRes = []*regexp.Regexp{
		regexp.MustCompile(`totalMs=(-?\d+)ms`),
		regexp.MustCompile(`handlerMs=(-?\d+)ms`),
		regexp.MustCompile(`writeMs=(-?\d+)ms`),
		regexp.MustCompile(`bodyMs=(-?\d+)ms`),
		regexp.MustCompile(`authMs=(-?\d+)ms`),
		regexp.MustCompile(`afterOutputMs=(-?\d+)ms`),
		regexp.MustCompile(`从队列进入到中间件时间间隔Ms=(-?\d+)ms`),
	}
)

type DetailLogEntry struct {
	Time      string   `json:"time"`
	Level     string   `json:"level"`
	TraceId   string   `json:"traceId"`
	ReqId     string   `json:"reqId"`
	AuthId    string   `json:"authId"`
	Url       string   `json:"url"`
	ElapsedMs *float64 `json:"elapsedMs,omitempty"`
	Message   string   `json:"message"`
	Raw       string   `json:"raw"`
}

type AccessLogEntry struct {
	Time       string  `json:"time"`
	TraceId    string  `json:"traceId"`
	StatusCode int     `json:"statusCode"`
	Method     string  `json:"method"`
	Url        string  `json:"url"`
	HandlerMs  float64 `json:"handlerMs"`
	Ip         string  `json:"ip"`
	UserAgent  string  `json:"userAgent"`
	Raw        string  `json:"raw"`
}

func parseDetailLogLine(line string) (*DetailLogEntry, bool) {
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
		Raw:     line,
	}
	if reqMatch := reqIdRe.FindStringSubmatch(message); len(reqMatch) > 1 {
		entry.ReqId = reqMatch[1]
	}
	entry.AuthId = extractAuthIdFromMessage(message)
	if entry.ReqId == "" || entry.AuthId == "" {
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

func extractElapsedMsFromMessage(message string) (float64, bool) {
	for _, re := range elapsedMsRes {
		if match := re.FindStringSubmatch(message); len(match) > 1 {
			value, err := strconv.ParseFloat(match[1], 64)
			if err != nil {
				continue
			}
			return value, true
		}
	}
	return 0, false
}

func extractAuthIdFromMessage(message string) string {
	if authMatch := authIdRe.FindStringSubmatch(message); len(authMatch) > 1 {
		return authMatch[1]
	}
	if userMatch := userIdRe.FindStringSubmatch(message); len(userMatch) > 1 {
		return userMatch[1]
	}
	if playerMatch := playerRe.FindStringSubmatch(message); len(playerMatch) > 1 {
		return playerMatch[1]
	}
	return ""
}

func extractIdsFromLogHeaders(message string) (reqId, authId string) {
	m := headersRe.FindStringSubmatch(message)
	if len(m) < 2 {
		return "", ""
	}
	var headers map[string][]string
	if err := json.Unmarshal([]byte(m[1]), &headers); err != nil {
		return "", ""
	}
	for key, values := range headers {
		if len(values) == 0 || values[0] == "" {
			continue
		}
		switch strings.ToLower(key) {
		case "reqid":
			if reqId == "" {
				reqId = values[0]
			}
		case "authorization":
			if authId == "" {
				authId = strings.SplitN(values[0], ".", 2)[0]
			}
		case "authid":
			if authId == "" {
				authId = values[0]
			}
		}
	}
	return reqId, authId
}

func messageMatchesAuthId(entry *DetailLogEntry, authId string) bool {
	if entry == nil || authId == "" {
		return false
	}
	if fuzzyMatch(entry.AuthId, authId) {
		return true
	}
	if fuzzyMatch(entry.Message, "authId="+authId) {
		return true
	}
	if fuzzyMatch(entry.Message, "userid="+authId) {
		return true
	}
	if fuzzyMatch(entry.Message, "玩家="+authId) {
		return true
	}
	return fuzzyMatch(entry.Raw, authId)
}

func parseAccessLogLine(line string) (*AccessLogEntry, bool) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, false
	}
	m := accessLogRe.FindStringSubmatch(line)
	if len(m) < 8 {
		return nil, false
	}
	statusCode, _ := strconv.Atoi(m[3])
	handlerSec, _ := strconv.ParseFloat(m[6], 64)
	entry := &AccessLogEntry{
		Time:       m[1],
		TraceId:    m[2],
		StatusCode: statusCode,
		Method:     m[4],
		Url:        m[5],
		HandlerMs:  handlerSec * 1000,
		Ip:         strings.TrimSpace(m[7]),
		UserAgent:  extractUserAgent(line),
		Raw:        line,
	}
	return entry, true
}

func parseAccessLogTime(timeStr string) (time.Time, bool) {
	return parseLogTime(timeStr)
}

func parseLogTime(timeStr string) (time.Time, bool) {
	layouts := []string{
		"2006-01-02T15:04:05.000",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05.000",
		"2006-01-02 15:04:05",
	}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, timeStr, time.Local); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func fuzzyMatch(value, needle string) bool {
	if needle == "" {
		return true
	}
	if value == "" {
		return false
	}
	return strings.Contains(strings.ToLower(value), strings.ToLower(needle))
}

func matchTraceId(entryTraceId, reqTraceId string) bool {
	if reqTraceId == "" {
		return false
	}
	if entryTraceId == reqTraceId {
		return true
	}
	return fuzzyMatch(entryTraceId, reqTraceId)
}

func extractUserAgent(line string) string {
	lastQuote := strings.LastIndex(line, `"`)
	if lastQuote <= 0 {
		return ""
	}
	prevQuote := strings.LastIndex(line[:lastQuote], `"`)
	if prevQuote < 0 || prevQuote+1 >= lastQuote {
		return ""
	}
	return line[prevQuote+1 : lastQuote]
}

func matchDetailEntry(entry *DetailLogEntry, traceId, reqId, authId, url, keyword string) bool {
	if entry == nil {
		return false
	}
	if traceId != "" && !matchTraceId(entry.TraceId, traceId) {
		return false
	}
	if reqId != "" && !fuzzyMatch(entry.ReqId, reqId) && !fuzzyMatch(entry.Message, reqId) {
		return false
	}
	if authId != "" && !messageMatchesAuthId(entry, authId) {
		return false
	}
	if url != "" && !fuzzyMatch(entry.Url, url) && !fuzzyMatch(entry.Message, url) {
		return false
	}
	if keyword != "" && !fuzzyMatch(entry.Raw, keyword) {
		return false
	}
	return true
}

func matchAccessEntry(entry *AccessLogEntry, traceId, url, ip string, statusCode int, minHandlerMs, maxHandlerMs *float64) bool {
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
	if minHandlerMs != nil && entry.HandlerMs < *minHandlerMs {
		return false
	}
	if maxHandlerMs != nil && entry.HandlerMs > *maxHandlerMs {
		return false
	}
	return true
}
