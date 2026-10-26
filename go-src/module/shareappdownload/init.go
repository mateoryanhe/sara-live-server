package shareappdownload

func Init() {
	reloadCfgMemory()
	// APK URL 扫描失败不得影响启动
	RefreshOfficialApkDownloadUrl()
}
