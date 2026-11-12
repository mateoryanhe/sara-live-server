package attrreport

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"regexp"
	"strings"
)

var nonDigitRE = regexp.MustCompile(`[^0-9]`)

// HashSHA256 对已标准化字符串做 SHA-256 小写 hex.
func HashSHA256(normalized string) string {
	normalized = strings.TrimSpace(normalized)
	if normalized == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(sum[:])
}

// NormalizeEmail 去前后空格并转小写(不哈希).
func NormalizeEmail(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

// HashEmail 标准化后哈希; 无效空串返回空.
func HashEmail(raw string) string {
	return HashSHA256(NormalizeEmail(raw))
}

// NormalizePhone 去符号/字母/前置零,并保证带区号数字串(不哈希).
// areaCode 为纯区号(如 "1"); phone 可为 "(650)555-1212".
func NormalizePhone(areaCode, phone string) string {
	digits := nonDigitRE.ReplaceAllString(phone, "")
	digits = strings.TrimLeft(digits, "0")
	if digits == "" {
		return ""
	}
	area := nonDigitRE.ReplaceAllString(areaCode, "")
	area = strings.TrimLeft(area, "0")
	if area != "" && !strings.HasPrefix(digits, area) {
		digits = area + digits
	}
	if len(digits) < 7 {
		return ""
	}
	return digits
}

// HashPhone 标准化后哈希; 无效返回空.
func HashPhone(areaCode, phone string) string {
	return HashSHA256(NormalizePhone(areaCode, phone))
}

// SplitClientIPs 从 client_ip 拆出 ipv4/ipv6(可空).
func SplitClientIPs(raw string) (ipv4, ipv6 string) {
	ipStr := strings.TrimSpace(raw)
	if ipStr == "" {
		return "", ""
	}
	// 兼容 host:port / [ipv6]:port
	if host, _, err := net.SplitHostPort(ipStr); err == nil {
		ipStr = host
	}
	ipStr = strings.Trim(ipStr, "[]")
	parsed := net.ParseIP(ipStr)
	if parsed == nil {
		return "", ""
	}
	if v4 := parsed.To4(); v4 != nil {
		return v4.String(), ""
	}
	return "", parsed.String()
}
