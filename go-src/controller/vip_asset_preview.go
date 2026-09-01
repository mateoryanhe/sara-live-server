package controller

import (
	"net/http"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/httpserver"
	"xr-game-server/module/vip"
)

const (
	vipAssetPreviewPath         = "/vipCfg/assetPreview"
	vipAssetPreviewResourcePath = "/vipCfg/assetPreview/resource"
)

func initVipAssetPreviewController() {
	httpserver.RegNonAuthGETHandler(vipAssetPreviewResourcePath, handleVipAssetPreviewResource)
	httpserver.RegNonAuthGETHandler(vipAssetPreviewPath, handleVipAssetPreviewList)
}

func handleVipAssetPreviewList(r *ghttp.Request) {
	html, err := vip.RenderVipAssetPreviewListHTML()
	if err != nil {
		r.Response.WriteStatus(http.StatusInternalServerError, "render vip asset preview failed")
		return
	}
	r.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	r.Response.Write(html)
}

func handleVipAssetPreviewResource(r *ghttp.Request) {
	cfgID := r.Get("id").Uint64()
	if cfgID == 0 {
		r.Response.WriteStatus(http.StatusBadRequest, "missing vip cfg id")
		return
	}
	field := r.Get("field").String()
	html, err := vip.RenderVipAssetPreviewResourceHTML(cfgID, field)
	if err != nil {
		r.Response.WriteStatus(http.StatusNotFound, "vip asset not found")
		return
	}
	r.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	r.Response.Write(html)
}
