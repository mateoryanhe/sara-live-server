package banner

import "strings"

// bannerHidesForPackage 当前包名在排除列表中时不展示该 Banner。
func bannerHidesForPackage(excludeRaw, packageName string) bool {
	packageName = strings.TrimSpace(strings.ToLower(packageName))
	if packageName == "" {
		return false
	}
	for _, part := range splitPackageNameList(excludeRaw) {
		if strings.TrimSpace(strings.ToLower(part)) == packageName {
			return true
		}
	}
	return false
}

func splitPackageNameList(raw string) []string {
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	raw = strings.ReplaceAll(raw, "\r", "\n")
	raw = strings.ReplaceAll(raw, "\n", ",")
	raw = strings.ReplaceAll(raw, ";", ",")
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func normalizeExcludePackageNames(raw string) string {
	parts := splitPackageNameList(raw)
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, ",")
}
