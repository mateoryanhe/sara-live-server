package core

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"xr-game-server/core/cfg"
	"xr-game-server/core/push"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
	"xr-game-server/core/xrcorn"
	"xr-game-server/core/xrpool"
	"xr-game-server/core/xrtimer"
)

func Init() {
	cfg.InitCfg()
	snowflake.InitSnowflake()
	xrpool.Init()
	xrtimer.Init()
	push.Init()
	xrcorn.Init()
	syndb.InitSynCache()
}
