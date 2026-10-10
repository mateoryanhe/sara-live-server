package lang

import "strings"

// Lang 语言代码(与 CMS SUPPORTED_LOCALES 对齐,另含 zh-TW 供 App/繁体客户端)
type Lang string

const (
	LangZHCN Lang = "zh-CN" // 简体中文
	LangZHTW Lang = "zh-TW" // 繁体中文(App, CMS 无单独项)
	LangEN   Lang = "en"    // 英文
	LangES   Lang = "es"    // 西班牙语
	LangPT   Lang = "pt"    // 葡萄牙语
	LangHI   Lang = "hi"    // 印地语
	LangID   Lang = "id"    // 印尼语
)

// CMSLocales 与 cms/src/i18n/locales/types.ts SUPPORTED_LOCALES 一致
var CMSLocales = []Lang{LangZHCN, LangEN, LangES, LangPT, LangHI, LangID}

// DefaultLang 默认语言(未识别 Accept-Language 时)
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
	case strings.HasPrefix(low, "id"):
		return LangID
	case strings.HasPrefix(low, "es"):
		return LangES
	case strings.HasPrefix(low, "pt"):
		return LangPT
	case strings.HasPrefix(low, "hi"):
		return LangHI
	case strings.HasPrefix(low, "en"):
		return LangEN
	default:
		return DefaultLang
	}
}
