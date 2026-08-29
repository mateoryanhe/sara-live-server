package upload

import "xr-game-server/core/cfg"

// Init 加载上传资源配置到内存
func Init() {
	cfg.RegisterCMSExportStoragePathOverride(GetStoragePath)
	cfg.RegisterCMSExportTtlOverride(GetCmsExportTtlMinutes)
	reloadResourceCfgMemory()
	registerStaticMappings()
}
