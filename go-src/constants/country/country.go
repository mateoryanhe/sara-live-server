package country

import "strings"

// Country 国家/地区静态信息。
// Code 为 ISO 3166-1 alpha-2（大写）；EU 为支付侧特例区域码，非主权国家。
type Country struct {
	Code     string // 简码,如 ID / US
	NameEn   string // 英文名称
	NameZh   string // 中文名称
	FlagIcon string // 国旗文件名,如 id.png(完整路径用 RelPath+当前 version)
}

const (
	// AssetDir 图片静态根下国旗目录(其下再按 version 分子目录)
	AssetDir    = "country-flags"
	FlagIconExt = ".png"
)

// FlagFileName 国旗文件名(小写简码.png)
func FlagFileName(code string) string {
	c := strings.ToLower(strings.TrimSpace(code))
	if c == "" {
		return ""
	}
	return c + FlagIconExt
}

// RelPath 相对图片根路径: country-flags/{version}/{code}.png
func RelPath(code, version string) string {
	name := FlagFileName(code)
	v := strings.TrimSpace(version)
	if name == "" || v == "" {
		return ""
	}
	return AssetDir + "/" + v + "/" + name
}

// Get 按简码查询(大小写不敏感);不存在返回 false。
func Get(code string) (Country, bool) {
	c := normalizeCode(code)
	if c == "" {
		return Country{}, false
	}
	v, ok := byCode[c]
	return v, ok
}

// MustGet 查询;不存在返回仅含 Code/FlagIcon 的占位(NameEn/NameZh 空)。
func MustGet(code string) Country {
	if v, ok := Get(code); ok {
		return v
	}
	c := normalizeCode(code)
	return Country{Code: c, FlagIcon: FlagFileName(c)}
}

// Exists 简码是否在静态表中。
func Exists(code string) bool {
	_, ok := Get(code)
	return ok
}

// All 返回静态表副本(稳定顺序,与定义顺序一致)。
func All() []Country {
	out := make([]Country, len(all))
	copy(out, all)
	return out
}

func normalizeCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func newCountry(code, nameEn, nameZh string) Country {
	code = normalizeCode(code)
	return Country{
		Code:     code,
		NameEn:   strings.TrimSpace(nameEn),
		NameZh:   strings.TrimSpace(nameZh),
		FlagIcon: FlagFileName(code),
	}
}
