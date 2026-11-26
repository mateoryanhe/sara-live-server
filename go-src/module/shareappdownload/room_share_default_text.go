package shareappdownload

import "xr-game-server/constants/lang"

// 直播间分享默认文案仅在 Go 维护; CMS 某语言留空时 App 与预览均回落此处.
const (
	defaultLiveRoomShareTextEn = "Join me in this live room on Sofie — tap the link below!"
	defaultLiveRoomShareTextEs = "¡Únete a mi sala en vivo en Sofie! Toca el enlace de abajo."
	defaultLiveRoomShareTextPt = "Entre na minha sala ao vivo no Sofie — toque no link abaixo."
	defaultLiveRoomShareTextHi = "Sofie पर मेरे लाइव रूम में आइए — नीचे लिंक टैप करें।"
	defaultLiveRoomShareTextId = "Gabung ruang live saya di Sofie — ketuk tautan di bawah."
)

func defaultLiveRoomShareTextForLang(clientLang lang.Lang) string {
	switch clientLang {
	case lang.LangES:
		return defaultLiveRoomShareTextEs
	case lang.LangPT:
		return defaultLiveRoomShareTextPt
	case lang.LangHI:
		return defaultLiveRoomShareTextHi
	case lang.LangID:
		return defaultLiveRoomShareTextId
	default:
		return defaultLiveRoomShareTextEn
	}
}

// DefaultLiveRoomShareText CMS/未配置时按客户端语言取内置文案(简繁与其它未维护语言回落英文).
func DefaultLiveRoomShareText(clientLang lang.Lang) string {
	return defaultLiveRoomShareTextForLang(clientLang)
}

// BuiltinLiveRoomShareTexts 供 CMS 查询展示默认文案(与 App 回落一致).
func BuiltinLiveRoomShareTexts() (en, es, pt, hi, id string) {
	return defaultLiveRoomShareTextEn,
		defaultLiveRoomShareTextEs,
		defaultLiveRoomShareTextPt,
		defaultLiveRoomShareTextHi,
		defaultLiveRoomShareTextId
}