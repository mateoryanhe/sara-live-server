package luckygiftpoolcfg

import (
	"context"
	"sync/atomic"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/event"
	"xr-game-server/core/xrpool"
	"xr-game-server/errercode"
)

const luckyGiftPoolQueueSize = 2048

type jobKind uint8

const (
	jobDraw jobKind = iota + 1
	jobApplyConfig
	jobRefundPayout
	jobMergePending
)

type applyConfigPayload struct {
	cfgID             uint64
	poolBalance       float64
	contributePercent float64
	winProbability    float64
	tiers             []*tierInput
}

type tierInput struct {
	Multiplier uint32
	Weight     uint32
	Type       string
}

type poolJob struct {
	kind         jobKind
	unitPrice    float64
	count        int
	apply        *applyConfigPayload
	refundAmount float64
	resp         chan poolJobResult
}

type poolJobResult struct {
	draws []DrawResult
	cfgID uint64
	err   error
}

var (
	poolJobCh       chan *poolJob
	poolWorkerOn    = make(chan struct{})
	acceptingJobsOn atomic.Bool
)

func startPoolQueueWorker() {
	poolJobCh = make(chan *poolJob, luckyGiftPoolQueueSize)
	acceptingJobsOn.Store(true)
	// 启动至延时合并完成前：强制不中奖，进池走临时表
	setSettleNoWin(true)
	resetStuckPendingOnStart()

	xrpool.AddWithRecover(gctx.New(), func(ctx context.Context) {
		close(poolWorkerOn)
		for job := range poolJobCh {
			if job == nil {
				continue
			}
			handlePoolJob(job)
		}
	})
	<-poolWorkerOn

	event.Sub(event.PrepareRestart, onPrepareRestartDrainQueue)
	event.Sub(event.ColdStart, onServerStartScheduleMerge)
	event.Sub(event.HotStart, onServerStartScheduleMerge)
}

func onPrepareRestartDrainQueue(_ any) {
	acceptingJobsOn.Store(false)
	setSettleNoWin(true)
	drainPoolQueueAsNoWin()
}

func onServerStartScheduleMerge(_ any) {
	schedulePendingMergeAfterStartOnce()
}

func handlePoolJob(job *poolJob) {
	res := poolJobResult{}
	switch job.kind {
	case jobDraw:
		res.draws = processBatchDrawLocked(job.unitPrice, job.count)
	case jobApplyConfig:
		if job.apply == nil {
			res.err = errercode.CreateCode(errercode.InvalidParam)
		} else {
			res.cfgID, res.err = processApplyConfigLocked(job.apply)
		}
	case jobRefundPayout:
		processRefundPayoutLocked(job.refundAmount)
	case jobMergePending:
		processMergePendingLocked()
	default:
		res.err = errercode.CreateCode(errercode.InvalidParam)
	}
	if job.resp != nil {
		job.resp <- res
		close(job.resp)
	}
}

func enqueueAndWait(job *poolJob) poolJobResult {
	if job == nil {
		return poolJobResult{err: errercode.CreateCode(errercode.InvalidParam)}
	}
	// 合并任务在停收后仍允许，保证启动延时合并可执行
	if job.kind != jobMergePending && (!acceptingJobsOn.Load() || poolJobCh == nil) {
		return poolJobResult{err: errercode.CreateCode(errercode.InvalidParam)}
	}
	if poolJobCh == nil {
		return poolJobResult{err: errercode.CreateCode(errercode.InvalidParam)}
	}
	job.resp = make(chan poolJobResult, 1)
	select {
	case poolJobCh <- job:
	default:
		poolJobCh <- job
	}
	return <-job.resp
}

func errInvalidParam() error {
	return errercode.CreateCode(errercode.InvalidParam)
}

// EnqueueBatchDraw 串行入队：count 次独立抽奖后返回明细(再由 gift 侧合并)
func EnqueueBatchDraw(unitPrice float64, count int) []DrawResult {
	if count <= 0 || unitPrice < 0 {
		return nil
	}
	res := enqueueAndWait(&poolJob{
		kind:      jobDraw,
		unitPrice: unitPrice,
		count:     count,
	})
	if res.err != nil {
		return nil
	}
	return res.draws
}

// EnqueueRefundPayout 钱包加钻失败时退回已扣奖池金额
func EnqueueRefundPayout(amount float64) {
	if amount <= 0 {
		return
	}
	_ = enqueueAndWait(&poolJob{
		kind:         jobRefundPayout,
		refundAmount: amount,
	})
}

func enqueueApplyConfig(p *applyConfigPayload) (uint64, error) {
	res := enqueueAndWait(&poolJob{
		kind:  jobApplyConfig,
		apply: p,
	})
	return res.cfgID, res.err
}

func processRefundPayoutLocked(amount float64) {
	if amount <= 0 {
		return
	}
	if isSettleNoWin() {
		enqueuePendingPoolCredit(amount)
		runtimeMu.Lock()
		ensureRuntimeCfgLocked().SubTotalPayoutAmount(amount)
		runtimeMu.Unlock()
		return
	}
	runtimeMu.Lock()
	defer runtimeMu.Unlock()
	setRuntimePoolBalanceLocked(runtimePoolBalanceLocked() + amount)
	ensureRuntimeCfgLocked().SubTotalPayoutAmount(amount)
	storePublicSnapshotLocked()
}

// drainPoolQueueAsNoWin 排空在途抽奖：强制不中奖，进池进临时表
func drainPoolQueueAsNoWin() {
	if poolJobCh == nil {
		return
	}
	for {
		select {
		case job := <-poolJobCh:
			if job == nil {
				continue
			}
			switch job.kind {
			case jobDraw:
				res := poolJobResult{draws: processBatchDrawLocked(job.unitPrice, job.count)}
				if job.resp != nil {
					job.resp <- res
					close(job.resp)
				}
			case jobRefundPayout:
				enqueuePendingPoolCredit(job.refundAmount)
				if job.resp != nil {
					job.resp <- poolJobResult{}
					close(job.resp)
				}
			default:
				handlePoolJob(job)
			}
		default:
			return
		}
	}
}
