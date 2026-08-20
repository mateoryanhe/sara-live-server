package game

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/module/game/sign"
)

const (
	vendorGameListAPI            = "open/game/list"
	vendorGameListPath           = "/open/game/list"
	vendorGameFetchPageSize      = 30
	vendorGameListRequestTimeout = 30 * time.Second
	vendorGameListRetryInterval  = 3 * time.Second
	vendorGameListMaxAttempts    = 10
)

type vendorGameListResp struct {
	Code int `json:"code"`
	Data struct {
		List     []*vendorGameListItem `json:"list"`
		Data     []*vendorGameListItem `json:"data"`
		Total    int                   `json:"total"`
		Page     int                   `json:"page"`
		PageSize int                   `json:"pageSize"`
	} `json:"data"`
	Message string `json:"message"`
}

type vendorGameListItem struct {
	GameCode string `json:"gameCode"`
	Name     string `json:"name"`
	NameEn   string `json:"nameEn"`
	Category string `json:"category"`
	Cover    string `json:"cover"`
	Platform string `json:"platform"`
}

func fetchAllVendorGames(ctx context.Context) ([]*VendorGame, error) {
	if err := waitGamePlatformCfgReady(ctx); err != nil {
		return nil, err
	}
	cfg := cfgdao.GetGamePlatformCfgFromMemory()
	if cfg == nil {
		return nil, fmt.Errorf("game platform cfg not found")
	}

	baseURL := cfgdao.GetVendorUrlFromMemory()
	operatorToken := strings.TrimSpace(cfg.Token)
	secretKey := strings.TrimSpace(cfg.SecretKey)

	all := make([]*VendorGame, 0, vendorGameFetchPageSize)
	page := 1
	for {
		query := map[string]string{
			"page":     strconv.Itoa(page),
			"pageSize": strconv.Itoa(vendorGameFetchPageSize),
		}
		items, _, err := fetchVendorGamePage(ctx, baseURL, query, operatorToken, secretKey)
		if err != nil {
			return all, err
		}
		if len(items) == 0 {
			break
		}
		all = append(all, items...)
		if len(items) < vendorGameFetchPageSize {
			break
		}
		page++
	}
	return all, nil
}

func waitGamePlatformCfgReady(ctx context.Context) error {
	var lastErr error
	for attempt := 1; attempt <= vendorGameListMaxAttempts; attempt++ {
		cfgdao.ReloadGamePlatformCfgCache()
		if cfgdao.GamePlatformCfgReady() {
			return nil
		}
		lastErr = fmt.Errorf("game platform cfg not ready")
		if attempt >= vendorGameListMaxAttempts {
			break
		}
		vendorDetailLog().Warningf(ctx, "game platform cfg not ready retry attempt=%d/%d wait=%s",
			attempt, vendorGameListMaxAttempts, vendorGameListRetryInterval)
		if err := sleepWithContext(ctx, vendorGameListRetryInterval); err != nil {
			return err
		}
	}
	return lastErr
}

func fetchVendorGamePage(
	ctx context.Context,
	baseURL string,
	query map[string]string,
	operatorToken, secretKey string,
) ([]*VendorGame, int, error) {
	var lastErr error
	for attempt := 1; attempt <= vendorGameListMaxAttempts; attempt++ {
		if attempt > 1 {
			if ctx == nil {
				ctx = gctx.New()
			}
			vendorDetailLog().Warningf(ctx, "vendor game list retry attempt=%d/%d wait=%s err=%v",
				attempt, vendorGameListMaxAttempts, vendorGameListRetryInterval, lastErr)
			if err := sleepWithContext(ctx, vendorGameListRetryInterval); err != nil {
				return nil, 0, err
			}
		}
		items, total, err := fetchVendorGamePageOnce(ctx, baseURL, query, operatorToken, secretKey)
		if err == nil {
			return items, total, nil
		}
		lastErr = err
	}
	return nil, 0, lastErr
}

func fetchVendorGamePageOnce(
	ctx context.Context,
	baseURL string,
	query map[string]string,
	operatorToken, secretKey string,
) ([]*VendorGame, int, error) {
	if ctx == nil {
		ctx = gctx.New()
	}
	reqCtx, cancel := context.WithTimeout(ctx, vendorGameListRequestTimeout)
	defer cancel()

	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signParams := sign.MergeParams(query, nil, operatorToken, timestamp)
	signValue := sign.Build(signParams, secretKey)
	headers := buildVendorRequestHeaders(operatorToken, timestamp, signValue)

	client := gclient.New()
	for key, value := range headers {
		client.SetHeader(key, value)
	}

	url := buildVendorGameListURL(baseURL, query)
	logVendorAPIRequest(ctx, vendorGameListAPI, "GET", url, query, headers, signParams, secretKey)
	start := time.Now()
	resp, err := client.Get(reqCtx, url)
	if err != nil {
		logVendorAPIError(ctx, vendorGameListAPI, "GET", url, err, vendorAPICostMs(start))
		return nil, 0, err
	}
	defer resp.Close()

	body := resp.ReadAll()
	logVendorAPIResponse(ctx, vendorGameListAPI, "GET", url, resp.StatusCode, body, vendorAPICostMs(start))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, 0, fmt.Errorf("vendor game list http status=%d body=%s", resp.StatusCode, string(body))
	}

	var result vendorGameListResp
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, 0, err
	}
	if result.Code != 0 {
		msg := strings.TrimSpace(result.Message)
		if msg == "" {
			msg = "vendor game list failed"
		}
		return nil, 0, fmt.Errorf("vendor game list code=%d msg=%s", result.Code, msg)
	}

	items := result.Data.List
	if len(items) == 0 {
		items = result.Data.Data
	}
	games := make([]*VendorGame, 0, len(items))
	for _, item := range items {
		if item == nil || strings.TrimSpace(item.GameCode) == "" {
			continue
		}
		games = append(games, &VendorGame{
			GameCode: strings.TrimSpace(item.GameCode),
			Name:     strings.TrimSpace(item.Name),
			NameEn:   strings.TrimSpace(item.NameEn),
			Category: strings.TrimSpace(item.Category),
			Cover:    normalizeVendorGameCover(item.Cover),
			Platform: strings.TrimSpace(item.Platform),
		})
	}
	return games, result.Data.Total, nil
}

func buildVendorGameListURL(baseURL string, query map[string]string) string {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	values := make([]string, 0, len(query))
	keys := make([]string, 0, len(query))
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		values = append(values, k+"="+query[k])
	}
	if len(values) == 0 {
		return baseURL + vendorGameListPath
	}
	return baseURL + vendorGameListPath + "?" + strings.Join(values, "&")
}

func sleepWithContext(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
