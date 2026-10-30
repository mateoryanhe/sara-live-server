package preloadcfgdto

import "github.com/gogf/gf/v2/frame/g"

type GetPreloadCfgReq struct {
	g.Meta `path:"/getPreloadCfg" method:"post" summary:"查询服务器运行配置" tags:"服务器运行配置"`
}

type PreloadCfgItem struct {
	ID                              string  `json:"id"`
	RecentLoginLimit                int     `json:"recentLoginLimit"`
	HotRestartAuth                  string  `json:"hotRestartAuth"`
	HotRestartFlushTimeout          int     `json:"hotRestartFlushTimeout"`
	HotRestartExitTimeout           int     `json:"hotRestartExitTimeout"`
	MemoryLimitM                    int     `json:"memoryLimitM"`
	IpGeoDbPath                     string  `json:"ipGeoDbPath"`
	ResourceMetricCollectionEnabled bool    `json:"resourceMetricCollectionEnabled"`
	CreatedAt                       string  `json:"createdAt"`
	UpdatedAt                       string  `json:"updatedAt"`
}

type GetPreloadCfgRes struct {
	Cfg *PreloadCfgItem `json:"cfg"`
}

type SavePreloadCfgReq struct {
	g.Meta                          `path:"/savePreloadCfg" method:"post" summary:"保存服务器运行配置" tags:"服务器运行配置"`
	ID                              uint64  `json:"id" dc:"配置ID,新建传0"`
	RecentLoginLimit                int     `json:"recentLoginLimit" v:"required|min:1|max:10000#预热数量不能为空|预热数量需大于0|预热数量不能超过10000" dc:"最近登录用户预热数量"`
	HotRestartAuth                  string  `json:"hotRestartAuth" v:"required|length:8,128#热重启密钥不能为空|热重启密钥长度需在8-128之间" dc:"热重启接口密钥"`
	HotRestartFlushTimeout          int     `json:"hotRestartFlushTimeout" v:"required|min:1|max:3600#第一阶段时间不能为空|第一阶段时间最小1秒|第一阶段时间最大3600秒" dc:"热重启第一阶段固定等待秒"`
	HotRestartExitTimeout           int     `json:"hotRestartExitTimeout" v:"required|min:1|max:3600#第二阶段时间不能为空|第二阶段时间最小1秒|第二阶段时间最大3600秒" dc:"热重启第二阶段固定等待秒"`
	MemoryLimitM                    int     `json:"memoryLimitM" v:"required|min:64|max:32768#内存上限不能为空|内存上限最小64MB|内存上限最大32768MB" dc:"Go堆内存软上限MB"`
	IpGeoDbPath                     string  `json:"ipGeoDbPath" v:"required|length:1,512#IP库路径不能为空|IP库路径过长" dc:"GeoLite2-Country.mmdb绝对路径"`
	ResourceMetricCollectionEnabled bool    `json:"resourceMetricCollectionEnabled" dc:"是否采集服务器资源指标"`
}

type SavePreloadCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
