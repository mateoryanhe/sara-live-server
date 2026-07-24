package userinfo

import (
	"xr-game-server/core/event"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/gameevent"
	"xr-game-server/module/randomnick"
)

func initNicknameEvent() {
	event.Sub(gameevent.RegisterEvent, onRegisterAssignNickname)
}

func onRegisterAssignNickname(data any) {
	val, ok := data.(*gameevent.RegisterEventData)
	if !ok || val == nil || val.UserId == 0 {
		return
	}
	user := userinfodao.GetUserInfoByUserId(val.UserId)
	if user == nil || user.Nickname != "" {
		return
	}
	lang := val.NickLang
	if lang == 0 {
		lang = randomnick.DefaultLang
	}
	user.SetNickname(randomnick.PickRandom(lang))
}
