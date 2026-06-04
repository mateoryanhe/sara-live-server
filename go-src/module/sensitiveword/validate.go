package sensitiveword

import (
	"strings"

	"xr-game-server/errercode"
)

// RequireTextCompliant 校验文本合规;空字符串跳过。不合规返回 TextSensitiveWord 错误
func RequireTextCompliant(texts ...string) error {
	for _, text := range texts {
		if strings.TrimSpace(text) == "" {
			continue
		}
		if !CheckTextCompliant(text) {
			return errercode.CreateCode(errercode.TextSensitiveWord)
		}
	}
	return nil
}
