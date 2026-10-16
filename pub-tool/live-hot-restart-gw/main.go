// live-hot-restart-gw：独立进程，部署在直播服本机。
// 监听 18900，手机/外网 GET ?auth= 后 curl 本机 https://127.0.0.1/internal/hotRestart 触发热重启。
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
	Listen         string `json:"listen"`
	Auth           string `json:"auth"`
	HotRestartURL  string `json:"hotRestartURL"`
	MinIntervalSec int    `json:"minIntervalSec"`
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

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleTrigger)
	mux.HandleFunc("/hotRestart", handleTrigger)
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	})

	log.Printf("live-hot-restart-gw listen http://%s/  config=%s", cfg.Listen, configPath)
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
	return nil
}

type triggerResp struct {
	OK       bool   `json:"ok"`
	Reason   string `json:"reason,omitempty"`
	Upstream string `json:"upstream,omitempty"`
}

func handleTrigger(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	token := strings.TrimSpace(r.URL.Query().Get("auth"))
	if token == "" {
		token = strings.TrimSpace(r.URL.Query().Get("key"))
	}
	if subtle.ConstantTimeCompare([]byte(token), []byte(cfg.Auth)) != 1 {
		writeJSON(w, http.StatusForbidden, triggerResp{OK: false, Reason: "unauthorized"})
		return
	}

	triggerMu.Lock()
	defer triggerMu.Unlock()
	if since := time.Since(lastTrigger); since < time.Duration(cfg.MinIntervalSec)*time.Second {
		writeJSON(w, http.StatusTooManyRequests, triggerResp{
			OK:     false,
			Reason: fmt.Sprintf("rate limited, retry after %ds", cfg.MinIntervalSec-int(since.Seconds())),
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
	defer cancel()

	upstream, err := callGameHotRestart(ctx, token)
	if err != nil {
		log.Printf("hot restart upstream error: %v body=%s", err, upstream)
		writeJSON(w, http.StatusBadGateway, triggerResp{OK: false, Reason: err.Error(), Upstream: upstream})
		return
	}
	lastTrigger = time.Now()
	ok := strings.Contains(upstream, `"ok":true`) || strings.Contains(upstream, `"ok": true`)
	status := http.StatusOK
	if !ok {
		status = http.StatusBadGateway
	}
	log.Printf("hot restart trigger ok=%v upstream=%s remote=%s", ok, upstream, r.RemoteAddr)
	writeJSON(w, status, triggerResp{OK: ok, Reason: upstream, Upstream: upstream})
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
