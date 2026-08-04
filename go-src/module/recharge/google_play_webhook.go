package recharge

import (
	"io"
	"net/http"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/httpserver"
)

const googlePlayWebhookPrefix = "/webhook/googlePlay"

func initGooglePlayWebhook() {
	httpserver.RegNonAuthHandler(googlePlayWebhookPrefix, "/rtdn", handleGooglePlayRTDNHTTP)
}

func handleGooglePlayRTDNHTTP(r *ghttp.Request) {
	ctx := gctx.New()
	body, err := io.ReadAll(r.Request.Body)
	if err != nil {
		logGooglePlayError(ctx, "rtdn read body failed err=%v", err)
		r.Response.WriteStatus(http.StatusBadRequest)
		return
	}
	err = HandleGooglePlayRTDN(ctx, body, r.GetHeader("Authorization"))
	if err != nil {
		if isGooglePlayUnauthorized(err) {
			r.Response.WriteStatus(http.StatusUnauthorized)
			return
		}
		logGooglePlayError(ctx, "rtdn handle failed err=%v", err)
		r.Response.WriteStatus(http.StatusInternalServerError)
		return
	}
	r.Response.WriteStatus(http.StatusOK)
}
