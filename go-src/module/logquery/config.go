package logquery

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/errercode"
)

const (
	logTypeAccess = "access"
	logTypeDetail = "detail"
	logTypeError  = "error"

	defaultExportSubDir     = "log-query-export"
	defaultExportTTLMinutes = 30
	defaultMaxMatchLines    = 100000
	defaultMaxPageSize      = 200
)

type logQueryConfig struct {
	LogDir           string
	AccessPrefix     string
	DetailPrefix     string
	ErrorPrefix      string
	ExportSubDir     string
	ExportTTLMinutes int
	MaxMatchLines    int
	MaxPageSize      int
	ServerRoot       string
}

func loadLogQueryConfig() logQueryConfig {
	ctx := gctx.New()
	logDir := strings.TrimSpace(g.Cfg().MustGet(ctx, "logger.detail.path").String())
	if logDir == "" {
		logDir = strings.TrimSpace(g.Cfg().MustGet(ctx, "logger.error.path").String())
	}
	if logDir == "" {
		logDir = strings.TrimSpace(g.Cfg().MustGet(ctx, "server.logPath").String())
	}
	return logQueryConfig{
		LogDir:           filepath.Clean(logDir),
		AccessPrefix:     logFilePrefixFromPattern(g.Cfg().MustGet(ctx, "server.accessLogPattern").String()),
		DetailPrefix:     logFilePrefixFromPattern(g.Cfg().MustGet(ctx, "logger.detail.file").String()),
		ErrorPrefix:      logFilePrefixFromPattern(g.Cfg().MustGet(ctx, "logger.error.file").String()),
		ExportSubDir:     defaultExportSubDir,
		ExportTTLMinutes: defaultExportTTLMinutes,
		MaxMatchLines:    defaultMaxMatchLines,
		MaxPageSize:      defaultMaxPageSize,
		ServerRoot:       filepath.Clean(g.Cfg().MustGet(ctx, "server.serverRoot").String()),
	}
}

func logFilePrefixFromPattern(pattern string) string {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		return ""
	}
	idx := strings.Index(pattern, "{")
	prefix := pattern
	if idx > 0 {
		prefix = pattern[:idx]
	} else {
		prefix = strings.TrimSuffix(pattern, ".log")
	}
	return strings.TrimRight(prefix, "-")
}

func (c logQueryConfig) normalized() logQueryConfig {
	if c.AccessPrefix == "" {
		c.AccessPrefix = "access"
	}
	if c.DetailPrefix == "" {
		c.DetailPrefix = "detail"
	}
	if c.ErrorPrefix == "" {
		c.ErrorPrefix = "error"
	}
	return c
}

func (c logQueryConfig) prefixForType(logType string) string {
	c = c.normalized()
	switch logType {
	case logTypeAccess:
		return c.AccessPrefix
	case logTypeError:
		return c.ErrorPrefix
	default:
		return c.DetailPrefix
	}
}

func (c logQueryConfig) exportAbsDir() string {
	return filepath.Join(c.ServerRoot, c.ExportSubDir)
}

func (c logQueryConfig) exportURLPrefix() string {
	sub := strings.Trim(strings.ReplaceAll(c.ExportSubDir, "\\", "/"), "/")
	if sub == "" {
		return "/log-query-export"
	}
	return "/" + sub
}

func formatLogFileName(pattern string, dateStr string) string {
	t, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	if err != nil {
		return pattern
	}
	s := pattern
	s = strings.ReplaceAll(s, "{Y-m-d}", t.Format("2006-01-02"))
	s = strings.ReplaceAll(s, "{Ymd}", t.Format("20060102"))
	return s
}

const (
	logQueryDateLayout      = "2006-01-02"
	logQueryDateTimeLayout  = "2006-01-02 15:04:05"
	logQueryDateTimeTLayout = "2006-01-02T15:04:05"
)

func parseLogQueryDateTime(value string) (time.Time, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, false
	}
	for _, layout := range []string{logQueryDateTimeLayout, logQueryDateTimeTLayout, logQueryDateLayout} {
		if t, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func isLogQueryDateOnly(value string) bool {
	value = strings.TrimSpace(value)
	return len(value) == len(logQueryDateLayout)
}

func hasLogQueryTimeComponent(value string) bool {
	return len(strings.TrimSpace(value)) > len(logQueryDateLayout)
}

func shouldApplyLogTimeFilter(startDate, endDate string) bool {
	return hasLogQueryTimeComponent(startDate) || hasLogQueryTimeComponent(endDate)
}

func normalizeLogQueryLogRange(startDate, endDate string) (startLog, endLog string, ok bool) {
	start, ok1 := parseLogQueryDateTime(startDate)
	end, ok2 := parseLogQueryDateTime(endDate)
	if !ok1 || !ok2 {
		return "", "", false
	}
	if isLogQueryDateOnly(endDate) {
		end = end.Add(24*time.Hour - time.Second)
	}
	return start.Format(logQueryDateTimeTLayout), end.Format(logQueryDateTimeTLayout), true
}

func listDates(startDate, endDate string) []string {
	start, ok1 := parseLogQueryDateTime(startDate)
	end, ok2 := parseLogQueryDateTime(endDate)
	if !ok1 || !ok2 {
		return nil
	}
	startDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	endDay := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())
	var dates []string
	for d := startDay; !d.After(endDay); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format(logQueryDateLayout))
	}
	return dates
}

func validateLogQueryDateTime(value string) error {
	if _, ok := parseLogQueryDateTime(value); !ok {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	return nil
}

func validateLogQueryDateRange(startDate, endDate string) error {
	if err := validateLogQueryDateTime(startDate); err != nil {
		return err
	}
	if err := validateLogQueryDateTime(endDate); err != nil {
		return err
	}
	start, _ := parseLogQueryDateTime(startDate)
	end, _ := parseLogQueryDateTime(endDate)
	if start.After(end) {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	return nil
}
