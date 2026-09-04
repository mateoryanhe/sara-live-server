package stat

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/container/gqueue"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/event"
	"xr-game-server/core/hotrestart"
	"xr-game-server/core/xrtimer"
)

const statConsumeInterval = time.Second

type jobKind uint8

const (
	jobLogin jobKind = iota + 1
	jobRegister
	jobRecharge
	jobCurrency
	jobAudience
)

// statJob 统计队列任务;Payload 为原始事件数据,校验与清洗在消费端完成.
type statJob struct {
	Kind    jobKind
	At      time.Time // 入队时捕获的业务时间(可选)
	Payload any
}

var statQueue = gqueue.NewTQueue[*statJob]()

func initStatQueue() {
	hotrestart.RegisterStatQueueFlush(drainStatQueue, statQueueIdle)
	event.Sub(event.HotStart, onStatHotStart)
	event.Sub(event.PrepareRestart, onStatPrepareRestart)
}

func onStatHotStart(_ any) {
	xrtimer.AddSingleton(gctx.New(), statConsumeInterval, func(ctx context.Context) {
		drainStatQueue()
	})
}

func onStatPrepareRestart(_ any) {
	drainStatQueue()
}

func enqueue(job *statJob) {
	if job == nil {
		return
	}
	statQueue.Push(job)
}

func statQueueIdle() bool {
	return statQueue.Len() == 0
}

// drainStatQueue 非阻塞尽量拉空队列并处理.
func drainStatQueue() {
	for {
		select {
		case job := <-statQueue.C:
			if job != nil {
				handleStatJob(job)
			}
		default:
			return
		}
	}
}

func handleStatJob(job *statJob) {
	switch job.Kind {
	case jobLogin:
		consumeLoginJob(job)
	case jobRegister:
		consumeRegisterJob(job)
	case jobRecharge:
		consumeRechargeJob(job)
	case jobCurrency:
		consumeCurrencyJob(job)
	case jobAudience:
		consumeAudienceJob(job)
	}
}
