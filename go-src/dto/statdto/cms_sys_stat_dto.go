package statdto

import "github.com/gogf/gf/v2/frame/g"

// CMSSysStatReq CMS获取系统总数据
type CMSSysStatReq struct {
	g.Meta `path:"/getSysStat" method:"post" summary:"CMS获取系统总数据" tags:"仪表盘"`
}

// CMSSysStatRes CMS系统总数据
type CMSSysStatRes struct {
	TotalGold         float64 `json:"totalGold"`
	TotalRecharge     float64 `json:"totalRecharge"`
	TotalWithdraw     float64 `json:"totalWithdraw"`
	TotalRegisterUser uint64  `json:"totalRegisterUser,string"`
	TodayRegisterUser uint64  `json:"todayRegisterUser,string"`
}
