package hotrestart

import (
	"context"
	"crypto/subtle"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/common"
	"xr-game-server/core/cfg"
	"xr-game-server/core/event"
	"xr-game-server/core/shutdown"
	"xr-game-server/core/syndb"
	"xr-game-server/core/xrlog"
	"xr-game-server/core/xrtimer"
)

const (
	defaultHotRestartFlushTimeout = 60
	defaultHotRestartExitTimeout  = 120
	defaultHotRestartAuth         = "nGH66S4TjBjQqCKyWJAM"
)

var (
	hotRestartOnce         sync.Once
	hotRestartTriggered    atomic.Bool
	oldProcessAfterRestart atomic.Bool
	enterRestartPhaseFunc  func()
	hotRestartAuthProvider func() string
)

// RegisterHotRestartAuthProvider 由 preload 注册,避免包循环依赖.
func RegisterHotRestartAuthProvider(fn func() string) {
	hotRestartAuthProvider = fn
}

// RegisterEnterRestartPhase 由 httpserver 注册,避免包循环依赖.
func RegisterEnterRestartPhase(fn func()) {
	enterRestartPhaseFunc = fn
}

// TryTriggerHotRestart 校验密钥后异步触发热重启.
func TryTriggerHotRestart(auth string) (accepted bool, reason string) {
	expected := hotRestartAuth()
	if expected == "" {
		return false, "hot restart disabled"
	}
	if subtle.ConstantTimeCompare([]byte(auth), []byte(expected)) != 1 {
		return false, "unauthorized"
	}
	if !hotRestartTriggered.CompareAndSwap(false, true) {
		return false, "already in progress"
	}
	go runHotRestart()
	return true, "accepted"
}

func runHotRestart() {
	hotRestartOnce.Do(func() {
		ctx := gctx.New()
		xrlog.DetailLog.Warning(ctx, "热重启:收到请求,停止定时器并开始刷盘(HTTP 继续服务)")

		event.Pub(event.PrepareRestart, nil)
		xrtimer.Pause()

		flushTimeout := hotRestartFlushTimeoutSec()
		xrlog.DetailLog.Warningf(ctx, "热重启:开始无条件刷盘,timeout=%vs", flushTimeout)

		idle := syndb.FlushUntilIdle(time.Duration(flushTimeout) * time.Second)
		if idle {
			xrlog.DetailLog.Info(ctx, "热重启:数据库入库队列已空,准备执行 Restart")
		} else {
			xrlog.DetailLog.Warning(ctx, "热重启:刷盘超时,强制执行 Restart")
		}

		shutdown.EnableSyncLoggers()

		xrlog.DetailLog.Warning(ctx, "热重启:调用 GoFrame RestartAllServer")
		if err := ghttp.RestartAllServer(ctx, ""); err != nil {
			xrlog.DetailLog.Errorf(ctx, "热重启:RestartAllServer 失败,err=%v", err)
			return
		}

		oldProcessAfterRestart.Store(true)
		if enterRestartPhaseFunc != nil {
			enterRestartPhaseFunc()
		}
		xrlog.DetailLog.Warning(ctx, "热重启:RestartAllServer 完成,旧进程开始自行退出")

		exitOldProcess(ctx, hotRestartExitTimeoutSec())
	})
}

// exitOldProcess 旧进程自行关闭 HTTP 并退出,不依赖新进程发消息;超时后强制退出.
func exitOldProcess(ctx context.Context, exitTimeoutSec int) {
	xrlog.DetailLog.Warningf(ctx, "热重启:旧进程关闭 HTTP,最长等待=%vs", exitTimeoutSec)

	shutdownCtx, cancel := context.WithTimeout(ctx, time.Duration(exitTimeoutSec)*time.Second)
	defer cancel()

	if err := ghttp.ShutdownAllServer(shutdownCtx); err != nil {
		xrlog.DetailLog.Warningf(ctx, "热重启:旧进程 ShutdownAllServer,err=%v", err)
	}

	if shutdownCtx.Err() != nil {
		xrlog.DetailLog.Warningf(ctx, "热重启:旧进程退出超时(%vs),强制退出", exitTimeoutSec)
	} else {
		xrlog.DetailLog.Info(ctx, "热重启:旧进程 HTTP 已关闭")
	}

	xrlog.DetailLog.Warning(ctx, "热重启:旧进程成功退出,佛祖保佑,成功重启")
	os.Exit(0)
}

// NotifyOldProcessExit 非热重启路径下 HTTP Run 返回时调用(热重启旧进程走 exitOldProcess).
func NotifyOldProcessExit() {
	if !oldProcessAfterRestart.Load() {
		return
	}
}

func hotRestartAuth() string {
	if hotRestartAuthProvider != nil {
		if auth := strings.TrimSpace(hotRestartAuthProvider()); auth != "" {
			return auth
		}
	}
	return defaultHotRestartAuth
}

func hotRestartFlushTimeoutSec() int {
	serverCfg := cfg.GetServerCfg()
	if serverCfg == nil || serverCfg.HotRestartFlushTimeout <= common.Zero {
		return defaultHotRestartFlushTimeout
	}
	return serverCfg.HotRestartFlushTimeout
}

func hotRestartExitTimeoutSec() int {
	serverCfg := cfg.GetServerCfg()
	if serverCfg == nil || serverCfg.HotRestartExitTimeout <= common.Zero {
		return defaultHotRestartExitTimeout
	}
	return serverCfg.HotRestartExitTimeout
}
