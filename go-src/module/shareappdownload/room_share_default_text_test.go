package shareappdownload

import (
	"testing"

	"xr-game-server/constants/lang"
)

func TestLiveRoomShareTextForLangUsesBuiltinWhenCMSEmpty(t *testing.T) {
	snap := &liveRoomShareSnapshot{}
	if got := liveRoomShareTextForLang(snap, lang.LangEN); got != defaultLiveRoomShareTextEn {
		t.Fatalf("en: got %q", got)
	}
	if got := liveRoomShareTextForLang(snap, lang.LangES); got != defaultLiveRoomShareTextEs {
		t.Fatalf("es: got %q", got)
	}
	if got := liveRoomShareTextForLang(snap, lang.LangZHCN); got != defaultLiveRoomShareTextEn {
		t.Fatalf("zh-cn: want en default, got %q", got)
	}
}

func TestLiveRoomShareTextForLangPrefersCMS(t *testing.T) {
	snap := &liveRoomShareSnapshot{ShareTextEn: "Custom EN"}
	if got := liveRoomShareTextForLang(snap, lang.LangEN); got != "Custom EN" {
		t.Fatalf("got %q", got)
	}
}

func TestLiveRoomShareTextForLangLocaleFallbackThenBuiltin(t *testing.T) {
	snap := &liveRoomShareSnapshot{ShareTextEn: "CMS EN only"}
	if got := liveRoomShareTextForLang(snap, lang.LangES); got != "CMS EN only" {
		t.Fatalf("es with empty cms es should fall back en cms: got %q", got)
	}
	snap.ShareTextEs = "CMS ES"
	if got := liveRoomShareTextForLang(snap, lang.LangES); got != "CMS ES" {
		t.Fatalf("got %q", got)
	}
}