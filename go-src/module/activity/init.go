package activity

import (
	"xr-game-server/core/event"
	"xr-game-server/gameevent"
)

func Init() {
	ReloadFirstRechargeActivityCache()
	event.Sub(gameevent.FirstRechargeCompletedEvent, onFirstRechargeCompleted)
}
