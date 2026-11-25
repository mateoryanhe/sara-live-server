package shareappdownloaddto

import "github.com/gogf/gf/v2/frame/g"

type GetSofieShareAppDownloadCfgReq struct {
	g.Meta `path:"/getSofieShareAppDownloadCfg" method:"post" summary:"查询 Sofie 分享下载配置" tags:"Sofie分享下载"`
}

type SofieShareAppDownloadCfgItem struct {
	ID               string `json:"id"`
	PackageName      string `json:"packageName" dc:"适用包名,与 X-App-Package 一致时优先"`
	GooglePlayUrl    string `json:"googlePlayUrl"`
	ShortUrl         string `json:"shortUrl" dc:"只读:分享短链"`
	ShareText        string `json:"shareText"`
	PreviewShareCode string `json:"previewShareCode" dc:"只读:最终分享文案预览"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}

type GetSofieShareAppDownloadCfgRes struct {
	Cfg *SofieShareAppDownloadCfgItem `json:"cfg"`
}

type SaveSofieShareAppDownloadCfgReq struct {
	g.Meta        `path:"/saveSofieShareAppDownloadCfg" method:"post" summary:"保存 Sofie 分享下载配置" tags:"Sofie分享下载"`
	ID            uint64 `json:"id" dc:"配置ID,首次保存可为0"`
	PackageName   string `json:"packageName" v:"max-length:128#包名长度不能超过128"`
	GooglePlayUrl string `json:"googlePlayUrl" v:"required|max-length:512#请填写 Google Play 链接|链接长度不能超过512"`
	ShareText     string `json:"shareText" v:"max-length:4000#分享文案长度不能超过4000"`
}

type SaveSofieShareAppDownloadCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}