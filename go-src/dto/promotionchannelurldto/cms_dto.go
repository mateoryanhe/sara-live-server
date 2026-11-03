package promotionchannelurldto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type PromotionChannelUrlListReq struct {
	g.Meta `path:"/promotionChannelUrlList" method:"post" summary:"推广渠道URL列表" tags:"推广"`
	httpserver.CMSQueryReq
	Key string `json:"key" dc:"搜索渠道名称或编码"`
}

type PromotionChannelUrlItem struct {
	ID             string `json:"id"`
	ChannelName    string `json:"channelName"`
	ChannelCode    int    `json:"channelCode"`
	H5AppChannelCode  string `json:"h5AppChannelCode" dc:"H5-App渠道码"`
	ShareURL       string `json:"shareUrl" dc:"分享链接(由H5直播域名+ch参数生成)"`
	URL            string `json:"url" dc:"兼容字段,与shareUrl一致"`
	ApkFileName    string `json:"apkFileName" dc:"渠道APK云桶对象名"`
	ApkDownloadURL string `json:"apkDownloadUrl" dc:"lookupApkDownload 渠道包URL"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

type PromotionChannelUrlListRes struct {
	Total         int                        `json:"total"`
	Data          []*PromotionChannelUrlItem `json:"data"`
	ShareBaseHost string                     `json:"shareBaseHost" dc:"H5直播部署域名(来自domain_site_mappings/h5-live)"`
}

type CreatePromotionChannelUrlReq struct {
	g.Meta      `path:"/createPromotionChannelUrl" method:"post" summary:"新增推广渠道URL" tags:"推广"`
	ChannelName      string `json:"channelName" v:"required#渠道名称不能为空" dc:"渠道名称"`
	ChannelCode      int    `json:"channelCode" v:"required|min:1#渠道编码不能为空|渠道编码须大于0" dc:"渠道编码(整数,唯一)"`
	H5AppChannelCode string `json:"h5AppChannelCode" v:"required#H5-App渠道码不能为空" dc:"H5-App渠道码(唯一,与APK渠道码配对)"`
}

type CreatePromotionChannelUrlRes struct {
	ID string `json:"id"`
}

type UpdatePromotionChannelUrlReq struct {
	g.Meta      `path:"/updatePromotionChannelUrl" method:"post" summary:"修改推广渠道URL" tags:"推广"`
	ID          uint64 `json:"id" v:"required#配置ID不能为空" dc:"配置ID"`
	ChannelName      string `json:"channelName" v:"required#渠道名称不能为空" dc:"渠道名称"`
	ChannelCode      int    `json:"channelCode" v:"required|min:1#渠道编码不能为空|渠道编码须大于0" dc:"渠道编码(整数,唯一)"`
	H5AppChannelCode string `json:"h5AppChannelCode" v:"required#H5-App渠道码不能为空" dc:"H5-App渠道码(唯一,与APK渠道码配对)"`
}

type UpdatePromotionChannelUrlRes struct {
	Success bool `json:"success"`
}

type DeletePromotionChannelUrlReq struct {
	g.Meta `path:"/deletePromotionChannelUrl" method:"post" summary:"删除推广渠道URL" tags:"推广"`
	ID     uint64 `json:"id" v:"required#配置ID不能为空" dc:"配置ID"`
}

type DeletePromotionChannelUrlRes struct {
	Success bool `json:"success"`
}

type PromotionChannelDailyStatListReq struct {
	g.Meta `path:"/promotionChannelDailyStatList" method:"post" summary:"推广URL日统计(deviceInfo.ch)" tags:"推广"`
	httpserver.CMSQueryReq
	StartDate   string `json:"startDate" v:"required|date-format:Y-m-d#开始日期不能为空|开始日期格式应为Y-m-d" dc:"开始日期(含)"`
	EndDate     string `json:"endDate" v:"required|date-format:Y-m-d#结束日期不能为空|结束日期格式应为Y-m-d" dc:"结束日期(含)"`
	ChannelCode int    `json:"channelCode" dc:"渠道编码(0=全部)"`
}

type PromotionChannelDailyStatItem struct {
	StatDate              string  `json:"statDate"`
	ChannelName           string  `json:"channelName"`
	ChannelCode           int     `json:"channelCode"`
	RegisterCount         uint64  `json:"registerCount"`
	OrderCount            uint64  `json:"orderCount"`
	PayUserCount          uint64  `json:"payUserCount"`
	OrderTotalAmount      float64 `json:"orderTotalAmount"`
	RechargeSuccessAmount float64 `json:"rechargeSuccessAmount"`
	H5DownloadApkCount    uint64  `json:"h5DownloadApkCount" dc:"H5下载APK次数"`
}

type PromotionChannelDailyStatListRes struct {
	Total int                              `json:"total"`
	Data  []*PromotionChannelDailyStatItem `json:"data"`
}

type InitPromotionChannelUrlApkUploadReq struct {
	g.Meta    `path:"/initPromotionChannelUrlApkUpload" method:"post" summary:"初始化推广渠道APK分片上传" tags:"推广"`
	ChannelID uint64 `json:"channelId" v:"required#渠道ID不能为空" dc:"推广渠道URL配置ID"`
	FileName  string `json:"fileName" v:"required#文件名不能为空" dc:"APK文件名"`
	FileSize  int64  `json:"fileSize" v:"required|min:1#文件大小无效" dc:"文件大小(字节)"`
}

type InitPromotionChannelUrlApkUploadRes struct {
	UploadId   string `json:"uploadId"`
	ChunkSize  int64  `json:"chunkSize"`
	TotalChunk int64  `json:"totalChunk"`
}

type CompletePromotionChannelUrlApkUploadReq struct {
	g.Meta    `path:"/completePromotionChannelUrlApkUpload" method:"post" summary:"完成推广渠道APK分片上传" tags:"推广"`
	UploadId  string `json:"uploadId" v:"required#uploadId不能为空"`
	ChannelID uint64 `json:"channelId" v:"required#渠道ID不能为空"`
}

type CompletePromotionChannelUrlApkUploadRes struct {
	ApkDownloadURL string `json:"apkDownloadUrl"`
	ApkFileName    string `json:"apkFileName"`
}

type AbortPromotionChannelUrlApkUploadReq struct {
	g.Meta   `path:"/abortPromotionChannelUrlApkUpload" method:"post" summary:"取消推广渠道APK分片上传" tags:"推广"`
	UploadId string `json:"uploadId" v:"required#uploadId不能为空"`
}

type AbortPromotionChannelUrlApkUploadRes struct{}

type H5AppPromotionDailyStatListReq struct {
	g.Meta `path:"/h5AppPromotionDailyStatList" method:"post" summary:"H5-App推广日统计(apkCh配对URL渠道)" tags:"推广"`
	httpserver.CMSQueryReq
	StartDate   string `json:"startDate" v:"required|date-format:Y-m-d#开始日期不能为空|开始日期格式应为Y-m-d" dc:"开始日期(含)"`
	EndDate     string `json:"endDate" v:"required|date-format:Y-m-d#结束日期不能为空|结束日期格式应为Y-m-d" dc:"结束日期(含)"`
	ChannelCode int    `json:"channelCode" dc:"(兼容)推广URL整数编码(0=全部)"`
	H5AppChannelCode string `json:"h5AppChannelCode" dc:"H5-App渠道码(空=全部)"`
}

type H5AppPromotionDailyStatItem struct {
	StatDate              string  `json:"statDate"`
	ChannelName           string  `json:"channelName"`
	H5AppChannelCode      string  `json:"h5AppChannelCode"`
	RegisterCount         uint64  `json:"registerCount"`
	OrderCount            uint64  `json:"orderCount"`
	PayUserCount          uint64  `json:"payUserCount"`
	OrderTotalAmount      float64 `json:"orderTotalAmount"`
	RechargeSuccessAmount float64 `json:"rechargeSuccessAmount"`
}

type H5AppPromotionDailyStatListRes struct {
	Total int                            `json:"total"`
	Data  []*H5AppPromotionDailyStatItem `json:"data"`
}
