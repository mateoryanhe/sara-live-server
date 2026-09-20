package recharge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"xr-game-server/constants/country"
	"xr-game-server/core/xrlog"
	"xr-game-server/dao/cfgdao"
	rechargeentity "xr-game-server/entity/recharge"
	"xr-game-server/module/fxrate"
)

const (
	haiPayProviderName      = "haipay"
	haiPayCollectNotifyPath = "/webhook/haipay/collect/notify"
	haiPayMinUsdAmount      = 0.99
	haiPayHTTPTimeout       = 30 * time.Second
	haiPayLogBodyMax        = 4096
)

type haiPayProvider struct{}

func (p *haiPayProvider) Name() string { return haiPayProviderName }

func (p *haiPayProvider) Enabled() bool { return cfgdao.HaiPayEnabled() }

func (p *haiPayProvider) resolveRegion(regionHint string) (string, error) {
	regionHint = strings.TrimSpace(regionHint)
	if regionHint == "" {
		return "", fmt.Errorf("haipay region required")
	}
	return haiPayResolveRegion(regionHint)
}

func (p *haiPayProvider) QuotePay(ctx context.Context, priceUsd float64, regionHint string, payChannel uint8) (payCurrency string, payAmount float64, err error) {
	if isHaiPayGlobalCashierBiz(haiPayCollectionBizType(payChannel)) {
		return p.quoteGlobalCashier(ctx, priceUsd, regionHint)
	}
	if priceUsd < haiPayMinUsdAmount {
		return "", 0, fmt.Errorf("haipay usd amount must be >= %.2f", haiPayMinUsdAmount)
	}
	region, err := p.resolveRegion(regionHint)
	if err != nil {
		return "", 0, err
	}
	currency, _, enabled := haiPayCollectionCountrySelection(region, payChannel)
	if !enabled {
		return "", 0, fmt.Errorf("haipay collection region disabled region=%s", region)
	}
	if currency == "" {
		return "", 0, fmt.Errorf("haipay collection currency missing for region=%s", region)
	}
	conversion, err := fxrate.ConvertUSD(ctx, priceUsd, currency)
	if err != nil {
		return "", 0, err
	}
	// 订单保存值必须和本地代收接口实际提交精度一致，避免查单对账出现尾差。
	amount := haiPayRoundCollectionAmount(currency, conversion.TargetAmount)
	if amount <= 0 {
		return "", 0, fmt.Errorf("haipay collection amount invalid region=%s currency=%s priceUsd=%v rate=%v", region, currency, priceUsd, conversion.Rate)
	}
	xrlog.DetailLog.Infof(ctx, "channelPay fx rate source=%s cached=%t USD/%s=%v", conversion.Source, conversion.Cached, conversion.TargetCurrency, conversion.Rate)
	return currency, amount, nil
}

func (p *haiPayProvider) CreatePay(ctx context.Context, req *ChannelPayCreateReq) (*ChannelPayCreateRes, error) {
	if isHaiPayGlobalCashierBiz(haiPayCollectionBizType(req.PayChannel)) {
		return p.createGlobalCashierPay(ctx, req)
	}
	cfg := cfgdao.GetHaiPayCfgCached()
	if cfg == nil || !cfgdao.HaiPayEnabled() {
		return nil, fmt.Errorf("haipay not configured")
	}
	region, err := p.resolveRegion(req.Region)
	if err != nil {
		return nil, err
	}
	currency := strings.ToUpper(strings.TrimSpace(req.Currency))
	wantCurrency, methods, enabled := haiPayCollectionCountrySelection(region, req.PayChannel)
	if !enabled {
		return nil, fmt.Errorf("haipay collection region disabled region=%s", region)
	}
	if currency == "" || currency != wantCurrency {
		return nil, fmt.Errorf("haipay collection currency mismatch region=%s currency=%s want=%s", region, currency, wantCurrency)
	}
	amount := req.Amount
	if amount <= 0 {
		return nil, fmt.Errorf("haipay collection amount invalid region=%s currency=%s amount=%v", region, currency, amount)
	}
	if currency == "USD" && amount < haiPayMinUsdAmount {
		return nil, fmt.Errorf("haipay usd amount must be >= %.2f got=%v", haiPayMinUsdAmount, amount)
	}
	method, err := selectHaiPayCollectionMethod(methods, req.PayType, req.InBankCode)
	if err != nil {
		return nil, fmt.Errorf("haipay collection payment method region=%s currency=%s: %w", region, currency, err)
	}
	if err = validateHaiPayCollectionMethodAmount(region, currency, method, amount); err != nil {
		return nil, err
	}
	if req.PayChannel != rechargeentity.RechargeCfgTypeCoinMerchant {
		return nil, fmt.Errorf("haipay unsupported local collection channel=%d", req.PayChannel)
	}
	if _, credentialErr := haiPayCoinMerchantCollectionCredential(region, currency); credentialErr != nil {
		return nil, credentialErr
	}
	appID, ok := country.LookupHaiPayAppID(currency)
	if !ok {
		return nil, fmt.Errorf("haipay appId enum missing currency=%s", currency)
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
	phone := strings.TrimSpace(req.Phone)
	if phone == "" {
		return nil, fmt.Errorf("haipay local collect phone missing userId=%d region=%s", req.UserID, region)
	}

	amountStr := haiPayFormatCollectionAmount(currency, amount)
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
		"appId":           appID,
		"orderId":         req.OrderID,
		"name":            name,
		"phone":           phone,
		"email":           email,
		"amount":          amountStr,
		"currency":        currency,
		"payType":         method.PayType,
		"inBankCode":      method.InBankCode,
		"callBackUrl":     returnURL,
		"callBackFailUrl": failURL,
		"notifyUrl":       notifyURL,
		"subject":         subject,
		"partnerUserId":   strconv.FormatUint(req.UserID, 10),
	}
	if remark := strings.TrimSpace(req.OrderID); remark != "" {
		body["body"] = "order:" + remark
	}

	// 签名前记录脱敏业务字段，便于与 HaiPay 对账。
	xrlog.DetailLog.Infof(ctx, "haipay apply sign-before(request) params=%s", haiPaySafeParams(body))
	sign, err := haiPaySign(ctx, body, cfg.MerchantSecretKey, cfg.MerchantPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("haipay sign: %w", err)
	}
	body["sign"] = sign

	var res haiPayApplyAPIRes
	url := strings.TrimRight(cfg.ApiHost, "/") + haiPayLocalCollectPath(currency, "apply")
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
	// 普通充值与币商充值的下单应答统一不验签；最终支付状态由回调触发主动查单确认。
	payURL := strings.TrimSpace(res.Data.PayUrl)
	if payURL == "" {
		payURL = strings.TrimSpace(res.Data.QrCode)
	}
	if payURL == "" {
		return nil, fmt.Errorf("haipay empty payUrl")
	}
	return &ChannelPayCreateRes{
		PayURL:       payURL,
		ThirdOrderID: strings.TrimSpace(res.Data.OrderNo),
	}, nil
}

type haiPayCollectionMethod struct {
	PayType    string
	InBankCode string
}

// haiPayCollectionCountrySelection 返回 CMS 配置的本地代收币种和可用支付方式。
// 国家/地区是否可用只由 App 可见开关决定；币种只决定本地代收接口。
// 不匹配当前币种的历史支付方式会被忽略，并在创建订单时以“未配置支付方式”明确报错。
func haiPayCollectionCountrySelection(region string, payChannel uint8) (currency string, methods []haiPayCollectionMethod, enabled bool) {
	if payChannel == rechargeentity.RechargeCfgTypeCoinMerchant {
		config := cfgdao.GetHaiPayCoinMerchantCollectionCfgCached(region)
		if !cfgdao.HaiPayCoinMerchantCollectionCfgComplete(config) {
			return country.HaiPayCollectionCurrency(region), nil, false
		}
		currency = strings.ToUpper(strings.TrimSpace(config.CurrencyCode))
		payType := strings.ToUpper(strings.TrimSpace(config.PayType))
		canonicalCode, ok := country.ResolveHaiPayCollectionPaymentMethodCode(
			region, currency, payType, config.InBankCode,
		)
		if !ok || canonicalCode == "" {
			return currency, nil, false
		}
		return currency, []haiPayCollectionMethod{{PayType: payType, InBankCode: canonicalCode}}, config.Enabled
	}
	bizType := haiPayCollectionBizType(payChannel)
	aggregate := cfgdao.GetHaiPayCollectionCountryCfgCached(bizType, region)
	if aggregate == nil || aggregate.Config == nil {
		return country.HaiPayCollectionCurrency(region), nil, false
	}
	currency = strings.ToUpper(strings.TrimSpace(aggregate.Config.CurrencyCode))
	enabled = aggregate.Config.Enabled
	return currency, nil, enabled
}

func haiPayCoinMerchantCollectionCredential(region, currency string) (*rechargeentity.HaiPayCoinMerchantCollectionCfg, error) {
	config := cfgdao.GetHaiPayCoinMerchantCollectionCfgCached(region)
	if !cfgdao.HaiPayCoinMerchantCollectionCfgComplete(config) {
		return nil, fmt.Errorf("haipay coin merchant collection config missing region=%s", strings.ToUpper(strings.TrimSpace(region)))
	}
	wantCurrency := strings.ToUpper(strings.TrimSpace(currency))
	if strings.ToUpper(strings.TrimSpace(config.CurrencyCode)) != wantCurrency {
		return nil, fmt.Errorf("haipay coin merchant collection currency mismatch region=%s currency=%s", region, wantCurrency)
	}
	return config, nil
}

func haiPayCollectionBizType(payChannel uint8) rechargeentity.HaiPayBizType {
	if payChannel == rechargeentity.RechargeCfgTypeCoinMerchant {
		return rechargeentity.HaiPayBizTypeCoinMerchantCollection
	}
	return rechargeentity.HaiPayBizTypeNormalCollection
}

func selectHaiPayCollectionMethod(methods []haiPayCollectionMethod, payType, inBankCode string) (haiPayCollectionMethod, error) {
	payType = strings.ToUpper(strings.TrimSpace(payType))
	inBankCode = strings.TrimSpace(inBankCode)
	if len(methods) == 0 {
		return haiPayCollectionMethod{}, fmt.Errorf("no configured method")
	}
	if len(methods) > 1 {
		return haiPayCollectionMethod{}, fmt.Errorf("multiple configured methods for active currency")
	}
	configured := methods[0]
	if payType == "" && inBankCode == "" {
		// App 不上报支付类型和支付编码，服务端使用当前国家和代收币种唯一配置的一条。
		return configured, nil
	}
	if payType == "" || inBankCode == "" {
		return haiPayCollectionMethod{}, fmt.Errorf("payType and inBankCode must be supplied together")
	}
	if configured.PayType == payType && strings.EqualFold(configured.InBankCode, inBankCode) {
		return configured, nil
	}
	return haiPayCollectionMethod{}, fmt.Errorf("method is not enabled payType=%s inBankCode=%s", payType, inBankCode)
}

func validateHaiPayCollectionMethodAmount(region, currency string, selected haiPayCollectionMethod, amount float64) error {
	for _, method := range country.ListHaiPayCollectionPaymentMethods(region) {
		if !method.Available || method.CurrencyCode != currency || method.PayType != selected.PayType || !strings.EqualFold(method.InBankCode, selected.InBankCode) {
			continue
		}
		minAmount, minErr := strconv.ParseFloat(method.MinAmount, 64)
		maxAmount, maxErr := strconv.ParseFloat(method.MaxAmount, 64)
		if minErr == nil && amount < minAmount {
			return fmt.Errorf("haipay collection amount below limit amount=%v min=%s currency=%s method=%s/%s", amount, method.MinAmount, currency, selected.PayType, selected.InBankCode)
		}
		if maxErr == nil && amount > maxAmount {
			return fmt.Errorf("haipay collection amount above limit amount=%v max=%s currency=%s method=%s/%s", amount, method.MaxAmount, currency, selected.PayType, selected.InBankCode)
		}
		return nil
	}
	return fmt.Errorf("haipay collection method catalog missing region=%s currency=%s method=%s/%s", region, currency, selected.PayType, selected.InBankCode)
}

func haiPayLocalCollectPath(currency, action string) string {
	return "/" + strings.ToLower(strings.TrimSpace(currency)) + "/collect/" + strings.TrimSpace(action)
}

func haiPayFormatCollectionAmount(currency string, amount float64) string {
	switch strings.ToUpper(strings.TrimSpace(currency)) {
	case "IDR", "VND", "KRW", "JPY", "CLP", "UGX", "XAF":
		return fmt.Sprintf("%.0f", math.Round(amount))
	default:
		return fmt.Sprintf("%.2f", amount)
	}
}

func haiPayRoundCollectionAmount(currency string, amount float64) float64 {
	switch strings.ToUpper(strings.TrimSpace(currency)) {
	case "IDR", "VND", "KRW", "JPY", "CLP", "UGX", "XAF":
		return math.Round(amount)
	default:
		return math.Round(amount*100) / 100
	}
}

func haiPayResolveRegion(currencyOrRegion string) (string, error) {
	code := strings.ToUpper(strings.TrimSpace(currencyOrRegion))
	if code == "" {
		return "", fmt.Errorf("empty region")
	}
	if region := country.HaiPayCollectionCountryFromCurrency(code); region != "" {
		return region, nil
	}
	if country.IsHaiPayRegion(code) {
		return code, nil
	}
	if region := country.HaiPayGlobalCashierCountryFromCurrency(code); region != "" {
		return region, nil
	}
	if country.IsHaiPayGlobalCashierRegion(code) {
		return code, nil
	}
	return "", fmt.Errorf("unsupported region currency=%s", code)
}

type haiPayApplyAPIRes struct {
	Status string `json:"status"`
	Error  string `json:"error"`
	Msg    string `json:"msg"`
	Data   *struct {
		OrderId  string `json:"orderId"`
		OrderNo  string `json:"orderNo"`
		PayUrl   string `json:"payUrl"`
		BankNo   string `json:"bankNo"`
		BankCode string `json:"bankCode"`
		QrCode   string `json:"qrCode"`
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
	for _, k := range []string{
		"sign", "merchantSecretKey", "merchantPrivateKey",
		"name", "phone", "email", "accountNo", "identifyType",
		"address1", "address2", "address3", "postalCode",
	} {
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
