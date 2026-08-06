package gameplatform

import "strings"

// Platform 游戏平台类型
type Platform string

const (
	ZY Platform = "ZY"
)

func (p Platform) String() string {
	return string(p)
}

// ParsePlatform 解析平台类型(忽略大小写).
func ParsePlatform(raw string) (Platform, bool) {
	p := Platform(strings.ToUpper(strings.TrimSpace(raw)))
	if IsValid(p) {
		return p, true
	}
	return "", false
}

// IsValid 是否为合法平台类型
func IsValid(p Platform) bool {
	switch p {
	case ZY:
		return true
	default:
		return false
	}
}

// IsValidString 是否为合法平台类型字符串
func IsValidString(s string) bool {
	_, ok := ParsePlatform(s)
	return ok
}
