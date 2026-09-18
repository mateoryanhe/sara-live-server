package game

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/xrpool"
	"xr-game-server/dto/gameplatformdto"
)

const vendorGameSyncJobTimeout = 30 * time.Minute

type vendorGameSyncJobState struct {
	mu           sync.RWMutex
	running      bool
	success      bool
	count        int
	errorMessage string
}

type vendorGameSyncJobSnapshot struct {
	running      bool
	success      bool
	count        int
	errorMessage string
}

var currentVendorGameSyncJob vendorGameSyncJobState

func startVendorGameSyncJob() (*gameplatformdto.ReloadVendorGameCacheRes, bool) {
	currentVendorGameSyncJob.mu.Lock()
	if currentVendorGameSyncJob.running {
		snapshot := currentVendorGameSyncJob.snapshotLocked()
		currentVendorGameSyncJob.mu.Unlock()
		return vendorGameSyncJobResponse(snapshot, false), false
	}
	currentVendorGameSyncJob.running = true
	currentVendorGameSyncJob.success = false
	currentVendorGameSyncJob.count = 0
	currentVendorGameSyncJob.errorMessage = ""
	snapshot := currentVendorGameSyncJob.snapshotLocked()
	currentVendorGameSyncJob.mu.Unlock()

	xrpool.AddWithRecover(gctx.New(), func(ctx context.Context) {
		runVendorGameSyncJob(ctx)
	})
	return vendorGameSyncJobResponse(snapshot, true), true
}

func runVendorGameSyncJob(ctx context.Context) {
	jobCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), vendorGameSyncJobTimeout)
	defer cancel()

	vendorDetailLog().Infof(jobCtx, "sync vendor game library background job started")
	count := 0
	var syncErr error
	defer func() {
		if recovered := recover(); recovered != nil {
			syncErr = fmt.Errorf("panic: %v", recovered)
		}
		finishVendorGameSyncJob(jobCtx, count, syncErr)
	}()
	count, syncErr = SyncVendorGameLibraryFromVendor(jobCtx)
}

func finishVendorGameSyncJob(ctx context.Context, count int, syncErr error) {
	currentVendorGameSyncJob.mu.Lock()
	currentVendorGameSyncJob.running = false
	currentVendorGameSyncJob.count = count
	currentVendorGameSyncJob.success = syncErr == nil
	if syncErr != nil {
		currentVendorGameSyncJob.errorMessage = syncErr.Error()
	} else {
		currentVendorGameSyncJob.errorMessage = ""
	}
	currentVendorGameSyncJob.mu.Unlock()

	if syncErr != nil {
		vendorDetailLog().Errorf(ctx, "sync vendor game library background job failed err=%v", syncErr)
		return
	}
	vendorDetailLog().Infof(ctx, "sync vendor game library background job completed total=%d", count)
}

func getVendorGameSyncJobResponse() *gameplatformdto.ReloadVendorGameCacheRes {
	currentVendorGameSyncJob.mu.RLock()
	snapshot := currentVendorGameSyncJob.snapshotLocked()
	currentVendorGameSyncJob.mu.RUnlock()
	return vendorGameSyncJobResponse(snapshot, false)
}

func (s *vendorGameSyncJobState) snapshotLocked() vendorGameSyncJobSnapshot {
	return vendorGameSyncJobSnapshot{
		running:      s.running,
		success:      s.success,
		count:        s.count,
		errorMessage: s.errorMessage,
	}
}

func vendorGameSyncJobResponse(snapshot vendorGameSyncJobSnapshot, started bool) *gameplatformdto.ReloadVendorGameCacheRes {
	return &gameplatformdto.ReloadVendorGameCacheRes{
		Success:      snapshot.success,
		Count:        snapshot.count,
		Running:      snapshot.running,
		Started:      started,
		ErrorMessage: snapshot.errorMessage,
	}
}
