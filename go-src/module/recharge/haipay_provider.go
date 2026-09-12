package recharge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"xr-game-server/constants/country"
	"xr-game-server/core/xrlog"
	"xr-game-server/dao/cfgdao"
)

const (
	haiPayProviderName           = "haipay"
	haiPayGlobalCashierApplyPath = "/global/cashier/collect/apply"
	haiPayCollectNotifyPath      = "/webhook/haipay/collect/notify"
	haiPayMinUsdAmount           = 0.99
	haiPayHTTPTimeout            = 30 * time.Second
	haiPayLogBodyMax             = 4096
)

// 法币代码 → 全球收银台 region（ISO 3166-1 alpha-2）
var haiPayCurrencyToRegion = map[string]string{
	"AED": "AE", "BHD": "BH", "BRL": "BR", "EGP": "EG", "EUR": "EU",
	"GBP": "GB", "HKD": "HK", "IDR": "ID", "INR": "IN", "JPY": "JP",
	"KRW": "KR", "KWD": "KW", "MYR": "MY", "OMR": "OM", "PHP": "PH",
	"PKR": "PK", "QAR": "QA", "SAR": "SA", "SGD": "SG", "THB": "TH",
	"TRY": "TR", "TWD": "TW", "USD": "US", "VND": "VN",
}

type haiPayProvider struct{}

func (p *haiPayProvider) Name() string { return haiPayProviderName }

func (p *haiPayProvider) Enabled() bool { return cfgdao.HaiPayEnabled() }

func (p *haiPayProvider) QuotePay(priceUsd float64) (payCurrency string, payAmount float64, err error) {
	if priceUsd < haiPayMinUsdAmount {
		return "", 0, fmt.Errorf("haipay usd amount must be >= %.2f", haiPayMinUsdAmount)
	}
	return "USD", priceUsd, nil
}

func (p *haiPayProvider) CreatePay(ctx context.Context, req *ChannelPayCreateReq) (*ChannelPayCreateRes, error) {
	cfg := cfgdao.GetHaiPayCfgCached()
	if cfg == nil || !cfgdao.HaiPayEnabled() {
		return nil, fmt.Errorf("haipay not configured")
	}
	regionHint := strings.TrimSpace(req.Region)
	if regionHint == "" {
		regionHint = strings.TrimSpace(cfg.DefaultRegion)
	}
	if regionHint == "" {
		regionHint = "ID"
	}
	region, err := haiPayResolveRegion(regionHint)
	if err != nil {
		return nil, err
	}
	amountUsd := req.PriceUsd
	if amountUsd <= 0 {
		amountUsd = req.Amount
	}
	if amountUsd < haiPayMinUsdAmount {
		return nil, fmt.Errorf("haipay usd amount must be >= %.2f got=%v", haiPayMinUsdAmount, amountUsd)
	}

	name := strings.TrimSpace(req.PlayerName)
	if name == "" {
		name = fmt.Sprintf("User %d", req.UserID)
	}
	if !strings.Contains(name, " ") {
		name = name + " User"
	}
	email := strings.TrimSpace(req.Email)
	if email == "" {
		email = fmt.Sprintf("user%d@noreply.local", req.UserID)
	}

	amountStr := fmt.Sprintf("%.2f", amountUsd)
	notifyURL := strings.TrimRight(cfg.CallbackBaseUrl, "/") + haiPayCollectNotifyPath
	returnURL := strings.TrimSpace(cfg.ReturnUrl)
	if returnURL == "" {
		returnURL = cfg.CallbackBaseUrl
	}
	failURL := strings.TrimSpace(cfg.FailReturnUrl)
	if failURL == "" {
		failURL = returnURL
	}
	subject := strings.TrimSpace(cfg.Subject)
	if subject == "" {
		subject = "Recharge"
	}

	body := map[string]any{
		"appId":          cfg.AppId,
		"orderId":        req.OrderID,
		"name":           name,
		"email":          email,
		"amount":         amountStr,
		"currency":       "USD",
		"callBackUrl":    returnURL,
		"callBackFailUrl": failURL,
		"notifyUrl":      notifyURL,
		"subject":        subject,
		"region":         region,
		"partnerUserId":  strconv.FormatUint(req.UserID, 10),
	}
	if cancel := strings.TrimSpace(cfg.CancelUrl); cancel != "" {
		body["cancelUrl"] = cancel
	}
	if methods := strings.TrimSpace(cfg.PaymentMethods); methods != "" {
		body["paymentMethods"] = methods
	}
	if remark := strings.TrimSpace(req.OrderID); remark != "" {
		body["body"] = "order:" + remark
	}
	if phone := strings.TrimSpace(req.Phone); phone != "" {
		body["phone"] = phone
	}

	// 签名前：业务字段 + 拼串（含 key），便于与 HaiPay 对账
	xrlog.DetailLog.Infof(ctx, "haipay apply sign-before(request) params=%s", haiPaySafeParams(body))
	sign, err := haiPaySign(ctx, body, cfg.MerchantSecretKey, cfg.MerchantPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("haipay sign: %w", err)
	}
	body["sign"] = sign

	var res haiPayApplyAPIRes
	url := strings.TrimRight(cfg.ApiHost, "/") + haiPayGlobalCashierApplyPath
	if err = haiPayPostJSON(ctx, url, body, &res); err != nil {
		return nil, err
	}
	if res.Status != "1" {
		msg := strings.TrimSpace(res.Msg)
		if msg == "" {
			msg = res.Error
		}
		return nil, fmt.Errorf("haipay apply failed status=%s error=%s msg=%s", res.Status, res.Error, msg)
	}
	if res.Data == nil {
		return nil, fmt.Errorf("haipay apply empty data")
	}
	// 验签只对应答 data 内实际字段;验签失败只打日志不阻断(HaiPay 已出单可付款,回调仍严格验签)
	dataMap := map[string]any{}
	rawData, _ := json.Marshal(res.Data)
	_ = json.Unmarshal(rawData, &dataMap)
	xrlog.DetailLog.Infof(ctx, "haipay apply sign-before(response) data=%s", haiPaySafeParams(dataMap))
	if err = haiPayVerify(ctx, dataMap, cfg.MerchantSecretKey, cfg.HaiPayPublicKey, res.Data.Sign); err != nil {
		xrlog.DetailLog.Warningf(ctx, "haipay apply verify failed(skip block) orderId=%s err=%v", req.OrderID, err)
	}
	payURL := strings.TrimSpace(res.Data.PayUrl)
	if payURL == "" {
		return nil, fmt.Errorf("haipay empty payUrl")
	}
	return &ChannelPayCreateRes{
		PayURL:       payURL,
		ThirdOrderID: strings.TrimSpace(res.Data.OrderNo),
	}, nil
}

func haiPayResolveRegion(currencyOrRegion string) (string, error) {
	code := strings.ToUpper(strings.TrimSpace(currencyOrRegion))
	if code == "" {
		return "", fmt.Errorf("empty region")
	}
	if region, ok := haiPayCurrencyToRegion[code]; ok {
		return region, nil
	}
	if country.IsHaiPayRegion(code) {
		return code, nil
	}
	return "", fmt.Errorf("unsupported region currency=%s", code)
}

type haiPayApplyAPIRes struct {
	Status string `json:"status"`
	Error  string `json:"error"`
	Msg    string `json:"msg"`
	Data   *struct {
		OrderId string `json:"orderId"`
		OrderNo string `json:"orderNo"`
		PayUrl  string `json:"payUrl"`
		Sign    string `json:"sign"`
	} `json:"data"`
}

func haiPayTruncateLog(s string, max int) string {
	if max <= 0 || len(s) <= max {
		return s
	}
	return s[:max] + "...(truncated)"
}

func haiPayPostJSON(ctx context.Context, fullURL string, payload any, out any) error {
	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	xrlog.DetailLog.Infof(ctx, "haipay http request url=%s body=%s", fullURL, haiPaySafeJSON(raw))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: haiPayHTTPTimeout}
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		costMs := time.Since(start).Milliseconds()
		xrlog.DetailLog.Errorf(ctx, "haipay http do failed url=%s costMs=%d err=%v", fullURL, costMs, err)
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	costMs := time.Since(start).Milliseconds()
	if err != nil {
		xrlog.DetailLog.Errorf(ctx, "haipay http read failed url=%s status=%d costMs=%d err=%v", fullURL, resp.StatusCode, costMs, err)
		return err
	}
	xrlog.DetailLog.Infof(ctx, "haipay http response url=%s status=%d costMs=%d body=%s", fullURL, resp.StatusCode, costMs, haiPaySafeJSON(body))
	if len(body) == 0 {
		return fmt.Errorf("haipay empty response status=%d", resp.StatusCode)
	}
	if err = json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("haipay decode failed status=%d: %w", resp.StatusCode, err)
	}
	return nil
}

func haiPaySafeJSON(raw []byte) string {
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return haiPayTruncateLog(string(raw), haiPayLogBodyMax)
	}
	return haiPaySafeParams(m)
}

// haiPaySafeParams 打日志用：脱敏 sign/密钥，保留签名前业务字段
func haiPaySafeParams(m map[string]any) string {
	if m == nil {
		return "{}"
	}
	cp := make(map[string]any, len(m))
	for k, v := range m {
		cp[k] = v
	}
	for _, k := range []string{"sign", "merchantSecretKey", "merchantPrivateKey", "haiPayPublicKey"} {
		if _, ok := cp[k]; ok {
			cp[k] = "***"
		}
	}
	b, err := json.Marshal(cp)
	if err != nil {
		return "{}"
	}
	return haiPayTruncateLog(string(b), haiPayLogBodyMax)
}
