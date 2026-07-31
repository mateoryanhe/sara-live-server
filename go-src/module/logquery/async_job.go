package logquery

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"xr-game-server/core/xrlog"
	"xr-game-server/core/xrpool"
	"xr-game-server/core/xrtimer"
	"xr-game-server/dto/logquerydto"
	"xr-game-server/errercode"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gogf/gf/v2/util/gutil"
)

const (
	logQueryJobTimeout   = 30 * time.Minute
	logQueryJobTTL       = 40 * time.Minute // 需大于任务超时,避免运行中 job 被清理
	logQueryJobQueueSize = 32

	logQueryJobStatusPending = "pending"
	logQueryJobStatusRunning = "running"
	logQueryJobStatusDone    = "done"
	logQueryJobStatusFailed  = "failed"
)

const (
	LogQueryJobTypeDetail      = "detail"
	LogQueryJobTypeAccess      = "access"
	LogQueryJobTypeError       = "error"
	LogQueryJobTypeTrace       = "trace"
	LogQueryJobTypeAccessStats = "accessStats"
	LogQueryJobTypeAccessTrend = "accessTrend"
)

type logQueryJob struct {
	id        string
	queryType string
	payload   json.RawMessage
	createdAt time.Time
}

type logQueryJobState struct {
	mu            sync.RWMutex
	id            string
	queryType     string
	status        string
	queuePosition int
	errorMessage  string
	result        any
	createdAt     time.Time
	updatedAt     time.Time
}

var (
	logQueryJobStates    sync.Map
	logQueryJobInitOnce  sync.Once
	logQueryJobSerialMu  sync.Mutex
	logQueryPendingMu    sync.Mutex
	logQueryPendingCount int
)

func initLogQueryJobWorker() {
	logQueryJobInitOnce.Do(func() {
		xrtimer.AddSingleton(gctx.New(), time.Minute, cleanupLogQueryJobStates)
	})
}

func incLogQueryPending() int {
	logQueryPendingMu.Lock()
	defer logQueryPendingMu.Unlock()
	logQueryPendingCount++
	return logQueryPendingCount
}

func decLogQueryPending() {
	logQueryPendingMu.Lock()
	defer logQueryPendingMu.Unlock()
	if logQueryPendingCount > 0 {
		logQueryPendingCount--
	}
}

func SubmitLogQueryJob(_ context.Context, req *logquerydto.CMSSubmitLogQueryJobReq) (*logquerydto.CMSSubmitLogQueryJobRes, error) {
	if req == nil || req.QueryType == "" || len(req.Payload) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if !isSupportedLogQueryJobType(req.QueryType) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	initLogQueryJobWorker()

	queuePosition := incLogQueryPending()
	if queuePosition > logQueryJobQueueSize {
		decLogQueryPending()
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	jobID := guid.S()
	now := time.Now()
	state := &logQueryJobState{
		id:            jobID,
		queryType:     req.QueryType,
		status:        logQueryJobStatusPending,
		queuePosition: queuePosition,
		createdAt:     now,
		updatedAt:     now,
	}
	logQueryJobStates.Store(jobID, state)

	job := &logQueryJob{
		id:        jobID,
		queryType: req.QueryType,
		payload:   append(json.RawMessage(nil), req.Payload...),
		createdAt: now,
	}

	xrpool.AddWithRecover(gctx.New(), func(ctx context.Context) {
		defer decLogQueryPending()
		logQueryJobSerialMu.Lock()
		defer logQueryJobSerialMu.Unlock()
		runLogQueryJob(job)
	})

	return &logquerydto.CMSSubmitLogQueryJobRes{
		JobId:         jobID,
		QueuePosition: queuePosition,
	}, nil
}

func GetLogQueryJob(_ context.Context, req *logquerydto.CMSGetLogQueryJobReq) (*logquerydto.CMSGetLogQueryJobRes, error) {
	if req == nil || req.JobId == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	state, ok := loadLogQueryJobState(req.JobId)
	if !ok {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	state.mu.RLock()
	defer state.mu.RUnlock()
	return &logquerydto.CMSGetLogQueryJobRes{
		JobId:         state.id,
		QueryType:     state.queryType,
		Status:        state.status,
		QueuePosition: state.queuePosition,
		ErrorMessage:  state.errorMessage,
		Result:        state.result,
	}, nil
}

func isSupportedLogQueryJobType(queryType string) bool {
	switch queryType {
	case LogQueryJobTypeDetail, LogQueryJobTypeAccess, LogQueryJobTypeError,
		LogQueryJobTypeTrace, LogQueryJobTypeAccessStats, LogQueryJobTypeAccessTrend:
		return true
	default:
		return false
	}
}

func runLogQueryJob(job *logQueryJob) {
	if job == nil {
		return
	}
	state, ok := loadLogQueryJobState(job.id)
	if !ok {
		return
	}
	updateLogQueryJobState(state, logQueryJobStatusRunning, "", nil)

	ctx, cancel := context.WithTimeout(context.Background(), logQueryJobTimeout)
	defer cancel()

	// TryCatch: panic 时更新任务状态; xrpool.AddWithRecover 写 error 日志并防止进程退出
	gutil.TryCatch(ctx, func(try context.Context) {
		result, err := executeLogQueryJob(try, job.queryType, job.payload)
		if err != nil {
			updateLogQueryJobState(state, logQueryJobStatusFailed, err.Error(), nil)
			return
		}
		updateLogQueryJobState(state, logQueryJobStatusDone, "", result)
	}, func(catch context.Context, exception error) {
		xrlog.ErrorWithErr(catch, "LogQueryJob", job.queryType, exception)
		msg := "日志查询异常"
		if exception != nil {
			msg = exception.Error()
		}
		updateLogQueryJobState(state, logQueryJobStatusFailed, msg, nil)
	})
}

func executeLogQueryJob(ctx context.Context, queryType string, payload json.RawMessage) (any, error) {
	switch queryType {
	case LogQueryJobTypeDetail:
		var req logquerydto.CMSQueryDetailLogsReq
		if err := json.Unmarshal(payload, &req); err != nil {
			return nil, err
		}
		return QueryDetailLogs(ctx, &req)
	case LogQueryJobTypeAccess:
		var req logquerydto.CMSQueryAccessLogsReq
		if err := json.Unmarshal(payload, &req); err != nil {
			return nil, err
		}
		return QueryAccessLogs(ctx, &req)
	case LogQueryJobTypeError:
		var req logquerydto.CMSQueryErrorLogsReq
		if err := json.Unmarshal(payload, &req); err != nil {
			return nil, err
		}
		return QueryErrorLogs(ctx, &req)
	case LogQueryJobTypeTrace:
		var req logquerydto.CMSGetTraceLogsReq
		if err := json.Unmarshal(payload, &req); err != nil {
			return nil, err
		}
		return GetTraceLogs(ctx, &req)
	case LogQueryJobTypeAccessStats:
		var req logquerydto.CMSGetAccessStatsReq
		if err := json.Unmarshal(payload, &req); err != nil {
			return nil, err
		}
		return GetAccessStats(ctx, &req)
	case LogQueryJobTypeAccessTrend:
		var req logquerydto.CMSGetAccessTrendReq
		if err := json.Unmarshal(payload, &req); err != nil {
			return nil, err
		}
		return GetAccessTrend(ctx, &req)
	default:
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
}

func loadLogQueryJobState(jobID string) (*logQueryJobState, bool) {
	value, ok := logQueryJobStates.Load(jobID)
	if !ok {
		return nil, false
	}
	state, ok := value.(*logQueryJobState)
	return state, ok
}

func updateLogQueryJobState(state *logQueryJobState, status, errMsg string, result any) {
	if state == nil {
		return
	}
	state.mu.Lock()
	defer state.mu.Unlock()
	state.status = status
	state.errorMessage = errMsg
	state.result = result
	state.updatedAt = time.Now()
}

func cleanupLogQueryJobStates(_ context.Context) {
	expireBefore := time.Now().Add(-logQueryJobTTL)
	logQueryJobStates.Range(func(key, value any) bool {
		state, ok := value.(*logQueryJobState)
		if !ok || state == nil {
			logQueryJobStates.Delete(key)
			return true
		}
		state.mu.RLock()
		expired := state.updatedAt.Before(expireBefore)
		state.mu.RUnlock()
		if expired {
			logQueryJobStates.Delete(key)
		}
		return true
	})
}
