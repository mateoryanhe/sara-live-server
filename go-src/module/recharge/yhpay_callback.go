package recharge

import (
	"context"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/constants/currency"
	"xr-game-server/core/xrlog"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dao/rechargeorderdao"
	"xr-game-server/entity/recharge"
	"xr-game-server/errercode"
)

type yhPayCallbackAck struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

func writeYhPayAck(r *ghttp.Request, code int, message string) {
	r.Response.WriteJson(&yhPayCallbackAck{ErrorCode: code, Message: message})
}

// HandleYhPayPayinCallback DepositV2 入款回调(form-urlencoded)
func HandleYhPayPayinCallback(r *ghttp.Request) {
	ctx := r.Context()
	cfg := cfgdao.GetYhPayCfgCached()
	if cfg == nil || !cfg.Enabled {
		writeYhPayAck(r, 4, "Operation Failed")
		return
	}

	itemID := strings.TrimSpace(r.GetForm("ItemID").String())
	statusStr := strings.TrimSpace(r.GetForm("status").String())
	amount := strings.TrimSpace(r.GetForm("Amount").String())
	currencyCode := strings.TrimSpace(r.GetForm("Currency").String())
	transaction := strings.TrimSpace(r.GetForm("transaction").String())
	createdAt := strings.TrimSpace(r.GetForm("created_at").String())
	signature2 := strings.TrimSpace(r.GetForm("signature2").String())

	xrlog.DetailLog.Infof(ctx, "yhpay payin callback recv itemID=%s status=%s amount=%s currency=%s transaction=%s created_at=%s",
		itemID, statusStr, amount, currencyCode, transaction, createdAt)

	expect := yhPayHMAC(cfg.ApiKey, transaction+statusStr+currencyCode+amount+createdAt)
	if !strings.EqualFold(expect, signature2) {
		xrlog.DetailLog.Warningf(ctx, "yhpay payin bad signature itemID=%s", itemID)
		writeYhPayAck(r, 3, "Invalid Parameter")
		return
	}

	status, _ := strconv.Atoi(statusStr)
	if status != 1 && status != 3 {
		xrlog.DetailLog.Infof(ctx, "yhpay payin ignored status=%d itemID=%s", status, itemID)
		writeYhPayAck(r, 0, "Operation Success")
		return
	}

	if err := completeYhPayOrder(ctx, itemID, transaction); err != nil {
		xrlog.DetailLog.Errorf(ctx, "yhpay payin complete failed itemID=%s err=%v", itemID, err)
		writeYhPayAck(r, 4, "Operation Failed")
		return
	}
	xrlog.DetailLog.Infof(ctx, "yhpay payin complete ok itemID=%s transaction=%s status=%d", itemID, transaction, status)
	writeYhPayAck(r, 0, "Operation Success")
}

type yhPayCryptoCallbackBody struct {
	InvoiceNo           string  `json:"InvoiceNo"`
	MerchantCode        string  `json:"MerchantCode"`
	SenderWalletAddress string  `json:"SenderWalletAddress"`
	Reference           string  `json:"Reference"`
	ExchangeRate        float64 `json:"ExchangeRate"`
	FiatCurrency        string  `json:"FiatCurrency"`
	CryptoCurrency      string  `json:"CryptoCurrency"`
	FiatAmount          float64 `json:"FiatAmount"`
	CryptoAmount        float64 `json:"CryptoAmount"`
	Status              int     `json:"Status"`
	RefId               string  `json:"RefId"`
	PaymentDate         string  `json:"PaymentDate"`
	Signature           string  `json:"Signature"`
	StatusMessage       string  `json:"StatusMessage"`
	DeclineReason       string  `json:"DeclineReason"`
}

// HandleYhPayCryptoCallback 加密货币入款回调(JSON)
func HandleYhPayCryptoCallback(r *ghttp.Request) {
	ctx := r.Context()
	cfg := cfgdao.GetYhPayCfgCached()
	if cfg == nil || !cfg.Enabled {
		writeYhPayAck(r, 4, "Operation Failed")
		return
	}

	var body yhPayCryptoCallbackBody
	if err := r.Parse(&body); err != nil {
		xrlog.DetailLog.Warningf(ctx, "yhpay crypto callback parse failed err=%v", err)
		writeYhPayAck(r, 3, "Invalid Parameter")
		return
	}
	xrlog.DetailLog.Infof(ctx, "yhpay crypto callback recv refId=%s invoice=%s status=%d fiat=%s/%.2f crypto=%s/%.8f msg=%s",
		body.RefId, body.InvoiceNo, body.Status, body.FiatCurrency, body.FiatAmount, body.CryptoCurrency, body.CryptoAmount, body.StatusMessage)

	sig := strings.TrimSpace(body.Signature)
	payloadA := body.InvoiceNo + body.SenderWalletAddress + body.FiatCurrency +
		formatYhPayAmount(body.FiatAmount) + body.CryptoCurrency +
		strconv.FormatFloat(body.CryptoAmount, 'f', -1, 64) + body.PaymentDate
	payloadB := body.InvoiceNo + body.SenderWalletAddress + body.FiatCurrency +
		strconv.FormatFloat(body.FiatAmount, 'f', -1, 64) + body.CryptoCurrency +
		strconv.FormatFloat(body.CryptoAmount, 'f', -1, 64) + body.PaymentDate
	if !strings.EqualFold(yhPayHMAC(cfg.ApiKey, payloadA), sig) &&
		!strings.EqualFold(yhPayHMAC(cfg.ApiKey, payloadB), sig) {
		xrlog.DetailLog.Warningf(ctx, "yhpay crypto bad signature refId=%s", body.RefId)
		writeYhPayAck(r, 3, "Invalid Parameter")
		return
	}

	if body.Status != 1 && body.Status != 3 {
		xrlog.DetailLog.Infof(ctx, "yhpay crypto ignored status=%d refId=%s", body.Status, body.RefId)
		writeYhPayAck(r, 0, "Operation Success")
		return
	}

	if err := completeYhPayOrder(ctx, strings.TrimSpace(body.RefId), strings.TrimSpace(body.InvoiceNo)); err != nil {
		xrlog.DetailLog.Errorf(ctx, "yhpay crypto complete failed refId=%s err=%v", body.RefId, err)
		writeYhPayAck(r, 4, "Operation Failed")
		return
	}
	xrlog.DetailLog.Infof(ctx, "yhpay crypto complete ok refId=%s invoice=%s status=%d", body.RefId, body.InvoiceNo, body.Status)
	writeYhPayAck(r, 0, "Operation Success")
}

func completeYhPayOrder(ctx context.Context, orderIdStr, thirdOrderId string) error {
	orderId, err := strconv.ParseUint(strings.TrimSpace(orderIdStr), 10, 64)
	if err != nil || orderId == 0 {
		return errercode.CreateCode(errercode.RechargeOrderNonExist)
	}
	thirdOrderId = strings.TrimSpace(thirdOrderId)
	if thirdOrderId != "" {
		gmlock.Lock(rechargeThirdOrderLockKey(thirdOrderId))
		defer gmlock.Unlock(rechargeThirdOrderLockKey(thirdOrderId))
		if existing := rechargeorderdao.GetByThirdOrderId(thirdOrderId); existing != nil {
			if existing.ID == orderId && existing.Status == entity.RechargeOrderStatusCompleted {
				return nil
			}
			if existing.ID != orderId {
				xrlog.DetailLog.Warningf(ctx, "yhpay thirdOrder reused third=%s orderId=%d existing=%d", thirdOrderId, orderId, existing.ID)
				return errercode.CreateCode(errercode.RechargeOrderStateInvalid)
			}
		}
	}

	order := rechargeorderdao.GetById(orderId)
	if order == nil {
		return errercode.CreateCode(errercode.RechargeOrderNonExist)
	}
	if order.Status == entity.RechargeOrderStatusCompleted {
		return nil
	}
	if order.PayChannel != entity.RechargeCfgTypeChannel {
		return errercode.CreateCode(errercode.RechargeOrderStateInvalid)
	}
	if thirdOrderId != "" && order.ThirdOrderId != thirdOrderId {
		order.SetThirdOrderId(thirdOrderId)
	}
	_, err = completeOrder(order, currency.ReasonRecharge)
	return err
}
