package country

import "strings"

// HaiPayRegionCodes HaiPay 全球收银台支持的 region 简码(有序列表,供 App/CMS 选型).
// 与 HaiPay 文档 region 白名单对齐;增减只改此处.
var HaiPayRegionCodes = []string{
	"ID", "PH", "MY", "IN", "TH", "VN", "SG", "HK", "TW",
	"JP", "KR", "PK", "BR", "US", "GB", "EU", "IT", "AT", "BE",
	"NL", "PL", "TR", "AE", "SA", "QA", "KW", "BH", "OM", "EG",
}

var haiPayRegionSet map[string]struct{}

func init() {
	haiPayRegionSet = make(map[string]struct{}, len(HaiPayRegionCodes))
	for _, code := range HaiPayRegionCodes {
		haiPayRegionSet[normalizeCode(code)] = struct{}{}
	}
}

// IsHaiPayRegion 是否为 HaiPay 收银台合法 region 简码.
func IsHaiPayRegion(code string) bool {
	_, ok := haiPayRegionSet[normalizeCode(code)]
	return ok
}

// ListHaiPayRegions 返回 HaiPay 可选国家/地区(按 HaiPayRegionCodes 顺序;缺表项跳过).
func ListHaiPayRegions() []Country {
	out := make([]Country, 0, len(HaiPayRegionCodes))
	for _, code := range HaiPayRegionCodes {
		if c, ok := Get(code); ok {
			out = append(out, c)
			continue
		}
		// 表缺失时仍返回简码占位,避免选型列表空洞
		code = normalizeCode(code)
		out = append(out, Country{
			Code:     code,
			NameEn:   code,
			FlagIcon: FlagFileName(code),
		})
	}
	return out
}

// NormalizeHaiPayRegionHint 将 2 位 region 规范化;非法返回空.
func NormalizeHaiPayRegionHint(code string) string {
	code = strings.ToUpper(strings.TrimSpace(code))
	if !IsHaiPayRegion(code) {
		return ""
	}
	return code
}
