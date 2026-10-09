package metapixeldto

import "github.com/gogf/gf/v2/frame/g"

type GetMetaPixelCfgReq struct {
	g.Meta `path:"/getMetaPixelCfg" method:"post" summary:"查询Meta Pixel服务端上报配置" tags:"Meta Pixel配置"`
}

type MetaPixelCfgItem struct {
	ID            string `json:"id"`
	Enabled       uint8  `json:"enabled" dc:"0关闭1启用"`
	PixelId       string `json:"pixelId"`
	AccessToken   string `json:"accessToken"`
	TestEventCode string `json:"testEventCode"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

type GetMetaPixelCfgRes struct {
	Cfg *MetaPixelCfgItem `json:"cfg"`
}

type SaveMetaPixelCfgReq struct {
	g.Meta        `path:"/saveMetaPixelCfg" method:"post" summary:"保存Meta Pixel服务端上报配置" tags:"Meta Pixel配置"`
	ID            uint64 `json:"id" dc:"配置ID,新建传0"`
	Enabled       uint8  `json:"enabled" v:"in:0,1#启用状态仅支持0或1" dc:"0关闭1启用"`
	PixelId       string `json:"pixelId" v:"max-length:32#Pixel ID最长32字符" dc:"Meta Pixel ID"`
	AccessToken   string `json:"accessToken" dc:"Conversions API Access Token"`
	TestEventCode string `json:"testEventCode" v:"max-length:64#测试事件代码最长64字符" dc:"测试事件代码(可选)"`
}

type SaveMetaPixelCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
