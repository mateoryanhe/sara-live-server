package lang

import "strings"

// Lang 语言代码(项目通用)
type Lang string

const (
	LangZHCN Lang = "zh-CN" // 简体中文
	LangZHTW Lang = "zh-TW" // 繁体中文
	LangEN   Lang = "en"    // 英文(默认)
)

// DefaultLang 默认语言
const DefaultLang = LangEN

// Parse 解析 Accept-Language 风格字符串(如 "zh-CN,zh;q=0.9,en;q=0.8"),
// 不识别时回落到默认语言
func Parse(s string) Lang {
	if s == "" {
		return DefaultLang
	}
	first := strings.SplitN(s, ",", 2)[0]
	first = strings.SplitN(first, ";", 2)[0]
	low := strings.ToLower(strings.TrimSpace(first))

	switch {
	case strings.HasPrefix(low, "zh-tw"),
		strings.HasPrefix(low, "zh-hk"),
		strings.HasPrefix(low, "zh-hant"):
		return LangZHTW
	case strings.HasPrefix(low, "zh"):
		return LangZHCN
	case strings.HasPrefix(low, "en"):
		return LangEN
	default:
		return DefaultLang
	}
}
