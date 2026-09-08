package upload

import (
	"xr-game-server/core/cfg"
	"xr-game-server/module/fileexport"
)

// Init 加载上传资源配置到内存
func Init() {
	cfg.RegisterCMSExportStoragePathOverride(GetStoragePath)
	cfg.RegisterCMSExportTtlOverride(GetCmsExportTtlMinutes)
	cfg.RegisterCMSExportURLBuilder(BuildExportFileURL)
	cfg.RegisterCMSExportURLPrefixProvider(GetFileAccessDomain)
	fileexport.RegisterPublisher(func(rec *fileexport.Record) error {
		if rec == nil {
			return nil
		}
		return PublishLocalFileToS3(StoreCatExport, rec.FileName, rec.AbsPath)
	})
	fileexport.RegisterRemover(func(fileName string) {
		DeleteExportStored(fileName)
	})
	reloadResourceCfgMemory()
	registerStaticMappings()
}
