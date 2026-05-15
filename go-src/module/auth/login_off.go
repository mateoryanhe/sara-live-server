package auth

import (
	"xr-game-server/core/cfg"
	"xr-game-server/core/event"
	"xr-game-server/dao/globalcfgdao"
	"xr-game-server/gameevent"
	"xr-game-server/module/globalcfg"
)

const (
	LoginOff = "loginOff"
)

func initLoginOff() {
	flag := globalcfgdao.GetBool(globalcfg.Auth, LoginOff, false)
	cfg.GetAuthCfg().LoginOff = flag
	event.Sub(gameevent.GlobalCfgEvent, func(val any) {
		flagRet := globalcfgdao.GetBool(globalcfg.Auth, LoginOff, false)
		cfg.GetAuthCfg().LoginOff = flagRet
	})
}
