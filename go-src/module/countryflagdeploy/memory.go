package countryflagdeploy

import (
	"sync/atomic"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity/sys"
)

type cfgSnapshot struct {
	ID        uint64
	Version   string
	UpdatedAt string
}

var cfgCache atomic.Pointer[cfgSnapshot]

func ReloadCountryFlagCache() {
	cfgCache.Store(toCfgSnapshot(cfgdao.LoadCountryFlagCfg()))
}

func getCfgCache() *cfgSnapshot {
	if v := cfgCache.Load(); v != nil {
		return v
	}
	return &cfgSnapshot{}
}

// CurrentVersion 当前生效国旗版本(空表示尚未部署)
func CurrentVersion() string {
	return getCfgCache().Version
}

func toCfgSnapshot(row *entity.CountryFlagCfg) *cfgSnapshot {
	if row == nil {
		return &cfgSnapshot{}
	}
	updated := ""
	if !row.UpdatedAt.IsZero() {
		updated = row.UpdatedAt.Format("2006-01-02 15:04:05")
	}
	return &cfgSnapshot{
		ID:        row.ID,
		Version:   row.Version,
		UpdatedAt: updated,
	}
}
