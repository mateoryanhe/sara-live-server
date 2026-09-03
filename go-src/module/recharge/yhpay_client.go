package recharge

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"xr-game-server/core/xrlog"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity/recharge"
)

const yhPaySupportedCurrency = "IDR"

const (
	yhPayManualDepositPath    = "/Payin/Manual/Deposit"
	yhPayPayinCallbackPath    = "/webhook/yhpay/payin"
)

func yhPayHMAC(secret, payload string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

func normalizeYhPayRedirect(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	const prefix = "redirectlink="
	if strings.HasPrefix(strings.ToLower(raw), prefix) {
		return strings.TrimSpace(raw[len(prefix):])
	}
	return raw
}

func formatYhPayAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}

const yhPayLogBodyMax = 4096

func truncateYhPayLog(s string, max int) string {
	if max <= 0 || len(s) <= max {
		return s
	}
	return s[:max] + fmt.Sprintf("...(%d bytes)", len(s))
}

// yhPayLogSafeJSON 脱敏 APIKey/Hash/Signature 等字段，便于排障又不落密钥全文。
func yhPayLogSafeJSON(raw []byte) string {
	if len(raw) == 0 {
		return ""
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return truncateYhPayLog(string(raw), yhPayLogBodyMax)
	}
	for k := range m {
		switch strings.ToLower(k) {
		case "apikey", "api_key", "hash", "signature", "signature2":
			m[k] = "***"
		}
	}
	b, err := json.Marshal(m)
	if err != nil {
		return truncateYhPayLog(string(raw), yhPayLogBodyMax)
	}
	return truncateYhPayLog(string(b), yhPayLogBodyMax)
}

type yhPayManualDepositReq struct {
	MerchantCode   string `json:"merchant_code"`
	RefId          string `json:"ref_id"`
	CurrencyCode   string `json:"currency_code"`
	Amount         string `json:"amount"`
	Hash           string `json:"hash"`
	CallbackURL    string `json:"callback_url"`
	ReturnURL      string `json:"return_url"`
	PlayerUsername string `json:"player_username"`
	PlayerIP       string `json:"player_ip"`
}

type yhPayManualDepositRes struct {
	ErrorCode     int     `json:"error_code"`
	ErrorMessage  string  `json:"error_message"`
	InvoiceNumber string  `json:"invoice_number"`
	RefId         string  `json:"ref_id"`
	CurrencyCode  string  `json:"currency_code"`
	Amount        float64 `json:"amount"`
	Token         string  `json:"token"`
	RedirectTo    string  `json:"redirect_to"`
	BankCode      string  `json:"bank_code"`
	BankName      string  `json:"bank_name"`
	AccountName   string  `json:"account_name"`
	AccountNumber string  `json:"account_number"`
	StatusCode    int     `json:"status_code"`
	StatusMessage string  `json:"status_message"`
}

func yhPayPostJSON(ctx context.Context, fullURL string, body any, out any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}
	xrlog.DetailLog.Infof(ctx, "yhpay http request url=%s body=%s", fullURL, yhPayLogSafeJSON(payload))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, strings.NewReader(string(payload)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		xrlog.DetailLog.Errorf(ctx, "yhpay http do failed url=%s err=%v", fullURL, err)
		return err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		xrlog.DetailLog.Errorf(ctx, "yhpay http read failed url=%s status=%d err=%v", fullURL, resp.StatusCode, err)
		return err
	}
	xrlog.DetailLog.Infof(ctx, "yhpay http response url=%s status=%d body=%s", fullURL, resp.StatusCode, yhPayLogSafeJSON(raw))
	if len(raw) == 0 {
		return fmt.Errorf("yhpay empty response status=%d", resp.StatusCode)
	}
	if err = json.Unmarshal(raw, out); err != nil {
		xrlog.DetailLog.Warningf(ctx, "yhpay decode failed url=%s status=%d body=%s err=%v", fullURL, resp.StatusCode, truncateYhPayLog(string(raw), yhPayLogBodyMax), err)
		return fmt.Errorf("yhpay decode failed status=%d: %w", resp.StatusCode, err)
	}
	return nil
}

func getYhPayActiveCfg() (*entity.YhPayCfg, error) {
	cfg := cfgdao.GetYhPayCfgCached()
	if cfg == nil || !cfg.Enabled || cfg.MerchantCode == "" || cfg.ApiKey == "" {
		return nil, fmt.Errorf("yhpay not configured")
	}
	return cfg, nil
}

func yhPayFiatApiHost(cfg *entity.YhPayCfg) string {
	if cfg == nil {
		return ""
	}
	return strings.TrimRight(strings.TrimSpace(cfg.ApiHost), "/")
}

func createYhPayManualDeposit(ctx context.Context, cfg *entity.YhPayCfg, orderId, playerName, playerIP, currency string, amount float64) (payUrl string, err error) {
	host := yhPayFiatApiHost(cfg)
	if host == "" {
		return "", fmt.Errorf("yhpay api host not configured")
	}
	currency = strings.ToUpper(strings.TrimSpace(currency))
	if currency != yhPaySupportedCurrency {
		return "", fmt.Errorf("yhpay unsupported currency=%s", currency)
	}

	amountStr := formatYhPayAmount(amount)
	hash := yhPayHMAC(cfg.ApiKey, cfg.MerchantCode+orderId+currency+amountStr)
	callbackURL := strings.TrimRight(cfg.CallbackBaseUrl, "/") + yhPayPayinCallbackPath
	returnURL := cfg.ReturnUrl
	if returnURL == "" {
		returnURL = cfg.CallbackBaseUrl
	}

	reqBody := &yhPayManualDepositReq{
		MerchantCode:   cfg.MerchantCode,
		RefId:          orderId,
		CurrencyCode:   currency,
		Amount:         amountStr,
		Hash:           hash,
		CallbackURL:    callbackURL,
		ReturnURL:      returnURL,
		PlayerUsername: playerName,
		PlayerIP:       playerIP,
	}
	var res yhPayManualDepositRes
	url := host + yhPayManualDepositPath
	if err = yhPayPostJSON(ctx, url, reqBody, &res); err != nil {
		return "", err
	}
	if res.ErrorCode != 0 {
		msg := strings.TrimSpace(res.ErrorMessage)
		if msg == "" {
			msg = strings.TrimSpace(res.StatusMessage)
		}
		if msg == "" {
			msg = fmt.Sprintf("error_code=%d", res.ErrorCode)
		}
		return "", fmt.Errorf("yhpay manual deposit failed: %s", msg)
	}
	payUrl = normalizeYhPayRedirect(res.RedirectTo)
	if payUrl == "" {
		return "", fmt.Errorf("yhpay manual deposit empty redirect invoice=%s", res.InvoiceNumber)
	}
	xrlog.DetailLog.Infof(ctx, "yhpay manual deposit ok orderId=%s currency=%s amount=%s invoice=%s bank=%s payUrl=%s",
		orderId, currency, amountStr, res.InvoiceNumber, res.BankCode, payUrl)
	return payUrl, nil
}
