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
	if !cfgdao.HaiPayEnabled() {
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

	orderId, _ := haiPayValueToSignString(params["orderId"])
	notifyOrderNo, _ := haiPayValueToSignString(params["orderNo"])
	notifyStatus, _ := haiPayValueToSignString(params["status"])
	xrlog.DetailLog.Infof(ctx, "haipay payout notify trigger orderId=%s orderNo=%s notifyStatus=%s",
		orderId, notifyOrderNo, notifyStatus)

	settlementID, err := strconv.ParseUint(strings.TrimSpace(orderId), 10, 64)
	if err != nil || settlementID == 0 {
		xrlog.DetailLog.Warningf(ctx, "haipay payout notify invalid settlement orderId=%s err=%v", orderId, err)
		r.Response.WriteStatus(200)
		r.Response.Write([]byte("OK"))
		return
	}
	row := liveroomdao.GetGuildIncomeSettlementLogById(settlementID)
	if row == nil {
		xrlog.DetailLog.Warningf(ctx, "haipay payout notify settlement not found orderId=%s", orderId)
		r.Response.WriteStatus(200)
		r.Response.Write([]byte("OK"))
		return
	}
	if row.Status == liveentity.GuildIncomeSettlementStatusTransferred {
		r.Response.WriteStatus(200)
		r.Response.Write([]byte("OK"))
		return
	}
	query, err := haiPayQueryPayout(ctx, row)
	if err != nil {
		xrlog.DetailLog.Warningf(ctx, "haipay payout notify query failed orderId=%s err=%v", orderId, err)
		r.Response.WriteStatus(500)
		r.Response.Write([]byte("query fail"))
		return
	}

	switch query.Status {
	case 2: // 放款成功
		now := time.Now()
		if orderNo := strings.TrimSpace(query.OrderNo); orderNo != "" {
			row.SetTransferPayout(row.TransferOrderId, orderNo, row.TransferCurrency, row.TransferLocalAmount)
		}
		row.SetStatus(liveentity.GuildIncomeSettlementStatusTransferred)
		row.SetTransferAt(&now)
		row.SetTransferFailMsg("")
		liveroomdao.PublishGuildIncomeSettlementLog(row)
	case 3: // 放款失败 → 回审核通过可重试
		msg := strings.TrimSpace(query.ErrorMsg)
		if msg == "" {
			msg = "payout failed"
		}
		row.SetStatus(liveentity.GuildIncomeSettlementStatusApproved)
		row.SetTransferFailMsg(msg)
		liveroomdao.PublishGuildIncomeSettlementLog(row)
	default:
		xrlog.DetailLog.Infof(ctx, "haipay payout notify query pending status=%d orderId=%s", query.Status, orderId)
	}

	r.Response.WriteStatus(200)
	r.Response.Write([]byte("OK"))
}
