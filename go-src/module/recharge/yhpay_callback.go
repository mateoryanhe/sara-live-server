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

type yhPayManualPayinCallbackBody struct {
	MerchantCode   string  `json:"MerchantCode"`
	InvoiceNo      string  `json:"InvoiceNo"`
	Currency       string  `json:"Currency"`
	Amount         string  `json:"Amount"`
	OriginalAmount float64 `json:"OriginalAmount"`
	PaymentDate    string  `json:"PaymentDate"`
	RefId          string  `json:"RefId"`
	ReceiverBank   string  `json:"ReceiverBank"`
	ReceiverAccount string `json:"ReceiverAccount"`
	Status         int     `json:"Status"`
	Signature      string  `json:"Signature"`
	StatusMessage  string  `json:"StatusMessage"`
	DeclineReason  string  `json:"DeclineReason"`
}

// HandleYhPayPayinCallback IDR 手动提交入款回调(JSON)
func HandleYhPayPayinCallback(r *ghttp.Request) {
	ctx := r.Context()
	cfg := cfgdao.GetYhPayCfgCached()
	if cfg == nil || !cfg.Enabled {
		writeYhPayAck(r, 4, "Operation Failed")
		return
	}

	var body yhPayManualPayinCallbackBody
	if err := r.Parse(&body); err != nil {
		xrlog.DetailLog.Warningf(ctx, "yhpay manual payin callback parse failed err=%v", err)
		writeYhPayAck(r, 3, "Invalid Parameter")
		return
	}

	xrlog.DetailLog.Infof(ctx, "yhpay manual payin callback recv refId=%s invoice=%s status=%d currency=%s amount=%s msg=%s reason=%s",
		body.RefId, body.InvoiceNo, body.Status, body.Currency, body.Amount, body.StatusMessage, body.DeclineReason)

	sig := strings.TrimSpace(body.Signature)
	expect := yhPayHMAC(cfg.ApiKey, body.InvoiceNo+body.ReceiverBank+body.ReceiverAccount+
		strings.ToUpper(strings.TrimSpace(body.Currency))+strings.TrimSpace(body.Amount)+body.PaymentDate)
	if !strings.EqualFold(expect, sig) {
		xrlog.DetailLog.Warningf(ctx, "yhpay manual payin bad signature refId=%s", body.RefId)
		writeYhPayAck(r, 3, "Invalid Parameter")
		return
	}

	// 3=批准 Approved, 2=拒绝 Rejected；仅最终状态会回调
	if body.Status != 3 {
		xrlog.DetailLog.Infof(ctx, "yhpay manual payin ignored status=%d refId=%s", body.Status, body.RefId)
		writeYhPayAck(r, 0, "Operation Success")
		return
	}

	if err := completeYhPayOrder(ctx, strings.TrimSpace(body.RefId), strings.TrimSpace(body.InvoiceNo)); err != nil {
		xrlog.DetailLog.Errorf(ctx, "yhpay manual payin complete failed refId=%s err=%v", body.RefId, err)
		writeYhPayAck(r, 4, "Operation Failed")
		return
	}
	xrlog.DetailLog.Infof(ctx, "yhpay manual payin complete ok refId=%s invoice=%s status=%d", body.RefId, body.InvoiceNo, body.Status)
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
