package controller

import (
	"io"
	"net/http"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/httpserver"
	"xr-game-server/module/recharge"
)

const (
	// GooglePlayWebhookUrl Sara / 旧 App RTDN
	GooglePlayWebhookUrl = "/webhook/googlePlay"
	// SofieGooglePlayWebhookUrl Sofie 专用 RTDN(Pub/Sub Push 与 Sara 分开填)
	SofieGooglePlayWebhookUrl = "/webhook/sofie/googlePlay"
)

func initGooglePlayWebhookController() {
	httpserver.RegNonAuthHandler(GooglePlayWebhookUrl, "/rtdn", handleGooglePlayRTDNHTTP)
	httpserver.RegNonAuthHandler(SofieGooglePlayWebhookUrl, "/rtdn", handleGooglePlayRTDNHTTP)
}

func handleGooglePlayRTDNHTTP(r *ghttp.Request) {
	ctx := gctx.New()
	body, err := io.ReadAll(r.Request.Body)
	if err != nil {
		g.Log().Errorf(ctx, "google play rtdn read body failed path=%s err=%v", r.URL.Path, err)
		r.Response.WriteStatus(http.StatusBadRequest)
		return
	}
	if err = recharge.HandleGooglePlayRTDN(ctx, body); err != nil {
		g.Log().Errorf(ctx, "google play rtdn handle failed path=%s err=%v", r.URL.Path, err)
		r.Response.WriteStatus(http.StatusInternalServerError)
		return
	}
	r.Response.WriteStatus(http.StatusOK)
}
