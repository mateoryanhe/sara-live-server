package controller

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/httpserver"
	"xr-game-server/module/recharge"
)

const HaiPayWebhookUrl = "/webhook/haipay"

func initHaiPayWebhookController() {
	httpserver.RegNonAuthHandler(HaiPayWebhookUrl, "/collect/notify", handleHaiPayCollectNotifyHTTP)
	httpserver.RegNonAuthHandler(HaiPayWebhookUrl, "/payout/notify", handleHaiPayPayoutNotifyHTTP)
}

func handleHaiPayCollectNotifyHTTP(r *ghttp.Request) {
	recharge.HandleHaiPayCollectNotify(r)
}

func handleHaiPayPayoutNotifyHTTP(r *ghttp.Request) {
	recharge.HandleHaiPayPayoutNotify(r)
}
