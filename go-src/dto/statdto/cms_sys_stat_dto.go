package statdto

import "github.com/gogf/gf/v2/frame/g"

// CMSSysStatReq CMS获取系统总数据
type CMSSysStatReq struct {
	g.Meta `path:"/getSysStat" method:"post" summary:"CMS获取系统总数据" tags:"仪表盘"`
}

// CMSSysStatRes CMS系统总数据
type CMSSysStatRes struct {
	TotalGold           float64 `json:"totalGold"`
	TotalRecharge       float64 `json:"totalRecharge"`
	TotalWithdraw       float64 `json:"totalWithdraw"`
	TotalRegisterUser   uint64  `json:"totalRegisterUser,string"`
	TotalGoldConsume    float64 `json:"totalGoldConsume"    dc:"金币总消费"`
	TotalDiamondConsume float64 `json:"totalDiamondConsume" dc:"钻石总消费"`
	TodayRecharge       float64 `json:"todayRecharge"       dc:"今日充值金额(USD)"`
	TodayGoldConsume    float64 `json:"todayGoldConsume"    dc:"今日金币消费金额"`
	TodayDiamondConsume float64 `json:"todayDiamondConsume" dc:"今日钻石消费金额"`
	TodayRegisterUser   uint64  `json:"todayRegisterUser,string" dc:"今日注册用户数"`
}
