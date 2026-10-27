package main

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"
)

// processWatchState tracks continuous "process missing" window for watchdog cold start.
type processWatchState struct {
	mu           sync.Mutex
	missingSince time.Time
	inFlight     bool
}

func (s *processWatchState) notePresent() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.missingSince = time.Time{}
}

// noteMissing records first miss time; returns when continuous missing started.
func (s *processWatchState) noteMissing(now time.Time) time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.missingSince.IsZero() {
		s.missingSince = now
	}
	return s.missingSince
}

func (s *processWatchState) tryBeginCold() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.inFlight {
		return false
	}
	s.inFlight = true
	return true
}

func (s *processWatchState) endCold(resetMissing bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.inFlight = false
	if resetMissing {
		s.missingSince = time.Time{}
	}
}

// shouldWatchdogColdStart is true when process has been missing continuously for >= threshold.
func shouldWatchdogColdStart(missingSince, now time.Time, threshold time.Duration) bool {
	if missingSince.IsZero() || threshold <= 0 {
		return false
	}
	return !now.Before(missingSince.Add(threshold))
}

func processWatchEnabled() bool {
	if cfg.ProcessWatchEnabled != nil {
		return *cfg.ProcessWatchEnabled
	}
	// Default on when cold start script is configured (prod recovery).
	return strings.TrimSpace(cfg.ColdStartScript) != ""
}

func startProcessWatch() {
	if !processWatchEnabled() {
		log.Println("process watch: disabled")
		return
	}
	if strings.TrimSpace(cfg.ColdStartScript) == "" {
		log.Println("process watch: skipped (coldStartScript empty)")
		return
	}
	intervalSec := cfg.ProcessWatchIntervalSec
	if intervalSec <= 0 {
		intervalSec = 3
	}
	missingSec := cfg.ProcessWatchMissingSec
	if missingSec <= 0 {
		missingSec = 180 // 3 minutes
	}
	interval := time.Duration(intervalSec) * time.Second
	threshold := time.Duration(missingSec) * time.Second
	state := &processWatchState{}
	log.Printf("process watch: every %s, cold start if %s missing continuously", interval, threshold)
	go runProcessWatchLoop(state, interval, threshold)
}

func runProcessWatchLoop(state *processWatchState, interval, threshold time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		probeOnce(state, threshold)
	}
}

func probeOnce(state *processWatchState, threshold time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	appName := strings.TrimSpace(cfg.AppName)
	if appName == "" {
		appName = "xr-game-server"
	}
	_, count, err := pgrepApp(ctx, appName)
	if err != nil {
		log.Printf("process watch: pgrep %s: %v", appName, err)
		return
	}
	now := time.Now()
	if count > 0 {
		state.notePresent()
		return
	}

	missingSince := state.noteMissing(now)
	if !shouldWatchdogColdStart(missingSince, now, threshold) {
		gone := now.Sub(missingSince).Truncate(time.Second)
		log.Printf("process watch: %s not running for %s (cold after %s)", appName, gone, threshold)
		return
	}
	if !state.tryBeginCold() {
		return
	}

	// Share rate limit with HTTP triggers so we don't thrash start.sh.
	triggerMu.Lock()
	if since := time.Since(lastTrigger); since < time.Duration(cfg.MinIntervalSec)*time.Second {
		left := cfg.MinIntervalSec - int(since.Seconds())
		triggerMu.Unlock()
		state.endCold(false) // keep missingSince; retry after rate window
		log.Printf("process watch: cold start deferred (rate limited, %ds left)", left)
		return
	}
	triggerMu.Unlock()

	log.Printf("process watch: %s missing >= %s, starting cold start", appName, threshold)
	detectedAt := time.Now()
	triggerMu.Lock()
	lastTrigger = detectedAt // reserve slot even if start.sh fails (avoid thrash)
	triggerMu.Unlock()

	coldCtx, coldCancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer coldCancel()
	result, coldErr := doColdStart(coldCtx, false, 0)
	coldOK := coldErr == nil
	recordProcessDownEvent(appName, missingSince, detectedAt, coldOK, result.PID, formatDownDetail(coldOK, result, coldErr))
	if coldErr != nil {
		state.endCold(false) // keep missing clock; retry after minInterval
		log.Printf("process watch: cold start failed: %v output=%s", coldErr, result.Output)
		return
	}
	state.endCold(true)
	if result.Output == "already running" {
		log.Printf("process watch: cold start skipped, already running pid=%d", result.PID)
		return
	}
	log.Printf("process watch: cold start ok pid=%d", result.PID)
}
