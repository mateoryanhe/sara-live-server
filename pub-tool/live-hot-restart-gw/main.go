// live-hot-restart-gw：独立进程，部署在直播服本机。
// 监听 18900：带 auth 的 HTTP 触发热重启；Go 未运行时走 start.sh 冷启（与 go-prod/restart.bat 一致）。
package main

import (
	"context"
	"crypto/subtle"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type config struct {
	Listen                  string `json:"listen"`
	Auth                    string `json:"auth"`
	HotRestartURL           string `json:"hotRestartURL"`
	MinIntervalSec          int    `json:"minIntervalSec"`
	ColdStartScript         string `json:"coldStartScript"`
	AppName                 string `json:"appName"`
	HotRestartFallbackCold  bool   `json:"hotRestartFallbackCold"`
	ColdStartWaitSec        int    `json:"coldStartWaitSec"`
	ProcessWatchEnabled     *bool  `json:"processWatchEnabled"`     // nil=auto on when coldStartScript set
	ProcessWatchIntervalSec int    `json:"processWatchIntervalSec"` // default 3
	ProcessWatchMissingSec  int    `json:"processWatchMissingSec"`  // default 180 (3min continuous miss → cold)
	MysqlDSN                string `json:"mysqlDSN"`                // 进程挂掉记录入库; 例 user:pass@tcp(127.0.0.1:13307)/live_db?parseTime=true&loc=Local
}

var (
	cfg         config
	configPath  string
	lastTrigger time.Time
	triggerMu   sync.Mutex
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	if err := loadConfig(); err != nil {
		log.Fatalf("config: %v", err)
	}
	if strings.TrimSpace(cfg.Auth) == "" {
		log.Fatal("config.json: auth 不能为空（与 CMS 热重启密钥 hotRestartAuth 一致）")
	}
	if cfg.Listen == "" {
		cfg.Listen = "0.0.0.0:18900"
	}
	if strings.TrimSpace(cfg.HotRestartURL) == "" {
		cfg.HotRestartURL = "https://127.0.0.1/internal/hotRestart"
	}
	if cfg.MinIntervalSec <= 0 {
		cfg.MinIntervalSec = 60
	}
	if strings.TrimSpace(cfg.AppName) == "" {
		cfg.AppName = "xr-game-server"
	}
	if cfg.HotRestartFallbackCold && strings.TrimSpace(cfg.ColdStartScript) == "" {
		log.Println("warn: hotRestartFallbackCold=true 但未配置 coldStartScript，进程挂掉时无法冷启")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleSmartRestart)
	mux.HandleFunc("/hotRestart", handleSmartRestart)
	mux.HandleFunc("/coldStart", handleColdStartOnly)
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	})

	initProcessDownDB()
	startProcessWatch()

	log.Printf("live-hot-restart-gw listen http://%s/  config=%s coldStart=%q",
		cfg.Listen, configPath, cfg.ColdStartScript)
	if err := http.ListenAndServe(cfg.Listen, mux); err != nil {
		log.Fatal(err)
	}
}

func loadConfig() error {
	configPath = strings.TrimSpace(os.Getenv("LIVE_HR_GW_CONFIG"))
	if configPath == "" {
		exe, err := os.Executable()
		if err != nil {
			configPath = "config.json"
		} else {
			configPath = filepath.Join(filepath.Dir(exe), "config.json")
		}
	}
	raw, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("read %s: %w (见 config.example.json)", configPath, err)
	}
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return err
	}
	if v := strings.TrimSpace(os.Getenv("LIVE_HR_GW_LISTEN")); v != "" {
		cfg.Listen = v
	}
	if v := strings.TrimSpace(os.Getenv("LIVE_HR_GW_AUTH")); v != "" {
		cfg.Auth = v
	}
	if v := strings.TrimSpace(os.Getenv("LIVE_HR_GW_HOT_RESTART_URL")); v != "" {
		cfg.HotRestartURL = v
	}
	if v := strings.TrimSpace(os.Getenv("LIVE_HR_GW_COLD_START_SCRIPT")); v != "" {
		cfg.ColdStartScript = v
	}
	if v := strings.TrimSpace(os.Getenv("LIVE_HR_GW_MYSQL_DSN")); v != "" {
		cfg.MysqlDSN = v
	}
	return nil
}

type triggerResp struct {
	OK       bool   `json:"ok"`
	Mode     string `json:"mode,omitempty"`
	PID      int    `json:"pid,omitempty"`
	Reason   string `json:"reason,omitempty"`
	Upstream string `json:"upstream,omitempty"`
}

func handleSmartRestart(w http.ResponseWriter, r *http.Request) {
	if !authorizeRequest(w, r) {
		return
	}
	if !acquireTriggerSlot(w) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 120*time.Second)
	defer cancel()

	pid, count, err := pgrepApp(ctx, cfg.AppName)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, triggerResp{OK: false, Reason: err.Error()})
		return
	}
	if count == 0 {
		respondColdStart(w, r, ctx, false, 0, "no running process")
		return
	}

	upstream, hotErr := callGameHotRestart(ctx, requestAuthToken(r))
	if hotErr == nil {
		ok := strings.Contains(upstream, `"ok":true`) || strings.Contains(upstream, `"ok": true`)
		if ok {
			lastTrigger = time.Now()
			log.Printf("hot restart ok remote=%s pid=%d", r.RemoteAddr, pid)
			writeJSON(w, http.StatusOK, triggerResp{OK: true, Mode: "hot", PID: pid, Reason: upstream, Upstream: upstream})
			return
		}
	}
	if !cfg.HotRestartFallbackCold || strings.TrimSpace(cfg.ColdStartScript) == "" {
		reason := upstream
		if hotErr != nil {
			reason = hotErr.Error()
		}
		log.Printf("hot restart failed (no cold fallback): %v upstream=%s", hotErr, upstream)
		writeJSON(w, http.StatusBadGateway, triggerResp{OK: false, Mode: "hot", Reason: reason, Upstream: upstream})
		return
	}
	if !shouldColdFallbackHot(hotErr, upstream) {
		reason := upstream
		if hotErr != nil {
			reason = hotErr.Error()
		}
		writeJSON(w, http.StatusBadGateway, triggerResp{OK: false, Mode: "hot", Reason: reason, Upstream: upstream})
		return
	}
	respondColdStart(w, r, ctx, true, pid, fmt.Sprintf("hot failed: %v", hotErr))
}

func handleColdStartOnly(w http.ResponseWriter, r *http.Request) {
	if !authorizeRequest(w, r) {
		return
	}
	if !acquireTriggerSlot(w) {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 120*time.Second)
	defer cancel()
	pid, count, _ := pgrepApp(ctx, cfg.AppName)
	kill := count > 0
	respondColdStart(w, r, ctx, kill, pid, "coldStart endpoint")
}

func respondColdStart(w http.ResponseWriter, r *http.Request, ctx context.Context, killExisting bool, knownPID int, why string) {
	if strings.TrimSpace(cfg.ColdStartScript) == "" {
		writeJSON(w, http.StatusServiceUnavailable, triggerResp{
			OK:     false,
			Mode:   "cold",
			Reason: "coldStartScript not configured",
		})
		return
	}
	result, err := doColdStart(ctx, killExisting, knownPID)
	if err != nil {
		log.Printf("cold start failed (%s): %v output=%s remote=%s", why, err, result.Output, r.RemoteAddr)
		writeJSON(w, http.StatusBadGateway, triggerResp{
			OK:     false,
			Mode:   "cold",
			Reason: err.Error(),
			PID:    result.PID,
		})
		return
	}
	lastTrigger = time.Now()
	if result.Output == "already running" {
		log.Printf("cold start skipped (%s): already running pid=%d remote=%s", why, result.PID, r.RemoteAddr)
		writeJSON(w, http.StatusOK, triggerResp{OK: true, Mode: "cold", PID: result.PID, Reason: "already running"})
		return
	}
	log.Printf("cold start ok (%s) pid=%d remote=%s", why, result.PID, r.RemoteAddr)
	writeJSON(w, http.StatusOK, triggerResp{OK: true, Mode: "cold", PID: result.PID, Reason: "started"})
}

func requestAuthToken(r *http.Request) string {
	token := strings.TrimSpace(r.URL.Query().Get("auth"))
	if token == "" {
		token = strings.TrimSpace(r.URL.Query().Get("key"))
	}
	return token
}

func authorizeRequest(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	token := requestAuthToken(r)
	if subtle.ConstantTimeCompare([]byte(token), []byte(cfg.Auth)) != 1 {
		writeJSON(w, http.StatusForbidden, triggerResp{OK: false, Reason: "unauthorized"})
		return false
	}
	return true
}

func acquireTriggerSlot(w http.ResponseWriter) bool {
	triggerMu.Lock()
	defer triggerMu.Unlock()
	if since := time.Since(lastTrigger); since < time.Duration(cfg.MinIntervalSec)*time.Second {
		writeJSON(w, http.StatusTooManyRequests, triggerResp{
			OK:     false,
			Reason: fmt.Sprintf("rate limited, retry after %ds", cfg.MinIntervalSec-int(since.Seconds())),
		})
		return false
	}
	return true
}

func callGameHotRestart(ctx context.Context, auth string) (string, error) {
	base := strings.TrimSpace(cfg.HotRestartURL)
	sep := "?"
	if strings.Contains(base, "?") {
		sep = "&"
	}
	url := base + sep + "auth=" + auth
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	text := strings.TrimSpace(string(body))
	if resp.StatusCode != http.StatusOK {
		return text, fmt.Errorf("upstream HTTP %d", resp.StatusCode)
	}
	return text, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(v)
}
