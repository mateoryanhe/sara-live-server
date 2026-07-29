package cfg

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
)

const dbLoggerQueryOnlyCfg = "database.logger.queryOnly"

func initDbQueryOnlyLogger() {
	ctx := gctx.New()
	if !g.Cfg().MustGet(ctx, dbLoggerQueryOnlyCfg).Bool() {
		return
	}
	db := g.DB()
	logger, ok := db.GetLogger().(*glog.Logger)
	if !ok {
		return
	}
	logger.SetHandlers(dbQueryOnlyLogHandler)
}

func dbQueryOnlyLogHandler(ctx context.Context, in *glog.HandlerInput) {
	if in.Level == glog.LEVEL_DEBU && !isDbAllowedSqlLog(in.Content, in.Values) {
		return
	}
	in.Next(ctx)
}

func isDbAllowedSqlLog(content string, values []any) bool {
	sqlText := strings.TrimSpace(content)
	if sqlText == "" && len(values) > 0 {
		sqlText = strings.TrimSpace(gconv.String(values[0]))
	}
	sqlText = extractSqlFromGdbLog(sqlText)
	if sqlText == "" {
		return false
	}
	upper := strings.ToUpper(sqlText)
	return strings.HasPrefix(upper, "SELECT") ||
		strings.HasPrefix(upper, "UPDATE") ||
		strings.HasPrefix(upper, "DELETE")
}

func extractSqlFromGdbLog(content string) string {
	idx := strings.LastIndex(content, "[rows:")
	if idx < 0 {
		return strings.TrimSpace(content)
	}
	rest := strings.TrimSpace(content[idx:])
	closeIdx := strings.Index(rest, "]")
	if closeIdx < 0 {
		return strings.TrimSpace(content)
	}
	rest = strings.TrimSpace(rest[closeIdx+1:])
	if strings.HasPrefix(rest, "[txid:") {
		if end := strings.Index(rest, "]"); end >= 0 {
			rest = strings.TrimSpace(rest[end+1:])
		}
	}
	return rest
}
