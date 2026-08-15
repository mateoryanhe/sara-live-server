package agora

import (
	"sync/atomic"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity/live"
)

const (
	defaultTokenExpireSeconds   uint32 = 24 * 60 * 60
	minTokenExpireSeconds       uint32 = 4 * 60 * 60
	maxTokenExpireSeconds       uint32 = 24 * 60 * 60
	minTokenRefreshSeconds      uint32 = 2 * 60 * 60
	tokenRefreshAheadGapSeconds uint32 = 2 * 60 * 60
	defaultCloudPlayerRegion           = "cn"
)

type agoraCfgSnapshot struct {
	AppId               string
	AppCertificate      string
	RestCustomerId      string
	RestCustomerSecret  string
	CloudPlayerRegion   string
	TokenExpireSeconds  uint32
	TokenRefreshSeconds uint32
}

var (
	agoraCfgCache         atomic.Value // *agoraCfgSnapshot
	emptyAgoraCfgSnapshot = &agoraCfgSnapshot{
		TokenExpireSeconds:  defaultTokenExpireSeconds,
		TokenRefreshSeconds: defaultTokenExpireSeconds - tokenRefreshAheadGapSeconds,
	}
)

func reloadAgoraCfgMemory() {
	agoraCfgCache.Store(toAgoraCfgSnapshot(cfgdao.LoadAgoraCfg()))
}

func getAgoraCfgCache() *agoraCfgSnapshot {
	v := agoraCfgCache.Load()
	if v == nil {
		return emptyAgoraCfgSnapshot
	}
	cfg, ok := v.(*agoraCfgSnapshot)
	if !ok || cfg == nil {
		return emptyAgoraCfgSnapshot
	}
	return cfg
}

func toAgoraCfgSnapshot(row *entity.AgoraCfg) *agoraCfgSnapshot {
	if row == nil {
		return emptyAgoraCfgSnapshot
	}
	expireSeconds, refreshSeconds := normalizeAgoraTokenCfg(row.TokenExpireSeconds, row.TokenRefreshSeconds)
	return &agoraCfgSnapshot{
		AppId:               row.AppId,
		AppCertificate:      row.AppCertificate,
		RestCustomerId:      row.RestCustomerId,
		RestCustomerSecret:  row.RestCustomerSecret,
		CloudPlayerRegion:   normalizeCloudPlayerRegion(row.CloudPlayerRegion),
		TokenExpireSeconds:  expireSeconds,
		TokenRefreshSeconds: refreshSeconds,
	}
}

func normalizeCloudPlayerRegion(region string) string {
	switch region {
	case "cn", "ap", "eu", "na":
		return region
	default:
		return defaultCloudPlayerRegion
	}
}

func normalizeAgoraTokenCfg(expireSeconds, refreshSeconds uint32) (uint32, uint32) {
	expireSeconds = normalizeTokenExpireSeconds(expireSeconds)
	if refreshSeconds == 0 {
		refreshSeconds = expireSeconds - tokenRefreshAheadGapSeconds
	}
	refreshSeconds = normalizeTokenRefreshSeconds(expireSeconds, refreshSeconds)
	return expireSeconds, refreshSeconds
}

func normalizeTokenExpireSeconds(seconds uint32) uint32 {
	if seconds == 0 {
		return defaultTokenExpireSeconds
	}
	if seconds < minTokenExpireSeconds {
		return minTokenExpireSeconds
	}
	if seconds > maxTokenExpireSeconds {
		return maxTokenExpireSeconds
	}
	return seconds
}

func normalizeTokenRefreshSeconds(expireSeconds, refreshSeconds uint32) uint32 {
	maxRefresh := expireSeconds - tokenRefreshAheadGapSeconds
	if refreshSeconds < minTokenRefreshSeconds {
		refreshSeconds = minTokenRefreshSeconds
	}
	if refreshSeconds > maxRefresh {
		return maxRefresh
	}
	return refreshSeconds
}

func isValidAgoraTokenCfg(expireSeconds, refreshSeconds uint32) bool {
	if expireSeconds < minTokenExpireSeconds || expireSeconds > maxTokenExpireSeconds {
		return false
	}
	maxRefresh := expireSeconds - tokenRefreshAheadGapSeconds
	if maxRefresh < minTokenRefreshSeconds {
		return false
	}
	return refreshSeconds >= minTokenRefreshSeconds && refreshSeconds <= maxRefresh
}

func getChannelTokenExpireSeconds() uint32 {
	return getAgoraCfgCache().TokenExpireSeconds
}
