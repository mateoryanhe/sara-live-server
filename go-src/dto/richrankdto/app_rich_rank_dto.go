package richrankdto

import "github.com/gogf/gf/v2/frame/g"

const (
	RichRankPeriodToday  = 1 // 当天
	RichRankPeriodLast7  = 2 // 最近7天
	RichRankPeriodLast30 = 3 // 最近30天
)

// AppRichRankListReq App端查询富豪榜
type AppRichRankListReq struct {
	g.Meta   `path:"/appRichRankList" method:"post" summary:"App查询富豪榜" tags:"富豪榜"`
	Period   int `json:"period"   v:"required|in:1,2,3#统计维度不能为空|统计维度无效" dc:"统计维度 1当天 2最近7天 3最近30天"`
	Page     int `json:"page"     v:"min:1#页码从1开始" dc:"页码(从1开始,默认1)"`
	PageSize int `json:"pageSize" v:"max:100#单页最多100条" dc:"每页数量(默认20,最大100)"`
}

// AppRichRankItem 富豪榜条目
type AppRichRankItem struct {
	Rank          int     `json:"rank"          dc:"排名"`
	UserId        string  `json:"userId"        dc:"用户ID"`
	Nickname      string  `json:"nickname"      dc:"用户昵称"`
	Avatar        string  `json:"avatar"        dc:"用户头像URL"`
	ConsumeAmount float64 `json:"consumeAmount,string" dc:"钻石消费总额"`
}

// AppRichRankListRes App端富豪榜响应
type AppRichRankListRes struct {
	Period    int                `json:"period"    dc:"统计维度"`
	Total     int                `json:"total"     dc:"榜单总数(最多500)"`
	Page      int                `json:"page"      dc:"当前页码"`
	PageSize  int                `json:"pageSize"  dc:"每页数量"`
	UpdatedAt int64              `json:"updatedAt" dc:"榜单刷新时间(秒)"`
	List      []*AppRichRankItem `json:"list"      dc:"榜单列表"`
}
