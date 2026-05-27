package agora

import (
	"sync/atomic"
	"xr-game-server/dao/agoracfgdao"
	"xr-game-server/entity"
)

const defaultTokenExpireSeconds uint32 = 21600

type agoraCfgSnapshot struct {
	AppId              string
	AppCertificate     string
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
	expireSeconds := row.TokenExpireSeconds
	if expireSeconds == 0 {
		expireSeconds = defaultTokenExpireSeconds
	}
	return &agoraCfgSnapshot{
		AppId:              row.AppId,
		AppCertificate:     row.AppCertificate,
		TokenExpireSeconds: expireSeconds,
	}
}
