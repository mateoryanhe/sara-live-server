package message

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/event"
	"xr-game-server/core/xrpool"
	"xr-game-server/entity/message"
	"xr-game-server/gameevent"
)

const welcomePersonalMessageParams = ""

type personalSystemMessageDisplay struct {
	IconEn, IconEs, IconPt, IconHi, IconId                string
	TitleEn, TitleEs, TitlePt, TitleHi, TitleId           string
	ContentEn, ContentEs, ContentPt, ContentHi, ContentId string
}

var welcomePersonalMessageContent = personalSystemMessageDisplay{
	TitleEn:   "Welcome",
	TitleEs:   "Bienvenido",
	TitlePt:   "Bem-vindo",
	TitleHi:   "स्वागत है",
	TitleId:   "Selamat datang",
	ContentEn: "Welcome to Sara Live! We're glad to have you here.",
	ContentEs: "¡Bienvenido a Sara Live! Nos alegra tenerte aquí.",
	ContentPt: "Bem-vindo ao Sara Live! Ficamos felizes em ter você aqui.",
	ContentHi: "Sara Live में आपका स्वागत है! हमें खुशी है कि आप यहाँ हैं।",
	ContentId: "Selamat datang di Sara Live! Kami senang Anda ada di sini.",
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

		addPersonalSystemMessage(val.UserId, entity.PersonalSystemMessageTypeWelcome, welcomePersonalMessageParams)
	})
}

func welcomePersonalMessageDisplay() personalSystemMessageDisplay {
	return welcomePersonalMessageContent
}
