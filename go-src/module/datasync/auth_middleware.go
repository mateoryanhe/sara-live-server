package datasync

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/errercode"
)

const HeaderDataSyncToken = "X-Data-Sync-Token"

func MiddlewareDataSyncAuth(r *ghttp.Request) {
	cfg := cfgdao.GetDataSyncCfg()
	if cfg == nil || cfg.Token == "" {
		httpserver.WriteFailJson(r, int(errercode.Token))
		return
	}
	if r.GetHeader(HeaderDataSyncToken) != cfg.Token {
		httpserver.WriteFailJson(r, int(errercode.Token))
		return
	}
	r.Middleware.Next()
}
