package currency

import "strings"

// Lang 语言代码
type Lang string

const (
	LangZHCN Lang = "zh-CN" // 简体中文
	LangZHTW Lang = "zh-TW" // 繁体中文
	LangEN   Lang = "en"    // 英文(默认)
)

// DefaultLang 默认语言
const DefaultLang = LangEN

// ParseLang 解析前端传入的语言标识(支持 Accept-Language 风格,如 "zh-CN,zh;q=0.9,en;q=0.8"),
// 不识别时回落到默认语言
func ParseLang(s string) Lang {
	if s == "" {
		return DefaultLang
	}
	// 取第一段(忽略 q=权重)
	first := strings.SplitN(s, ",", 2)[0]
	first = strings.SplitN(first, ";", 2)[0]
	first = strings.TrimSpace(first)
	low := strings.ToLower(first)

	switch {
	case strings.HasPrefix(low, "zh-tw"), strings.HasPrefix(low, "zh-hk"), strings.HasPrefix(low, "zh-hant"):
		return LangZHTW
	case strings.HasPrefix(low, "zh"):
		return LangZHCN
	case strings.HasPrefix(low, "en"):
		return LangEN
	default:
		return DefaultLang
	}
}
