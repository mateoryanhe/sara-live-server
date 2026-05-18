package currency

import "xr-game-server/constants/lang"

// Lang 语言代码(转发至通用 constants/lang 包,保持本包 API 兼容)
type Lang = lang.Lang

const (
	LangZHCN = lang.LangZHCN
	LangZHTW = lang.LangZHTW
	LangEN   = lang.LangEN
)

// DefaultLang 默认语言
const DefaultLang = lang.DefaultLang

// ParseLang 解析 Accept-Language 风格字符串
func ParseLang(s string) Lang { return lang.Parse(s) }
