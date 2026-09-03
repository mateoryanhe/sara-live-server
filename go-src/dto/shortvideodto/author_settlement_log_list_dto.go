package shortvideodto

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// CMSAuthorSettlementLogListReq CMS 分页查询短视频作者结算日志
type CMSAuthorSettlementLogListReq struct {
	g.Meta `path:"/cmsAuthorSettlementLogList" method:"post" summary:"CMS查询短视频作者结算日志" tags:"短视频"`
	httpserver.CMSQueryReq
	UserId    string `json:"userId" dc:"用户ID(可选,留空查全部)"`
	StartTime int64  `json:"startTime" dc:"开始时间(Unix秒,可选)"`
	EndTime   int64  `json:"endTime" dc:"结束时间(Unix秒,可选)"`
}

// CMSAuthorSettlementLogItem CMS 短视频作者结算日志列表项
type CMSAuthorSettlementLogItem struct {
	Id                 uint64     `json:"id,string"`
	UserId             uint64     `json:"userId,string"`
	UnsettledIncome    float64    `json:"unsettledIncome"`
	SettlementDiamond  float64    `json:"settlementDiamond"`
	AnchorSharePercent float64    `json:"anchorSharePercent"`
	CreatedAt          *time.Time `json:"createdAt"`
}
