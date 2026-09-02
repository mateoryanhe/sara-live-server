package entity

import (
	"xr-game-server/core/migrate"
	activityentity "xr-game-server/entity/activity"
	callentity "xr-game-server/entity/call"
	cmsentity "xr-game-server/entity/cms"
	gameentity "xr-game-server/entity/game"
	liveentity "xr-game-server/entity/live"
	msgentity "xr-game-server/entity/message"
	rechargeentity "xr-game-server/entity/recharge"
	fiatentity "xr-game-server/entity/fiat"
	sventity "xr-game-server/entity/shortvideo"
	statentity "xr-game-server/entity/stat"
	sysentity "xr-game-server/entity/sys"
	userentity "xr-game-server/entity/user"
)

func Init() {
	userentity.Init()
	cmsentity.Init()
	statentity.Init()
	liveentity.Init()
	sventity.Init()
	rechargeentity.Init()
	fiatentity.Init()
	activityentity.Init()
	initAppPkg()
	gameentity.Init()
	sysentity.Init()
	msgentity.Init()
	callentity.Init()
	migrate.Close()
}
