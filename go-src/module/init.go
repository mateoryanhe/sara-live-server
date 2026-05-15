package module

import (
	"xr-game-server/module/auth"
	"xr-game-server/module/name"
	"xr-game-server/module/rank"
)

func Init() {
	auth.InitAuth()
	rank.Init()
	name.Init()
}
