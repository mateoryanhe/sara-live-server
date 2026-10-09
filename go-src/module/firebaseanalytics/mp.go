package firebaseanalytics

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	sysentity "xr-game-server/entity/sys"
)

type mpCollectRequest struct {
	ClientID string    `json:"client_id"`
	UserID   string    `json:"user_id,omitempty"`
	Events   []mpEvent `json:"events"`
}

type mpEvent struct {
	Name   string                 `json:"name"`
	Params map[string]interface{} `json:"params"`
}

func sendSignUpEvent(ctx context.Context, userId uint64, countryCode string, at time.Time) {
	cfg := GetCfgCached()
	if cfg == nil {
		return
	}
	measurementID, err := resolveMeasurementID(cfg)
	if err != nil {
		g.Log().Warningf(ctx, "firebase analytics measurement id missing: %v", err)
		return
	}
	apiSecret := strings.TrimSpace(cfg.MeasurementApiSecret)
	if apiSecret == "" {
		return
	}
	userIDStr := strconv.FormatUint(userId, 10)
	params := map[string]interface{}{
		"engagement_time_msec": 100,
	}
	countryCode = strings.ToUpper(strings.TrimSpace(countryCode))
	if countryCode != "" {
		params["country_code"] = countryCode
	}
	events := []mpEvent{{
		Name: "sign_up",
		Params: params,
	}}
	sendMeasurementProtocol(ctx, measurementID, apiSecret, userIDStr, clientIDForUser(userId), events, at)
}

func resolveMeasurementID(cfg *sysentity.FirebaseAnalyticsCfg) (string, error) {
	if cfg == nil {
		return "", fmt.Errorf("cfg is nil")
	}
	parsed, err := parseClientConfig(cfg.ProjectId, cfg.ClientConfigJson)
	if err != nil {
		return "", err
	}
	id := strings.TrimSpace(parsed.MeasurementId)
	if id == "" {
		return "", fmt.Errorf("measurementId empty")
	}
	return id, nil
}

func sendMeasurementProtocol(ctx context.Context, measurementID, apiSecret, userID, clientID string, events []mpEvent, at time.Time) {
	if measurementID == "" || apiSecret == "" || clientID == "" || len(events) == 0 {
		return
	}
	_ = at
	reqBody := mpCollectRequest{
		ClientID: clientID,
		UserID:   userID,
		Events:   events,
	}
	url := fmt.Sprintf("https://www.google-analytics.com/mp/collect?measurement_id=%s&api_secret=%s",
		measurementID, apiSecret)
	resp, err := g.Client().ContentJson().Post(ctx, url, reqBody)
	if err != nil {
		g.Log().Warningf(ctx, "firebase analytics mp request failed: %v", err)
		return
	}
	defer resp.Close()
	body := resp.ReadAllString()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		g.Log().Warningf(ctx, "firebase analytics mp bad status=%d body=%s", resp.StatusCode, truncateMPLog(body))
		return
	}
	if strings.TrimSpace(body) != "" {
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(body), &parsed); err == nil && len(parsed) > 0 {
			g.Log().Warningf(ctx, "firebase analytics mp response: %s", truncateMPLog(body))
		}
	}
}

func truncateMPLog(s string) string {
	const max = 512
	s = strings.TrimSpace(s)
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
