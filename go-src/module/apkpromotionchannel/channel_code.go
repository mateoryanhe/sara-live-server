package apkpromotionchannel

import (
	"regexp"
	"strings"
)

var channelCodePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`)

// 与根路径短链冲突的保留段,禁止作为渠道码.
var reservedChannelCodes = map[string]struct{}{
	"dl": {},
}

func normalizeChannelCode(raw string) (string, bool) {
	code := strings.TrimSpace(raw)
	if code == "" || !channelCodePattern.MatchString(code) {
		return "", false
	}
	if _, reserved := reservedChannelCodes[strings.ToLower(code)]; reserved {
		return "", false
	}
	return code, true
}

// HasCachedChannel 渠道是否已在内存注册(含尚未上传 APK).
func HasCachedChannel(channelCode string) bool {
	code, ok := normalizeChannelCode(channelCode)
	if !ok {
		return false
	}
	_, ok = channelURLCache.Load(code)
	return ok
}
