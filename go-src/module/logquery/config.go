package logquery

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cfg"
	"xr-game-server/errercode"
)

const (
	logTypeAccess = "access"
	logTypeDetail = "detail"
	logTypeError  = "error"

	defaultExportStaticPrefix = "/cms"
	defaultExportSubDir       = "log-query-export"
	defaultExportTTLMinutes   = 30
	defaultMaxMatchLines      = 100000
	defaultMaxPageSize        = 200
)

type logQueryConfig struct {
	LogDir             string
	AccessPrefix       string
	DetailPrefix       string
	ErrorPrefix        string
	ExportStaticPrefix string
	ExportSubDir       string
	ExportRoot         string
	ExportTTLMinutes   int
	MaxMatchLines      int
	MaxPageSize        int
}

func buildLogQueryConfig() logQueryConfig {
	ctx := gctx.New()
	logDir := cfgGetString(ctx, "logger.access.path")
	if logDir == "" {
		logDir = cfgGetString(ctx, "logger.detail.path")
	}
	if logDir == "" {
		logDir = cfgGetString(ctx, "logger.error.path")
	}
	if logDir == "" {
		logDir = cfgGetString(ctx, "server.logPath")
	}

	exportStaticPrefix := normalizeLogQueryURLPrefix(cfgGetString(ctx, "logQuery.exportStaticPrefix"))
	if exportStaticPrefix == "" {
		exportStaticPrefix = defaultExportStaticPrefix
	}
	exportSubDir := cfgGetString(ctx, "logQuery.exportSubDir")
	if exportSubDir == "" {
		exportSubDir = defaultExportSubDir
	}

	exportRoot := strings.TrimSpace(cfg.GetStaticPathRoot(exportStaticPrefix))
	if exportRoot == "" {
		if serverCfg := cfg.GetServerCfg(); serverCfg != nil {
			if root := strings.TrimSpace(serverCfg.StaticResPath); root != "" {
				exportRoot = filepath.Join(root, strings.Trim(exportStaticPrefix, "/"))
			}
		}
	}
	if exportRoot == "" {
		exportRoot = "."
	}

	accessPattern := cfgGetString(ctx, "logger.access.file")
	if accessPattern == "" {
		accessPattern = cfgGetString(ctx, "server.accessLogPattern")
	}

	return logQueryConfig{
		LogDir:             filepath.Clean(logDir),
		AccessPrefix:       logFilePrefixFromPattern(accessPattern),
		DetailPrefix:       logFilePrefixFromPattern(cfgGetString(ctx, "logger.detail.file")),
		ErrorPrefix:        logFilePrefixFromPattern(cfgGetString(ctx, "logger.error.file")),
		ExportStaticPrefix: exportStaticPrefix,
		ExportSubDir:       exportSubDir,
		ExportRoot:         filepath.Clean(exportRoot),
		ExportTTLMinutes:   defaultExportTTLMinutes,
		MaxMatchLines:      defaultMaxMatchLines,
		MaxPageSize:        defaultMaxPageSize,
	}
}

func cfgGetString(ctx context.Context, key string) string {
	value, err := g.Cfg().Get(ctx, key)
	if err != nil || value.IsNil() {
		return ""
	}
	return strings.TrimSpace(value.String())
}

func normalizeLogQueryURLPrefix(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return ""
	}
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	return strings.TrimRight(prefix, "/")
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
	return filepath.Join(c.ExportRoot, c.ExportSubDir)
}

func (c logQueryConfig) exportURLPrefix() string {
	prefix := normalizeLogQueryURLPrefix(c.ExportStaticPrefix)
	if prefix == "" {
		prefix = defaultExportStaticPrefix
	}
	sub := strings.Trim(strings.ReplaceAll(c.ExportSubDir, "\\", "/"), "/")
	if sub == "" {
		sub = defaultExportSubDir
	}
	return prefix + "/" + sub
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
