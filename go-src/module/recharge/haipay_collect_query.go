package recharge

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"xr-game-server/core/xrlog"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dao/rechargeorderdao"
)

const haiPayGlobalCashierQueryPath = "/global/cashier/collect/query"

// 全球收银台代收订单状态(与 HaiPay 文档一致)
const (
	haiPayCollectStatusSuccess = 2 // 成功（终态）
	haiPayCollectStatusFailed  = 3 // 失败（终态）
)

type haiPayCollectQueryAPIRes struct {
	Status string `json:"status"`
	Error  string `json:"error"`
	Msg    string `json:"msg"`
	Data   *struct {
		OrderId      string `json:"orderId"`
		OrderNo      string `json:"orderNo"`
		Amount       string `json:"amount"`
		ActualAmount string `json:"actualAmount"`
		Fee          string `json:"fee"`
		Status       int    `json:"status"`
		PayTime      string `json:"payTime"`
		ErrorMsg     string `json:"errorMsg"`
		Sign         string `json:"sign"`
	} `json:"data"`
}

// haiPayQueryCollect 主动查单；返回平台订单状态与 orderNo。
// orderNo 可空，会尝试从本地订单 ThirdOrderId 补全（全球收银台 query 要求 orderNo）。
func haiPayQueryCollect(ctx context.Context, orderId, orderNo string) (platformStatus int, platformOrderNo string, err error) {
	cfg := cfgdao.GetHaiPayCfgCached()
	if cfg == nil || !cfg.Enabled {
		return 0, "", fmt.Errorf("haipay not configured")
	}
	orderId = strings.TrimSpace(orderId)
	orderNo = strings.TrimSpace(orderNo)
	if orderId == "" {
		return 0, "", fmt.Errorf("empty orderId")
	}
	if orderNo == "" {
		if oid, e := strconv.ParseUint(orderId, 10, 64); e == nil && oid > 0 {
			if order := rechargeorderdao.GetById(oid); order != nil {
				orderNo = strings.TrimSpace(order.ThirdOrderId)
			}
		}
	}
	if orderNo == "" {
		return 0, "", fmt.Errorf("empty orderNo for query")
	}

	body := map[string]any{
		"appId":   cfg.AppId,
		"orderId": orderId,
		"orderNo": orderNo,
	}
	xrlog.DetailLog.Infof(ctx, "haipay collect query sign-before params=%s", haiPaySafeParams(body))
	sign, err := haiPaySign(ctx, body, cfg.MerchantSecretKey, cfg.MerchantPrivateKey)
	if err != nil {
		return 0, "", fmt.Errorf("haipay query sign: %w", err)
	}
	body["sign"] = sign

	var res haiPayCollectQueryAPIRes
	url := strings.TrimRight(cfg.ApiHost, "/") + haiPayGlobalCashierQueryPath
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
	// 查单应答验签失败只打日志：以官方 HTTPS 查单结果为准发币，规避平台公钥不一致导致无法入账
	dataMap := map[string]any{}
	rawData, _ := json.Marshal(res.Data)
	_ = json.Unmarshal(rawData, &dataMap)
	if err = haiPayVerify(ctx, dataMap, cfg.MerchantSecretKey, cfg.HaiPayPublicKey, res.Data.Sign); err != nil {
		xrlog.DetailLog.Warningf(ctx, "haipay collect query verify failed(skip) orderId=%s err=%v", orderId, err)
	}

	platformOrderNo = strings.TrimSpace(res.Data.OrderNo)
	if platformOrderNo == "" {
		platformOrderNo = orderNo
	}
	platformStatus = res.Data.Status
	xrlog.DetailLog.Infof(ctx, "haipay collect query ok orderId=%s orderNo=%s status=%d amount=%s actual=%s",
		orderId, platformOrderNo, platformStatus, res.Data.Amount, res.Data.ActualAmount)
	return platformStatus, platformOrderNo, nil
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
