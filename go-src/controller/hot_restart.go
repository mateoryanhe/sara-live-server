package controller

import (
	"xr-game-server/core/hotrestart"
	"xr-game-server/core/httpserver"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

const HotRestartUrl = "/internal/hotRestart"

func initHotRestartController() {
	httpserver.RegInternalHandler(HotRestartUrl, handleHotRestart)
}

func handleHotRestart(r *ghttp.Request) {
	auth := r.Get("auth").String()
	accepted, reason := hotrestart.TryTriggerHotRestart(auth)
	if !accepted {
		r.Response.WriteStatus(403)
		r.Response.WriteJson(g.Map{"ok": false, "reason": reason})
		return
	}
	r.Response.WriteJson(g.Map{"ok": true, "reason": reason})
}
