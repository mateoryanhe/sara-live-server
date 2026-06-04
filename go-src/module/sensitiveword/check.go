package sensitiveword

import (
	"strings"

	"xr-game-server/constants/lang"
)

// CheckTextCompliant 校验文本在所有已加载词库下均合规(用于聊天、昵称等 UGC)
func CheckTextCompliant(text string) bool {
	text = strings.TrimSpace(text)
	if text == "" {
		return true
	}
	for _, l := range supportedLangs {
		if !CheckCompliant(l, text) {
			return false
		}
	}
	return true
}

// CheckTextCompliantForLang 仅按指定语言词库校验
func CheckTextCompliantForLang(l lang.Lang, text string) bool {
	return CheckCompliant(l, text)
}
