package shareactivitydeploy

import (
	"sync/atomic"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity/sys"
)

type cfgSnapshot struct {
	ID        uint64
	UpdatedAt string
}

var cfgCache atomic.Value

func reloadCfgMemory() {
	cfgCache.Store(toCfgSnapshot(cfgdao.LoadShareActivityDeployCfg()))
}

func ReloadShareActivityDeployCache() {
	reloadCfgMemory()
}

func getCfgCache() *cfgSnapshot {
	if cfgCache.Load() == nil {
		reloadCfgMemory()
	}
	v := cfgCache.Load()
	if v == nil {
		return &cfgSnapshot{}
	}
	snap, ok := v.(*cfgSnapshot)
	if !ok || snap == nil {
		return &cfgSnapshot{}
	}
	return snap
}

func toCfgSnapshot(row *entity.ShareActivityDeployCfg) *cfgSnapshot {
	if row == nil || row.ID == 0 {
		return &cfgSnapshot{}
	}
	return &cfgSnapshot{
		ID:        row.ID,
		UpdatedAt: formatCfgTime(row.UpdatedAt),
	}
}
