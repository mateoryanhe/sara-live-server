package apkpromotionchannel

import (
	"strings"
	"sync"
)

type channelURLCacheEntry struct {
	apkURL string
}

var channelURLCache sync.Map // channelCode -> *channelURLCacheEntry

// SetCachedApkURL 写入/刷新渠道 APK 公开 URL(空串表示尚未上传).
func SetCachedApkURL(channelCode, apkURL string) {
	code, ok := normalizeChannelCode(channelCode)
	if !ok {
		return
	}
	channelURLCache.Store(code, &channelURLCacheEntry{apkURL: strings.TrimSpace(apkURL)})
}

// GetCachedApkURL 短链 302 用,只读内存.
func GetCachedApkURL(channelCode string) string {
	code, ok := normalizeChannelCode(channelCode)
	if !ok {
		return ""
	}
	v, ok := channelURLCache.Load(code)
	if !ok {
		return ""
	}
	entry, _ := v.(*channelURLCacheEntry)
	if entry == nil {
		return ""
	}
	return entry.apkURL
}

// RemoveCachedChannel 删除渠道缓存项.
func RemoveCachedChannel(channelCode string) {
	code, ok := normalizeChannelCode(channelCode)
	if !ok {
		return
	}
	channelURLCache.Delete(code)
}

// RenameCachedChannel 改渠道码时迁移缓存.
func RenameCachedChannel(oldCode, newCode string) {
	oldCode, ok1 := normalizeChannelCode(oldCode)
	newCode, ok2 := normalizeChannelCode(newCode)
	if !ok1 {
		return
	}
	if !ok2 {
		RemoveCachedChannel(oldCode)
		return
	}
	url := GetCachedApkURL(oldCode)
	RemoveCachedChannel(oldCode)
	SetCachedApkURL(newCode, url)
}
