package attrreport

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/xrlog"
	"xr-game-server/core/xrpool"
	"xr-game-server/dao/cfgdao"
)

const (
	EventNamePurchase             = "Purchase"
	EventNameCompleteRegistration = "CompleteRegistration"

	TypePWA = "pwa"
	TypeH5  = "h5"
	TypeAPK = "apk"

	DefaultActionSource = "website"
	// AcceptPath CMS baseApi 只填主机(如 https://cdntk.showsme.click),上报时固定拼此路径.
	AcceptPath = "/tj/api/accept.php"
)

// AcceptPayload 归因 accept.php POST body(字段名与对方约定一致).
type AcceptPayload struct {
	APIToken         string `json:"api_token"`
	Type             string `json:"type"`
	Ch               string `json:"ch"`
	TjUserId         string `json:"tjuserid"`
	UserId           string `json:"userid"`
	ClientIPAddress  string `json:"client_ip_address"`
	ClientIPv4       string `json:"client_ipv4"`
	ClientIPv6       string `json:"client_ipv6"`
	EventID          string `json:"event_id"`
	EventName        string `json:"event_name"`
	// EventTime Unix 秒(GMT/UTC epoch); Normalize 为空时补 Now().Unix().
	EventTime        int64  `json:"event_time"`
	EventSourceURL   string `json:"event_source_url"`
	Em               string `json:"em"`
	Ph               string `json:"ph"`
	ActionSource     string `json:"action_source"`
	ClientUserAgent  string `json:"client_user_agent"`
	PhoneModel       string `json:"Phone_Model"`
	SystemVersion    string `json:"System_Version"`
	// Country ISO 3166-1 alpha-2(如 ID); 所有事件均带此字段,解析不到则为空串.
	Country string `json:"country"`
	// Currency / Value 充值 Purchase 用; 无则 omitempty 不上报(注册等事件留空).
	Currency string  `json:"currency,omitempty"`
	Value    float64 `json:"value,omitempty"`
}

// Normalize 清洗常用字段; api_token 为空时由 PostAccept 用配置密钥补齐.
func (p *AcceptPayload) Normalize() {
	if p == nil {
		return
	}
	p.APIToken = strings.TrimSpace(p.APIToken)
	p.Type = strings.TrimSpace(p.Type)
	p.Ch = strings.TrimSpace(p.Ch)
	p.TjUserId = strings.TrimSpace(p.TjUserId)
	p.UserId = strings.TrimSpace(p.UserId)
	p.ClientIPAddress = strings.TrimSpace(p.ClientIPAddress)
	p.ClientIPv4 = strings.TrimSpace(p.ClientIPv4)
	p.ClientIPv6 = strings.TrimSpace(p.ClientIPv6)
	p.EventID = strings.TrimSpace(p.EventID)
	p.EventName = strings.TrimSpace(p.EventName)
	p.EventSourceURL = strings.TrimSpace(p.EventSourceURL)
	p.Em = strings.TrimSpace(p.Em)
	p.Ph = strings.TrimSpace(p.Ph)
	p.ActionSource = strings.TrimSpace(p.ActionSource)
	p.ClientUserAgent = strings.TrimSpace(p.ClientUserAgent)
	p.PhoneModel = strings.TrimSpace(p.PhoneModel)
	p.SystemVersion = strings.TrimSpace(p.SystemVersion)
	p.Country = strings.ToUpper(strings.TrimSpace(p.Country))
	p.Currency = strings.ToUpper(strings.TrimSpace(p.Currency))
	if p.ActionSource == "" {
		p.ActionSource = DefaultActionSource
	}
	if p.EventTime <= 0 {
		// Unix 秒即 UTC/GMT epoch
		p.EventTime = time.Now().Unix()
	}
}

// PostAccept 同步 POST 到 CMS 配置的 Base API(accept.php); 未配置则跳过.
// 各业务模块可直接调用; 热路径请用 PostAcceptAsync.
func PostAccept(ctx context.Context, payload *AcceptPayload) {
	if payload == nil {
		return
	}
	cfg := cfgdao.ResolveAttrReportCfg()
	if cfg == nil || !cfg.IsConfigured() {
		return
	}
	payload.Normalize()
	if payload.APIToken == "" {
		payload.APIToken = strings.TrimSpace(cfg.SecretKey)
	}
	if payload.EventName == "" || payload.APIToken == "" {
		xrlog.DetailLog.Warningf(ctx, "attr report skip empty event_name/api_token event=%s user=%s", payload.EventName, payload.UserId)
		return
	}
	url := buildAcceptURL(cfg.BaseApi)
	if url == "" {
		xrlog.DetailLog.Warningf(ctx, "attr report skip empty accept url event=%s user=%s", payload.EventName, payload.UserId)
		return
	}
	resp, err := g.Client().ContentJson().Post(ctx, url, payload)
	if err != nil {
		xrlog.DetailLog.Warningf(ctx, "attr report post failed url=%s event=%s user=%s err=%v", url, payload.EventName, payload.UserId, err)
		return
	}
	defer resp.Close()
	body := resp.ReadAllString()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		xrlog.DetailLog.Warningf(ctx, "attr report bad status=%d url=%s event=%s user=%s body=%s",
			resp.StatusCode, url, payload.EventName, payload.UserId, truncateAttrLog(body))
		return
	}
	xrlog.DetailLog.Infof(ctx, "attr report ok event=%s user=%s type=%s ch=%s event_id=%s status=%d body=%s",
		payload.EventName, payload.UserId, payload.Type, payload.Ch, payload.EventID, resp.StatusCode, truncateAttrLog(body))
}

// PostAcceptAsync 异步上报(xrpool); payload 按值拷贝,避免调用方改脏.
func PostAcceptAsync(ctx context.Context, payload AcceptPayload) {
	if !cfgdao.ResolveAttrReportCfg().IsConfigured() {
		return
	}
	copied := payload
	runCtx := ctx
	if runCtx == nil {
		runCtx = gctx.New()
	}
	xrpool.AddWithRecover(runCtx, func(asyncCtx context.Context) {
		PostAccept(asyncCtx, &copied)
	})
}

// buildAcceptURL CMS 存主机根地址,自动拼 /tj/api/accept.php; 若已带该路径则不重复拼.
func buildAcceptURL(baseApi string) string {
	base := strings.TrimSpace(baseApi)
	if base == "" {
		return ""
	}
	base = strings.TrimRight(base, "/")
	if strings.HasSuffix(strings.ToLower(base), strings.ToLower(AcceptPath)) {
		return base
	}
	return base + AcceptPath
}

func truncateAttrLog(s string) string {
	const max = 512
	s = strings.TrimSpace(s)
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
