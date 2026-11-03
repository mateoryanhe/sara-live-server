package shareappdownloaddto

import "github.com/gogf/gf/v2/frame/g"

type GetShareAppDownloadCfgReq struct {
	g.Meta `path:"/getShareAppDownloadCfg" method:"post" summary:"查询分享App下载文案配置" tags:"分享下载文案"`
}

type ShareAppDownloadCfgItem struct {
	ID               string `json:"id"`
	DownloadUrl      string `json:"downloadUrl" dc:"只读:官网assets下体积最大apk的公开URL"`
	ShortUrl         string `json:"shortUrl" dc:"只读:分享短链,GET 302 到 downloadUrl"`
	ShareText        string `json:"shareText"`
	PreviewShareCode string `json:"previewShareCode" dc:"只读:下发给App的最终文案预览(邀请码为占位符)"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}

type GetShareAppDownloadCfgRes struct {
	Cfg *ShareAppDownloadCfgItem `json:"cfg"`
}

type SaveShareAppDownloadCfgReq struct {
	g.Meta    `path:"/saveShareAppDownloadCfg" method:"post" summary:"保存分享App下载文案配置" tags:"分享下载文案"`
	ID        uint64 `json:"id" dc:"配置ID,首次保存可为0"`
	ShareText string `json:"shareText" v:"max-length:4000#分享文案长度不能超过4000" dc:"分享文案"`
}

type SaveShareAppDownloadCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

// GetApkDownloadUrlReq App 查询 APK 下载链接(免登录). 可选 ch 匹配 CMS APK 推广渠道包.
type GetApkDownloadUrlReq struct {
	g.Meta `path:"/getApkDownloadUrl" method:"post" summary:"查询APK下载链接" tags:"分享下载"`
	Ch     string `json:"ch" dc:"推广渠道编码(可选,与 channel-url-list 渠道编码一致,如 \"2\");匹配 lookupApkDownload 渠道包"`
}

type GetApkDownloadUrlRes struct {
	DownloadUrl string `json:"downloadUrl" dc:"渠道包或官网assets最大apk公开URL,均无则空"`
}
