package h5livedeploy

import "xr-game-server/core/httpserver"

// Init 加载 H5 直播部署配置到内存
func Init() {
	reloadCfgMemory()
	httpserver.SetH5DeploySecretProvider(GetDeploySecret)
}
