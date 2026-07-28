package agora

// Init 服务启动时加载声网配置到内存,并启动云播放器 Token 续期任务
func Init() {
	reloadAgoraCfgMemory()
	initCloudPlayerTokenRefresh()
}
