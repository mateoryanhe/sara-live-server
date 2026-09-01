package logquery

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"xr-game-server/dto/logquerydto"
	"xr-game-server/module/fileexport"
)

type shellExportResult struct {
	ExportID  string `json:"exportId"`
	FileName  string `json:"fileName"`
	FileUrl   string `json:"fileUrl"`
	Total     int    `json:"total"`
	PageIndex int    `json:"pageIndex"`
	PageSize  int    `json:"pageSize"`
}

func toShellExportResult(rec *fileexport.Record) *shellExportResult {
	if rec == nil {
		return nil
	}
	return &shellExportResult{
		ExportID: rec.ExportID,
		FileName: rec.FileName,
		FileUrl:  rec.FileURL,
	}
}

func createShellExport(logType string, patterns []string, startDate, endDate string, pageIndex, pageSize int, minHandlerMs, maxHandlerMs float64) (*shellExportResult, error) {
	if err := ensureLinuxLogQuery(); err != nil {
		return nil, err
	}

	cfg := loadLogQueryConfig().normalized()
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 50
	}
	if pageSize > cfg.MaxPageSize {
		pageSize = cfg.MaxPageSize
	}

	files := listLogFilesByPrefix(cfg.LogDir, cfg.prefixForType(logType), startDate, endDate)
	rec, err := fileexport.Create(".log")
	if err != nil {
		return nil, err
	}
	exportID := rec.ExportID
	workDir := filepath.Join(cfg.exportAbsDir(), ".work")
	rawPath := filepath.Join(workDir, exportID+".raw")
	pagePath := filepath.Join(workDir, exportID+".page")
	exportPath := rec.AbsPath

	fail := func(err error) (*shellExportResult, error) {
		_ = fileexport.Delete(exportID)
		removeFile(rawPath)
		removeFile(pagePath)
		return nil, err
	}

	if logType == logTypeError {
		if err := grepFilesWithLogContinuationToFile(patterns, files, cfg.MaxMatchLines, rawPath, true); err != nil {
			return fail(err)
		}
	} else if err := grepFilesToReversedFile(patterns, files, cfg.MaxMatchLines, rawPath); err != nil {
		return fail(err)
	}
	var filterErr error
	rawPath, filterErr = applyLogTimeRangeFilter(rawPath, startDate, endDate, logType == logTypeError)
	if filterErr != nil {
		return fail(filterErr)
	}
	if logType == logTypeAccess {
		rawPath, filterErr = applyAccessLogQueryExcludeFilter(rawPath)
		if filterErr != nil {
			return fail(filterErr)
		}
		rawPath, filterErr = applyAccessHandlerMsFilter(rawPath, minHandlerMs, maxHandlerMs)
		if filterErr != nil {
			return fail(filterErr)
		}
	}
	if err := paginateFile(rawPath, pagePath, pageIndex, pageSize); err != nil {
		return fail(err)
	}
	if err := copyFile(pagePath, exportPath); err != nil {
		return fail(err)
	}
	removeFile(rawPath)
	removeFile(pagePath)

	result := toShellExportResult(rec)
	result.PageIndex = pageIndex
	result.PageSize = pageSize
	return result, nil
}

func createTraceShellExport(traceId, startDate, endDate string) (*shellExportResult, error) {
	if err := ensureLinuxLogQuery(); err != nil {
		return nil, err
	}
	cfg := loadLogQueryConfig().normalized()
	patterns := buildTracePatterns(traceId)
	rec, err := fileexport.Create(".trace.log")
	if err != nil {
		return nil, err
	}
	exportID := rec.ExportID
	workDir := filepath.Join(cfg.exportAbsDir(), ".work")
	exportPath := rec.AbsPath

	sections := []struct {
		tag    string
		prefix string
	}{
		{"detail", cfg.DetailPrefix},
		{"access", cfg.AccessPrefix},
		{"error", cfg.ErrorPrefix},
	}
	var builder strings.Builder
	for _, section := range sections {
		files := listLogFilesByPrefix(cfg.LogDir, section.prefix, startDate, endDate)
		rawPath := filepath.Join(workDir, exportID+"-"+section.tag+".raw")
		var grepErr error
		if section.tag == "error" {
			grepErr = grepFilesWithLogContinuationToFile(patterns, files, cfg.MaxMatchLines, rawPath, false)
		} else {
			grepErr = grepFilesToFile(patterns, files, cfg.MaxMatchLines, rawPath)
		}
		if grepErr != nil {
			_ = fileexport.Delete(exportID)
			return nil, grepErr
		}
		data, _ := os.ReadFile(rawPath)
		removeFile(rawPath)
		builder.WriteString("@" + section.tag + "\n")
		builder.Write(data)
		if len(data) > 0 && data[len(data)-1] != '\n' {
			builder.WriteByte('\n')
		}
	}
	if err := writeFile(exportPath, []byte(builder.String())); err != nil {
		_ = fileexport.Delete(exportID)
		return nil, err
	}
	return toShellExportResult(rec), nil
}

func createAccessStatsExport(startDate, endDate string, topN int) (*shellExportResult, error) {
	if err := ensureLinuxLogQuery(); err != nil {
		return nil, err
	}
	cfg := loadLogQueryConfig().normalized()
	if topN <= 0 {
		topN = 20
	}
	files := listLogFilesByPrefix(cfg.LogDir, cfg.AccessPrefix, startDate, endDate)
	rec, err := fileexport.Create(".stats.tsv")
	if err != nil {
		return nil, err
	}
	exportID := rec.ExportID
	workDir := filepath.Join(cfg.exportAbsDir(), ".work")
	rawPath := filepath.Join(workDir, exportID+".raw")
	exportPath := rec.AbsPath

	fail := func(err error) (*shellExportResult, error) {
		_ = fileexport.Delete(exportID)
		removeFile(rawPath)
		removeFile(filepath.Join(workDir, exportID+".url"))
		removeFile(filepath.Join(workDir, exportID+".ip"))
		return nil, err
	}

	if err := grepFilesToFile(nil, files, cfg.MaxMatchLines, rawPath); err != nil {
		return fail(err)
	}
	var filterErr error
	rawPath, filterErr = applyLogTimeRangeFilter(rawPath, startDate, endDate, false)
	if filterErr != nil {
		return fail(filterErr)
	}
	// access 日志格式: time {traceId} status "METHOD scheme host /path HTTP/1.1" ...
	script := `awk -F'"' 'NF>=2 {split($2,p," "); if (length(p)>=4) print p[4]}' ` + shellQuote(rawPath) +
		` | grep -v '^/logQuery/' | sort | uniq -c | sort -rn | head -n ` + strconv.Itoa(topN) +
		` > ` + shellQuote(filepath.Join(workDir, exportID+".url"))
	ipScript := `awk 'match($0, /" [0-9]+\.[0-9]+, ([^,]+),/, m) {gsub(/^ +| +$/,"",m[1]); if (m[1]!="") print m[1]}' ` + shellQuote(rawPath) +
		` | sort | uniq -c | sort -rn | head -n ` + strconv.Itoa(topN) +
		` > ` + shellQuote(filepath.Join(workDir, exportID+".ip"))
	_ = runShellScript(script)
	_ = runShellScript(ipScript)

	var builder strings.Builder
	builder.WriteString("urlTop\n")
	if urlData, err := os.ReadFile(filepath.Join(workDir, exportID+".url")); err == nil {
		builder.Write(urlData)
	}
	builder.WriteString("ipTop\n")
	if ipData, err := os.ReadFile(filepath.Join(workDir, exportID+".ip")); err == nil {
		builder.Write(ipData)
	}
	if err := writeFile(exportPath, []byte(builder.String())); err != nil {
		return fail(err)
	}
	removeFile(rawPath)
	removeFile(filepath.Join(workDir, exportID+".url"))
	removeFile(filepath.Join(workDir, exportID+".ip"))

	return toShellExportResult(rec), nil
}

func createAccessTrendExport(req *logquerydto.CMSGetAccessTrendReq) (*shellExportResult, error) {
	if err := ensureLinuxLogQuery(); err != nil {
		return nil, err
	}
	cfg := loadLogQueryConfig().normalized()
	intervalMinutes := resolveTrendIntervalMinutes(req.StartDate, req.EndDate, req.IntervalMinutes)
	patterns := buildAccessPatterns(req.TraceId, req.AuthId, req.Url, req.Ip, req.StatusCode)

	files := listLogFilesByPrefix(cfg.LogDir, cfg.AccessPrefix, req.StartDate, req.EndDate)
	rec, err := fileexport.Create(".trend.tsv")
	if err != nil {
		return nil, err
	}
	exportID := rec.ExportID
	workDir := filepath.Join(cfg.exportAbsDir(), ".work")
	rawPath := filepath.Join(workDir, exportID+".raw")
	bucketPath := filepath.Join(workDir, exportID+".bucket")
	exportPath := rec.AbsPath

	fail := func(err error) (*shellExportResult, error) {
		_ = fileexport.Delete(exportID)
		removeFile(rawPath)
		removeFile(bucketPath)
		return nil, err
	}

	if err := grepFilesToFile(patterns, files, cfg.MaxMatchLines, rawPath); err != nil {
		return fail(err)
	}
	var filterErr error
	rawPath, filterErr = applyLogTimeRangeFilter(rawPath, req.StartDate, req.EndDate, false)
	if filterErr != nil {
		return fail(filterErr)
	}
	rawPath, filterErr = applyAccessLogQueryExcludeFilter(rawPath)
	if filterErr != nil {
		return fail(filterErr)
	}

	minMs := 0.0
	maxMs := 0.0
	if req.MinHandlerMs != nil {
		minMs = *req.MinHandlerMs
	}
	if req.MaxHandlerMs != nil {
		maxMs = *req.MaxHandlerMs
	}

	awkScript := `awk -v interval=` + strconv.Itoa(intervalMinutes) +
		` -v minMs=` + strconv.FormatFloat(minMs, 'f', -1, 64) +
		` -v maxMs=` + strconv.FormatFloat(maxMs, 'f', -1, 64) + ` '
function handler_ms(line,   m) {
  if (match(line, /HTTP\/[0-9.]+\"[[:space:]]+([0-9.]+),/, m)) return m[1] * 1000
  return -1
}
{
  ms = handler_ms($0)
  if (minMs > 0 && (ms < 0 || ms <= minMs)) next
  if (maxMs > 0 && (ms < 0 || ms > maxMs)) next
  if (match($0, /^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}/)) {
    ts = substr($0, RSTART, RLENGTH)
    split(substr(ts, 12), parts, ":")
    hour = parts[1]
    minute = int(parts[2] / interval) * interval
    bucket = substr(ts, 1, 11) hour ":" sprintf("%02d", minute)
    count[bucket]++
    total++
  }
}
END {
  peak = 0
  peakTime = ""
  for (b in count) {
    if (count[b] > peak) { peak = count[b]; peakTime = b }
    print b "\t" count[b]
  }
  print "@meta"
  print "intervalMinutes\t" interval
  print "totalCount\t" total
  print "peakTime\t" peakTime
  print "peakCount\t" peak
}' ` + shellQuote(rawPath) + ` > ` + shellQuote(bucketPath)
	if err := runShellScript(awkScript); err != nil {
		return fail(err)
	}

	bucketData, err := os.ReadFile(bucketPath)
	if err != nil {
		return fail(err)
	}
	var builder strings.Builder
	builder.WriteString("points\n")
	for _, line := range splitLines(bucketData) {
		if line == "" || strings.HasPrefix(line, "@meta") {
			break
		}
		builder.WriteString(line)
		builder.WriteByte('\n')
	}
	builder.WriteString("meta\n")
	metaStarted := false
	for _, line := range splitLines(bucketData) {
		if line == "@meta" {
			metaStarted = true
			continue
		}
		if metaStarted && line != "" {
			builder.WriteString(line)
			builder.WriteByte('\n')
		}
	}
	if err := writeFile(exportPath, []byte(builder.String())); err != nil {
		return fail(err)
	}
	removeFile(rawPath)
	removeFile(bucketPath)

	return toShellExportResult(rec), nil
}

func resolveTrendIntervalMinutes(startDate, endDate string, requested int) int {
	if requested > 0 {
		return normalizeTrendInterval(requested)
	}
	start, ok1 := parseLogQueryDateTime(startDate)
	end, ok2 := parseLogQueryDateTime(endDate)
	if !ok1 || !ok2 {
		return 15
	}
	startDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	endDay := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())
	days := int(endDay.Sub(startDay).Hours()/24) + 1
	switch {
	case days <= 1:
		return 1
	case days <= 3:
		return 5
	default:
		return 15
	}
}

func normalizeTrendInterval(minutes int) int {
	for _, item := range []int{1, 5, 15, 60} {
		if minutes <= item {
			return item
		}
	}
	return 60
}

func deleteExport(exportID string) error {
	return fileexport.Delete(exportID)
}
