package shareappdownload

import (
	"strings"
	"sync/atomic"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity/user"
)

type cfgSnapshot struct {
	ShareText string
}

var (
	cfgCache         atomic.Value // *cfgSnapshot
	emptyCfgSnapshot = &cfgSnapshot{}
	apkDownloadURL   atomic.Value // string
)

func reloadCfgMemory() {
	cfgCache.Store(toCfgSnapshot(cfgdao.LoadShareAppDownloadCfg()))
}

func getCfgCache() *cfgSnapshot {
	v := cfgCache.Load()
	if v == nil {
		return emptyCfgSnapshot
	}
	cfg, ok := v.(*cfgSnapshot)
	if !ok || cfg == nil {
		return emptyCfgSnapshot
	}
	return cfg
}

func toCfgSnapshot(row *entity.ShareAppDownloadCfg) *cfgSnapshot {
	if row == nil {
		return emptyCfgSnapshot
	}
	return &cfgSnapshot{
		ShareText: strings.TrimSpace(row.ShareText),
	}
}

func storeApkDownloadURL(apkURL string) {
	apkDownloadURL.Store(strings.TrimSpace(apkURL))
}

// GetDownloadUrl 官网 assets 下最大 APK 的公开 URL(内存).
func GetDownloadUrl() string {
	v := apkDownloadURL.Load()
	if v == nil {
		return ""
	}
	s, _ := v.(string)
	return strings.TrimSpace(s)
}

// GetShareText 内存中的分享文案
func GetShareText() string {
	return getCfgCache().ShareText
}
