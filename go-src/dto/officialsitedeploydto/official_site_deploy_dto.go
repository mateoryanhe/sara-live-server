package officialsitedeploydto

import "github.com/gogf/gf/v2/frame/g"

const (
	OfficialSiteStaticPrefix         = "/web"
	ThirdPayOfficialSiteStaticPrefix = "/hom"
	OfficialSiteSiteKey              = "official-site"
	ThirdPayOfficialSiteSiteKey      = "third-pay-official-site"
	DefaultOfficialSiteDeployPath    = "/home/ec2-user/cdn/official-site"
	DefaultThirdPayOfficialSitePath  = "/home/ec2-user/cdn/third-pay-official-site"
)

type SiteDeployInfoItem struct {
	Domain       string `json:"domain" dc:"站点域名"`
	UrlPrefix    string `json:"urlPrefix" dc:"静态访问前缀"`
	DeployPath   string `json:"deployPath" dc:"解压目标物理目录"`
	AcceptExt    string `json:"acceptExt" dc:"允许上传的扩展名"`
	LastUploadAt string `json:"lastUploadAt" dc:"最后一次成功上传并解压时间"`
}

type GetOfficialSiteDeployInfoReq struct {
	g.Meta `path:"/getOfficialSiteDeployInfo" method:"post" summary:"获取官网部署信息" tags:"官网部署"`
}

type GetOfficialSiteDeployInfoRes struct {
	Info *SiteDeployInfoItem `json:"info"`
}

type SaveOfficialSiteDeployCfgReq struct {
	g.Meta     `path:"/saveOfficialSiteDeployCfg" method:"post" summary:"保存官网部署配置" tags:"官网部署"`
	Domain     string `json:"domain" v:"required#域名不能为空" dc:"官网站点域名,不含协议与路径"`
	DeployPath string `json:"deployPath" v:"required#部署目录不能为空" dc:"官网静态文件绝对目录"`
}

type SaveOfficialSiteDeployCfgRes struct {
	Success bool `json:"success"`
}

type GetThirdPayOfficialSiteDeployInfoReq struct {
	g.Meta `path:"/getThirdPayOfficialSiteDeployInfo" method:"post" summary:"获取第三方支付官网部署信息" tags:"第三方支付官网部署"`
}

type GetThirdPayOfficialSiteDeployInfoRes struct {
	Info *SiteDeployInfoItem `json:"info"`
}

type SaveThirdPayOfficialSiteDeployCfgReq struct {
	g.Meta     `path:"/saveThirdPayOfficialSiteDeployCfg" method:"post" summary:"保存第三方支付官网部署配置" tags:"第三方支付官网部署"`
	Domain     string `json:"domain" v:"required#域名不能为空" dc:"第三方支付官网域名,不含协议与路径"`
	DeployPath string `json:"deployPath" v:"required#部署目录不能为空" dc:"第三方支付官网静态文件绝对目录"`
}

type SaveThirdPayOfficialSiteDeployCfgRes struct {
	Success bool `json:"success"`
}

type DeploySiteZipRes struct {
	FileCount  int    `json:"fileCount" dc:"解压写入的文件数"`
	DirCount   int    `json:"dirCount" dc:"创建的目录数"`
	DeployPath string `json:"deployPath" dc:"解压目标物理目录"`
	UrlPrefix  string `json:"urlPrefix" dc:"静态访问前缀"`
}

type DeployOfficialSiteZipReq struct {
	g.Meta `path:"/deployZip" method:"post" mime:"multipart/form-data" summary:"上传zip并部署官网静态资源" tags:"官网部署"`
}

type DeployThirdPayOfficialSiteZipReq struct {
	g.Meta `path:"/deployZip" method:"post" mime:"multipart/form-data" summary:"上传zip并部署第三方支付官网静态资源" tags:"第三方支付官网部署"`
}

type InitOfficialSiteFileUploadReq struct {
	g.Meta   `path:"/initFileUpload" method:"post" summary:"初始化官网文件分片上传" tags:"官网部署"`
	FileName string `json:"fileName" dc:"文件名称，仅支持.zip或.apk"`
	FileSize int64  `json:"fileSize" dc:"文件字节数"`
}

type InitOfficialSiteFileUploadRes struct {
	UploadId    string `json:"uploadId" dc:"上传会话ID"`
	ChunkSize   int64  `json:"chunkSize" dc:"每个分片的字节数"`
	TotalChunks int64  `json:"totalChunks" dc:"分片总数"`
}

type UploadOfficialSiteFileChunkRes struct {
	ChunkIndex int64 `json:"chunkIndex" dc:"分片序号，从0开始"`
	FileSize   int64 `json:"fileSize" dc:"本次写入的字节数"`
}

type CompleteOfficialSiteFileUploadReq struct {
	g.Meta   `path:"/completeFileUpload" method:"post" summary:"合并官网文件分片并完成部署" tags:"官网部署"`
	UploadId string `json:"uploadId" dc:"上传会话ID"`
}

type CompleteOfficialSiteFileUploadRes struct {
	FileType    string `json:"fileType" dc:"文件类型：zip或apk"`
	FileName    string `json:"fileName" dc:"文件名称"`
	FileSize    int64  `json:"fileSize" dc:"文件字节数"`
	FileCount   int    `json:"fileCount" dc:"ZIP解压写入的文件数"`
	DirCount    int    `json:"dirCount" dc:"ZIP创建的目录数"`
	DeployPath  string `json:"deployPath" dc:"部署目录或APK文件路径"`
	UrlPrefix   string `json:"urlPrefix" dc:"静态访问前缀"`
	DownloadUrl string `json:"downloadUrl" dc:"APK下载地址"`
}

type AbortOfficialSiteFileUploadReq struct {
	g.Meta   `path:"/abortFileUpload" method:"post" summary:"取消官网文件分片上传" tags:"官网部署"`
	UploadId string `json:"uploadId" dc:"上传会话ID"`
}

type AbortOfficialSiteFileUploadRes struct{}

type InitThirdPayOfficialSiteFileUploadReq struct {
	g.Meta   `path:"/initFileUpload" method:"post" summary:"初始化第三方支付官网文件分片上传" tags:"第三方支付官网部署"`
	FileName string `json:"fileName" dc:"文件名称，仅支持.zip或.apk"`
	FileSize int64  `json:"fileSize" dc:"文件字节数"`
}

type InitThirdPayOfficialSiteFileUploadRes = InitOfficialSiteFileUploadRes

type CompleteThirdPayOfficialSiteFileUploadReq struct {
	g.Meta   `path:"/completeFileUpload" method:"post" summary:"合并第三方支付官网文件分片并完成部署" tags:"第三方支付官网部署"`
	UploadId string `json:"uploadId" dc:"上传会话ID"`
}

type CompleteThirdPayOfficialSiteFileUploadRes = CompleteOfficialSiteFileUploadRes

type AbortThirdPayOfficialSiteFileUploadReq struct {
	g.Meta   `path:"/abortFileUpload" method:"post" summary:"取消第三方支付官网文件分片上传" tags:"第三方支付官网部署"`
	UploadId string `json:"uploadId" dc:"上传会话ID"`
}

type AbortThirdPayOfficialSiteFileUploadRes = AbortOfficialSiteFileUploadRes
