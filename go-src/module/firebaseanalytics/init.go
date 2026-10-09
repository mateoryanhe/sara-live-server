package firebaseanalytics

import (
	"sync"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/dao/cfgdao"
	sysentity "xr-game-server/entity/sys"
)

var eventsOnce sync.Once

// Init 加载埋点配置缓存;配置有效时订阅业务事件并初始化 Admin SDK.
func Init() {
	cfgdao.InitFirebaseAnalyticsCfgDao()
	cfgdao.ReloadFirebaseAnalyticsCfgCache()
	cfg := cfgdao.LoadFirebaseAnalyticsCfgRow()
	if isRuntimeCfgActive(cfg) {
		ensureRuntimeInitialized()
	}
	warmInitAdmin(gctx.New())
}

func ensureRuntimeInitialized() {
	eventsOnce.Do(initEvents)
}

// IsEnabled 当前是否启用 Firebase Analytics 埋点.
func IsEnabled() bool {
	return isRuntimeCfgActive(cfgdao.ResolveFirebaseAnalyticsCfg())
}

func isRuntimeCfgActive(cfg *sysentity.FirebaseAnalyticsCfg) bool {
	return cfg != nil && cfg.IsActive()
}

// OnCfgSaved CMS 保存配置后刷新运行时状态.
func OnCfgSaved(cfg *sysentity.FirebaseAnalyticsCfg) {
	ResetFirebaseAdmin()
	if isRuntimeCfgActive(cfg) {
		ensureRuntimeInitialized()
		warmInitAdmin(gctx.New())
	}
}

// GetCfgCached 供后续 App 接口或服务端上报读取.
func GetCfgCached() *sysentity.FirebaseAnalyticsCfg {
	cfg := cfgdao.ResolveFirebaseAnalyticsCfg()
	if !isRuntimeCfgActive(cfg) {
		return nil
	}
	return cfg
}
