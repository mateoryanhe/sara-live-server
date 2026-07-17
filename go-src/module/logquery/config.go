package logquery

import (
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

type logPaths struct {
	DetailLogDir     string
	DetailLogPattern string
	AccessLogDir     string
	AccessLogPattern string
	ErrorLogDir      string
	ErrorLogPattern  string
}

func loadLogPaths() logPaths {
	ctx := gctx.New()
	detailDir := strings.TrimSpace(g.Cfg().MustGet(ctx, "logger.path").String())
	detailPattern := strings.TrimSpace(g.Cfg().MustGet(ctx, "logger.file").String())
	accessDir := strings.TrimSpace(g.Cfg().MustGet(ctx, "server.logPath").String())
	accessPattern := strings.TrimSpace(g.Cfg().MustGet(ctx, "server.accessLogPattern").String())
	errorPattern := strings.TrimSpace(g.Cfg().MustGet(ctx, "server.errorLogPattern").String())

	if detailPattern == "" {
		detailPattern = "{Y-m-d}.log"
	}
	if accessDir == "" {
		accessDir = detailDir
	}
	if accessPattern == "" {
		accessPattern = "access-{Y-m-d}.log"
	}
	if errorPattern == "" {
		errorPattern = "error-{Ymd}.log"
	}

	return logPaths{
		DetailLogDir:     normalizeDir(detailDir),
		DetailLogPattern: detailPattern,
		AccessLogDir:     normalizeDir(accessDir),
		AccessLogPattern: accessPattern,
		ErrorLogDir:      normalizeDir(accessDir),
		ErrorLogPattern:  errorPattern,
	}
}

func normalizeDir(dir string) string {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		return ""
	}
	return filepath.Clean(dir)
}

func formatLogFileName(pattern string, dateStr string) string {
	t, err := parseDate(dateStr)
	if err != nil {
		return pattern
	}
	s := pattern
	s = strings.ReplaceAll(s, "{Y-m-d}", t.Format("2006-01-02"))
	s = strings.ReplaceAll(s, "{Ymd}", t.Format("20060102"))
	s = strings.ReplaceAll(s, "{Y-m}", t.Format("2006-01"))
	s = strings.ReplaceAll(s, "{Y}", t.Format("2006"))
	s = strings.ReplaceAll(s, "{m}", t.Format("01"))
	s = strings.ReplaceAll(s, "{d}", t.Format("02"))
	return s
}

func resolveDetailLogFile(dateStr string) string {
	paths := loadLogPaths()
	if paths.DetailLogDir == "" {
		return ""
	}
	return filepath.Join(paths.DetailLogDir, formatLogFileName(paths.DetailLogPattern, dateStr))
}

func resolveAccessLogFile(dateStr string) string {
	paths := loadLogPaths()
	if paths.AccessLogDir == "" {
		return ""
	}
	return filepath.Join(paths.AccessLogDir, formatLogFileName(paths.AccessLogPattern, dateStr))
}
