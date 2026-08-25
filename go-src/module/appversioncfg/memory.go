package appversioncfg

import (
	"sync/atomic"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/appversioncfgdto"
	"xr-game-server/entity/cms"
)

type cfgSnapshot struct {
	VersionQueryEnabled bool
	Version             string
	DownloadUrl         string
	UpdateDetails       []*appversioncfgdto.AppVersionUpdateDetailItem
}

var cfgCache atomic.Value

func reloadCfgMemory() {
	cfg := cfgdao.LoadAppVersionCfg()
	details := cfgdao.LoadAppVersionUpdateDetails()
	cfgCache.Store(toCfgSnapshot(cfg, details))
}

func getCfgCache() *cfgSnapshot {
	v := cfgCache.Load()
	if v == nil {
		return &cfgSnapshot{
			UpdateDetails: []*appversioncfgdto.AppVersionUpdateDetailItem{},
		}
	}
	snap, ok := v.(*cfgSnapshot)
	if !ok || snap == nil {
		return &cfgSnapshot{
			UpdateDetails: []*appversioncfgdto.AppVersionUpdateDetailItem{},
		}
	}
	return snap
}

func toCfgSnapshot(row *entity.AppVersionCfg, details []*entity.AppVersionUpdateDetail) *cfgSnapshot {
	if row == nil {
		return &cfgSnapshot{
			UpdateDetails: toUpdateDetailItems(details),
		}
	}
	return &cfgSnapshot{
		VersionQueryEnabled: row.VersionQueryEnabled,
		Version:             row.Version,
		DownloadUrl:         row.DownloadUrl,
		UpdateDetails:       toUpdateDetailItems(details),
	}
}

// GetVersionQuerySnapshot 从内存缓存读取 App 版本配置(App 端使用,不查库)
func GetVersionQuerySnapshot() (enabled bool, version, downloadUrl string, updateDetails []*appversioncfgdto.AppVersionUpdateDetailItem) {
	snap := getCfgCache()
	return snap.VersionQueryEnabled, snap.Version, snap.DownloadUrl, snap.UpdateDetails
}
