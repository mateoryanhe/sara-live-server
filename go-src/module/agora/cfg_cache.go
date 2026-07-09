package agora

import (
	"sync/atomic"
	"xr-game-server/dao/agoracfgdao"
	"xr-game-server/entity"
)

const (
	defaultTokenExpireSeconds uint32 = 6 * 60 * 60
	minTokenExpireSeconds     uint32 = 6 * 60 * 60
	maxTokenExpireSeconds     uint32 = 18 * 60 * 60
)

type agoraCfgSnapshot struct {
	AppId              string
	AppCertificate     string
	RestCustomerId     string
	RestCustomerSecret string
	TokenExpireSeconds uint32
}

var (
	agoraCfgCache         atomic.Value // *agoraCfgSnapshot
	emptyAgoraCfgSnapshot = &agoraCfgSnapshot{
		TokenExpireSeconds: defaultTokenExpireSeconds,
	}
)

func reloadAgoraCfgMemory() {
	agoraCfgCache.Store(toAgoraCfgSnapshot(agoracfgdao.Load()))
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
	expireSeconds := normalizeTokenExpireSeconds(row.TokenExpireSeconds)
	return &agoraCfgSnapshot{
		AppId:              row.AppId,
		AppCertificate:     row.AppCertificate,
		RestCustomerId:     row.RestCustomerId,
		RestCustomerSecret: row.RestCustomerSecret,
		TokenExpireSeconds: expireSeconds,
	}
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

func isValidTokenExpireSeconds(seconds uint32) bool {
	return seconds >= minTokenExpireSeconds && seconds <= maxTokenExpireSeconds
}
