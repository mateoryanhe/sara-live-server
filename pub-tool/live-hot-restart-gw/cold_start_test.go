package main

import (
	"errors"
	"net"
	"testing"
	"time"
)

func TestShouldColdFallbackHot(t *testing.T) {
	if shouldColdFallbackHot(nil, `{"ok":true}`) {
		t.Fatal("ok body should not fallback")
	}
	if !shouldColdFallbackHot(errors.New("dial tcp 127.0.0.1:443: connect: connection refused"), "") {
		t.Fatal("connection refused should fallback")
	}
	var ne net.OpError
	_ = ne
	if !shouldColdFallbackHot(errors.New("upstream HTTP 502"), "") {
		t.Fatal("502 should fallback")
	}
}

func TestShouldWatchdogColdStart(t *testing.T) {
	now := time.Date(2026, 10, 26, 12, 0, 0, 0, time.UTC)
	threshold := 5 * time.Minute

	if shouldWatchdogColdStart(time.Time{}, now, threshold) {
		t.Fatal("zero missingSince should not trigger")
	}
	if shouldWatchdogColdStart(now.Add(-4*time.Minute), now, threshold) {
		t.Fatal("4m missing should not trigger")
	}
	if !shouldWatchdogColdStart(now.Add(-5*time.Minute), now, threshold) {
		t.Fatal("exactly 5m missing should trigger")
	}
	if !shouldWatchdogColdStart(now.Add(-6*time.Minute), now, threshold) {
		t.Fatal("6m missing should trigger")
	}
	if shouldWatchdogColdStart(now.Add(-10*time.Minute), now, 0) {
		t.Fatal("zero threshold should not trigger")
	}
}

func TestProcessWatchStateMissingWindow(t *testing.T) {
	st := &processWatchState{}
	t0 := time.Date(2026, 10, 26, 12, 0, 0, 0, time.UTC)
	first := st.noteMissing(t0)
	if !first.Equal(t0) {
		t.Fatalf("first miss want %v got %v", t0, first)
	}
	second := st.noteMissing(t0.Add(2 * time.Minute))
	if !second.Equal(t0) {
		t.Fatalf("continuous miss should keep first time, got %v", second)
	}
	st.notePresent()
	again := st.noteMissing(t0.Add(3 * time.Minute))
	if !again.Equal(t0.Add(3 * time.Minute)) {
		t.Fatalf("after present, miss clock should reset, got %v", again)
	}
}

func TestProcessWatchEnabledDefault(t *testing.T) {
	old := cfg
	t.Cleanup(func() { cfg = old })

	cfg = config{ColdStartScript: "/home/ec2-user/xgameserver/start.sh"}
	if !processWatchEnabled() {
		t.Fatal("should auto-enable when coldStartScript set")
	}
	off := false
	cfg.ProcessWatchEnabled = &off
	if processWatchEnabled() {
		t.Fatal("explicit false should disable")
	}
	on := true
	cfg = config{ProcessWatchEnabled: &on}
	if !processWatchEnabled() {
		t.Fatal("explicit true should enable even without script")
	}
	cfg = config{}
	if processWatchEnabled() {
		t.Fatal("empty config should stay off")
	}
}
