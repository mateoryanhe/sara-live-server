package module

import (
	"xr-game-server/module/agora"
	"xr-game-server/module/anchorrank"
	"xr-game-server/module/auth"
	"xr-game-server/module/banner"
	"xr-game-server/module/currencylog"
	"xr-game-server/module/game"
	"xr-game-server/module/liveroom"
	"xr-game-server/module/message"
	"xr-game-server/module/name"
	"xr-game-server/module/rank"
	"xr-game-server/module/recharge"
	"xr-game-server/module/richrank"
	"xr-game-server/module/sensitiveword"
	"xr-game-server/module/shortvideo"
	"xr-game-server/module/stat"
	"xr-game-server/module/userinfo"
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
	richrank.Init()
	anchorrank.Init()
	banner.Init()
	vip.Init()
	game.Init()
	agora.Init()
	verification_code.Init()
	sensitiveword.Init()
	liveroom.Init()
	userinfo.Init()
	stat.Init()
	shortvideo.Init()
}
