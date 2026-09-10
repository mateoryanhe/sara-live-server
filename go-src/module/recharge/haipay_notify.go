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

// HandleHaiPayCollectNotify 全球收银台代收异步通知（Map 解析以兼容字段扩展）
func HandleHaiPayCollectNotify(r *ghttp.Request) {
	ctx := r.Context()
	cfg := cfgdao.GetHaiPayCfgCached()
	if cfg == nil || !cfg.Enabled {
		r.Response.WriteStatus(503)
		r.Response.Write([]byte("not configured"))
		return
	}

	raw := r.GetBody()
	var params map[string]any
	if err := json.Unmarshal(raw, &params); err != nil {
		xrlog.DetailLog.Warningf(ctx, "haipay notify parse failed err=%v body=%s", err, haiPayTruncateLog(string(raw), 1024))
		r.Response.WriteStatus(400)
		r.Response.Write([]byte("bad request"))
		return
	}

	sign, _ := params["sign"].(string)
	if err := haiPayVerify(params, cfg.MerchantSecretKey, cfg.HaiPayPublicKey, sign); err != nil {
		xrlog.DetailLog.Warningf(ctx, "haipay notify verify failed err=%v", err)
		r.Response.WriteStatus(400)
		r.Response.Write([]byte("bad sign"))
		return
	}

	orderId, _ := haiPayValueToSignString(params["orderId"])
	orderNo, _ := haiPayValueToSignString(params["orderNo"])
	statusVal, _ := haiPayValueToSignString(params["status"])
	status, _ := strconv.Atoi(statusVal)

	xrlog.DetailLog.Infof(ctx, "haipay notify recv orderId=%s orderNo=%s status=%d amount=%v currency=%v",
		orderId, orderNo, status, params["amount"], params["currency"])

	// 2=成功；4/5 部分/超额也视为已收款(美金包装下按成功发币)
	if status != 2 && status != 4 && status != 5 {
		xrlog.DetailLog.Infof(ctx, "haipay notify ignored status=%d orderId=%s", status, orderId)
		r.Response.WriteStatus(200)
		r.Response.Write([]byte("OK"))
		return
	}

	if err := CompleteChannelPayOrder(ctx, strings.TrimSpace(orderId), strings.TrimSpace(orderNo)); err != nil {
		xrlog.DetailLog.Errorf(ctx, "haipay notify complete failed orderId=%s err=%v", orderId, err)
		r.Response.WriteStatus(500)
		r.Response.Write([]byte("fail"))
		return
	}
	notifyName, _ := haiPayValueToSignString(params["name"])
	notifyEmail, _ := haiPayValueToSignString(params["email"])
	if isLikelyRealPayName(notifyName) || isLikelyRealPayEmail(notifyEmail) {
		if oid, e := strconv.ParseUint(strings.TrimSpace(orderId), 10, 64); e == nil && oid > 0 {
			if order := rechargeorderdao.GetById(oid); order != nil && order.UserId > 0 {
				if requireExistingAppUser(order.UserId) == nil {
					channelpaydao.UpsertPayerInfo(order.UserId, notifyName, notifyEmail)
				}
			}
		}
	}
	r.Response.WriteStatus(200)
	r.Response.Write([]byte("OK"))
}
