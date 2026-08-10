package privacypolicy

import (
	"strings"

	"xr-game-server/errercode"
)

func normalizeApiBase(raw string) (string, error) {
	base := strings.TrimSpace(raw)
	if base == "" {
		return "", nil
	}
	base = strings.TrimRight(base, "/")
	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		return "", errercode.CreateCode(errercode.InvalidParam)
	}
	return base, nil
}

// resolvePolicyUrl 已是 http(s) 则原样返回; 否则与 apiBase 拼接为完整 URL。
func resolvePolicyUrl(apiBase, raw string) (string, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", nil
	}
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		return value, nil
	}
	base, err := normalizeApiBase(apiBase)
	if err != nil {
		return "", err
	}
	if base == "" {
		return "", errercode.CreateCode(errercode.InvalidParam)
	}
	if !strings.HasPrefix(value, "/") {
		value = "/" + value
	}
	return base + value, nil
}

func stripApiBase(apiBase, full string) string {
	value := strings.TrimSpace(full)
	base := strings.TrimRight(strings.TrimSpace(apiBase), "/")
	if value == "" || base == "" {
		return value
	}
	prefix := base + "/"
	if strings.HasPrefix(value, prefix) {
		return value[len(base):]
	}
	if value == base {
		return "/"
	}
	return value
}
