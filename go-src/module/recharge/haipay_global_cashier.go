package recharge

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"xr-game-server/constants/country"
	"xr-game-server/core/xrlog"
	"xr-game-server/dao/cfgdao"
	rechargeentity "xr-game-server/entity/recharge"
	"xr-game-server/module/fxrate"
)

const (
	haiPayGlobalCollectApplyPath   = "/global/cashier/collect/apply"
	haiPayGlobalCollectQueryV2Path = "/global/cashier/collect/query/v2"
)

func (p *haiPayProvider) quoteGlobalCashier(ctx context.Context, priceUSD float64, regionHint string) (string, float64, error) {
	if priceUSD < haiPayMinUsdAmount {
		return "", 0, fmt.Errorf("haipay usd amount must be >= %.2f", haiPayMinUsdAmount)
	}
	region, err := haiPayResolveGlobalCashierRegion(regionHint)
	if err != nil {
		return "", 0, err
	}
	currency, enabled := haiPayGlobalCashierCountrySelection(region)
	if !enabled {
		return "", 0, fmt.Errorf("haipay global cashier region disabled region=%s", region)
	}
	conversion, err := fxrate.ConvertUSD(ctx, priceUSD, currency)
	if err != nil {
		return "", 0, err
	}
	amount := math.Round(conversion.TargetAmount*100) / 100
	if amount <= 0 {
		return "", 0, fmt.Errorf("haipay global cashier amount invalid region=%s currency=%s priceUsd=%v rate=%v", region, currency, priceUSD, conversion.Rate)
	}
	xrlog.DetailLog.Infof(ctx, "channelPay global cashier fx rate source=%s cached=%t USD/%s=%v", conversion.Source, conversion.Cached, conversion.TargetCurrency, conversion.Rate)
	return currency, amount, nil
}

func (p *haiPayProvider) createGlobalCashierPay(ctx context.Context, req *ChannelPayCreateReq) (*ChannelPayCreateRes, error) {
	cfg := cfgdao.GetHaiPayCfgCached()
	if cfg == nil || !cfgdao.HaiPayEnabled() {
		return nil, fmt.Errorf("haipay not configured")
	}
	region, err := haiPayResolveGlobalCashierRegion(req.Region)
	if err != nil {
		return nil, err
	}
	configuredCurrency, enabled := haiPayGlobalCashierCountrySelection(region)
	if !enabled {
		return nil, fmt.Errorf("haipay global cashier region disabled region=%s", region)
	}
	currency := strings.ToUpper(strings.TrimSpace(req.Currency))
	if currency == "" || currency != configuredCurrency {
		return nil, fmt.Errorf("haipay global cashier currency mismatch region=%s currency=%s want=%s", region, currency, configuredCurrency)
	}
	amount := req.Amount
	if amount <= 0 {
		return nil, fmt.Errorf("haipay global cashier amount invalid region=%s currency=%s amount=%v", region, currency, amount)
	}
	if currency == "USD" && amount < haiPayMinUsdAmount {
		return nil, fmt.Errorf("haipay usd amount must be >= %.2f got=%v", haiPayMinUsdAmount, amount)
	}
	name := strings.TrimSpace(req.PlayerName)
	if name == "" {
		name = fmt.Sprintf("User %d", req.UserID)
	}
	if !strings.Contains(name, " ") {
		name += " User"
	}
	email := strings.TrimSpace(req.Email)
	if email == "" {
		email = fmt.Sprintf("user%d@noreply.local", req.UserID)
	}
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
		"appId":           int64(country.HaiPayAppIDCashier),
		"orderId":         req.OrderID,
		"name":            name,
		"email":           email,
		"amount":          fmt.Sprintf("%.2f", amount),
		"currency":        currency,
		"callBackUrl":     returnURL,
		"callBackFailUrl": failURL,
		"notifyUrl":       strings.TrimRight(cfg.CallbackBaseUrl, "/") + haiPayCollectNotifyPath,
		"subject":         subject,
		"region":          region,
		"partnerUserId":   strconv.FormatUint(req.UserID, 10),
	}
	if orderID := strings.TrimSpace(req.OrderID); orderID != "" {
		body["body"] = "order:" + orderID
	}

	xrlog.DetailLog.Infof(ctx, "haipay global cashier apply sign-before params=%s", haiPaySafeParams(body))
	sign, err := haiPaySign(ctx, body, cfg.MerchantSecretKey, cfg.MerchantPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("haipay global cashier sign: %w", err)
	}
	body["sign"] = sign

	var res haiPayApplyAPIRes
	url := strings.TrimRight(cfg.ApiHost, "/") + haiPayGlobalCollectApplyPath
	if err = haiPayPostJSON(ctx, url, body, &res); err != nil {
		return nil, err
	}
	if res.Status != "1" {
		msg := strings.TrimSpace(res.Msg)
		if msg == "" {
			msg = res.Error
		}
		return nil, fmt.Errorf("haipay global cashier apply failed status=%s error=%s msg=%s", res.Status, res.Error, msg)
	}
	if res.Data == nil {
		return nil, fmt.Errorf("haipay global cashier apply empty data")
	}
	payURL := strings.TrimSpace(res.Data.PayUrl)
	if payURL == "" {
		return nil, fmt.Errorf("haipay global cashier empty payUrl")
	}
	return &ChannelPayCreateRes{PayURL: payURL, ThirdOrderID: strings.TrimSpace(res.Data.OrderNo)}, nil
}

func haiPayGlobalCashierCountrySelection(region string) (string, bool) {
	aggregate := cfgdao.GetHaiPayCollectionCountryCfgCached(rechargeentity.HaiPayBizTypeNormalCollection, region)
	if aggregate == nil || aggregate.Config == nil || !aggregate.Config.Enabled {
		return country.HaiPayGlobalCashierCurrency(region), false
	}
	currency := strings.ToUpper(strings.TrimSpace(aggregate.Config.CurrencyCode))
	if !containsHaiPayString(country.HaiPayGlobalCashierCurrencies(region), currency) {
		return country.HaiPayGlobalCashierCurrency(region), false
	}
	return currency, true
}

func haiPayResolveGlobalCashierRegion(regionHint string) (string, error) {
	code := strings.ToUpper(strings.TrimSpace(regionHint))
	if code == "" {
		return "", fmt.Errorf("haipay global cashier region required")
	}
	if country.IsHaiPayGlobalCashierRegion(code) {
		return code, nil
	}
	if region := country.HaiPayGlobalCashierCountryFromCurrency(code); region != "" {
		return region, nil
	}
	return "", fmt.Errorf("unsupported haipay global cashier region=%s", code)
}
