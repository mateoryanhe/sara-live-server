package controller

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/httpserver"
	"xr-game-server/module/recharge"
)

const YhPayWebhookUrl = "/webhook/yhpay"

func initYhPayWebhookController() {
	httpserver.RegNonAuthHandler(YhPayWebhookUrl, "/payin", handleYhPayPayinHTTP)
}

func handleYhPayPayinHTTP(r *ghttp.Request) {
	recharge.HandleYhPayPayinCallback(r)
}
