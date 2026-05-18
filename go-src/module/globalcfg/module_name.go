package globalcfg

import "xr-game-server/entity"

const (
	Auth     entity.GlobalCfgModule = "Auth"
	Resource entity.GlobalCfgModule = "Resource"
)

// Resource 模块下的 key
const (
	// ResourceKeyDomain 静态资源域名(如 https://cdn.example.com)
	ResourceKeyDomain = "Domain"
)
