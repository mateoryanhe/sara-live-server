package metapixel

import (
	"sync"

	"xr-game-server/dao/cfgdao"
	sysentity "xr-game-server/entity/sys"
)

var eventsOnce sync.Once

func Init() {
	cfg := cfgdao.LoadMetaPixelCfgRow()
	if !isRuntimeCfgActive(cfg) {
		return
	}
	ensureRuntimeInitialized()
}

func ensureRuntimeInitialized() {
	cfgdao.EnsureMetaPixelCfgDao()
	cfgdao.ReloadMetaPixelCfgCache()
	eventsOnce.Do(initEvents)
}

func isRuntimeCfgActive(cfg *sysentity.MetaPixelCfg) bool {
	return cfg != nil && cfg.IsActive()
}

// OnMetaPixelCfgSaved CMS 保存配置后: 仅在有有效配置时初始化上报运行时.
func OnMetaPixelCfgSaved(cfg *sysentity.MetaPixelCfg) {
	if !isRuntimeCfgActive(cfg) {
		return
	}
	ensureRuntimeInitialized()
}
