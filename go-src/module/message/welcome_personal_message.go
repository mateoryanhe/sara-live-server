package message

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/event"
	"xr-game-server/core/xrpool"
	"xr-game-server/gameevent"
)

const welcomePersonalMessageParams = ""

var welcomePersonalMessage = struct {
	titleEn, titleEs, titlePt, titleHi         string
	contentEn, contentEs, contentPt, contentHi string
}{
	titleEn:   "Welcome",
	titleEs:   "Bienvenido",
	titlePt:   "Bem-vindo",
	titleHi:   "स्वागत है",
	contentEn: "Welcome to Sara Live! We're glad to have you here.",
	contentEs: "¡Bienvenido a Sara Live! Nos alegra tenerte aquí.",
	contentPt: "Bem-vindo ao Sara Live! Ficamos felizes em ter você aqui.",
	contentHi: "Sara Live में आपका स्वागत है! हमें खुशी है कि आप यहाँ हैं।",
}

func initWelcomePersonalMessageEvent() {
	event.Sub(gameevent.RegisterEvent, onRegisterWelcomePersonalMessage)
}

func onRegisterWelcomePersonalMessage(data any) {
	xrpool.AddWithRecover(gctx.New(), func(ctx context.Context) {
		val, ok := data.(*gameevent.RegisterEventData)
		if !ok || val == nil || val.UserId == 0 {
			g.Log().Errorf(gctx.New(), "RegisterEvent payload type error: %T", data)
			return
		}

		addPersonalSystemMessage(
			val.UserId,
			"",
			welcomePersonalMessageParams,
			welcomePersonalMessage.titleEn,
			welcomePersonalMessage.titleEs,
			welcomePersonalMessage.titlePt,
			welcomePersonalMessage.titleHi,
			welcomePersonalMessage.contentEn,
			welcomePersonalMessage.contentEs,
			welcomePersonalMessage.contentPt,
			welcomePersonalMessage.contentHi,
		)
	})

}
