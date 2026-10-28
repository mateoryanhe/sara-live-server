package userappsubmissiondto

import (
	"github.com/gogf/gf/v2/frame/g"
)

type CMSListUserFeedbackRecordReq struct {
	g.Meta    `path:"/listUserFeedbackRecord" method:"post" summary:"CMS用户反馈记录列表" tags:"用户反馈"`
	UserId    uint64 `json:"userId" dc:"用户ID,可选"`
	StartTime int64  `json:"startTime" dc:"开始时间 Unix 秒"`
	EndTime   int64  `json:"endTime" dc:"结束时间 Unix 秒"`
	PageIndex int    `json:"pageIndex" dc:"页码"`
	PageSize  int    `json:"pageSize" dc:"每页条数"`
}

type CMSUserFeedbackRecordItem struct {
	ID        uint64 `json:"id"`
	UserId    uint64 `json:"userId"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

type CMSListUserReportRecordReq struct {
	g.Meta    `path:"/listUserReportRecord" method:"post" summary:"CMS用户举报记录列表" tags:"用户举报"`
	UserId    uint64 `json:"userId" dc:"用户ID,可选"`
	StartTime int64  `json:"startTime" dc:"开始时间 Unix 秒"`
	EndTime   int64  `json:"endTime" dc:"结束时间 Unix 秒"`
	PageIndex int    `json:"pageIndex" dc:"页码"`
	PageSize  int    `json:"pageSize" dc:"每页条数"`
}

type CMSUserReportRecordItem struct {
	ID        uint64 `json:"id"`
	UserId    uint64 `json:"userId"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}
