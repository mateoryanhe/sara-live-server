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
	liveentity "xr-game-server/entity/live"
	rechargeentity "xr-game-server/entity/recharge"
)

const haiPayPayoutNotifyPath = "/webhook/haipay/payout/notify"

// haiPayParseKVMap 解析 "IDR:16000,PHP:58" 或 "IDR=16000" 形式键值表
func haiPayParseKVMap(raw string) map[string]string {
	out := make(map[string]string)
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return out
	}
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		sep := ":"
		if !strings.Contains(part, ":") && strings.Contains(part, "=") {
			sep = "="
		}
		kv := strings.SplitN(part, sep, 2)
		if len(kv) != 2 {
			continue
		}
		k := strings.ToUpper(strings.TrimSpace(kv[0]))
		v := strings.TrimSpace(kv[1])
		if k == "" || v == "" {
			continue
		}
		out[k] = v
	}
	return out
}

func HaiPayPayoutAppId(cfg *rechargeentity.HaiPayCfg, currency string) (int64, error) {
	if cfg == nil {
		return 0, fmt.Errorf("haipay cfg nil")
	}
	currency = strings.ToUpper(strings.TrimSpace(currency))
	m := haiPayParseKVMap(cfg.PayoutAppIds)
	raw, ok := m[currency]
	if !ok || raw == "" {
		return 0, fmt.Errorf("payout appId missing for currency=%s", currency)
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid payout appId currency=%s raw=%s", currency, raw)
	}
	return id, nil
}

func haiPayFormatPayoutAmount(currency string, amount float64) string {
	currency = strings.ToUpper(strings.TrimSpace(currency))
	switch currency {
	case "IDR", "VND", "KRW", "JPY", "CLP", "UGX", "XAF":
		return fmt.Sprintf("%.0f", math.Round(amount))
	default:
		return fmt.Sprintf("%.2f", amount)
	}
}

type HaiPayPayoutApplyReq struct {
	OrderID       string
	Currency      string
	Amount        float64
	AccountType   string
	BankCode      string
	AccountNo     string
	Name          string
	Phone         string
	Email         string
	PartnerUserID string
	Subject       string
	Body          string
}

type HaiPayPayoutApplyRes struct {
	OrderID string
	OrderNo string
}

type haiPayPayoutAPIRes struct {
	Status string `json:"status"`
	Error  string `json:"error"`
	Msg    string `json:"msg"`
	Data   *struct {
		OrderId string `json:"orderId"`
		OrderNo string `json:"orderNo"`
	} `json:"data"`
}

type haiPayPayoutQueryData struct {
	OrderId  string `json:"orderId"`
	OrderNo  string `json:"orderNo"`
	Amount   string `json:"amount"`
	Fee      string `json:"fee"`
	Status   int    `json:"status"`
	PayTime  string `json:"payTime"`
	ErrorMsg string `json:"errorMsg"`
}

type haiPayPayoutQueryAPIRes struct {
	Status string                 `json:"status"`
	Error  string                 `json:"error"`
	Msg    string                 `json:"msg"`
	Data   *haiPayPayoutQueryData `json:"data"`
}

// HaiPayApplyPayout 调用 /{currency}/pay/apply 代付
func HaiPayApplyPayout(ctx context.Context, req *HaiPayPayoutApplyReq) (*HaiPayPayoutApplyRes, error) {
	if req == nil {
		return nil, fmt.Errorf("nil payout req")
	}
	cfg := cfgdao.GetHaiPayCfgCached()
	if cfg == nil || !cfg.PayoutEnabled {
		return nil, fmt.Errorf("haipay payout not enabled")
	}
	currency := strings.ToUpper(strings.TrimSpace(req.Currency))
	if currency == "" {
		return nil, fmt.Errorf("empty currency")
	}
	appId, err := HaiPayPayoutAppId(cfg, currency)
	if err != nil {
		return nil, err
	}
	if req.Amount <= 0 {
		return nil, fmt.Errorf("invalid amount")
	}
	accountType := strings.TrimSpace(req.AccountType)
	if accountType == "" {
		accountType = liveentity.GuildTransferAccountTypeBank
	}
	accountType = strings.ToUpper(accountType)
	if !country.IsHaiPayPayoutAccountType(currency, accountType) {
		return nil, fmt.Errorf("unsupported payout accountType=%s currency=%s", accountType, currency)
	}
	name := strings.TrimSpace(req.Name)
	email := strings.ToLower(strings.TrimSpace(req.Email))
	phone := strings.TrimSpace(req.Phone)
	bankCode := strings.TrimSpace(req.BankCode)
	accountNo := strings.TrimSpace(req.AccountNo)
	if name == "" || email == "" || phone == "" || bankCode == "" || accountNo == "" {
		return nil, fmt.Errorf("missing payout payee fields")
	}
	if !strings.Contains(name, " ") {
		name = name + " User"
	}
	subject := strings.TrimSpace(req.Subject)
	if subject == "" {
		subject = strings.TrimSpace(cfg.PayoutSubject)
	}
	if subject == "" {
		subject = "GuildSettlement"
	}
	notifyURL := strings.TrimRight(cfg.CallbackBaseUrl, "/") + haiPayPayoutNotifyPath
	amountStr := haiPayFormatPayoutAmount(currency, req.Amount)

	body := map[string]any{
		"appId":         appId,
		"orderId":       strings.TrimSpace(req.OrderID),
		"amount":        amountStr,
		"accountType":   accountType,
		"bankCode":      bankCode,
		"accountNo":     accountNo,
		"name":          name,
		"phone":         phone,
		"email":         email,
		"partnerUserId": strings.TrimSpace(req.PartnerUserID),
		"notifyUrl":     notifyURL,
		"subject":       subject,
	}
	if b := strings.TrimSpace(req.Body); b != "" {
		body["body"] = b
	}

	xrlog.DetailLog.Infof(ctx, "haipay payout apply sign-before currency=%s params=%s", currency, haiPaySafeParams(body))
	sign, err := haiPaySign(ctx, body, cfg.MerchantSecretKey, cfg.MerchantPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("haipay payout sign: %w", err)
	}
	body["sign"] = sign

	var res haiPayPayoutAPIRes
	url := strings.TrimRight(cfg.ApiHost, "/") + "/" + strings.ToLower(currency) + "/pay/apply"
	if err = haiPayPostJSON(ctx, url, body, &res); err != nil {
		return nil, err
	}
	if res.Status != "1" {
		msg := strings.TrimSpace(res.Msg)
		if msg == "" {
			msg = res.Error
		}
		return nil, fmt.Errorf("haipay payout failed status=%s error=%s msg=%s", res.Status, res.Error, msg)
	}
	if res.Data == nil {
		return nil, fmt.Errorf("haipay payout empty data")
	}
	return &HaiPayPayoutApplyRes{
		OrderID: strings.TrimSpace(res.Data.OrderId),
		OrderNo: strings.TrimSpace(res.Data.OrderNo),
	}, nil
}

// haiPayQueryPayout 主动查询代付订单。回调只用于触发查单，最终状态以该接口结果为准。
func haiPayQueryPayout(ctx context.Context, row *liveentity.GuildIncomeSettlementLog) (*haiPayPayoutQueryData, error) {
	if row == nil || row.ID == 0 {
		return nil, fmt.Errorf("empty payout settlement")
	}
	cfg := cfgdao.GetHaiPayCfgCached()
	if cfg == nil || !cfg.PayoutEnabled {
		return nil, fmt.Errorf("haipay payout not enabled")
	}
	currency := strings.ToUpper(strings.TrimSpace(row.TransferCurrency))
	orderID := strings.TrimSpace(row.TransferOrderId)
	if currency == "" || orderID == "" {
		return nil, fmt.Errorf("haipay payout local order incomplete")
	}
	appID, err := HaiPayPayoutAppId(cfg, currency)
	if err != nil {
		return nil, err
	}
	body := map[string]any{
		"appId":   appID,
		"orderId": orderID,
	}
	if orderNo := strings.TrimSpace(row.TransferPlatformNo); orderNo != "" {
		body["orderNo"] = orderNo
	}
	sign, err := haiPaySign(ctx, body, cfg.MerchantSecretKey, cfg.MerchantPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("haipay payout query sign: %w", err)
	}
	body["sign"] = sign

	var res haiPayPayoutQueryAPIRes
	url := strings.TrimRight(cfg.ApiHost, "/") + "/" + strings.ToLower(currency) + "/pay/query"
	if err = haiPayPostJSON(ctx, url, body, &res); err != nil {
		return nil, err
	}
	if res.Status != "1" {
		msg := strings.TrimSpace(res.Msg)
		if msg == "" {
			msg = res.Error
		}
		return nil, fmt.Errorf("haipay payout query failed status=%s error=%s msg=%s", res.Status, res.Error, msg)
	}
	if res.Data == nil {
		return nil, fmt.Errorf("haipay payout query empty data")
	}
	queryOrderID := strings.TrimSpace(res.Data.OrderId)
	if queryOrderID == "" || queryOrderID != orderID {
		return nil, fmt.Errorf("haipay payout query order mismatch got=%s want=%s", queryOrderID, orderID)
	}
	localOrderNo := strings.TrimSpace(row.TransferPlatformNo)
	queryOrderNo := strings.TrimSpace(res.Data.OrderNo)
	if localOrderNo != "" && (queryOrderNo == "" || queryOrderNo != localOrderNo) {
		return nil, fmt.Errorf("haipay payout query platform order mismatch got=%s want=%s", queryOrderNo, localOrderNo)
	}
	if res.Data.Status == haiPayCollectStatusSuccess {
		queryAmount, amountErr := haiPayAmountCents(res.Data.Amount)
		if amountErr != nil {
			return nil, fmt.Errorf("invalid haipay payout query amount %q: %w", res.Data.Amount, amountErr)
		}
		expectedText := haiPayFormatPayoutAmount(currency, row.TransferLocalAmount)
		expectedAmount, amountErr := haiPayAmountCents(expectedText)
		if amountErr != nil {
			return nil, fmt.Errorf("invalid local payout amount %q: %w", expectedText, amountErr)
		}
		if queryAmount != expectedAmount {
			return nil, fmt.Errorf("haipay payout amount mismatch query=%s local=%s currency=%s", res.Data.Amount, expectedText, currency)
		}
	}
	xrlog.DetailLog.Infof(ctx, "haipay payout query ok orderId=%s orderNo=%s status=%d amount=%s currency=%s",
		orderID, res.Data.OrderNo, res.Data.Status, res.Data.Amount, currency)
	return res.Data, nil
}
