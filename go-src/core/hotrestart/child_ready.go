package hotrestart

import (
	"bytes"
	"context"
	"sync"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gproc"
	"xr-game-server/core/xrlog"
	"xr-game-server/core/xrpool"
)

const (
	hotRestartProcGroup     = "XR_HOT_RESTART"
	childStartupReadyPayload = "startup_ok"
	defaultChildReadyWaitSec = 180
)

var (
	childReadyListenerOnce sync.Once
	childReadyNotifyCh     = make(chan struct{}, 1)
)

// ensureChildReadyListener 旧进程在 fork 前注册 gproc 组,供子进程上报启动完成.
func ensureChildReadyListener() {
	childReadyListenerOnce.Do(func() {
		xrpool.AddWithRecover(gctx.New(), func(_ context.Context) {
			for {
				msg := gproc.Receive(hotRestartProcGroup)
				if msg == nil {
					continue
				}
				if !bytes.Equal(msg.Data, []byte(childStartupReadyPayload)) {
					continue
				}
				select {
				case childReadyNotifyCh <- struct{}{}:
				default:
				}
			}
		})
	})
}

func drainChildReadyNotify() {
	for {
		select {
		case <-childReadyNotifyCh:
		default:
			return
		}
	}
}

// NotifyChildStartupComplete 热重启子进程在 HTTP 就绪且业务初始化完成后通知父进程.
func NotifyChildStartupComplete() {
	if !gproc.IsChild() {
		return
	}
	ctx := gctx.New()
	ppid := gproc.PPid()
	if ppid <= 0 || ppid == gproc.Pid() {
		return
	}
	if err := gproc.Send(ppid, []byte(childStartupReadyPayload), hotRestartProcGroup); err != nil {
		xrlog.DetailLog.Errorf(ctx, "热重启子进程通知父进程启动完成失败,ppid=%d err=%v", ppid, err)
		xrlog.DetailLog.Errorf(ctx, "hot_restart child notify parent failed ppid=%d err=%v", ppid, err)
		return
	}
	xrlog.DetailLog.Info(ctx, "热重启子进程已通知父进程:启动完成")
	xrlog.DetailLog.Info(ctx, "hot_restart child startup ok")
}

func waitChildStartupReady(ctx context.Context) {
	ensureChildReadyListener()
	drainChildReadyNotify()

	timeout := time.Duration(childReadyMaxWaitSec()) * time.Second
	xrlog.DetailLog.Infof(ctx, "热重启等待新进程完全启动,最长 %ds", int(timeout.Seconds()))
	xrlog.DetailLog.Infof(ctx, "hot_restart waiting child startup max_sec=%d", int(timeout.Seconds()))

	select {
	case <-childReadyNotifyCh:
		xrlog.DetailLog.Info(ctx, "热重启新进程已完全启动,开始第二阶段计时")
		xrlog.DetailLog.Info(ctx, "hot_restart child ready phase2 timer start")
	case <-time.After(timeout):
		xrlog.DetailLog.Warning(ctx, "热重启等待新进程启动超时,仍进入第二阶段")
		xrlog.DetailLog.Warning(ctx, "hot_restart child ready timeout phase2 timer start")
	}
}

func childReadyMaxWaitSec() int {
	wait := hotRestartPhase1WaitSec() + hotRestartExitWaitSec() + 60
	if wait < defaultChildReadyWaitSec {
		return defaultChildReadyWaitSec
	}
	return wait
}
