package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func pgrepApp(ctx context.Context, appName string) (newestPID int, count int, err error) {
	appName = strings.TrimSpace(appName)
	if appName == "" {
		return 0, 0, errors.New("appName empty")
	}
	cmd := exec.CommandContext(ctx, "pgrep", "-x", appName)
	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			return 0, 0, nil
		}
		return 0, 0, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	pids := make([]int, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		pid, parseErr := strconv.Atoi(line)
		if parseErr != nil {
			continue
		}
		pids = append(pids, pid)
	}
	if len(pids) == 0 {
		return 0, 0, nil
	}
	maxPID := pids[0]
	for _, pid := range pids[1:] {
		if pid > maxPID {
			maxPID = pid
		}
	}
	return maxPID, len(pids), nil
}

func stopAppProcesses(ctx context.Context, appName string, preferPID int) error {
	if preferPID > 0 {
		_ = exec.CommandContext(ctx, "sudo", "kill", "-TERM", strconv.Itoa(preferPID)).Run()
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) {
			if err := exec.CommandContext(ctx, "kill", "-0", strconv.Itoa(preferPID)).Run(); err != nil {
				return nil
			}
			time.Sleep(500 * time.Millisecond)
		}
		_ = exec.CommandContext(ctx, "sudo", "kill", "-KILL", strconv.Itoa(preferPID)).Run()
		time.Sleep(time.Second)
	}
	_ = exec.CommandContext(ctx, "sudo", "pkill", "-9", "-x", appName).Run()
	time.Sleep(time.Second)
	return nil
}

func runColdStartScript(ctx context.Context, script string) (string, error) {
	script = strings.TrimSpace(script)
	if script == "" {
		return "", errors.New("coldStartScript not configured")
	}
	if runtime.GOOS != "linux" {
		return "", errors.New("cold start only supported on linux")
	}
	info, err := os.Stat(script)
	if err != nil {
		return "", fmt.Errorf("stat coldStartScript: %w", err)
	}
	if info.IsDir() {
		return "", fmt.Errorf("coldStartScript is a directory: %s", script)
	}
	cmd := exec.CommandContext(ctx, "sudo", script)
	cmd.Dir = filepath.Dir(script)
	out, err := cmd.CombinedOutput()
	text := strings.TrimSpace(string(out))
	if err != nil {
		if text != "" {
			return text, fmt.Errorf("sudo %s: %w: %s", script, err, text)
		}
		return text, fmt.Errorf("sudo %s: %w", script, err)
	}
	return text, nil
}

func waitForApp(ctx context.Context, appName string, timeout time.Duration) (pid int, err error) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if err := ctx.Err(); err != nil {
			return 0, err
		}
		pid, count, probeErr := pgrepApp(ctx, appName)
		if probeErr != nil {
			return 0, probeErr
		}
		if count > 0 && pid > 0 {
			return pid, nil
		}
		time.Sleep(time.Second)
	}
	return 0, fmt.Errorf("timed out waiting for %s after %s", appName, timeout)
}

type coldStartResult struct {
	PID    int
	Output string
}

func doColdStart(ctx context.Context, killExisting bool, knownPID int) (coldStartResult, error) {
	script := strings.TrimSpace(cfg.ColdStartScript)
	appName := strings.TrimSpace(cfg.AppName)
	if appName == "" {
		appName = "xr-game-server"
	}
	if killExisting {
		if err := stopAppProcesses(ctx, appName, knownPID); err != nil {
			return coldStartResult{}, err
		}
	} else {
		pid, count, err := pgrepApp(ctx, appName)
		if err != nil {
			return coldStartResult{}, err
		}
		if count > 0 {
			return coldStartResult{PID: pid, Output: "already running"}, nil
		}
	}
	out, err := runColdStartScript(ctx, script)
	if err != nil {
		return coldStartResult{Output: out}, err
	}
	waitSec := cfg.ColdStartWaitSec
	if waitSec <= 0 {
		waitSec = 45
	}
	pid, err := waitForApp(ctx, appName, time.Duration(waitSec)*time.Second)
	if err != nil {
		return coldStartResult{Output: out}, err
	}
	return coldStartResult{PID: pid, Output: out}, nil
}

func shouldColdFallbackHot(err error, upstream string) bool {
	body := strings.ToLower(upstream)
	if strings.Contains(body, "unauthorized") || strings.Contains(body, "already in progress") {
		return false
	}
	if err == nil {
		if strings.Contains(body, `"ok":true`) || strings.Contains(body, `"ok": true`) {
			return false
		}
		return true
	}
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return true
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "http 403") || strings.Contains(msg, "http 401") {
		return false
	}
	if strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "connect: connection") ||
		strings.Contains(msg, "no route") ||
		strings.Contains(msg, "upstream http 502") ||
		strings.Contains(msg, "upstream http 503") ||
		strings.Contains(msg, "upstream http 504") {
		return true
	}
	return false
}
