package logquery

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
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

func listDates(startDate, endDate string) []string {
	start, err1 := time.ParseInLocation("2006-01-02", startDate, time.Local)
	end, err2 := time.ParseInLocation("2006-01-02", endDate, time.Local)
	if err1 != nil || err2 != nil {
		return nil
	}
	var dates []string
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format("2006-01-02"))
	}
	return dates
}
