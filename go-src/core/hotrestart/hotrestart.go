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
	"xr-game-server/core/event"
	"xr-game-server/core/shutdown"
	"xr-game-server/core/syndb"
	"xr-game-server/core/xrlog"
	"xr-game-server/core/xrpool"
	"xr-game-server/core/xrtimer"
)

const (
	defaultHotRestartPhase1Wait = 60
	defaultHotRestartExitWait   = 60
	defaultHotRestartAuth       = "nGH66S4TjBjQqCKyWJAM"
)

var (
	hotRestartTriggered    atomic.Bool
	hotRestartTriggerAt    atomic.Int64 // 毫秒时间戳,用于回收卡死的 in-progress 标志
	oldProcessAfterRestart atomic.Bool
	enterRestartPhaseFunc  func()
	hotRestartAuthProvider       func() string
	hotRestartPhase1WaitProvider func() int
	hotRestartExitWaitProvider   func() int
)

const staleHotRestartTriggerAfter = 2 * time.Minute

// RegisterHotRestartAuthProvider 由 preload 注册,避免包循环依赖.
func RegisterHotRestartAuthProvider(fn func() string) {
	hotRestartAuthProvider = fn
}

// RegisterHotRestartTimeoutProviders 由 preload 注册,避免包循环依赖.
func RegisterHotRestartTimeoutProviders(phase1WaitSec, exitWaitSec func() int) {
	hotRestartPhase1WaitProvider = phase1WaitSec
	hotRestartExitWaitProvider = exitWaitSec
}

// RegisterEnterRestartPhase 由 httpserver 注册,避免包循环依赖.
func RegisterEnterRestartPhase(fn func()) {
	enterRestartPhaseFunc = fn
}

// clearStaleHotRestartTrigger 回收长时间未进入 fork 的 in-progress 标志(如 xrpool 未调度到 runHotRestart).
func clearStaleHotRestartTrigger(ctx context.Context) {
	if !hotRestartTriggered.Load() || oldProcessAfterRestart.Load() {
		return
	}
	at := hotRestartTriggerAt.Load()
	if at <= 0 {
		return
	}
	if time.Since(time.UnixMilli(at)) <= staleHotRestartTriggerAfter {
		return
	}
	hotRestartTriggered.Store(false)
	hotRestartTriggerAt.Store(0)
	xrlog.DetailLog.Warning(ctx, "热重启 in-progress 标志超时,已自动回收")
	xrlog.DetailLog.Warning(ctx, "hot_restart stale trigger cleared")
}

// TryTriggerHotRestart 校验密钥后异步触发热重启.
func TryTriggerHotRestart(auth string) (accepted bool, reason string) {
	ctx := gctx.New()
	expected := hotRestartAuth()
	if expected == "" {
		return false, "hot restart disabled"
	}
	if subtle.ConstantTimeCompare([]byte(auth), []byte(expected)) != 1 {
		return false, "unauthorized"
	}
	clearStaleHotRestartTrigger(ctx)
	if !hotRestartTriggered.CompareAndSwap(false, true) {
		return false, "already in progress"
	}
	hotRestartTriggerAt.Store(time.Now().UnixMilli())
	xrpool.AddWithRecover(gctx.New(), func(_ context.Context) {
		runHotRestart()
	})
	return true, "accepted"
}

func runHotRestart() {
	ctx := gctx.New()
	defer func() {
		if !oldProcessAfterRestart.Load() {
			hotRestartTriggered.Store(false)
			hotRestartTriggerAt.Store(0)
		}
	}()

	logPhase1Start(ctx)

	event.Pub(event.PrepareRestart, nil)
	xrtimer.Pause()

	startPhase1FlushBackground(ctx)
	waitHotRestartPhase1(ctx, time.Duration(hotRestartPhase1WaitSec())*time.Second)
	if syndb.AllCachesIdle() {
		logPhase1QueueEmpty(ctx)
	}
	logPhase1End(ctx)

	shutdown.EnableSyncLoggers()
	ensureChildReadyListener()

	if err := ghttp.RestartAllServer(ctx, ""); err != nil {
		xrlog.DetailLog.Errorf(ctx, "热更新---RestartAllServer失败,err=%v", err)
		xrlog.DetailLog.Errorf(ctx, "hot_restart RestartAllServer failed err=%v", err)
		return
	}
	xrlog.DetailLog.Info(ctx, "热更新---RestartAllServer成功,旧进程进入第二阶段")
	xrlog.DetailLog.Info(ctx, "hot_restart RestartAllServer ok")

	oldProcessAfterRestart.Store(true)
	waitChildStartupReady(ctx)
	if enterRestartPhaseFunc != nil {
		enterRestartPhaseFunc()
	}
	logPhase2Start(ctx)

	tryExitOldProcess()
}

var oldProcessExitOnce sync.Once

// tryExitOldProcess 旧进程固定等待后写第二阶段结束日志并退出;Run 返回与热重启协程均可触发,只执行一次.
func tryExitOldProcess() {
	if !oldProcessAfterRestart.Load() {
		return
	}
	oldProcessExitOnce.Do(func() {
		ctx := gctx.New()
		time.Sleep(time.Duration(hotRestartExitWaitSec()) * time.Second)

		shutdown.EnableSyncLoggers()
		logPhase2WaitEndFlush(ctx)
		stopOldProcessContinuousFlush(ctx)
		if syndb.AllCachesIdle() {
			logPhase2QueueEmpty(ctx)
		}
		logPhase2End(ctx)
		os.Exit(0)
	})
}

// NotifyOldProcessExit httpServer.Run 返回时触发(与热重启协程 tryExitOldProcess 互为补充).
func NotifyOldProcessExit() {
	tryExitOldProcess()
}

func hotRestartAuth() string {
	if hotRestartAuthProvider != nil {
		if auth := strings.TrimSpace(hotRestartAuthProvider()); auth != "" {
			return auth
		}
	}
	return defaultHotRestartAuth
}

func hotRestartPhase1WaitSec() int {
	if hotRestartPhase1WaitProvider != nil {
		if sec := hotRestartPhase1WaitProvider(); sec > 0 {
			return sec
		}
	}
	return defaultHotRestartPhase1Wait
}

func hotRestartExitWaitSec() int {
	if hotRestartExitWaitProvider != nil {
		if sec := hotRestartExitWaitProvider(); sec > 0 {
			return sec
		}
	}
	return defaultHotRestartExitWait
}
