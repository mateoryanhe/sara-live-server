package recharge

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/xrlog"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dao/liveroomdao"
	liveentity "xr-game-server/entity/live"
)

// HandleHaiPayPayoutNotify 代付异步通知
func HandleHaiPayPayoutNotify(r *ghttp.Request) {
	ctx := r.Context()
	cfg := cfgdao.GetHaiPayCfgCached()
	if cfg == nil || !cfg.PayoutEnabled {
		r.Response.WriteStatus(503)
		r.Response.Write([]byte("not configured"))
		return
	}

	raw := r.GetBody()
	xrlog.DetailLog.Infof(ctx, "haipay payout notify raw body=%s", haiPayTruncateLog(string(raw), haiPayLogBodyMax))
	var params map[string]any
	if err := json.Unmarshal(raw, &params); err != nil {
		xrlog.DetailLog.Warningf(ctx, "haipay payout notify parse failed err=%v", err)
		r.Response.WriteStatus(400)
		r.Response.Write([]byte("bad request"))
		return
	}

	sign, _ := params["sign"].(string)
	if err := haiPayVerify(ctx, params, cfg.MerchantSecretKey, cfg.HaiPayPublicKey, sign); err != nil {
		xrlog.DetailLog.Warningf(ctx, "haipay payout notify verify failed err=%v body=%s", err, haiPaySafeJSON(raw))
		// 与代收一致:验签失败仍记日志;代付资金敏感,暂拒绝以免误标成功
		r.Response.WriteStatus(400)
		r.Response.Write([]byte("bad sign"))
		return
	}

	orderId, _ := haiPayValueToSignString(params["orderId"])
	orderNo, _ := haiPayValueToSignString(params["orderNo"])
	statusVal, _ := haiPayValueToSignString(params["status"])
	status, _ := strconv.Atoi(statusVal)
	errMsg, _ := haiPayValueToSignString(params["errorMsg"])

	xrlog.DetailLog.Infof(ctx, "haipay payout notify orderId=%s orderNo=%s status=%d errMsg=%s",
		orderId, orderNo, status, errMsg)

	row := liveroomdao.GetGuildIncomeSettlementLogByTransferOrderId(strings.TrimSpace(orderId))
	if row == nil {
		xrlog.DetailLog.Warningf(ctx, "haipay payout notify settlement not found orderId=%s", orderId)
		r.Response.WriteStatus(200)
		r.Response.Write([]byte("OK"))
		return
	}

	switch status {
	case 2: // 放款成功
		if row.Status == liveentity.GuildIncomeSettlementStatusTransferred {
			r.Response.WriteStatus(200)
			r.Response.Write([]byte("OK"))
			return
		}
		now := time.Now()
		if orderNo != "" {
			row.SetTransferPayout(row.TransferOrderId, orderNo, row.TransferCurrency, row.TransferLocalAmount)
		}
		row.SetStatus(liveentity.GuildIncomeSettlementStatusTransferred)
		row.SetTransferAt(&now)
		row.SetTransferFailMsg("")
	case 3: // 放款失败 → 回审核通过可重试
		msg := strings.TrimSpace(errMsg)
		if msg == "" {
			msg = "payout failed"
		}
		row.SetStatus(liveentity.GuildIncomeSettlementStatusApproved)
		row.SetTransferFailMsg(msg)
		// 保留 TransferOrderId 便于排查;重试会换新单号覆盖
	default:
		xrlog.DetailLog.Infof(ctx, "haipay payout notify ignored status=%d orderId=%s", status, orderId)
	}

	r.Response.WriteStatus(200)
	r.Response.Write([]byte("OK"))
}
