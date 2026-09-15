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

// HandleHaiPayCollectNotify 全球收银台代收异步通知。
// 发币不以回调验签为准：解析 orderId/orderNo 后主动 /global/cashier/collect/query，
// 查单 status=2 成功则 CompleteChannelPayOrder；非成功则将待支付订单标为失败(StatusFailed)。
func HandleHaiPayCollectNotify(r *ghttp.Request) {
	ctx := r.Context()
	cfg := cfgdao.GetHaiPayCfgCached()
	if cfg == nil || !cfg.Enabled {
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

	// 回调验签仅打日志，不作为发币条件（平台公钥不一致时仍可靠查单入账）
	if sign, _ := params["sign"].(string); strings.TrimSpace(sign) != "" {
		if err := haiPayVerify(ctx, params, cfg.MerchantSecretKey, cfg.HaiPayPublicKey, sign); err != nil {
			xrlog.DetailLog.Warningf(ctx, "haipay notify verify failed(skip, use query) err=%v", err)
		}
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
		if isLikelyRealPayName(notifyName) || isLikelyRealPayEmail(notifyEmail) {
			if oid, e := strconv.ParseUint(orderId, 10, 64); e == nil && oid > 0 {
				if order := rechargeorderdao.GetById(oid); order != nil && order.UserId > 0 {
					if requireExistingAppUser(order.UserId) == nil {
						channelpaydao.UpsertPayerInfo(order.UserId, notifyName, notifyEmail)
					}
				}
			}
		}
		r.Response.WriteStatus(200)
		r.Response.Write([]byte("OK"))
		return
	}

	// 非成功(含失败终态 status=3)：待支付单标为失败，避免 CMS 一直显示待支付
	xrlog.DetailLog.Infof(ctx, "haipay notify not success orderId=%s orderNo=%s platformStatus=%d → mark failed",
		orderId, platformOrderNo, platformStatus)
	if err := FailChannelPayOrder(ctx, orderId, platformOrderNo, platformStatus); err != nil {
		xrlog.DetailLog.Infof(ctx, "haipay notify fail-mark failed orderId=%s err=%v", orderId, err)
		r.Response.WriteStatus(500)
		r.Response.Write([]byte("fail"))
		return
	}
	r.Response.WriteStatus(200)
	r.Response.Write([]byte("OK"))
}
