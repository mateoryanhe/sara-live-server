package gameconsumrankdto

import "github.com/gogf/gf/v2/frame/g"

const (
	GameConsumeRankPeriodToday  = 1 // 当天
	GameConsumeRankPeriodLast7  = 2 // 最近7天
	GameConsumeRankPeriodLast30 = 3 // 最近30天
)

// AppGameConsumeRankListReq App端查询游戏消费榜
type AppGameConsumeRankListReq struct {
	g.Meta   `path:"/appGameConsumeRankList" method:"post" summary:"App查询游戏消费榜" tags:"游戏消费榜"`
	Period   int `json:"period"   v:"required|in:1,2,3#统计维度不能为空|统计维度无效" dc:"统计维度 1当天 2最近7天 3最近30天"`
	Page     int `json:"page"     v:"min:1#页码从1开始" dc:"页码(从1开始,默认1)"`
	PageSize int `json:"pageSize" v:"max:100#单页最多100条" dc:"每页数量(默认20,最大100)"`
}

// AppGameConsumeRankItem 游戏消费榜条目
type AppGameConsumeRankItem struct {
	Rank          int     `json:"rank"          dc:"排名"`
	UserId        string  `json:"userId"        dc:"用户ID"`
	Nickname      string  `json:"nickname"      dc:"用户昵称"`
	Avatar        string  `json:"avatar"        dc:"用户头像URL"`
	ConsumeAmount float64 `json:"consumeAmount,string" dc:"游戏金币消费总额"`
	VipLevel      uint32  `json:"vipLevel"      dc:"VIP等级"`
	Gender        uint8   `json:"gender"        dc:"性别(0未知,1男,2女)"`
	Age           int     `json:"age"           dc:"年龄(未设置出生日期时为0)"`
}

// AppGameConsumeRankListRes App端游戏消费榜响应
type AppGameConsumeRankListRes struct {
	Period    int                       `json:"period"    dc:"统计维度"`
	MyRank    int                       `json:"myRank"    dc:"请求者在榜单中的排名,未上榜为-1"`
	Total     int                       `json:"total"     dc:"榜单总数(最多500)"`
	Page      int                       `json:"page"      dc:"当前页码"`
	PageSize  int                       `json:"pageSize"  dc:"每页数量"`
	UpdatedAt int64                     `json:"updatedAt" dc:"榜单刷新时间(秒)"`
	List      []*AppGameConsumeRankItem `json:"list"      dc:"榜单列表"`
}
