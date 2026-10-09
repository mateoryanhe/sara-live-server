package auth

import (
	"encoding/base64"
	"strings"
	"unicode/utf8"
)

// normalizeSmtpCredential 支持粘贴 SES CSV 的明文或 Base64 凭证.
func normalizeSmtpCredential(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil || len(decoded) == 0 || !utf8.Valid(decoded) {
		return raw
	}
	text := strings.TrimSpace(string(decoded))
	if text == "" || strings.ContainsRune(text, '\x00') {
		return raw
	}
	return text
}
