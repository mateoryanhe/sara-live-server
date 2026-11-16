package errercode

import "xr-game-server/constants/lang"

// 西/葡/印地语错误码表:先自英文完整复制,后续可逐条改为人工译文;阿拉伯语见 text_ar.go.
func init() {
	cloneCodeTextMap(lang.LangES, lang.LangEN)
	cloneCodeTextMap(lang.LangPT, lang.LangEN)
	cloneCodeTextMap(lang.LangHI, lang.LangEN)
}

func cloneCodeTextMap(dst, src lang.Lang) {
	srcMap, ok := codeTextMap[src]
	if !ok || len(srcMap) == 0 {
		return
	}
	dup := make(map[XRCode]string, len(srcMap))
	for code, text := range srcMap {
		dup[code] = text
	}
	codeTextMap[dst] = dup
}
