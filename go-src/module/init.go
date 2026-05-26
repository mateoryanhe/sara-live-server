package module

import (
	"xr-game-server/module/auth"
	"xr-game-server/module/banner"
	"xr-game-server/module/currencylog"
	"xr-game-server/module/game"
	"xr-game-server/module/liveroom"
	"xr-game-server/module/message"
	"xr-game-server/module/name"
	"xr-game-server/module/rank"
	"xr-game-server/module/recharge"
	"xr-game-server/module/shortvideo"
	"xr-game-server/module/verification_code"
	"xr-game-server/module/vip"
)

func Init() {
	auth.InitAuth()
	rank.Init()
	name.Init()
	currencylog.Init()
	message.Init()
	recharge.Init()
	banner.Init()
	shortvideo.Init()
	vip.Init()
	game.Init()
	verification_code.Init()
	liveroom.Init()
}
