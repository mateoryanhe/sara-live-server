package recharge

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/xrlog"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dao/channelpaydao"
	"xr-game-server/dao/rechargeorderdao"
)

// HandleHaiPayCollectNotify HaiPay 本地代收异步通知。
// 发币不以回调验签为准：解析 orderId/orderNo 后按订单币种主动查单，
// 查单 status=2 成功则 CompleteChannelPayOrder；status=3 才标记失败，其余状态保持待支付。
func HandleHaiPayCollectNotify(r *ghttp.Request) {
	ctx := r.Context()
	if !cfgdao.HaiPayEnabled() {
		r.Response.WriteStatus(503)
		r.Response.Write([]byte("not configured"))
		return
	}

	raw := r.GetBody()
	xrlog.DetailLog.Infof(ctx, "haipay notify raw body=%s", haiPayTruncateLog(string(raw), haiPayLogBodyMax))
	var params map[string]any
	if err := json.Unmarshal(raw, &params); err != nil {
		xrlog.DetailLog.Warningf(ctx, "haipay notify parse failed err=%v body=%s", err, haiPayTruncateLog(string(raw), 1024))
		r.Response.WriteStatus(400)
		r.Response.Write([]byte("bad request"))
		return
	}

	orderId, _ := haiPayValueToSignString(params["orderId"])
	orderNo, _ := haiPayValueToSignString(params["orderNo"])
	orderId = strings.TrimSpace(orderId)
	orderNo = haiPayResolveNotifyOrderNo(orderId, orderNo)
	notifyStatus, _ := haiPayValueToSignString(params["status"])

	xrlog.DetailLog.Infof(ctx, "haipay notify recv orderId=%s orderNo=%s notifyStatus=%s amount=%v currency=%v",
		orderId, orderNo, notifyStatus, params["amount"], params["currency"])

	if orderId == "" {
		r.Response.WriteStatus(400)
		r.Response.Write([]byte("missing orderId"))
		return
	}

	platformStatus, platformOrderNo, err := haiPayQueryCollect(ctx, orderId, orderNo)
	if err != nil {
		xrlog.DetailLog.Infof(ctx, "haipay notify query failed orderId=%s orderNo=%s err=%v", orderId, orderNo, err)
		r.Response.WriteStatus(500)
		r.Response.Write([]byte("query fail"))
		return
	}
	if platformStatus == haiPayCollectStatusSuccess {
		if err := CompleteChannelPayOrder(ctx, orderId, platformOrderNo); err != nil {
			xrlog.DetailLog.Infof(ctx, "haipay notify complete failed orderId=%s err=%v", orderId, err)
			r.Response.WriteStatus(500)
			r.Response.Write([]byte("fail"))
			return
		}
		notifyName, _ := haiPayValueToSignString(params["name"])
		notifyEmail, _ := haiPayValueToSignString(params["email"])
		notifyPhone, _ := haiPayValueToSignString(params["phone"])
		if isLikelyRealPayName(notifyName) || isLikelyRealPayEmail(notifyEmail) || strings.TrimSpace(notifyPhone) != "" {
			if oid, e := strconv.ParseUint(orderId, 10, 64); e == nil && oid > 0 {
				if order := rechargeorderdao.GetById(oid); order != nil && order.UserId > 0 {
					if requireExistingAppUser(order.UserId) == nil {
						channelpaydao.UpsertPayerInfo(order.UserId, notifyName, notifyEmail, notifyPhone)
					}
				}
			}
		}
		r.Response.WriteStatus(200)
		r.Response.Write([]byte("OK"))
		return
	}

	if platformStatus == haiPayCollectStatusFailed {
		xrlog.DetailLog.Infof(ctx, "haipay notify failed orderId=%s orderNo=%s platformStatus=%d → mark failed",
			orderId, platformOrderNo, platformStatus)
		if err := FailChannelPayOrder(ctx, orderId, platformOrderNo, platformStatus); err != nil {
			xrlog.DetailLog.Infof(ctx, "haipay notify fail-mark failed orderId=%s err=%v", orderId, err)
			r.Response.WriteStatus(500)
			r.Response.Write([]byte("fail"))
			return
		}
	} else {
		xrlog.DetailLog.Infof(ctx, "haipay notify pending orderId=%s orderNo=%s platformStatus=%d", orderId, platformOrderNo, platformStatus)
	}
	r.Response.WriteStatus(200)
	r.Response.Write([]byte("OK"))
}
