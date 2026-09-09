package coinmerchantdeploy

import "xr-game-server/core/httpserver"

// Init 加载币商 H5 部署配置到内存,并注册加解密密钥
func Init() {
	reloadCfgMemory()
	httpserver.SetCoinMerchantDeploySecretProvider(GetDeploySecret)
}
