package recharge

import (
	"testing"

	rechargeentity "xr-game-server/entity/recharge"
)

func TestAppRechargeCfgTypeFromClient(t *testing.T) {
	chromeWin := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36"
	if got := AppRechargeCfgTypeFromClient("web", chromeWin); got != rechargeentity.RechargeCfgTypeChannel {
		t.Fatalf("x-platform web got cfgType=%d want channel=%d", got, rechargeentity.RechargeCfgTypeChannel)
	}
	if got := AppRechargeCfgTypeFromClient("WEB", ""); got != rechargeentity.RechargeCfgTypeChannel {
		t.Fatalf("x-platform WEB got cfgType=%d want channel=%d", got, rechargeentity.RechargeCfgTypeChannel)
	}
	if got := AppRechargeCfgTypeFromClient("ios", chromeWin); got != rechargeentity.RechargeCfgTypeIOS {
		t.Fatalf("x-platform ios got cfgType=%d want ios=%d", got, rechargeentity.RechargeCfgTypeIOS)
	}
	if got := AppRechargeCfgTypeFromClient("android", ""); got != rechargeentity.RechargeCfgTypeGoogle {
		t.Fatalf("x-platform android got cfgType=%d want google=%d", got, rechargeentity.RechargeCfgTypeGoogle)
	}
	// 未传 X-Platform 时不因浏览器 UA 走渠道
	if got := AppRechargeCfgTypeFromClient("", chromeWin); got != rechargeentity.RechargeCfgTypeGoogle {
		t.Fatalf("no platform + chrome ua got cfgType=%d want google=%d", got, rechargeentity.RechargeCfgTypeGoogle)
	}
}

func TestAppRechargeCfgTypeFromUserAgent(t *testing.T) {
	okhttp := "okhttp/4.12.0"
	if got := AppRechargeCfgTypeFromUserAgent(okhttp); got != rechargeentity.RechargeCfgTypeGoogle {
		t.Fatalf("okhttp got cfgType=%d want google=%d", got, rechargeentity.RechargeCfgTypeGoogle)
	}

	iosNative := "MyApp/1.0 CFNetwork/1410.1 Darwin/22.6.0"
	if got := AppRechargeCfgTypeFromUserAgent(iosNative); got != rechargeentity.RechargeCfgTypeIOS {
		t.Fatalf("ios native got cfgType=%d want ios=%d", got, rechargeentity.RechargeCfgTypeIOS)
	}

	if got := AppRechargeCfgTypeFromUserAgent(""); got != rechargeentity.RechargeCfgTypeGoogle {
		t.Fatalf("empty ua got cfgType=%d want google=%d", got, rechargeentity.RechargeCfgTypeGoogle)
	}
}
