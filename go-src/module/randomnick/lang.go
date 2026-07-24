package randomnick

// 随机昵称语言(与 CMS 导入、注册分配一致)
const (
	LangEN uint8 = 1 // 英文(默认)
	LangES uint8 = 2 // 西班牙语
	LangHI uint8 = 3 // 印地语
	LangPT uint8 = 4 // 葡萄牙语
)

const DefaultLang = LangEN

var langCodes = map[uint8]string{
	LangEN: "en",
	LangES: "es",
	LangHI: "hi",
	LangPT: "pt",
}

var langLabels = map[uint8]string{
	LangEN: "English",
	LangES: "Español",
	LangHI: "हिन्दी",
	LangPT: "Português",
}

// LangCode 语言代码
func LangCode(lang uint8) string {
	if code, ok := langCodes[lang]; ok {
		return code
	}
	return langCodes[DefaultLang]
}

// LangLabel 语言展示名
func LangLabel(lang uint8) string {
	if label, ok := langLabels[lang]; ok {
		return label
	}
	return langLabels[DefaultLang]
}

// SupportedLangs 支持的语言列表
func SupportedLangs() []uint8 {
	return []uint8{LangEN, LangES, LangHI, LangPT}
}

// NormalizeLang 非法值回落英文
func NormalizeLang(lang uint8) uint8 {
	switch lang {
	case LangEN, LangES, LangHI, LangPT:
		return lang
	default:
		return DefaultLang
	}
}
