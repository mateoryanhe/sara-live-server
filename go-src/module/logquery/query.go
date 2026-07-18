package logquery

import (
	"context"
	"sort"
	"time"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/logquerydto"
	"xr-game-server/errercode"
)

const (
	maxDateRangeDays = 7
	maxMatchLines    = 5000
	defaultPageSize  = 50
	maxPageSize      = 200
	defaultTopN      = 20
	maxTopN          = 100
)

func GetLogPaths(_ context.Context, _ *logquerydto.CMSGetLogPathsReq) (*logquerydto.CMSGetLogPathsRes, error) {
	paths := loadLogPaths()
	return &logquerydto.CMSGetLogPathsRes{
		ServerTime:       time.Now().Format("2006-01-02 15:04:05.000"),
		DetailLogDir:     paths.DetailLogDir,
		DetailLogPattern: paths.DetailLogPattern,
		AccessLogDir:     paths.AccessLogDir,
		AccessLogPattern: paths.AccessLogPattern,
		ErrorLogDir:      paths.ErrorLogDir,
		ErrorLogPattern:  paths.ErrorLogPattern,
	}, nil
}

func QueryDetailLogs(_ context.Context, req *logquerydto.CMSQueryDetailLogsReq) (*httpserver.CMSQueryResp, error) {
	rangeStart, rangeEnd, err := buildTimeRange(req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}
	pageIndex, pageSize := normalizePage(req.PageIndex, req.PageSize)

	matched := queryDetailLogsMatched(req, rangeStart, rangeEnd)
	total := len(matched)
	start := (pageIndex - 1) * pageSize
	if start >= total {
		return httpserver.NewCMSQueryResp(total, []*DetailLogEntry{}), nil
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return httpserver.NewCMSQueryResp(total, matched[start:end]), nil
}

func QueryAccessLogs(_ context.Context, req *logquerydto.CMSQueryAccessLogsReq) (*httpserver.CMSQueryResp, error) {
	rangeStart, rangeEnd, err := buildTimeRange(req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}
	pageIndex, pageSize := normalizePage(req.PageIndex, req.PageSize)

	matched := queryAccessLogsMatched(req, rangeStart, rangeEnd)
	total := len(matched)
	start := (pageIndex - 1) * pageSize
	if start >= total {
		return httpserver.NewCMSQueryResp(total, []*AccessLogEntry{}), nil
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return httpserver.NewCMSQueryResp(total, matched[start:end]), nil
}

func QueryErrorLogs(_ context.Context, req *logquerydto.CMSQueryErrorLogsReq) (*httpserver.CMSQueryResp, error) {
	rangeStart, rangeEnd, err := buildTimeRange(req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}
	pageIndex, pageSize := normalizePage(req.PageIndex, req.PageSize)

	matched := queryErrorLogsMatched(req, rangeStart, rangeEnd)
	total := len(matched)
	start := (pageIndex - 1) * pageSize
	if start >= total {
		return httpserver.NewCMSQueryResp(total, []*ErrorLogEntry{}), nil
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return httpserver.NewCMSQueryResp(total, matched[start:end]), nil
}

func GetTraceLogs(_ context.Context, req *logquerydto.CMSGetTraceLogsReq) (*logquerydto.CMSGetTraceLogsRes, error) {
	if req.TraceId == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	rangeStart, rangeEnd, err := buildTimeRange(req.StartDate, req.EndDate)
	if err != nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	if res, ok := getTraceLogsWithGrep(req, rangeStart, rangeEnd); ok {
		return res, nil
	}

	var detailLogs []*DetailLogEntry
	for _, filePath := range listDetailLogFilesForRange(req.StartDate, req.EndDate) {
		_ = scanDetailLogFile(filePath, func(entry *DetailLogEntry) bool {
			if !logTimeInRange(entry.Time, rangeStart, rangeEnd) {
				return true
			}
			if matchTraceId(entry.TraceId, req.TraceId) {
				detailLogs = append(detailLogs, trimDetailLogEntryForQuery(entry))
			}
			return true
		})
	}

	var accessLogs []*AccessLogEntry
	for _, filePath := range listAccessLogFilesForRange(req.StartDate, req.EndDate) {
		_ = scanAccessLogFile(filePath, func(entry *AccessLogEntry) bool {
			if !logTimeInRange(entry.Time, rangeStart, rangeEnd) {
				return true
			}
			if matchTraceId(entry.TraceId, req.TraceId) {
				accessLogs = append(accessLogs, entry)
			}
			return true
		})
	}

	var errorLogs []*ErrorLogEntry
	scanErrorLogEntriesInRange(req.StartDate, req.EndDate, func(entry *ErrorLogEntry) bool {
		if !logTimeInRange(entry.Time, rangeStart, rangeEnd) {
			return true
		}
		if matchTraceId(entry.TraceId, req.TraceId) && !isLogQueryRelatedErrorEntry(entry) {
			errorLogs = append(errorLogs, entry)
		}
		return true
	})

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
	}, nil
}

func GetAccessStats(_ context.Context, req *logquerydto.CMSGetAccessStatsReq) (*logquerydto.CMSGetAccessStatsRes, error) {
	rangeStart, rangeEnd, err := buildTimeRange(req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}
	topN := normalizeTopN(req.TopN)

	urlCounter := make(map[string]int64)
	ipCounter := make(map[string]int64)
	for _, filePath := range listAccessLogFilesForRange(req.StartDate, req.EndDate) {
		_ = scanAccessLogFile(filePath, func(entry *AccessLogEntry) bool {
			if !logTimeInRange(entry.Time, rangeStart, rangeEnd) {
				return true
			}
			if entry.Url != "" {
				urlCounter[entry.Url]++
			}
			if entry.Ip != "" {
				ipCounter[entry.Ip]++
			}
			return true
		})
	}

	return &logquerydto.CMSGetAccessStatsRes{
		UrlTop: toTopStatItems(urlCounter, topN),
		IpTop:  toTopStatItems(ipCounter, topN),
	}, nil
}

func scanDetailLogFile(filePath string, fn func(entry *DetailLogEntry) bool) error {
	return scanLogFile(filePath, func(line string) bool {
		entry, ok := parseDetailLogLine(line)
		if !ok {
			return true
		}
		return fn(entry)
	})
}

func scanAccessLogFile(filePath string, fn func(entry *AccessLogEntry) bool) error {
	return scanLogFile(filePath, func(line string) bool {
		entry, ok := parseAccessLogLine(line)
		if !ok {
			return true
		}
		return fn(entry)
	})
}

func validateDateRange(startDate, endDate string) error {
	_, _, err := buildTimeRange(startDate, endDate)
	return err
}

func buildTimeRange(startDate, endDate string) (time.Time, time.Time, error) {
	start, err := parseDate(startDate)
	if err != nil {
		return time.Time{}, time.Time{}, errercode.CreateCode(errercode.InvalidParam)
	}
	end, err := parseDate(endDate)
	if err != nil {
		return time.Time{}, time.Time{}, errercode.CreateCode(errercode.InvalidParam)
	}
	if end.Before(start) {
		return time.Time{}, time.Time{}, errercode.CreateCode(errercode.InvalidParam)
	}
	days := int(end.Sub(start).Hours()/24) + 1
	if days > maxDateRangeDays {
		return time.Time{}, time.Time{}, errercode.CreateCode(errercode.InvalidParam)
	}
	rangeEnd := end.Add(24*time.Hour - time.Millisecond)
	return start, rangeEnd, nil
}

func logTimeInRange(timeStr string, rangeStart, rangeEnd time.Time) bool {
	t, ok := parseLogTime(timeStr)
	if !ok {
		return false
	}
	return !t.Before(rangeStart) && !t.After(rangeEnd)
}

func parseDate(dateStr string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", dateStr, time.Local)
}

func listDates(startDate, endDate string) []string {
	start, _ := parseDate(startDate)
	end, _ := parseDate(endDate)
	var dates []string
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format("2006-01-02"))
	}
	return dates
}

func normalizePage(pageIndex, pageSize int) (int, int) {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	return pageIndex, pageSize
}

func normalizeTopN(topN int) int {
	if topN <= 0 {
		return defaultTopN
	}
	if topN > maxTopN {
		return maxTopN
	}
	return topN
}

type topCounterItem struct {
	Key   string
	Count int64
}

func toTopStatItems(counter map[string]int64, topN int) []logquerydto.TopStatItem {
	items := topItems(counter, topN)
	ret := make([]logquerydto.TopStatItem, 0, len(items))
	for _, item := range items {
		ret = append(ret, logquerydto.TopStatItem{
			Key:   item.Key,
			Count: item.Count,
		})
	}
	return ret
}

func topItems(counter map[string]int64, topN int) []topCounterItem {
	items := make([]topCounterItem, 0, len(counter))
	for key, count := range counter {
		items = append(items, topCounterItem{Key: key, Count: count})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Count == items[j].Count {
			return items[i].Key < items[j].Key
		}
		return items[i].Count > items[j].Count
	})
	if len(items) > topN {
		items = items[:topN]
	}
	return items
}
