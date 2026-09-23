package recharge

import (
	"context"
	"fmt"
	"math"
	"strings"

	"xr-game-server/constants/country"
	"xr-game-server/core/xrlog"
	"xr-game-server/dao/cfgdao"
	liveentity "xr-game-server/entity/live"
)

const (
	haiPayPayoutNotifyPath = "/webhook/haipay/payout/notify"
	haiPayPayoutSubject    = "GuildSettlement"
)

func HaiPayPayoutAppID(currency string) (int64, error) {
	currency = strings.ToUpper(strings.TrimSpace(currency))
	appID, ok := country.LookupHaiPayAppID(currency)
	if !ok {
		return 0, fmt.Errorf("payout appId enum missing for currency=%s", currency)
	}
	return appID, nil
}

// HaiPayNormalizePayoutAmount 按 HaiPay 目标币种精度统一代付金额。
// 零小数币种在提交前取整，其余币种保留两位，调用方应保存该返回值作为实际代付金额。
func HaiPayNormalizePayoutAmount(currency string, amount float64) float64 {
	currency = strings.ToUpper(strings.TrimSpace(currency))
	switch currency {
	case "IDR", "VND", "KRW", "JPY", "CLP", "UGX", "XAF":
		return math.Round(amount)
	default:
		return math.Round(amount*100) / 100
	}
}

func haiPayFormatPayoutAmount(currency string, amount float64) string {
	amount = HaiPayNormalizePayoutAmount(currency, amount)
	if isHaiPayZeroDecimalCurrency(currency) {
		return fmt.Sprintf("%.0f", amount)
	}
	return fmt.Sprintf("%.2f", amount)
}

func isHaiPayZeroDecimalCurrency(currency string) bool {
	switch strings.ToUpper(strings.TrimSpace(currency)) {
	case "IDR", "VND", "KRW", "JPY", "CLP", "UGX", "XAF":
		return true
	default:
		return false
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
	Country       string
	IdentifyType  string
	Address1      string
	Address2      string
	Address3      string
	PostalCode    string
	PartnerUserID string
	Body          string
}

func haiPayPayoutExtraParams(req *HaiPayPayoutApplyReq) (map[string]any, error) {
	if req == nil {
		return nil, fmt.Errorf("nil payout req")
	}
	requirements := country.HaiPayPayoutExtraRequirementsFor(req.Currency, req.BankCode)
	params := make(map[string]any)
	identifyType := strings.TrimSpace(req.IdentifyType)
	if len(requirements.IdentifyTypeOptions) > 0 {
		identifyType = strings.ToUpper(identifyType)
	}
	if !requirements.IsAllowedIdentifyType(identifyType) {
		return nil, fmt.Errorf("payout identifyType missing or invalid currency=%s bankCode=%s", strings.ToUpper(strings.TrimSpace(req.Currency)), strings.TrimSpace(req.BankCode))
	}
	if requirements.IdentifyTypeRequired {
		params["identifyType"] = identifyType
	}
	if requirements.CountryRequired {
		countryCode := strings.ToUpper(strings.TrimSpace(req.Country))
		if countryCode == "" {
			return nil, fmt.Errorf("payout country missing currency=%s bankCode=%s", strings.ToUpper(strings.TrimSpace(req.Currency)), strings.TrimSpace(req.BankCode))
		}
		params["country"] = countryCode
	}
	if requirements.AddressRequired {
		address1 := strings.TrimSpace(req.Address1)
		address2 := strings.TrimSpace(req.Address2)
		address3 := strings.TrimSpace(req.Address3)
		postalCode := strings.TrimSpace(req.PostalCode)
		if address1 == "" || address2 == "" || address3 == "" || postalCode == "" {
			return nil, fmt.Errorf("payout address missing currency=%s bankCode=%s", strings.ToUpper(strings.TrimSpace(req.Currency)), strings.TrimSpace(req.BankCode))
		}
		params["address1"] = address1
		params["address2"] = address2
		params["address3"] = address3
		params["postalCode"] = postalCode
	}
	return params, nil
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
	if cfg == nil || !cfgdao.HaiPayEnabled() {
		return nil, fmt.Errorf("haipay config not ready")
	}
	callbackBaseURL := strings.TrimRight(strings.TrimSpace(cfg.CallbackBaseUrl), "/")
	if callbackBaseURL == "" {
		return nil, fmt.Errorf("haipay callback base url missing")
	}
	currency := strings.ToUpper(strings.TrimSpace(req.Currency))
	if currency == "" {
		return nil, fmt.Errorf("empty currency")
	}
	appId, err := HaiPayPayoutAppID(currency)
	if err != nil {
		return nil, err
	}
	normalizedAmount := HaiPayNormalizePayoutAmount(currency, req.Amount)
	if normalizedAmount <= 0 || math.IsNaN(normalizedAmount) || math.IsInf(normalizedAmount, 0) {
		return nil, fmt.Errorf("invalid amount")
	}
	accountType := strings.ToUpper(strings.TrimSpace(req.AccountType))
	method, ok := country.FindHaiPayPayoutMethod(currency, accountType, req.BankCode)
	if !ok {
		return nil, fmt.Errorf("unsupported payout method currency=%s accountType=%s bankCode=%s", currency, accountType, strings.TrimSpace(req.BankCode))
	}
	name := strings.TrimSpace(req.Name)
	email := strings.ToLower(strings.TrimSpace(req.Email))
	phone := strings.TrimSpace(req.Phone)
	bankCode := method.BankCode
	accountNo := strings.TrimSpace(req.AccountNo)
	if name == "" || email == "" || phone == "" || bankCode == "" || accountNo == "" {
		return nil, fmt.Errorf("missing payout payee fields")
	}
	if !strings.Contains(name, " ") {
		name = name + " User"
	}
	notifyURL := callbackBaseURL + haiPayPayoutNotifyPath
	amountStr := haiPayFormatPayoutAmount(currency, normalizedAmount)

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
		"subject":       haiPayPayoutSubject,
	}
	if b := strings.TrimSpace(req.Body); b != "" {
		body["body"] = b
	}
	extraParams, err := haiPayPayoutExtraParams(req)
	if err != nil {
		return nil, err
	}
	for key, value := range extraParams {
		body[key] = value
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
	return haiPayQueryPayoutOrder(ctx, row.TransferCurrency, row.TransferOrderId, row.TransferPlatformNo, row.TransferLocalAmount)
}

func haiPayQueryAnchorPayout(ctx context.Context, row *liveentity.AnchorIncomeSettlementLog) (*haiPayPayoutQueryData, error) {
	if row == nil || row.ID == 0 {
		return nil, fmt.Errorf("empty anchor payout settlement")
	}
	return haiPayQueryPayoutOrder(ctx, row.TransferCurrency, row.TransferOrderId, row.TransferPlatformNo, row.TransferLocalAmount)
}

func haiPayQueryPayoutOrder(ctx context.Context, transferCurrency, transferOrderId, transferPlatformNo string, transferLocalAmount float64) (*haiPayPayoutQueryData, error) {
	cfg := cfgdao.GetHaiPayCfgCached()
	if cfg == nil || !cfgdao.HaiPayEnabled() {
		return nil, fmt.Errorf("haipay config not ready")
	}
	currency := strings.ToUpper(strings.TrimSpace(transferCurrency))
	orderID := strings.TrimSpace(transferOrderId)
	if currency == "" || orderID == "" {
		return nil, fmt.Errorf("haipay payout local order incomplete")
	}
	appID, err := HaiPayPayoutAppID(currency)
	if err != nil {
		return nil, err
	}
	body := map[string]any{
		"appId":   appID,
		"orderId": orderID,
	}
	if orderNo := strings.TrimSpace(transferPlatformNo); orderNo != "" {
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
	localOrderNo := strings.TrimSpace(transferPlatformNo)
	queryOrderNo := strings.TrimSpace(res.Data.OrderNo)
	if localOrderNo != "" && (queryOrderNo == "" || queryOrderNo != localOrderNo) {
		return nil, fmt.Errorf("haipay payout query platform order mismatch got=%s want=%s", queryOrderNo, localOrderNo)
	}
	if res.Data.Status == haiPayCollectStatusSuccess {
		queryAmount, amountErr := haiPayAmountCents(res.Data.Amount)
		if amountErr != nil {
			return nil, fmt.Errorf("invalid haipay payout query amount %q: %w", res.Data.Amount, amountErr)
		}
		expectedText := haiPayFormatPayoutAmount(currency, transferLocalAmount)
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
