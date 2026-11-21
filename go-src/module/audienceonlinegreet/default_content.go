package audienceonlinegreet

import (
	"crypto/rand"
	"math/big"

	"xr-game-server/constants/lang"
	messageentity "xr-game-server/entity/message"
)

type greetContentsByLang struct {
	En, Es, Pt, Hi, Id []string
}

// defaultGreetContentsByLang 库无配置或列表为空时的回落(见 default_greet_lines.go).
var defaultGreetContentsByLang = defaultAudienceGreetLinesByLang()

func pickGreetContentForLang(cfg *messageentity.AudienceOnlineAnchorGreetCfg, clientLang lang.Lang) string {
	list := greetContentsForLang(cfg, clientLang)
	if len(list) == 0 {
		return ""
	}
	return pickRandomGreetContent(list)
}

func pickRandomGreetContent(list []string) string {
	if len(list) == 0 {
		return ""
	}
	if len(list) == 1 {
		return list[0]
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(list))))
	if err != nil {
		return list[0]
	}
	return list[n.Int64()]
}

func builtinGreetCfgContents() greetContentsByLang {
	return defaultAudienceGreetLinesByLang()
}