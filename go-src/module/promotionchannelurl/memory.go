package promotionchannelurl

import (
	"strconv"
	"strings"
	"sync"
)

type channelURLCacheEntry struct {
	apkURL string
}

var channelURLCache sync.Map // channelCode(int) -> *channelURLCacheEntry

func channelCacheKey(code int) string {
	return strconv.Itoa(code)
}

func SetCachedApkURL(channelCode int, apkURL string) {
	if channelCode <= 0 {
		return
	}
	channelURLCache.Store(channelCacheKey(channelCode), &channelURLCacheEntry{apkURL: strings.TrimSpace(apkURL)})
}

func GetCachedApkURL(channelCode int) string {
	if channelCode <= 0 {
		return ""
	}
	v, ok := channelURLCache.Load(channelCacheKey(channelCode))
	if !ok {
		return ""
	}
	entry, _ := v.(*channelURLCacheEntry)
	if entry == nil {
		return ""
	}
	return entry.apkURL
}

func HasCachedChannel(channelCode int) bool {
	if channelCode <= 0 {
		return false
	}
	_, ok := channelURLCache.Load(channelCacheKey(channelCode))
	return ok
}

func RemoveCachedChannel(channelCode int) {
	if channelCode <= 0 {
		return
	}
	channelURLCache.Delete(channelCacheKey(channelCode))
}

func RenameCachedChannel(oldCode, newCode int) {
	if oldCode <= 0 {
		return
	}
	url := GetCachedApkURL(oldCode)
	RemoveCachedChannel(oldCode)
	if newCode > 0 {
		SetCachedApkURL(newCode, url)
	}
}
