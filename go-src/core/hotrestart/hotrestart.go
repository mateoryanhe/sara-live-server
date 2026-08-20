package hotrestart

import (
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
	defaultHotRestartPhase1Wait = 60
	defaultHotRestartExitWait   = 60
	defaultHotRestartAuth       = "nGH66S4TjBjQqCKyWJAM"
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
		logPhase1Start(ctx)

		event.Pub(event.PrepareRestart, nil)
		xrtimer.Pause()

		startPhase1FlushBackground(ctx)
		time.Sleep(time.Duration(hotRestartPhase1WaitSec()) * time.Second)
		if syndb.AllCachesIdle() {
			logPhase1QueueEmpty(ctx)
		}
		logPhase1End(ctx)

		shutdown.EnableSyncLoggers()

		if err := ghttp.RestartAllServer(ctx, ""); err != nil {
			xrlog.DetailLog.Errorf(ctx, "热更新---RestartAllServer失败,err=%v", err)
			return
		}

		oldProcessAfterRestart.Store(true)
		if enterRestartPhaseFunc != nil {
			enterRestartPhaseFunc()
		}
		logPhase2Start(ctx)

		tryExitOldProcess()
	})
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
		flushPhase2Once(ctx)
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
	serverCfg := cfg.GetServerCfg()
	if serverCfg == nil || serverCfg.HotRestartFlushTimeout <= common.Zero {
		return defaultHotRestartPhase1Wait
	}
	return serverCfg.HotRestartFlushTimeout
}

func hotRestartExitWaitSec() int {
	serverCfg := cfg.GetServerCfg()
	if serverCfg == nil || serverCfg.HotRestartExitTimeout <= common.Zero {
		return defaultHotRestartExitWait
	}
	return serverCfg.HotRestartExitTimeout
}
