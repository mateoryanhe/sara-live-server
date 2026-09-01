package controller

import (
	"net/http"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/httpserver"
	"xr-game-server/module/gift"
)

const (
	giftAssetPreviewPath           = "/gift/assetPreview"
	giftAssetPreviewAnimationPath  = "/gift/assetPreview/animation"
)

func initGiftAssetPreviewController() {
	httpserver.RegNonAuthGETHandler(giftAssetPreviewAnimationPath, handleGiftAssetPreviewAnimation)
	httpserver.RegNonAuthGETHandler(giftAssetPreviewPath, handleGiftAssetPreviewList)
}

func handleGiftAssetPreviewList(r *ghttp.Request) {
	html, err := gift.RenderAssetPreviewListHTML()
	if err != nil {
		r.Response.WriteStatus(http.StatusInternalServerError, "render gift asset preview failed")
		return
	}
	r.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	r.Response.Write(html)
}

func handleGiftAssetPreviewAnimation(r *ghttp.Request) {
	giftID := r.Get("id").Uint64()
	if giftID == 0 {
		r.Response.WriteStatus(http.StatusBadRequest, "missing gift id")
		return
	}
	html, err := gift.RenderAssetPreviewAnimationHTML(giftID)
	if err != nil {
		r.Response.WriteStatus(http.StatusNotFound, "gift not found")
		return
	}
	r.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	r.Response.Write(html)
}
