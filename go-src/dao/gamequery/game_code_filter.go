package gamequery

import (
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
)

// NormalizeGameCodes 去空、去重，保持顺序。
func NormalizeGameCodes(codes []string) []string {
	if len(codes) == 0 {
		return nil
	}
	out := make([]string, 0, len(codes))
	seen := make(map[string]struct{}, len(codes))
	for _, code := range codes {
		code = strings.TrimSpace(code)
		if code == "" {
			continue
		}
		if _, ok := seen[code]; ok {
			continue
		}
		seen[code] = struct{}{}
		out = append(out, code)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// ApplyGameCodeFilter 多选优先 IN；否则单字段模糊匹配。
func ApplyGameCodeFilter(m *gdb.Model, column string, fuzzyCode string, codes []string) *gdb.Model {
	if m == nil {
		return m
	}
	normalized := NormalizeGameCodes(codes)
	if len(normalized) > 0 {
		return m.Where(column+" IN (?)", normalized)
	}
	if fuzzyCode = strings.TrimSpace(fuzzyCode); fuzzyCode != "" {
		return m.Where(column+" LIKE ?", "%"+fuzzyCode+"%")
	}
	return m
}

// MatchGameCode 内存过滤：codes 非空时精确命中其一；否则 fuzzy 子串匹配。
func MatchGameCode(gameCode string, fuzzyCode string, codes []string) bool {
	normalized := NormalizeGameCodes(codes)
	if len(normalized) > 0 {
		for _, code := range normalized {
			if gameCode == code {
				return true
			}
		}
		return false
	}
	fuzzyCode = strings.TrimSpace(fuzzyCode)
	if fuzzyCode == "" {
		return true
	}
	return strings.Contains(strings.ToLower(gameCode), strings.ToLower(fuzzyCode))
}
