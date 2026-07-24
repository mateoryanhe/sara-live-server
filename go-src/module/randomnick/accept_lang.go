package randomnick

import (
	"strings"
)

// ParseAcceptLanguage 从 Accept-Language 解析昵称语言;不识别时默认英文
func ParseAcceptLanguage(header string) uint8 {
	header = strings.TrimSpace(header)
	if header == "" {
		return DefaultLang
	}
	first := strings.SplitN(header, ",", 2)[0]
	first = strings.SplitN(first, ";", 2)[0]
	low := strings.ToLower(strings.TrimSpace(first))

	switch {
	case strings.HasPrefix(low, "es"):
		return LangES
	case strings.HasPrefix(low, "hi"):
		return LangHI
	case strings.HasPrefix(low, "pt"):
		return LangPT
	case strings.HasPrefix(low, "en"):
		return LangEN
	default:
		return DefaultLang
	}
}
