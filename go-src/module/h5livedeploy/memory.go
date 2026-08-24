package h5livedeploy

import (
	"strings"
	"sync/atomic"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity/sys"
)

type cfgSnapshot struct {
	ID           uint64
	DeploySecret string
	UpdatedAt    string
}

var cfgCache atomic.Value

func reloadCfgMemory() {
	cfgCache.Store(toCfgSnapshot(cfgdao.LoadH5LiveDeployCfg()))
}

func ReloadH5LiveDeployCache() {
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

func GetDeploySecret() string {
	return getCfgCache().DeploySecret
}

func toCfgSnapshot(row *entity.H5LiveDeployCfg) *cfgSnapshot {
	if row == nil || row.ID == 0 {
		return &cfgSnapshot{}
	}
	return &cfgSnapshot{
		ID:           row.ID,
		DeploySecret: strings.TrimSpace(row.DeploySecret),
		UpdatedAt:    formatCfgTime(row.UpdatedAt),
	}
}
