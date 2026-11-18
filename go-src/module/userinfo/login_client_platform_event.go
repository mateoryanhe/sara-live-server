package userinfo

import (
	"xr-game-server/core/event"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/gameevent"
)

func initLoginClientPlatformBackfill() {
	event.Sub(gameevent.LoginEvent, onLoginClientPlatformBackfill)
}

func onLoginClientPlatformBackfill(data any) {
	val, ok := data.(*gameevent.LoginEventData)
	if !ok || val == nil || val.UserId == 0 {
		return
	}
	userinfodao.TryBackfillRegisterClientPlatform(val.UserId, val.ClientPlatform)
}