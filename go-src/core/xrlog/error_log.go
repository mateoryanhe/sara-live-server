package xrlog

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

const Tag = "ErrorLog"

func Error(ctx context.Context, source, msg string) {
	g.Log().Errorf(ctx, "%s source=%s %s", Tag, source, msg)
}

func ErrorWithErr(ctx context.Context, source, msg string, err error) {
	if err == nil {
		Error(ctx, source, msg)
		return
	}
	g.Log().Errorf(ctx, "%s source=%s %s err=%v stack=%s", Tag, source, msg, err, formatStack(err))
}

func formatStack(err error) string {
	if err == nil {
		return ""
	}
	if stack := gerror.Stack(err); stack != "" {
		return stack
	}
	return fmt.Sprintf("%+v", err)
}
