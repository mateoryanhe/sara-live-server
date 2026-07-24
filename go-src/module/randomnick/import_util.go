package randomnick

import (
	"strings"
	"unicode/utf8"
)

const maxNicknameRunes = 32

// NormalizeNickname 清洗导入/存储的昵称
func NormalizeNickname(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return ""
	}
	if utf8.RuneCountInString(s) <= maxNicknameRunes {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxNicknameRunes])
}

// ParseImportText 解析批量导入文本(一行一个,忽略空行与重复)
func ParseImportText(text string) []string {
	lines := strings.Split(text, "\n")
	seen := make(map[string]struct{})
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		name := NormalizeNickname(line)
		if name == "" {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		out = append(out, name)
	}
	return out
}
