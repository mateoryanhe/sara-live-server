package anchorsalarycfgdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type AnchorSalaryCfgListReq struct {
	g.Meta `path:"/anchorSalaryCfgList" method:"post" summary:"获取主播结算薪资分档列表" tags:"主播结算薪资配置"`
	httpserver.CMSQueryReq
}

type AnchorSalaryCfgItem struct {
	ID                       string  `json:"id"`
	WeeklyWorkDays           uint64  `json:"weeklyWorkDays"`
	DailyLiveDurationMinutes uint64  `json:"dailyLiveDurationMinutes"`
	SalaryAmount             float64 `json:"salaryAmount"`
	Sort                     int     `json:"sort"`
	CreatedAt                string  `json:"createdAt"`
	UpdatedAt                string  `json:"updatedAt"`
}

type CreateAnchorSalaryCfgReq struct {
	g.Meta                   `path:"/createAnchorSalaryCfg" method:"post" summary:"创建主播结算薪资分档" tags:"主播结算薪资配置"`
	WeeklyWorkDays           uint64  `json:"weeklyWorkDays" dc:"每周工作天数门槛"`
	DailyLiveDurationMinutes uint64  `json:"dailyLiveDurationMinutes" dc:"每天直播时长门槛(分钟)"`
	SalaryAmount             float64 `json:"salaryAmount" dc:"薪资金额"`
	Sort                     int     `json:"sort" dc:"排序值(越大越靠前)"`
}

type CreateAnchorSalaryCfgRes struct {
	ID string `json:"id"`
}

type UpdateAnchorSalaryCfgReq struct {
	g.Meta                   `path:"/updateAnchorSalaryCfg" method:"post" summary:"修改主播结算薪资分档" tags:"主播结算薪资配置"`
	ID                       uint64  `json:"id" v:"required#配置ID不能为空" dc:"配置ID"`
	WeeklyWorkDays           uint64  `json:"weeklyWorkDays" dc:"每周工作天数门槛"`
	DailyLiveDurationMinutes uint64  `json:"dailyLiveDurationMinutes" dc:"每天直播时长门槛(分钟)"`
	SalaryAmount             float64 `json:"salaryAmount" dc:"薪资金额"`
	Sort                     int     `json:"sort" dc:"排序值(越大越靠前)"`
}

type UpdateAnchorSalaryCfgRes struct {
	Success bool `json:"success"`
}

type DeleteAnchorSalaryCfgReq struct {
	g.Meta `path:"/deleteAnchorSalaryCfg" method:"post" summary:"删除主播结算薪资分档" tags:"主播结算薪资配置"`
	ID     uint64 `json:"id" v:"required#配置ID不能为空" dc:"配置ID"`
}

type DeleteAnchorSalaryCfgRes struct {
	Success bool `json:"success"`
}
