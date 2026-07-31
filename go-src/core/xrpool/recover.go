package xrpool

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	"xr-game-server/core/xrlog"
)

// Recover 供 defer 使用,捕获当前协程 panic 并写 error 日志,不退出进程。
func Recover(ctx context.Context, source string) {
	p := recover()
	if p == nil {
		return
	}
	logPanic(ctx, source, p)
}

// RecoverMain 供 main 入口 defer 使用; 启动阶段 panic 写 error 日志后以非 0 退出。
func RecoverMain(ctx context.Context, source string) {
	p := recover()
	if p == nil {
		return
	}
	logPanic(ctx, source, p)
	os.Exit(1)
}

func logPanic(ctx context.Context, source string, p any) {
	if ctx == nil {
		ctx = context.Background()
	}
	err := fmt.Errorf("panic: %v\n%s", p, debug.Stack())
	xrlog.ErrorWithErr(ctx, source, "recover", err)
}
