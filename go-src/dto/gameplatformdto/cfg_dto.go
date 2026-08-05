package gameplatformdto

import "github.com/gogf/gf/v2/frame/g"

type GetGamePlatformCfgReq struct {
	g.Meta `path:"/getGamePlatformCfg" method:"post" summary:"查询游戏平台接入配置" tags:"游戏平台配置"`
}

type GamePlatformCfgItem struct {
	ID        string `json:"id"`
	VendorUrl string `json:"vendorUrl"`
	Token     string `json:"token"`
	SecretKey string `json:"secretKey"`
	IconUrl   string `json:"iconUrl"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type GetGamePlatformCfgRes struct {
	Cfg *GamePlatformCfgItem `json:"cfg"`
}

type SaveGamePlatformCfgReq struct {
	g.Meta    `path:"/saveGamePlatformCfg" method:"post" summary:"保存游戏平台接入配置" tags:"游戏平台配置"`
	ID        uint64 `json:"id" dc:"配置ID,新建传0"`
	VendorUrl string `json:"vendorUrl" v:"required|max-length:255#厂家URL不能为空|厂家URL最长255字符" dc:"厂家API根地址"`
	Token     string `json:"token" v:"required|max-length:512#Token不能为空|Token最长512字符" dc:"接入Token"`
	SecretKey string `json:"secretKey" v:"required|max-length:255#SecretKey不能为空|SecretKey最长255字符" dc:"SecretKey"`
	IconUrl   string `json:"iconUrl" v:"max-length:512#IconUrl最长512字符" dc:"游戏封面CDN根地址"`
}

type SaveGamePlatformCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
