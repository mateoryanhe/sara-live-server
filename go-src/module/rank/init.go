package rank

import (
	"xr-game-server/core/event"
	"xr-game-server/gameevent"
)

func Init() {
	event.Sub(gameevent.AddRankValEvent, onAddEvent)
	event.Sub(gameevent.ReduceRankValEvent, onReduceEvent)

	initSettlement()
	initRank(nil)
	initMgr()
}
