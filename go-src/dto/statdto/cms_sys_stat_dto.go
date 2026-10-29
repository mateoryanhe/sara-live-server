package statdto

import "github.com/gogf/gf/v2/frame/g"

// CMSSysStatReq CMS获取系统总数据
type CMSSysStatReq struct {
	g.Meta               `path:"/getSysStat" method:"post" summary:"CMS获取系统总数据" tags:"仪表盘"`
	BypassDashboardCache bool `json:"bypassDashboardCache" dc:"true 时跳过仪表盘读缓存, 重新查库/组装"`
}

// CMSSysStatRes CMS系统总数据
type CMSSysStatRes struct {
	TotalGold                 float64 `json:"totalGold"`
	TotalRecharge             float64 `json:"totalRecharge" dc:"全部美金入账累计(真实USD)"`
	TotalNormalUserRecharge   float64 `json:"totalNormalUserRecharge" dc:"普通用户美金入账累计(真实USD)"`
	TotalCoinMerchantRecharge float64 `json:"totalCoinMerchantRecharge" dc:"币商美金入账累计(真实USD)"`
	TotalAnchorPayout         float64 `json:"totalAnchorPayout" dc:"主播代付累计(USD)"`
	WeeklyAnchorPayout        float64 `json:"weeklyAnchorPayout" dc:"本周主播代付金额(USD)"`
	TotalVirtualRecharge      float64 `json:"totalVirtualRecharge" dc:"虚拟美金累计(充值白名单)"`
	TotalWithdraw             float64 `json:"totalWithdraw"`
	TotalRegisterUser         uint64  `json:"totalRegisterUser,string"`
	TotalGoldConsume          float64 `json:"totalGoldConsume"    dc:"金币总消费"`
	TotalDiamondConsume       float64 `json:"totalDiamondConsume" dc:"钻石总消费"`
	TodayRecharge             float64 `json:"todayRecharge"       dc:"今日充值金额(USD)"`
	TodayGoldConsume          float64 `json:"todayGoldConsume"    dc:"今日金币消费金额"`
	TodayDiamondConsume       float64 `json:"todayDiamondConsume" dc:"今日钻石消费金额"`
	TodayRegisterUser         uint64  `json:"todayRegisterUser,string" dc:"今日注册用户数"`
	OnlineCount               uint64  `json:"onlineCount,string" dc:"在线人数(TCP连接)"`
	GameOnlineCount           uint64  `json:"gameOnlineCount,string" dc:"游戏在线人数(瞬时)"`
	GameBetUserCount          uint64  `json:"gameBetUserCount,string" dc:"游戏下注人数(当前分钟累计去重)"`
	AudienceOnlineCount       uint64  `json:"audienceOnlineCount,string" dc:"观众在线人数"`
	LiveAnchorCount           uint64  `json:"liveAnchorCount,string" dc:"主播在线人数(开播中)"`
}
