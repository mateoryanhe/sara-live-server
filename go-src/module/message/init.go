package message

import (
	"xr-game-server/core/event"
	"xr-game-server/gameevent"
)

func Init() {
	event.Sub(gameevent.SystemMessageEvent, onSystemMessage)
}
