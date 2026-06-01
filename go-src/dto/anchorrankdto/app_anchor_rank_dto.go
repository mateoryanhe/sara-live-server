package anchorrankdto

import "github.com/gogf/gf/v2/frame/g"

const (
	AnchorRankPeriodToday  = 1 // 当天
	AnchorRankPeriodLast7  = 2 // 最近7天
	AnchorRankPeriodLast30 = 3 // 最近30天
)

// AppAnchorRankListReq App端查询主播红人榜
type AppAnchorRankListReq struct {
	g.Meta   `path:"/appAnchorRankList" method:"post" summary:"App查询主播红人榜" tags:"主播红人榜"`
	Period   int `json:"period"   v:"required|in:1,2,3#统计维度不能为空|统计维度无效" dc:"统计维度 1当天 2最近7天 3最近30天"`
	Page     int `json:"page"     v:"min:1#页码从1开始" dc:"页码(从1开始,默认1)"`
	PageSize int `json:"pageSize" v:"max:100#单页最多100条" dc:"每页数量(默认20,最大100)"`
}

// AppAnchorRankItem 主播红人榜条目
type AppAnchorRankItem struct {
	Rank          int    `json:"rank"          dc:"排名"`
	UserId        string `json:"userId"        dc:"主播用户ID"`
	Nickname      string `json:"nickname"      dc:"主播昵称"`
	Avatar        string `json:"avatar"        dc:"主播头像URL"`
	RevenueAmount uint64 `json:"revenueAmount" dc:"收益总额(钻石)"`
}

// AppAnchorRankListRes App端主播红人榜响应
type AppAnchorRankListRes struct {
	Period    int                  `json:"period"    dc:"统计维度"`
	Total     int                  `json:"total"     dc:"榜单总数(最多500)"`
	Page      int                  `json:"page"      dc:"当前页码"`
	PageSize  int                  `json:"pageSize"  dc:"每页数量"`
	UpdatedAt int64                `json:"updatedAt" dc:"榜单刷新时间(秒)"`
	List      []*AppAnchorRankItem `json:"list"      dc:"榜单列表"`
}
