package preloadcfgdto

import "github.com/gogf/gf/v2/frame/g"

type GetPreloadCfgReq struct {
	g.Meta `path:"/getPreloadCfg" method:"post" summary:"查询预热配置" tags:"预热配置"`
}

type PreloadCfgItem struct {
	ID               string `json:"id"`
	RecentLoginLimit int    `json:"recentLoginLimit"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}

type GetPreloadCfgRes struct {
	Cfg *PreloadCfgItem `json:"cfg"`
}

type SavePreloadCfgReq struct {
	g.Meta           `path:"/savePreloadCfg" method:"post" summary:"保存预热配置" tags:"预热配置"`
	ID               uint64 `json:"id" dc:"配置ID,新建传0"`
	RecentLoginLimit int    `json:"recentLoginLimit" v:"required|min:1|max:10000#预热数量不能为空|预热数量需大于0|预热数量不能超过10000" dc:"最近登录用户预热数量"`
}

type SavePreloadCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
