package game

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"xr-game-server/module/game/sign"
)

type vendorAPIRequestLog struct {
	API        string            `json:"api"`
	Method     string            `json:"method"`
	URL        string            `json:"url"`
	Query      map[string]string `json:"query"`
	Headers    map[string]string `json:"headers"`
	SignParams map[string]string `json:"signParams"`
	SignString string            `json:"signString"`
}

type vendorAPIResponseLog struct {
	API        string `json:"api"`
	Method     string `json:"method"`
	URL        string `json:"url"`
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
	Error      string `json:"error,omitempty"`
	CostMs     int64  `json:"costMs"`
}

func logVendorAPIRequest(ctx context.Context, api, method, url string, query map[string]string, headers map[string]string, signParams map[string]string, secretKey string) {
	payload := vendorAPIRequestLog{
		API:        api,
		Method:     method,
		URL:        url,
		Query:      cloneStringMap(query),
		Headers:    cloneStringMap(headers),
		SignParams: cloneStringMap(signParams),
		SignString: maskSignString(sign.BuildString(signParams, secretKey), secretKey),
	}
	logVendorAPIJSON(ctx, "vendor api request", payload)
}

func logVendorAPIResponse(ctx context.Context, api, method, url string, statusCode int, body []byte, costMs int64) {
	payload := vendorAPIResponseLog{
		API:        api,
		Method:     method,
		URL:        url,
		StatusCode: statusCode,
		Body:       string(body),
		CostMs:     costMs,
	}
	logVendorAPIJSON(ctx, "vendor api response", payload)
	vendorDetailLog().Infof(ctx, "vendor api finished api=%s method=%s costMs=%d statusCode=%d url=%s",
		api, method, costMs, statusCode, url)
}

func logVendorAPIError(ctx context.Context, api, method, url string, err error, costMs int64) {
	payload := vendorAPIResponseLog{
		API:    api,
		Method: method,
		URL:    url,
		Error:  err.Error(),
		CostMs: costMs,
	}
	logVendorAPIJSON(ctx, "vendor api error", payload)
	vendorDetailLog().Warningf(ctx, "vendor api failed api=%s method=%s costMs=%d url=%s err=%v",
		api, method, costMs, url, err)
}

func vendorAPICostMs(start time.Time) int64 {
	if start.IsZero() {
		return 0
	}
	return time.Since(start).Milliseconds()
}

func logVendorAPIJSON(ctx context.Context, title string, payload any) {
	raw, err := json.Marshal(payload)
	if err != nil {
		vendorDetailLog().Warningf(ctx, "%s marshal failed: %v", title, err)
		return
	}
	vendorDetailLog().Infof(ctx, "%s: %s", title, string(raw))
}

func buildVendorRequestHeaders(operatorToken, timestamp, signValue string) map[string]string {
	return map[string]string{
		"Accept":                     "application/json",
		sign.HTTPHeaderOperatorToken: operatorToken,
		sign.HTTPHeaderTimestamp:     timestamp,
		sign.HTTPHeaderSign:          signValue,
	}
}

func cloneStringMap(src map[string]string) map[string]string {
	if len(src) == 0 {
		return map[string]string{}
	}
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func maskSignString(signString, secretKey string) string {
	secretKey = strings.TrimSpace(secretKey)
	if secretKey == "" {
		return signString
	}
	return strings.Replace(signString, "&secret="+secretKey, "&secret=***", 1)
}
