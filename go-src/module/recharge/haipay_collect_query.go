package recharge

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"xr-game-server/constants/country"
	"xr-game-server/core/xrlog"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dao/rechargeorderdao"
	rechargeentity "xr-game-server/entity/recharge"
)

// HaiPay 代收订单状态（本地代收与全球收银台一致）
const (
	haiPayCollectStatusSuccess = 2 // 成功（终态）
	haiPayCollectStatusFailed  = 3 // 失败（终态）
)

type haiPayCollectQueryData struct {
	OrderId           string `json:"orderId"`
	OrderNo           string `json:"orderNo"`
	Currency          string `json:"currency"`
	Amount            string `json:"amount"`
	ActualAmount      string `json:"actualAmount"`
	Fee               string `json:"fee"`
	Status            int    `json:"status"`
	PayTime           string `json:"payTime"`
	ErrorMsg          string `json:"errorMsg"`
	OriginalCurrency  string `json:"originalCurrency"`
	OriginalAmount    string `json:"originalAmount"`
	FloatExchangeRate string `json:"floatExchangeRate"`
}

type haiPayCollectQueryAPIRes struct {
	Status string                  `json:"status"`
	Error  string                  `json:"error"`
	Msg    string                  `json:"msg"`
	Data   *haiPayCollectQueryData `json:"data"`
}

// haiPayQueryCollect 主动查单；返回平台订单状态与 orderNo。
// orderNo 会优先从本地订单 ThirdOrderId 补全；全球收银台 query v2 必须携带 orderNo。
func haiPayQueryCollect(ctx context.Context, orderId, orderNo string) (platformStatus int, platformOrderNo string, err error) {
	cfg := cfgdao.GetHaiPayCfgCached()
	if cfg == nil || !cfgdao.HaiPayEnabled() {
		return 0, "", fmt.Errorf("haipay not configured")
	}
	orderId = strings.TrimSpace(orderId)
	orderNo = strings.TrimSpace(orderNo)
	if orderId == "" {
		return 0, "", fmt.Errorf("empty orderId")
	}
	var localOrderPayAmount float64
	var localOrderCurrency string
	var localOrderPayRegion string
	var localOrderPayChannel uint8
	hasLocalOrder := false
	if oid, e := strconv.ParseUint(orderId, 10, 64); e == nil && oid > 0 {
		if order := rechargeorderdao.GetById(oid); order != nil {
			localOrderPayAmount = order.PayAmount
			localOrderCurrency = strings.ToUpper(strings.TrimSpace(order.Currency))
			localOrderPayRegion = strings.ToUpper(strings.TrimSpace(order.PayRegion))
			localOrderPayChannel = order.PayChannel
			hasLocalOrder = true
			if orderNo == "" {
				orderNo = strings.TrimSpace(order.ThirdOrderId)
			}
		}
	}
	if !hasLocalOrder || localOrderCurrency == "" {
		return 0, "", fmt.Errorf("haipay collect local order not found or currency missing orderId=%s", orderId)
	}
	globalCashier := isHaiPayGlobalCashierBiz(haiPayCollectionBizType(localOrderPayChannel))
	var appID int64
	if globalCashier {
		if orderNo == "" {
			return 0, "", fmt.Errorf("haipay global cashier query orderNo missing orderId=%s", orderId)
		}
		appID = int64(country.HaiPayAppIDCashier)
	} else {
		if localOrderPayChannel != rechargeentity.RechargeCfgTypeCoinMerchant {
			return 0, "", fmt.Errorf("haipay unsupported local collection channel=%d", localOrderPayChannel)
		}
		if localOrderPayRegion == "" {
			return 0, "", fmt.Errorf("haipay coin merchant collection region missing orderId=%s", orderId)
		}
		if _, credentialErr := haiPayCoinMerchantCollectionCredential(localOrderPayRegion, localOrderCurrency); credentialErr != nil {
			return 0, "", credentialErr
		}
		var ok bool
		appID, ok = country.LookupHaiPayAppID(localOrderCurrency)
		if !ok {
			return 0, "", fmt.Errorf("haipay appId enum missing currency=%s", localOrderCurrency)
		}
	}

	body := map[string]any{
		"appId":   appID,
		"orderId": orderId,
	}
	if orderNo != "" {
		body["orderNo"] = orderNo
	}
	xrlog.DetailLog.Infof(ctx, "haipay collect query sign-before params=%s", haiPaySafeParams(body))
	sign, err := haiPaySign(ctx, body, cfg.MerchantSecretKey, cfg.MerchantPrivateKey)
	if err != nil {
		return 0, "", fmt.Errorf("haipay query sign: %w", err)
	}
	body["sign"] = sign

	var res haiPayCollectQueryAPIRes
	queryPath := haiPayLocalCollectPath(localOrderCurrency, "query")
	if globalCashier {
		queryPath = haiPayGlobalCollectQueryV2Path
	}
	url := strings.TrimRight(cfg.ApiHost, "/") + queryPath
	if err = haiPayPostJSON(ctx, url, body, &res); err != nil {
		return 0, "", err
	}
	if res.Status != "1" {
		msg := strings.TrimSpace(res.Msg)
		if msg == "" {
			msg = res.Error
		}
		return 0, "", fmt.Errorf("haipay query failed status=%s error=%s msg=%s", res.Status, res.Error, msg)
	}
	if res.Data == nil {
		return 0, "", fmt.Errorf("haipay query empty data")
	}
	platformOrderNo = strings.TrimSpace(res.Data.OrderNo)
	if platformOrderNo == "" {
		platformOrderNo = orderNo
	}
	platformStatus = res.Data.Status
	if platformStatus == haiPayCollectStatusSuccess {
		// 本地代收查单响应不返回 currency；请求 URL 已绑定本地订单币种。
		if !globalCashier && strings.TrimSpace(res.Data.Currency) == "" {
			res.Data.Currency = localOrderCurrency
		}
		amountField, checkedAmount, checkErr := haiPayCheckCollectPayAmount(res.Data, localOrderCurrency, localOrderPayAmount)
		if checkErr != nil {
			xrlog.DetailLog.Warningf(ctx, "haipay collect amount check failed orderId=%s orderNo=%s queryCurrency=%s queryAmount=%s originalCurrency=%s originalAmount=%s orderCurrency=%s orderPayAmount=%v err=%v",
				orderId, platformOrderNo, res.Data.Currency, res.Data.Amount, res.Data.OriginalCurrency, res.Data.OriginalAmount,
				localOrderCurrency, localOrderPayAmount, checkErr)
			return 0, "", checkErr
		}
		xrlog.DetailLog.Infof(ctx, "haipay collect amount check ok orderId=%s field=%s amount=%s orderCurrency=%s orderPayAmount=%v",
			orderId, amountField, checkedAmount, localOrderCurrency, localOrderPayAmount)
	}
	xrlog.DetailLog.Infof(ctx, "haipay collect query ok orderId=%s orderNo=%s status=%d currency=%s amount=%s actual=%s originalCurrency=%s originalAmount=%s",
		orderId, platformOrderNo, platformStatus, res.Data.Currency, res.Data.Amount, res.Data.ActualAmount,
		res.Data.OriginalCurrency, res.Data.OriginalAmount)
	return platformStatus, platformOrderNo, nil
}

// haiPayCheckCollectPayAmount 选择与本地订单币种匹配的查单金额，再按最小单位（分）比较。
// 本地代收使用 amount；全球收银台按订单币种选择 originalAmount 或 amount。
func haiPayCheckCollectPayAmount(data *haiPayCollectQueryData, orderCurrency string, orderPayAmount float64) (field, queryAmount string, err error) {
	if data == nil {
		return "", "", fmt.Errorf("empty haipay query data")
	}
	orderCurrency = strings.ToUpper(strings.TrimSpace(orderCurrency))
	if orderCurrency == "" {
		return "", "", fmt.Errorf("empty order currency")
	}
	queryCurrency := strings.ToUpper(strings.TrimSpace(data.Currency))
	originalCurrency := strings.ToUpper(strings.TrimSpace(data.OriginalCurrency))
	switch {
	case originalCurrency == orderCurrency && strings.TrimSpace(data.OriginalAmount) != "":
		field = "originalAmount"
		queryAmount = data.OriginalAmount
	case queryCurrency == orderCurrency && strings.TrimSpace(data.Amount) != "":
		field = "amount"
		queryAmount = data.Amount
	default:
		return "", "", fmt.Errorf("haipay query has no amount for order currency=%s", orderCurrency)
	}

	queryCents, err := haiPayAmountCents(queryAmount)
	if err != nil {
		return field, queryAmount, fmt.Errorf("invalid haipay query %s %q: %w", field, queryAmount, err)
	}
	orderAmount := strconv.FormatFloat(orderPayAmount, 'f', 4, 64)
	orderCents, err := haiPayAmountCents(orderAmount)
	if err != nil {
		return field, queryAmount, fmt.Errorf("invalid order pay_amount %s: %w", orderAmount, err)
	}
	if queryCents != orderCents {
		return field, queryAmount, fmt.Errorf("haipay collect amount mismatch %s=%s order.pay_amount=%s", field, queryAmount, orderAmount)
	}
	return field, queryAmount, nil
}

func haiPayAmountCents(amount string) (int64, error) {
	amount = strings.TrimSpace(amount)
	if amount == "" {
		return 0, fmt.Errorf("empty amount")
	}
	value, ok := new(big.Rat).SetString(amount)
	if !ok || value.Sign() <= 0 {
		return 0, fmt.Errorf("amount must be a positive decimal")
	}
	value.Mul(value, big.NewRat(100, 1))
	if !value.IsInt() {
		return 0, fmt.Errorf("amount has more than two decimal places")
	}
	if !value.Num().IsInt64() {
		return 0, fmt.Errorf("amount exceeds supported range")
	}
	return value.Num().Int64(), nil
}

// haiPayResolveNotifyOrderNo 回调缺 orderNo 时从本地单补
func haiPayResolveNotifyOrderNo(orderId, orderNo string) string {
	orderNo = strings.TrimSpace(orderNo)
	if orderNo != "" {
		return orderNo
	}
	oid, err := strconv.ParseUint(strings.TrimSpace(orderId), 10, 64)
	if err != nil || oid == 0 {
		return ""
	}
	order := rechargeorderdao.GetById(oid)
	if order == nil {
		return ""
	}
	return strings.TrimSpace(order.ThirdOrderId)
}
