package module

import (
	"xr-game-server/module/auth"
	"xr-game-server/module/currencylog"
	"xr-game-server/module/name"
	"xr-game-server/module/rank"
	"xr-game-server/module/rechargecfg"
)

func Init() {
	auth.InitAuth()
	rank.Init()
	name.Init()
	currencylog.Init()
	rechargecfg.Init()
}
