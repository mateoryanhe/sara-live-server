package game

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/module/game/sign"
)

const (
	vendorGameStartAPI            = "open/game/start"
	vendorGameStartPath           = "/open/game/start"
	vendorGameStartRequestTimeout = 15 * time.Second
)

type vendorGameStartResp struct {
	Code int `json:"code"`
	Data struct {
		LaunchUrl  string `json:"launchUrl"`
		ExpireInMs int64  `json:"expireInMs"`
	} `json:"data"`
	Message string `json:"message"`
}

func fetchVendorGameStartURL(ctx context.Context, gameCode, platform, ops, lang string) (string, error) {
	if !cfgdao.GamePlatformCfgReady() {
		return "", fmt.Errorf("game platform cfg not ready")
	}
	cfg := cfgdao.GetGamePlatformCfgFromMemory()
	if cfg == nil {
		return "", fmt.Errorf("game platform cfg not found")
	}

	gameCode = strings.TrimSpace(gameCode)
	platform = strings.TrimSpace(platform)
	ops = strings.TrimSpace(ops)
	lang = strings.TrimSpace(lang)
	if gameCode == "" || platform == "" || ops == "" {
		return "", fmt.Errorf("game start params invalid")
	}
	if lang == "" {
		lang = "en"
	}

	baseURL := cfgdao.GetVendorUrlFromMemory()
	operatorToken := strings.TrimSpace(cfg.Token)
	secretKey := strings.TrimSpace(cfg.SecretKey)

	body := map[string]string{
		"gameId":   gameCode,
		"platform": platform,
		"ops":      ops,
		"lang":     lang,
	}
	return postVendorGameStart(ctx, baseURL, body, operatorToken, secretKey)
}

func postVendorGameStart(
	ctx context.Context,
	baseURL string,
	body map[string]string,
	operatorToken, secretKey string,
) (string, error) {
	if ctx == nil {
		ctx = gctx.New()
	}
	reqCtx, cancel := context.WithTimeout(ctx, vendorGameStartRequestTimeout)
	defer cancel()

	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signParams := sign.MergeParams(nil, body, operatorToken, timestamp)
	signValue := sign.Build(signParams, secretKey)
	headers := buildVendorRequestHeaders(operatorToken, timestamp, signValue)

	url := buildVendorAPIURL(baseURL, vendorGameStartPath)
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	logVendorAPIRequest(ctx, vendorGameStartAPI, "POST", url, body, headers, signParams, secretKey)

	client := gclient.New()
	client.SetHeader("Content-Type", "application/json")
	for key, value := range headers {
		client.SetHeader(key, value)
	}

	start := time.Now()
	resp, err := client.Post(reqCtx, url, bodyBytes)
	if err != nil {
		logVendorAPIError(ctx, vendorGameStartAPI, "POST", url, err, vendorAPICostMs(start))
		return "", err
	}
	defer resp.Close()

	respBody := resp.ReadAll()
	logVendorAPIResponse(ctx, vendorGameStartAPI, "POST", url, resp.StatusCode, respBody, vendorAPICostMs(start))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("vendor game start http status=%d body=%s", resp.StatusCode, string(respBody))
	}

	var result vendorGameStartResp
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}
	if result.Code != 0 {
		msg := strings.TrimSpace(result.Message)
		if msg == "" {
			msg = "vendor game start failed"
		}
		return "", fmt.Errorf("vendor game start code=%d msg=%s", result.Code, msg)
	}

	launchURL := strings.TrimSpace(result.Data.LaunchUrl)
	if launchURL == "" {
		return "", fmt.Errorf("vendor game start empty launchUrl")
	}
	return launchURL, nil
}

func buildVendorAPIURL(baseURL, path string) string {
	return strings.TrimRight(strings.TrimSpace(baseURL), "/") + path
}
