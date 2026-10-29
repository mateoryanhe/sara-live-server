package trackingdto

import "github.com/gogf/gf/v2/frame/g"

type ReportH5DownloadReq struct {
	g.Meta      `path:"/reportH5Download" method:"post" summary:"上报H5下载点击" tags:"埋点"`
	ChannelCode int `json:"channelCode" dc:"推广渠道编码(可选,>0 时计入推广渠道日统计 H5下载APK量)"`
}

type ReportH5DownloadRes struct {
	Success bool `json:"success"`
}

type ReportOfficialSiteDownloadReq struct {
	g.Meta `path:"/reportOfficialSiteDownload" method:"post" summary:"上报官网下载点击(免登录)" tags:"埋点"`
}

type ReportOfficialSiteDownloadRes struct {
	Success bool `json:"success"`
}

type CMSH5DownloadTrendReq struct {
	g.Meta `path:"/getH5DownloadTrend" method:"post" summary:"CMS H5下载量趋势" tags:"埋点"`
}

type CMSH5DownloadTrendRes struct {
	TodayCount uint64                        `json:"todayCount"`
	WeekCount  uint64                        `json:"weekCount"`
	MonthCount uint64                        `json:"monthCount"`
	Daily      []*CMSTrackingEventTrendPoint `json:"daily"`
	Weekly     []*CMSTrackingEventTrendPoint `json:"weekly"`
	Monthly    []*CMSTrackingEventTrendPoint `json:"monthly"`
}

type CMSOfficialSiteDownloadTrendReq struct {
	g.Meta `path:"/getOfficialSiteDownloadTrend" method:"post" summary:"CMS官网下载量趋势" tags:"埋点"`
}

type CMSOfficialSiteDownloadTrendRes struct {
	TodayCount uint64                        `json:"todayCount"`
	WeekCount  uint64                        `json:"weekCount"`
	MonthCount uint64                        `json:"monthCount"`
	Daily      []*CMSTrackingEventTrendPoint `json:"daily"`
	Weekly     []*CMSTrackingEventTrendPoint `json:"weekly"`
	Monthly    []*CMSTrackingEventTrendPoint `json:"monthly"`
}

type CMSShareShortLinkClickTrendReq struct {
	g.Meta `path:"/getShareShortLinkClickTrend" method:"post" summary:"CMS分享短链点击量趋势" tags:"埋点"`
}

type CMSShareShortLinkClickTrendRes struct {
	TodayCount uint64                        `json:"todayCount"`
	WeekCount  uint64                        `json:"weekCount"`
	MonthCount uint64                        `json:"monthCount"`
	Daily      []*CMSTrackingEventTrendPoint `json:"daily"`
	Weekly     []*CMSTrackingEventTrendPoint `json:"weekly"`
	Monthly    []*CMSTrackingEventTrendPoint `json:"monthly"`
}
