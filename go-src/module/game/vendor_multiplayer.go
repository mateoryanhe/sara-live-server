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
	vendorMultiplayerStartAPI            = "open/multiplayer/start"
	vendorMultiplayerStartPath           = "/open/multiplayer/start"
	vendorMultiplayerStartRequestTimeout = 15 * time.Second
	vendorMultiplayerDefaultPlatform     = "zy"
)

type vendorMultiplayerStartResp struct {
	Code int `json:"code"`
	Data struct {
		ConfigUrl  string `json:"configUrl"`
		ExpireInMs int64  `json:"expireInMs"`
	} `json:"data"`
	Message string `json:"message"`
}

func fetchVendorMultiplayerConfigURL(ctx context.Context, gameCode, platform string) (string, int64, error) {
	if err := waitGamePlatformCfgReady(ctx); err != nil {
		return "", 0, err
	}
	cfg := cfgdao.GetGamePlatformCfgFromMemory()
	if cfg == nil {
		return "", 0, fmt.Errorf("game platform cfg not found")
	}

	gameCode = strings.TrimSpace(gameCode)
	platform = strings.TrimSpace(platform)
	if gameCode == "" {
		return "", 0, fmt.Errorf("gameId required")
	}
	if platform == "" {
		platform = vendorMultiplayerDefaultPlatform
	}

	baseURL := cfgdao.GetVendorUrlFromMemory()
	operatorToken := strings.TrimSpace(cfg.Token)
	secretKey := strings.TrimSpace(cfg.SecretKey)

	body := map[string]string{
		"platform": platform,
		"gameId":   gameCode,
	}
	return postVendorMultiplayerStart(ctx, baseURL, body, operatorToken, secretKey)
}

func postVendorMultiplayerStart(
	ctx context.Context,
	baseURL string,
	body map[string]string,
	operatorToken, secretKey string,
) (string, int64, error) {
	if ctx == nil {
		ctx = gctx.New()
	}
	reqCtx, cancel := context.WithTimeout(ctx, vendorMultiplayerStartRequestTimeout)
	defer cancel()

	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signParams := sign.MergeParams(nil, body, operatorToken, timestamp)
	signValue := sign.Build(signParams, secretKey)
	headers := buildVendorRequestHeaders(operatorToken, timestamp, signValue)

	url := buildVendorAPIURL(baseURL, vendorMultiplayerStartPath)
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", 0, err
	}

	logVendorAPIRequest(ctx, vendorMultiplayerStartAPI, "POST", url, body, headers, signParams, secretKey)

	client := gclient.New()
	client.SetHeader("Content-Type", "application/json")
	for key, value := range headers {
		client.SetHeader(key, value)
	}

	start := time.Now()
	resp, err := client.Post(reqCtx, url, bodyBytes)
	if err != nil {
		logVendorAPIError(ctx, vendorMultiplayerStartAPI, "POST", url, err, vendorAPICostMs(start))
		return "", 0, err
	}
	defer resp.Close()

	respBody := resp.ReadAll()
	logVendorAPIResponse(ctx, vendorMultiplayerStartAPI, "POST", url, resp.StatusCode, respBody, vendorAPICostMs(start))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", 0, fmt.Errorf("vendor multiplayer start http status=%d body=%s", resp.StatusCode, string(respBody))
	}

	var result vendorMultiplayerStartResp
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", 0, err
	}
	if result.Code != 0 {
		msg := strings.TrimSpace(result.Message)
		if msg == "" {
			msg = "vendor multiplayer start failed"
		}
		return "", 0, fmt.Errorf("vendor multiplayer start code=%d msg=%s", result.Code, msg)
	}

	configURL := strings.TrimSpace(result.Data.ConfigUrl)
	if configURL == "" {
		return "", 0, fmt.Errorf("vendor multiplayer start empty configUrl")
	}
	return configURL, result.Data.ExpireInMs, nil
}
