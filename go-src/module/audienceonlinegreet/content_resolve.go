package audienceonlinegreet

import (
	"xr-game-server/constants/lang"
	messageentity "xr-game-server/entity/message"
)

// storedGreetContentsForLang 仅数据库已存文案(含历史单条字段); 无数据返回 nil.
func storedGreetContentsForLang(cfg *messageentity.AudienceOnlineAnchorGreetCfg, clientLang lang.Lang) []string {
	if cfg == nil {
		return nil
	}
	switch clientLang {
	case lang.LangZHCN, lang.LangZHTW:
		return cfg.ResolvedContentsEn()
	case lang.LangES:
		return cfg.ResolvedContentsEs()
	case lang.LangPT:
		return cfg.ResolvedContentsPt()
	case lang.LangHI:
		return cfg.ResolvedContentsHi()
	case lang.LangID:
		return cfg.ResolvedContentsId()
	default:
		return cfg.ResolvedContentsEn()
	}
}

func defaultGreetLinesForLang(clientLang lang.Lang) []string {
	def := defaultGreetContentsByLang
	switch clientLang {
	case lang.LangES:
		return def.Es
	case lang.LangPT:
		return def.Pt
	case lang.LangHI:
		return def.Hi
	case lang.LangID:
		return def.Id
	default:
		return def.En
	}
}

// greetContentsForLang 有库用库，无库用 Go 内置( default_greet_lines.go ).
func greetContentsForLang(cfg *messageentity.AudienceOnlineAnchorGreetCfg, clientLang lang.Lang) []string {
	if stored := storedGreetContentsForLang(cfg, clientLang); len(stored) > 0 {
		return append([]string(nil), stored...)
	}
	def := defaultGreetLinesForLang(clientLang)
	return append([]string(nil), def...)
}

func storedGreetContentsByLang(cfg *messageentity.AudienceOnlineAnchorGreetCfg) greetContentsByLang {
	if cfg == nil {
		return greetContentsByLang{}
	}
	return greetContentsByLang{
		En: cloneStrSlice(cfg.ResolvedContentsEn()),
		Es: cloneStrSlice(cfg.ResolvedContentsEs()),
		Pt: cloneStrSlice(cfg.ResolvedContentsPt()),
		Hi: cloneStrSlice(cfg.ResolvedContentsHi()),
		Id: cloneStrSlice(cfg.ResolvedContentsId()),
	}
}

func cloneStrSlice(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	return append([]string(nil), in...)
}