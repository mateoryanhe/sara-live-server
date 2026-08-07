package gamebetdto

import "github.com/gogf/gf/v2/frame/g"

// AppGameConsumeRankReq App 分页查询单场直播游戏消费榜
type AppGameConsumeRankReq struct {
	g.Meta       `path:"/appGameConsumeRank" method:"post" summary:"App分页查询单场直播游戏消费榜" tags:"游戏"`
	LiveRecordId uint64 `json:"liveRecordId" v:"required#直播记录ID不能为空" dc:"直播记录ID"`
	Page         int    `json:"page"         v:"min:1#页码从1开始" dc:"页码(从1开始,默认1)"`
	PageSize     int    `json:"pageSize"     v:"max:100#单页最多100条" dc:"每页数量(默认20,最大100)"`
}

// AppGameConsumeRankItem 游戏消费榜条目
type AppGameConsumeRankItem struct {
	Rank          int     `json:"rank"          dc:"排名"`
	UserId        string  `json:"userId"        dc:"用户ID"`
	Nickname      string  `json:"nickname"      dc:"昵称"`
	Avatar        string  `json:"avatar"        dc:"头像URL(已拼资源域名)"`
	ConsumeAmount float64 `json:"consumeAmount" dc:"本场直播游戏消费总额(金币下注)"`
	VipLevel      uint32  `json:"vipLevel"      dc:"VIP等级"`
	Gender        uint8   `json:"gender"        dc:"性别(0未知,1男,2女)"`
	Age           int     `json:"age"           dc:"年龄(未设置出生日期时为0)"`
}

// AppGameConsumeRankRes App 单场直播游戏消费榜响应
type AppGameConsumeRankRes struct {
	LiveRecordId  string                    `json:"liveRecordId"   dc:"直播记录ID"`
	MyRank        int                       `json:"myRank"         dc:"请求者在榜单中的排名,未上榜为-1"`
	ConsumeAmount float64                   `json:"consumeAmount"  dc:"请求者消费总额(未上榜为0)"`
	Total         int                       `json:"total"          dc:"榜单总数(最多500)"`
	Page          int                       `json:"page"           dc:"当前页码"`
	PageSize      int                       `json:"pageSize"       dc:"每页数量"`
	UpdatedAt     int64                     `json:"updatedAt"      dc:"查询时间(秒)"`
	List          []*AppGameConsumeRankItem `json:"list"           dc:"榜单列表"`
}
