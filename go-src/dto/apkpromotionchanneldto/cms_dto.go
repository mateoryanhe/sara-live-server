package apkpromotionchanneldto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type ApkPromotionChannelListReq struct {
	g.Meta `path:"/apkPromotionChannelList" method:"post" summary:"APK推广渠道列表" tags:"推广"`
	httpserver.CMSQueryReq
	Key string `json:"key" dc:"搜索渠道名称或编码"`
}

type ApkPromotionChannelItem struct {
	ID             string `json:"id"`
	ChannelName    string `json:"channelName"`
	ChannelCode    string `json:"channelCode"`
	ApkFileName    string `json:"apkFileName"`
	ApkDownloadURL string `json:"apkDownloadUrl"`
	ShareURL       string `json:"shareUrl" dc:"渠道专用下载短链"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

type ApkPromotionChannelListRes struct {
	Total         int                        `json:"total"`
	Data          []*ApkPromotionChannelItem `json:"data"`
	ShareBaseHost string                     `json:"shareBaseHost" dc:"官网静态域名(短链域名)"`
}

type CreateApkPromotionChannelReq struct {
	g.Meta      `path:"/createApkPromotionChannel" method:"post" summary:"新增APK推广渠道" tags:"推广"`
	ChannelName string `json:"channelName" v:"required#渠道名称不能为空" dc:"渠道名称"`
	ChannelCode string `json:"channelCode" v:"required#渠道编码不能为空" dc:"渠道编码(字符串,唯一)"`
}

type CreateApkPromotionChannelRes struct {
	ID string `json:"id"`
}

type UpdateApkPromotionChannelReq struct {
	g.Meta      `path:"/updateApkPromotionChannel" method:"post" summary:"修改APK推广渠道" tags:"推广"`
	ID          uint64 `json:"id" v:"required#配置ID不能为空" dc:"配置ID"`
	ChannelName string `json:"channelName" v:"required#渠道名称不能为空" dc:"渠道名称"`
	ChannelCode string `json:"channelCode" v:"required#渠道编码不能为空" dc:"渠道编码(字符串,唯一)"`
}

type UpdateApkPromotionChannelRes struct {
	Success bool `json:"success"`
}

type DeleteApkPromotionChannelReq struct {
	g.Meta `path:"/deleteApkPromotionChannel" method:"post" summary:"删除APK推广渠道" tags:"推广"`
	ID     uint64 `json:"id" v:"required#配置ID不能为空" dc:"配置ID"`
}

type DeleteApkPromotionChannelRes struct {
	Success bool `json:"success"`
}

type InitApkPromotionChannelApkUploadReq struct {
	g.Meta    `path:"/initApkPromotionChannelApkUpload" method:"post" summary:"初始化渠道APK分片上传" tags:"推广"`
	ChannelID uint64 `json:"channelId" v:"required#渠道ID不能为空" dc:"渠道配置ID"`
	FileName  string `json:"fileName" v:"required#文件名不能为空" dc:"APK文件名"`
	FileSize  int64  `json:"fileSize" v:"required|min:1#文件大小无效" dc:"文件大小(字节)"`
}

type InitApkPromotionChannelApkUploadRes struct {
	UploadId   string `json:"uploadId"`
	ChunkSize  int64  `json:"chunkSize"`
	TotalChunk int64  `json:"totalChunk"`
}

type CompleteApkPromotionChannelApkUploadReq struct {
	g.Meta    `path:"/completeApkPromotionChannelApkUpload" method:"post" summary:"完成渠道APK分片上传" tags:"推广"`
	UploadId  string `json:"uploadId" v:"required#uploadId不能为空" dc:"上传会话ID"`
	ChannelID uint64 `json:"channelId" v:"required#渠道ID不能为空" dc:"渠道配置ID"`
}

type CompleteApkPromotionChannelApkUploadRes struct {
	ApkDownloadURL string `json:"apkDownloadUrl"`
	ApkFileName    string `json:"apkFileName"`
}

type AbortApkPromotionChannelApkUploadReq struct {
	g.Meta   `path:"/abortApkPromotionChannelApkUpload" method:"post" summary:"取消渠道APK分片上传" tags:"推广"`
	UploadId string `json:"uploadId" v:"required#uploadId不能为空" dc:"上传会话ID"`
}

type AbortApkPromotionChannelApkUploadRes struct{}

type ApkPromotionChannelDailyStatListReq struct {
	g.Meta `path:"/apkPromotionChannelDailyStatList" method:"post" summary:"APK推广渠道日统计报表" tags:"推广"`
	httpserver.CMSQueryReq
	StartDate   string `json:"startDate" v:"required|date-format:Y-m-d#开始日期不能为空|开始日期格式应为Y-m-d" dc:"开始日期(含)"`
	EndDate     string `json:"endDate" v:"required|date-format:Y-m-d#结束日期不能为空|结束日期格式应为Y-m-d" dc:"结束日期(含)"`
	ChannelCode string `json:"channelCode" dc:"渠道编码(空=全部)"`
}

type ApkPromotionChannelDailyStatItem struct {
	StatDate              string  `json:"statDate"`
	ChannelName           string  `json:"channelName"`
	ChannelCode           string  `json:"channelCode"`
	RegisterCount         uint64  `json:"registerCount"`
	DownloadClickCount    uint64  `json:"downloadClickCount"`
	OrderCount            uint64  `json:"orderCount"`
	PayUserCount          uint64  `json:"payUserCount"`
	OrderTotalAmount      float64 `json:"orderTotalAmount"`
	RechargeSuccessAmount float64 `json:"rechargeSuccessAmount"`
}

type ApkPromotionChannelDailyStatListRes struct {
	Total int                                 `json:"total"`
	Data  []*ApkPromotionChannelDailyStatItem `json:"data"`
}
