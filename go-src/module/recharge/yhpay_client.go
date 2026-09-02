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

const (
	yhPayDepositV2Path      = "/Payin/DepositV2"
	yhPayDepositBanksPath   = "/Payin/DepositSenderBank/"
	yhPayCryptoDepositPath  = "/api/crypto/deposit"
	yhPayPayinCallbackPath  = "/webhook/yhpay/payin"
	yhPayCryptoCallbackPath = "/webhook/yhpay/crypto"
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

type yhPayDepositV2Req struct {
	MerchantCode      string `json:"MerchantCode"`
	ReturnURL         string `json:"ReturnURL"`
	FailedReturnURL   string `json:"FailedReturnURL"`
	HTTPPostURL       string `json:"HTTPPostURL"`
	Amount            string `json:"Amount"`
	Currency          string `json:"Currency"`
	ItemID            string `json:"ItemID"`
	ItemDescription   string `json:"ItemDescription"`
	PlayerId          string `json:"PlayerId"`
	Hash              string `json:"Hash"`
	BankCode          string `json:"BankCode"`
	SenderVerification int   `json:"SenderVerification"`
	ClientFullName    string `json:"ClientFullName"`
}

type yhPayDepositV2Res struct {
	Transaction any    `json:"transaction"`
	Status      int    `json:"status"`
	Token       string `json:"token"`
	RedirectTo  string `json:"redirect_to"`
	Amount      any    `json:"Amount"`
	Currency    string `json:"Currency"`
	Message     string `json:"message"`
}

type yhPayBankListRes struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
	Bank      []struct {
		BankCode string `json:"BankCode"`
		BankName string `json:"BankName"`
	} `json:"Bank"`
}

type yhPayCryptoDepositReq struct {
	MerchantCode       string `json:"merchant_code"`
	RefId              string `json:"ref_id"`
	PlayerUsername     string `json:"player_username"`
	PlayerIp           string `json:"player_ip"`
	FiatCurrencyCode   string `json:"fiat_currency_code"`
	FiatAmount         string `json:"fiat_amount"`
	CryptoCurrencyCode string `json:"crypto_currency_code"`
	Network            string `json:"network"`
	ClientUrl          string `json:"client_url"`
	Hash               string `json:"hash"`
	Lang               string `json:"lang"`
}

type yhPayCryptoDepositRes struct {
	ErrorCode     int     `json:"errorcode"`
	Message       string  `json:"message"`
	InvoiceNumber string  `json:"invoice_number"`
	CurrencyCode  string  `json:"currency_code"`
	Amount        float64 `json:"amount"`
	Token         string  `json:"token"`
	RedirectTo    string  `json:"redirect_to"`
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

func yhPayCryptoApiHost(cfg *entity.YhPayCfg) string {
	if cfg == nil {
		return ""
	}
	if host := strings.TrimRight(strings.TrimSpace(cfg.CryptoApiHost), "/"); host != "" {
		return host
	}
	// 兼容旧配置:未单独填 USDT host 时回退法币 host
	return yhPayFiatApiHost(cfg)
}

func pickYhPayBankCode(ctx context.Context, cfg *entity.YhPayCfg, currency string) (string, error) {
	host := yhPayFiatApiHost(cfg)
	if host == "" {
		return "", fmt.Errorf("yhpay fiat api host not configured")
	}
	url := host + yhPayDepositBanksPath
	var res yhPayBankListRes
	if err := yhPayPostJSON(ctx, url, map[string]string{
		"MerchantCode": cfg.MerchantCode,
		"Currency":     currency,
		"APIKey":       cfg.ApiKey,
	}, &res); err != nil {
		return "", err
	}
	if res.ErrorCode != 0 || len(res.Bank) == 0 {
		return "", fmt.Errorf("yhpay no bank available currency=%s code=%d msg=%s", currency, res.ErrorCode, res.Message)
	}
	bankCode := strings.TrimSpace(res.Bank[0].BankCode)
	xrlog.DetailLog.Infof(ctx, "yhpay bank picked currency=%s bankCode=%s bankCount=%d", currency, bankCode, len(res.Bank))
	return bankCode, nil
}

func createYhPayDepositV2(ctx context.Context, cfg *entity.YhPayCfg, orderId, playerId, playerName, currency string, amount float64) (payUrl string, err error) {
	host := yhPayFiatApiHost(cfg)
	if host == "" {
		return "", fmt.Errorf("yhpay fiat api host not configured")
	}
	bankCode, err := pickYhPayBankCode(ctx, cfg, currency)
	if err != nil {
		return "", err
	}
	amountStr := formatYhPayAmount(amount)
	itemID := orderId
	hash := yhPayHMAC(cfg.ApiKey, cfg.MerchantCode+itemID+currency+amountStr)
	callbackURL := strings.TrimRight(cfg.CallbackBaseUrl, "/") + yhPayPayinCallbackPath
	returnURL := cfg.ReturnUrl
	if returnURL == "" {
		returnURL = cfg.CallbackBaseUrl
	}
	failedURL := cfg.FailedReturnUrl
	if failedURL == "" {
		failedURL = returnURL
	}
	reqBody := &yhPayDepositV2Req{
		MerchantCode:       cfg.MerchantCode,
		ReturnURL:          returnURL,
		FailedReturnURL:    failedURL,
		HTTPPostURL:        callbackURL,
		Amount:             amountStr,
		Currency:           currency,
		ItemID:             itemID,
		ItemDescription:    "recharge " + orderId,
		PlayerId:           playerId,
		Hash:               hash,
		BankCode:           bankCode,
		SenderVerification: 0,
		ClientFullName:     playerName,
	}
	var res yhPayDepositV2Res
	url := host + yhPayDepositV2Path
	if err = yhPayPostJSON(ctx, url, reqBody, &res); err != nil {
		return "", err
	}
	payUrl = normalizeYhPayRedirect(res.RedirectTo)
	if payUrl == "" {
		if res.Message != "" {
			return "", fmt.Errorf("yhpay depositv2 failed: %s", res.Message)
		}
		return "", fmt.Errorf("yhpay depositv2 empty redirect status=%d", res.Status)
	}
	xrlog.DetailLog.Infof(ctx, "yhpay depositv2 ok orderId=%s currency=%s amount=%s bankCode=%s status=%d payUrl=%s",
		orderId, currency, amountStr, bankCode, res.Status, payUrl)
	return payUrl, nil
}

func createYhPayCryptoDeposit(ctx context.Context, cfg *entity.YhPayCfg, orderId, playerName, playerIP, cryptoCode string, usdAmount float64) (payUrl string, err error) {
	host := yhPayCryptoApiHost(cfg)
	if host == "" {
		return "", fmt.Errorf("yhpay crypto api host not configured")
	}
	fiatCurrency := "USD"
	fiatAmount := formatYhPayAmount(usdAmount)
	network := "TRC20"
	clientURL := cfg.CallbackBaseUrl
	if clientURL == "" {
		clientURL = cfg.ReturnUrl
	}
	hashPayload := cfg.MerchantCode + orderId + playerName + playerIP + fiatCurrency + fiatAmount + cryptoCode + clientURL
	reqBody := &yhPayCryptoDepositReq{
		MerchantCode:       cfg.MerchantCode,
		RefId:              orderId,
		PlayerUsername:     playerName,
		PlayerIp:           playerIP,
		FiatCurrencyCode:   fiatCurrency,
		FiatAmount:         fiatAmount,
		CryptoCurrencyCode: cryptoCode,
		Network:            network,
		ClientUrl:          clientURL,
		Hash:               yhPayHMAC(cfg.ApiKey, hashPayload),
		Lang:               "en",
	}
	var res yhPayCryptoDepositRes
	url := host + yhPayCryptoDepositPath
	if err = yhPayPostJSON(ctx, url, reqBody, &res); err != nil {
		return "", err
	}
	if res.ErrorCode != 0 {
		return "", fmt.Errorf("yhpay crypto deposit failed code=%d msg=%s", res.ErrorCode, res.Message)
	}
	payUrl = normalizeYhPayRedirect(res.RedirectTo)
	if payUrl == "" {
		return "", fmt.Errorf("yhpay crypto empty redirect")
	}
	xrlog.DetailLog.Infof(ctx, "yhpay crypto deposit ok orderId=%s crypto=%s usd=%s invoice=%s payUrl=%s",
		orderId, cryptoCode, fiatAmount, res.InvoiceNumber, payUrl)
	return payUrl, nil
}
