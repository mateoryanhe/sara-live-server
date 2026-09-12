package country

import "strings"

// FormatZhEn CMS 展示用:「中文名 / English」;未知简码原样返回(兼容历史中文入库).
func FormatZhEn(codeOrLegacy string) string {
	s := strings.TrimSpace(codeOrLegacy)
	if s == "" {
		return ""
	}
	c, ok := Get(s)
	if !ok {
		return s
	}
	zh := strings.TrimSpace(c.NameZh)
	en := strings.TrimSpace(c.NameEn)
	switch {
	case zh != "" && en != "":
		return zh + " / " + en
	case zh != "":
		return zh
	case en != "":
		return en
	default:
		return c.Code
	}
}
