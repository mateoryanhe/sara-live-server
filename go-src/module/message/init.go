package message

import (
	"xr-game-server/core/event"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/gameevent"
)

func Init() {
	cfgdao.InitActivityMessageDao()
	event.Sub(gameevent.SystemMessageEvent, onSystemMessage)
	initMessageCleanup()
}
