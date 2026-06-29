package liveroomdto

import "github.com/gogf/gf/v2/frame/g"

const (
	ContributionRankPeriodToday  = 1 // 当天
	ContributionRankPeriodLast7  = 2 // 最近7天
	ContributionRankPeriodLast30 = 3 // 最近30天
)

// GetContributionRankReq App端查询直播间观众贡献榜
type GetContributionRankReq struct {
	g.Meta   `path:"/contributionRank" method:"post" summary:"查询直播间观众贡献榜" tags:"直播间"`
	RoomId   uint64 `json:"roomId"   v:"required#直播间ID不能为空" dc:"直播间ID"`
	Period   int    `json:"period"   v:"required|in:1,2,3#统计维度不能为空|统计维度无效" dc:"统计维度 1当天 2最近7天 3最近30天"`
	Page     int    `json:"page"     v:"min:1#页码从1开始" dc:"页码(从1开始,默认1)"`
	PageSize int    `json:"pageSize" v:"max:100#单页最多100条" dc:"每页数量(默认20,最大100)"`
}

// ContributionRankItem 观众贡献榜条目
type ContributionRankItem struct {
	Rank               int     `json:"rank"              dc:"排名"`
	UserId             string  `json:"userId"            dc:"观众用户ID"`
	Nickname           string  `json:"nickname"          dc:"昵称"`
	Avatar             string  `json:"avatar"            dc:"头像URL(已拼资源域名)"`
	ContributionAmount float64 `json:"contributionAmount" dc:"贡献总额(钻石,礼物+付费弹幕)"`
	VipLevel           uint32  `json:"vipLevel"          dc:"VIP等级"`
	Gender             uint8   `json:"gender"            dc:"性别(0未知,1男,2女)"`
	Age                int     `json:"age"               dc:"年龄(未设置出生日期时为0)"`
}

// GetContributionRankRes App端直播间观众贡献榜响应
type GetContributionRankRes struct {
	RoomId             string                  `json:"roomId"             dc:"直播间ID"`
	Period             int                     `json:"period"             dc:"统计维度"`
	MyRank             int                     `json:"myRank"             dc:"请求者在榜单中的排名,未上榜为-1"`
	ContributionAmount float64                 `json:"contributionAmount" dc:"请求者贡献总额(钻石,礼物+付费弹幕,未上榜为0)"`
	Total              int                     `json:"total"              dc:"榜单总数(最多500)"`
	Page               int                     `json:"page"               dc:"当前页码"`
	PageSize           int                     `json:"pageSize"           dc:"每页数量"`
	UpdatedAt          int64                   `json:"updatedAt"          dc:"查询时间(秒)"`
	List               []*ContributionRankItem `json:"list"               dc:"榜单列表"`
}
