package metapixel

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/xrlog"
	sysentity "xr-game-server/entity/sys"
)

const metaGraphAPIVersion = "v21.0"

type capiRequest struct {
	Data          []capiEvent `json:"data"`
	TestEventCode string      `json:"test_event_code,omitempty"`
}

type capiEvent struct {
	EventName    string          `json:"event_name"`
	EventTime    int64           `json:"event_time"`
	EventID      string          `json:"event_id,omitempty"`
	ActionSource string          `json:"action_source"`
	UserData     capiUserData    `json:"user_data"`
	CustomData   *capiCustomData `json:"custom_data,omitempty"`
}

type capiUserData struct {
	ExternalID []string `json:"external_id,omitempty"`
}

type capiCustomData struct {
	Currency string  `json:"currency,omitempty"`
	Value    float64 `json:"value,omitempty"`
	OrderID  string  `json:"order_id,omitempty"`
}

func hashMetaExternalID(userId uint64) string {
	if userId == 0 {
		return ""
	}
	raw := strings.TrimSpace(strconv.FormatUint(userId, 10))
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func sendCAPIEvents(ctx context.Context, cfg *sysentity.MetaPixelCfg, events []capiEvent) {
	if cfg == nil || !cfg.IsActive() || len(events) == 0 {
		return
	}
	reqBody := capiRequest{Data: events}
	if code := strings.TrimSpace(cfg.TestEventCode); code != "" {
		reqBody.TestEventCode = code
	}
	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/events?access_token=%s",
		metaGraphAPIVersion, strings.TrimSpace(cfg.PixelId), strings.TrimSpace(cfg.AccessToken))
	eventSummary := formatCAPIEventSummary(events)
	resp, err := g.Client().ContentJson().Post(ctx, url, reqBody)
	if err != nil {
		xrlog.DetailLog.Warningf(ctx, "meta pixel capi request failed events=%s err=%v", eventSummary, err)
		return
	}
	defer resp.Close()
	body := resp.ReadAllString()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		xrlog.DetailLog.Warningf(ctx, "meta pixel capi bad status=%d events=%s body=%s", resp.StatusCode, eventSummary, truncateLog(body))
		return
	}
	var parsed struct {
		EventsReceived int      `json:"events_received"`
		Messages       []string `json:"messages"`
		FBTraceID      string   `json:"fbtrace_id"`
	}
	if err := json.Unmarshal([]byte(body), &parsed); err != nil {
		xrlog.DetailLog.Warningf(ctx, "meta pixel capi decode failed events=%s err=%v body=%s", eventSummary, err, truncateLog(body))
		return
	}
	if parsed.EventsReceived <= 0 {
		xrlog.DetailLog.Warningf(ctx, "meta pixel capi no events received events=%s messages=%v trace=%s", eventSummary, parsed.Messages, parsed.FBTraceID)
		return
	}
	xrlog.DetailLog.Infof(ctx, "meta pixel capi ok pixelId=%s events=%s received=%d trace=%s messages=%v testMode=%t",
		strings.TrimSpace(cfg.PixelId), eventSummary, parsed.EventsReceived, parsed.FBTraceID, parsed.Messages, reqBody.TestEventCode != "")
}

func formatCAPIEventSummary(events []capiEvent) string {
	if len(events) == 0 {
		return ""
	}
	parts := make([]string, 0, len(events))
	for _, e := range events {
		name := strings.TrimSpace(e.EventName)
		if e.EventID != "" {
			parts = append(parts, fmt.Sprintf("%s:%s", name, e.EventID))
			continue
		}
		parts = append(parts, name)
	}
	return strings.Join(parts, ",")
}

func truncateLog(s string) string {
	s = strings.TrimSpace(s)
	if len(s) <= 512 {
		return s
	}
	return s[:512] + "..."
}

func trackCompleteRegistration(ctx context.Context, userId uint64, at time.Time) {
	cfg := activeCfg()
	if cfg == nil {
		return
	}
	extID := hashMetaExternalID(userId)
	if extID == "" {
		return
	}
	eventTime := at.Unix()
	if eventTime <= 0 {
		eventTime = time.Now().Unix()
	}
	sendCAPIEvents(ctx, cfg, []capiEvent{{
		EventName:    "CompleteRegistration",
		EventTime:    eventTime,
		EventID:      fmt.Sprintf("reg-%d", userId),
		ActionSource: "system_generated",
		UserData:     capiUserData{ExternalID: []string{extID}},
	}})
}

func trackPurchase(ctx context.Context, userId uint64, orderId uint64, usdAmount float64, at time.Time) {
	cfg := activeCfg()
	if cfg == nil || usdAmount <= 0 || orderId == 0 || userId == 0 {
		return
	}
	extID := hashMetaExternalID(userId)
	if extID == "" {
		return
	}
	eventTime := at.Unix()
	if eventTime <= 0 {
		eventTime = time.Now().Unix()
	}
	sendCAPIEvents(ctx, cfg, []capiEvent{{
		EventName:    "Purchase",
		EventTime:    eventTime,
		EventID:      fmt.Sprintf("order-%d", orderId),
		ActionSource: "system_generated",
		UserData:     capiUserData{ExternalID: []string{extID}},
		CustomData: &capiCustomData{
			Currency: "USD",
			Value:    usdAmount,
			OrderID:  strconv.FormatUint(orderId, 10),
		},
	}})
}

func activeCfg() *sysentity.MetaPixelCfg {
	cfg := getCfgCached()
	if cfg == nil || !cfg.IsActive() {
		return nil
	}
	return cfg
}
