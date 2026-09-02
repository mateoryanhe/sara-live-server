package controller

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/httpserver"
	"xr-game-server/module/recharge"
)

const YhPayWebhookUrl = "/webhook/yhpay"

func initYhPayWebhookController() {
	httpserver.RegNonAuthHandler(YhPayWebhookUrl, "/payin", handleYhPayPayinHTTP)
	httpserver.RegNonAuthHandler(YhPayWebhookUrl, "/crypto", handleYhPayCryptoHTTP)
}

func handleYhPayPayinHTTP(r *ghttp.Request) {
	recharge.HandleYhPayPayinCallback(r)
}

func handleYhPayCryptoHTTP(r *ghttp.Request) {
	recharge.HandleYhPayCryptoCallback(r)
}
