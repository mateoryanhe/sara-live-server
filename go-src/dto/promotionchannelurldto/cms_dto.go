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
	ID            string `json:"id"`
	ChannelName   string `json:"channelName"`
	ChannelCode   int    `json:"channelCode"`
	ShareURL      string `json:"shareUrl" dc:"分享链接(由H5直播域名+ch参数生成)"`
	URL           string `json:"url" dc:"兼容字段,与shareUrl一致"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

type PromotionChannelUrlListRes struct {
	Total         int                        `json:"total"`
	Data          []*PromotionChannelUrlItem `json:"data"`
	ShareBaseHost string                     `json:"shareBaseHost" dc:"H5直播部署域名(来自domain_site_mappings/h5-live)"`
}

type CreatePromotionChannelUrlReq struct {
	g.Meta      `path:"/createPromotionChannelUrl" method:"post" summary:"新增推广渠道URL" tags:"推广"`
	ChannelName string `json:"channelName" v:"required#渠道名称不能为空" dc:"渠道名称"`
	ChannelCode int    `json:"channelCode" v:"required|min:1#渠道编码不能为空|渠道编码须大于0" dc:"渠道编码(整数,唯一)"`
}

type CreatePromotionChannelUrlRes struct {
	ID string `json:"id"`
}

type UpdatePromotionChannelUrlReq struct {
	g.Meta      `path:"/updatePromotionChannelUrl" method:"post" summary:"修改推广渠道URL" tags:"推广"`
	ID          uint64 `json:"id" v:"required#配置ID不能为空" dc:"配置ID"`
	ChannelName string `json:"channelName" v:"required#渠道名称不能为空" dc:"渠道名称"`
	ChannelCode int    `json:"channelCode" v:"required|min:1#渠道编码不能为空|渠道编码须大于0" dc:"渠道编码(整数,唯一)"`
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
	g.Meta `path:"/promotionChannelDailyStatList" method:"post" summary:"推广渠道日统计报表" tags:"推广"`
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
