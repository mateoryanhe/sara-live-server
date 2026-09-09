package upload

import (
	"os"

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
		if !IsS3Enabled() {
			return nil
		}
		if err := PublishLocalFileToS3(StoreCatExport, rec.FileName, rec.AbsPath); err != nil {
			return err
		}
		// 开云桶:导出只保留云对象,立刻删本地落盘,TTL/主动删除只清云
		if rec.AbsPath != "" {
			_ = os.Remove(rec.AbsPath)
		}
		return nil
	})
	fileexport.RegisterRemover(func(fileName string) {
		DeleteExportStored(fileName)
	})
	reloadResourceCfgMemory()
	registerStaticMappings()
}
