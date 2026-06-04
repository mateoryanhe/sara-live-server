package sensitiveword

import (
	"strings"

	"xr-game-server/constants/lang"
)

// CheckCompliant 校验文本是否合规: 未命中铭感词返回 true, 命中返回 false
func CheckCompliant(l lang.Lang, text string) bool {
	text = strings.TrimSpace(text)
	if text == "" {
		return true
	}
	words := getWords(l)
	if len(words) == 0 {
		return true
	}
	normalized := normalizeText(l, text)
	for _, w := range words {
		if w == "" {
			continue
		}
		if matchWord(l, normalized, w) {
			return false
		}
	}
	return true
}

func matchWord(l lang.Lang, normalizedText, word string) bool {
	if l == lang.LangEN {
		if strings.Contains(word, " ") {
			return strings.Contains(normalizedText, word)
		}
		return matchEnglishToken(normalizedText, word)
	}
	return strings.Contains(normalizedText, word)
}

// matchEnglishToken 英文单词按词边界匹配,降低误伤(如 class 不含 ass)
func matchEnglishToken(text, word string) bool {
	if text == word {
		return true
	}
	const sep = " "
	padded := sep + text + sep
	return strings.Contains(padded, sep+word+sep)
}

func normalizeText(l lang.Lang, text string) string {
	if l == lang.LangEN {
		return strings.ToLower(text)
	}
	return text
}
