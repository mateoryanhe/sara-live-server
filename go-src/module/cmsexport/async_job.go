package cmsexport

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"xr-game-server/core/httpserver"
	"xr-game-server/core/xrlog"
	"xr-game-server/core/xrpool"
	"xr-game-server/core/xrtimer"
	"xr-game-server/dto/cmsexportdto"
	"xr-game-server/errercode"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gogf/gf/v2/util/gutil"
)

const (
	exportJobTimeout   = 30 * time.Minute
	exportJobTTL       = 40 * time.Minute
	exportJobQueueSize = 32

	exportJobStatusPending = "pending"
	exportJobStatusRunning = "running"
	exportJobStatusDone    = "done"
	exportJobStatusFailed  = "failed"
)

type exportJob struct {
	id         string
	exportType string
	payload    json.RawMessage
	cmsUserId  uint64
	createdAt  time.Time
}

type exportJobState struct {
	mu            sync.RWMutex
	id            string
	exportType    string
	status        string
	queuePosition int
	errorMessage  string
	progress      *cmsexportdto.CMSExportJobProgress
	result        any
	createdAt     time.Time
	updatedAt     time.Time
}

var (
	exportJobStates    sync.Map
	exportJobInitOnce  sync.Once
	exportJobSerialMu  sync.Mutex
	exportPendingMu    sync.Mutex
	exportPendingCount int
)

func initExportJobWorker() {
	exportJobInitOnce.Do(func() {
		xrtimer.AddSingleton(gctx.New(), time.Minute, cleanupExportJobStates)
	})
}

func incExportPending() int {
	exportPendingMu.Lock()
	defer exportPendingMu.Unlock()
	exportPendingCount++
	return exportPendingCount
}

func decExportPending() {
	exportPendingMu.Lock()
	defer exportPendingMu.Unlock()
	if exportPendingCount > 0 {
		exportPendingCount--
	}
}

func SubmitExportJob(ctx context.Context, req *cmsexportdto.CMSSubmitExportJobReq) (*cmsexportdto.CMSSubmitExportJobRes, error) {
	if req == nil || req.ExportType == "" || len(req.Payload) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if !isSupportedExportType(req.ExportType) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := ensureExportReady(); err != nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	initExportJobWorker()

	queuePosition := incExportPending()
	if queuePosition > exportJobQueueSize {
		decExportPending()
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	jobID := guid.S()
	now := time.Now()
	state := &exportJobState{
		id:            jobID,
		exportType:    req.ExportType,
		status:        exportJobStatusPending,
		queuePosition: queuePosition,
		createdAt:     now,
		updatedAt:     now,
	}
	exportJobStates.Store(jobID, state)

	job := &exportJob{
		id:         jobID,
		exportType: req.ExportType,
		payload:    append(json.RawMessage(nil), req.Payload...),
		cmsUserId:  httpserver.GetAuthId(ctx),
		createdAt:  now,
	}

	xrpool.AddWithRecover(gctx.New(), func(ctx context.Context) {
		defer decExportPending()
		exportJobSerialMu.Lock()
		defer exportJobSerialMu.Unlock()
		runExportJob(job)
	})

	return &cmsexportdto.CMSSubmitExportJobRes{
		JobId:         jobID,
		QueuePosition: queuePosition,
	}, nil
}

func GetExportJob(_ context.Context, req *cmsexportdto.CMSGetExportJobReq) (*cmsexportdto.CMSGetExportJobRes, error) {
	if req == nil || req.JobId == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	state, ok := loadExportJobState(req.JobId)
	if !ok {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	state.mu.RLock()
	defer state.mu.RUnlock()
	return &cmsexportdto.CMSGetExportJobRes{
		JobId:         state.id,
		ExportType:    state.exportType,
		Status:        state.status,
		QueuePosition: state.queuePosition,
		ErrorMessage:  state.errorMessage,
		Progress:      state.progress,
		Result:        state.result,
	}, nil
}

func DeleteExport(_ context.Context, req *cmsexportdto.CMSDeleteExportReq) (*cmsexportdto.CMSDeleteExportRes, error) {
	if req == nil || req.ExportId == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := deleteExport(req.ExportId); err != nil {
		return nil, err
	}
	return &cmsexportdto.CMSDeleteExportRes{Success: true}, nil
}

func isSupportedExportType(exportType string) bool {
	switch exportType {
	case cmsexportdto.ExportTypeLiveRecord,
		cmsexportdto.ExportTypeLiveRevenueLog,
		cmsexportdto.ExportTypeVideoCallLog,
		cmsexportdto.ExportTypeAnchorIncomeSettlementLog,
		cmsexportdto.ExportTypeGuildIncomeSettlementLog,
		cmsexportdto.ExportTypeGuildAnchorIncomeSettlementLog,
		cmsexportdto.ExportTypeMyGuildAnchorIncomeSettlementLog,
		cmsexportdto.ExportTypeAnchorDailyEffectiveLive,
		cmsexportdto.ExportTypeGuildDailyEffectiveLive,
		cmsexportdto.ExportTypeGuildAnchorDailyEffectiveLive,
		cmsexportdto.ExportTypeMyGuildAnchorDailyEffectiveLive,
		cmsexportdto.ExportTypeLiveDailyEffectiveLive,
		cmsexportdto.ExportTypeLiveWeeklyUnsettledLive,
		cmsexportdto.ExportTypeCurrencyLog:
		return true
	default:
		return false
	}
}

func runExportJob(job *exportJob) {
	if job == nil {
		return
	}
	state, ok := loadExportJobState(job.id)
	if !ok {
		return
	}
	updateExportJobState(state, exportJobStatusRunning, "", nil, nil)

	ctx, cancel := context.WithTimeout(context.Background(), exportJobTimeout)
	defer cancel()

	gutil.TryCatch(ctx, func(try context.Context) {
		result, err := executeExportJob(try, job.exportType, job.cmsUserId, job.payload, func(exportedRows, totalRows int) {
			updateExportJobState(state, exportJobStatusRunning, "", &cmsexportdto.CMSExportJobProgress{
				ExportedRows: exportedRows,
				TotalRows:    totalRows,
			}, nil)
		})
		if err != nil {
			updateExportJobState(state, exportJobStatusFailed, err.Error(), nil, nil)
			return
		}
		updateExportJobState(state, exportJobStatusDone, "", nil, toExportDTO(result))
	}, func(catch context.Context, exception error) {
		xrlog.ErrorWithErr(catch, "CMSExportJob", job.exportType, exception)
		msg := "导出异常"
		if exception != nil {
			msg = exception.Error()
		}
		updateExportJobState(state, exportJobStatusFailed, msg, nil, nil)
	})
}

func executeExportJob(ctx context.Context, exportType string, cmsUserId uint64, payload json.RawMessage, onProgress func(exportedRows, totalRows int)) (*exportResult, error) {
	switch exportType {
	case cmsexportdto.ExportTypeLiveRecord:
		return exportLiveRecordCSV(ctx, payload, onProgress)
	case cmsexportdto.ExportTypeLiveRevenueLog:
		return exportLiveRevenueLogCSV(ctx, payload, onProgress)
	case cmsexportdto.ExportTypeVideoCallLog:
		return exportVideoCallLogCSV(ctx, payload, onProgress)
	case cmsexportdto.ExportTypeAnchorIncomeSettlementLog:
		return exportAnchorIncomeSettlementLogCSV(ctx, payload, onProgress)
	case cmsexportdto.ExportTypeGuildIncomeSettlementLog:
		return exportGuildIncomeSettlementLogCSV(ctx, payload, onProgress)
	case cmsexportdto.ExportTypeGuildAnchorIncomeSettlementLog:
		return exportGuildAnchorIncomeSettlementLogCSV(ctx, payload, onProgress)
	case cmsexportdto.ExportTypeMyGuildAnchorIncomeSettlementLog:
		return exportMyGuildAnchorIncomeSettlementLogCSV(ctx, cmsUserId, payload, onProgress)
	case cmsexportdto.ExportTypeAnchorDailyEffectiveLive:
		return exportAnchorDailyEffectiveLiveCSV(ctx, payload, onProgress)
	case cmsexportdto.ExportTypeGuildDailyEffectiveLive:
		return exportGuildDailyEffectiveLiveCSV(ctx, payload, onProgress)
	case cmsexportdto.ExportTypeGuildAnchorDailyEffectiveLive:
		return exportGuildAnchorDailyEffectiveLiveCSV(ctx, payload, onProgress)
	case cmsexportdto.ExportTypeMyGuildAnchorDailyEffectiveLive:
		return exportMyGuildAnchorDailyEffectiveLiveCSV(ctx, cmsUserId, payload, onProgress)
	case cmsexportdto.ExportTypeLiveDailyEffectiveLive:
		return exportLiveDailyEffectiveLiveCSV(ctx, payload, onProgress)
	case cmsexportdto.ExportTypeLiveWeeklyUnsettledLive:
		return exportLiveWeeklyUnsettledLiveCSV(ctx, payload, onProgress)
	case cmsexportdto.ExportTypeCurrencyLog:
		return exportCurrencyLogCSV(ctx, payload, onProgress)
	default:
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
}

func toExportDTO(result *exportResult) *cmsexportdto.CMSExportResult {
	if result == nil {
		return nil
	}
	return &cmsexportdto.CMSExportResult{
		ExportId: result.ExportID,
		FileName: result.FileName,
		FileUrl:  result.FileUrl,
		Total:    result.Total,
	}
}

func loadExportJobState(jobID string) (*exportJobState, bool) {
	value, ok := exportJobStates.Load(jobID)
	if !ok {
		return nil, false
	}
	state, ok := value.(*exportJobState)
	return state, ok
}

func updateExportJobState(state *exportJobState, status, errMsg string, progress *cmsexportdto.CMSExportJobProgress, result any) {
	if state == nil {
		return
	}
	state.mu.Lock()
	defer state.mu.Unlock()
	state.status = status
	state.errorMessage = errMsg
	if progress != nil {
		state.progress = progress
	}
	if result != nil {
		state.result = result
	}
	state.updatedAt = time.Now()
}

func cleanupExportJobStates(_ context.Context) {
	expireBefore := time.Now().Add(-exportJobTTL)
	exportJobStates.Range(func(key, value any) bool {
		state, ok := value.(*exportJobState)
		if !ok || state == nil {
			exportJobStates.Delete(key)
			return true
		}
		state.mu.RLock()
		expired := state.updatedAt.Before(expireBefore)
		state.mu.RUnlock()
		if expired {
			exportJobStates.Delete(key)
		}
		return true
	})
}
