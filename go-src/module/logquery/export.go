package logquery

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"xr-game-server/dto/logquerydto"

	"github.com/gogf/gf/v2/util/guid"
)

type exportRecord struct {
	exportID  string
	fileName  string
	absPath   string
	fileURL   string
	createdAt time.Time
}

var (
	exportRecords  sync.Map
	exportInitOnce sync.Once
)

func initExportCleanup() {
	exportInitOnce.Do(func() {
		go exportCleanupLoop()
	})
}

type shellExportResult struct {
	ExportID  string `json:"exportId"`
	FileName  string `json:"fileName"`
	FileUrl   string `json:"fileUrl"`
	Total     int    `json:"total"`
	PageIndex int    `json:"pageIndex"`
	PageSize  int    `json:"pageSize"`
}

func createShellExport(logType string, patterns []string, startDate, endDate string, pageIndex, pageSize int) (*shellExportResult, error) {
	if err := ensureLinuxLogQuery(); err != nil {
		return nil, err
	}
	initExportCleanup()

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
	exportID := guid.S()
	fileName := exportID + ".log"
	workDir := filepath.Join(cfg.exportAbsDir(), ".work")
	rawPath := filepath.Join(workDir, exportID+".raw")
	pagePath := filepath.Join(workDir, exportID+".page")
	exportPath := filepath.Join(cfg.exportAbsDir(), fileName)

	if err := grepFilesToFile(patterns, files, cfg.MaxMatchLines, rawPath); err != nil {
		return nil, err
	}
	total, err := countLines(rawPath)
	if err != nil {
		removeFile(rawPath)
		return nil, err
	}
	if err := paginateFile(rawPath, pagePath, pageIndex, pageSize); err != nil {
		removeFile(rawPath)
		return nil, err
	}
	if err := copyFile(pagePath, exportPath); err != nil {
		removeFile(rawPath)
		removeFile(pagePath)
		return nil, err
	}
	removeFile(rawPath)
	removeFile(pagePath)

	fileURL := cfg.exportURLPrefix() + "/" + fileName
	record := &exportRecord{
		exportID:  exportID,
		fileName:  fileName,
		absPath:   exportPath,
		fileURL:   fileURL,
		createdAt: time.Now(),
	}
	exportRecords.Store(exportID, record)

	return &shellExportResult{
		ExportID:  exportID,
		FileName:  fileName,
		FileUrl:   fileURL,
		Total:     total,
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}, nil
}

func createTraceShellExport(traceId, startDate, endDate string) (*shellExportResult, error) {
	if err := ensureLinuxLogQuery(); err != nil {
		return nil, err
	}
	initExportCleanup()
	cfg := loadLogQueryConfig().normalized()
	patterns := buildTracePatterns(traceId)
	exportID := guid.S()
	fileName := exportID + ".trace.log"
	workDir := filepath.Join(cfg.exportAbsDir(), ".work")
	exportPath := filepath.Join(cfg.exportAbsDir(), fileName)

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
		if err := grepFilesToFile(patterns, files, cfg.MaxMatchLines, rawPath); err != nil {
			return nil, err
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
		return nil, err
	}
	fileURL := cfg.exportURLPrefix() + "/" + fileName
	record := &exportRecord{
		exportID:  exportID,
		fileName:  fileName,
		absPath:   exportPath,
		fileURL:   fileURL,
		createdAt: time.Now(),
	}
	exportRecords.Store(exportID, record)
	return &shellExportResult{
		ExportID: exportID,
		FileName: fileName,
		FileUrl:  fileURL,
		Total:    0,
	}, nil
}

func createAccessStatsExport(startDate, endDate string, topN int) (*shellExportResult, error) {
	if err := ensureLinuxLogQuery(); err != nil {
		return nil, err
	}
	initExportCleanup()
	cfg := loadLogQueryConfig().normalized()
	if topN <= 0 {
		topN = 20
	}
	files := listLogFilesByPrefix(cfg.LogDir, cfg.AccessPrefix, startDate, endDate)
	exportID := guid.S()
	fileName := exportID + ".stats.tsv"
	workDir := filepath.Join(cfg.exportAbsDir(), ".work")
	rawPath := filepath.Join(workDir, exportID+".raw")
	exportPath := filepath.Join(cfg.exportAbsDir(), fileName)

	if err := grepFilesToFile(nil, files, cfg.MaxMatchLines, rawPath); err != nil {
		return nil, err
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
		return nil, err
	}
	removeFile(rawPath)
	removeFile(filepath.Join(workDir, exportID+".url"))
	removeFile(filepath.Join(workDir, exportID+".ip"))

	fileURL := cfg.exportURLPrefix() + "/" + fileName
	record := &exportRecord{exportID: exportID, fileName: fileName, absPath: exportPath, fileURL: fileURL, createdAt: time.Now()}
	exportRecords.Store(exportID, record)
	return &shellExportResult{ExportID: exportID, FileName: fileName, FileUrl: fileURL}, nil
}

func createAccessTrendExport(req *logquerydto.CMSGetAccessTrendReq) (*shellExportResult, error) {
	if err := ensureLinuxLogQuery(); err != nil {
		return nil, err
	}
	initExportCleanup()
	cfg := loadLogQueryConfig().normalized()
	intervalMinutes := resolveTrendIntervalMinutes(req.StartDate, req.EndDate, req.IntervalMinutes)
	patterns := buildAccessPatterns(req.TraceId, req.Url, req.Ip, req.StatusCode)

	files := listLogFilesByPrefix(cfg.LogDir, cfg.AccessPrefix, req.StartDate, req.EndDate)
	exportID := guid.S()
	fileName := exportID + ".trend.tsv"
	workDir := filepath.Join(cfg.exportAbsDir(), ".work")
	rawPath := filepath.Join(workDir, exportID+".raw")
	bucketPath := filepath.Join(workDir, exportID+".bucket")
	exportPath := filepath.Join(cfg.exportAbsDir(), fileName)

	if err := grepFilesToFile(patterns, files, cfg.MaxMatchLines, rawPath); err != nil {
		return nil, err
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
  if (minMs > 0 && ms >= 0 && ms < minMs) next
  if (maxMs > 0 && ms >= 0 && ms > maxMs) next
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
		removeFile(rawPath)
		return nil, err
	}

	bucketData, err := os.ReadFile(bucketPath)
	if err != nil {
		removeFile(rawPath)
		removeFile(bucketPath)
		return nil, err
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
		removeFile(rawPath)
		removeFile(bucketPath)
		return nil, err
	}
	removeFile(rawPath)
	removeFile(bucketPath)

	fileURL := cfg.exportURLPrefix() + "/" + fileName
	record := &exportRecord{exportID: exportID, fileName: fileName, absPath: exportPath, fileURL: fileURL, createdAt: time.Now()}
	exportRecords.Store(exportID, record)
	return &shellExportResult{ExportID: exportID, FileName: fileName, FileUrl: fileURL}, nil
}

func resolveTrendIntervalMinutes(startDate, endDate string, requested int) int {
	if requested > 0 {
		return normalizeTrendInterval(requested)
	}
	start, err1 := time.ParseInLocation("2006-01-02", startDate, time.Local)
	end, err2 := time.ParseInLocation("2006-01-02", endDate, time.Local)
	if err1 != nil || err2 != nil {
		return 15
	}
	days := int(end.Sub(start).Hours()/24) + 1
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
	value, ok := exportRecords.Load(exportID)
	if ok {
		if record, ok := value.(*exportRecord); ok {
			removeFile(record.absPath)
		}
		exportRecords.Delete(exportID)
	}
	cfg := loadLogQueryConfig()
	removeFile(filepath.Join(cfg.exportAbsDir(), exportID+".log"))
	removeFile(filepath.Join(cfg.exportAbsDir(), exportID+".trace.log"))
	removeFile(filepath.Join(cfg.exportAbsDir(), exportID+".stats.tsv"))
	removeFile(filepath.Join(cfg.exportAbsDir(), exportID+".trend.tsv"))
	return nil
}

func exportCleanupLoop() {
	cfg := loadLogQueryConfig()
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		expireBefore := time.Now().Add(-time.Duration(cfg.ExportTTLMinutes) * time.Minute)
		exportRecords.Range(func(key, value any) bool {
			record, ok := value.(*exportRecord)
			if !ok || record.createdAt.After(expireBefore) {
				return true
			}
			removeFile(record.absPath)
			exportRecords.Delete(key)
			return true
		})
		cleanExportDir(cfg.exportAbsDir(), expireBefore)
	}
}

func cleanExportDir(dir string, expireBefore time.Time) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil || info.ModTime().After(expireBefore) {
			continue
		}
		removeFile(filepath.Join(dir, entry.Name()))
	}
}

func shellQuote(path string) string {
	return "'" + strings.ReplaceAll(path, "'", `'\'\'\'`) + "'"
}

func runShellScript(script string) error {
	_, err := exec.Command("bash", "-lc", script).Output()
	return err
}
