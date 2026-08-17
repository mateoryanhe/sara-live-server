package guildsalarycfgdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type GuildSalaryCfgListReq struct {
	g.Meta `path:"/guildSalaryCfgList" method:"post" summary:"获取工会结算薪资分档列表" tags:"工会结算薪资配置"`
	httpserver.CMSQueryReq
}

type GuildSalaryCfgItem struct {
	ID                       string  `json:"id"`
	WeeklyWorkDays           uint64  `json:"weeklyWorkDays"`
	DailyLiveDurationMinutes uint64  `json:"dailyLiveDurationMinutes"`
	SalaryAmount             float64 `json:"salaryAmount"`
	Sort                     int     `json:"sort"`
	CreatedAt                string  `json:"createdAt"`
	UpdatedAt                string  `json:"updatedAt"`
}

type CreateGuildSalaryCfgReq struct {
	g.Meta                   `path:"/createGuildSalaryCfg" method:"post" summary:"创建工会结算薪资分档" tags:"工会结算薪资配置"`
	WeeklyWorkDays           uint64  `json:"weeklyWorkDays" dc:"每周工作天数门槛"`
	DailyLiveDurationMinutes uint64  `json:"dailyLiveDurationMinutes" dc:"每天直播时长门槛(分钟)"`
	SalaryAmount             float64 `json:"salaryAmount" dc:"薪资金额"`
	Sort                     int     `json:"sort" dc:"排序值(越大越靠前)"`
}

type CreateGuildSalaryCfgRes struct {
	ID string `json:"id"`
}

type UpdateGuildSalaryCfgReq struct {
	g.Meta                   `path:"/updateGuildSalaryCfg" method:"post" summary:"修改工会结算薪资分档" tags:"工会结算薪资配置"`
	ID                       uint64  `json:"id" v:"required#配置ID不能为空" dc:"配置ID"`
	WeeklyWorkDays           uint64  `json:"weeklyWorkDays" dc:"每周工作天数门槛"`
	DailyLiveDurationMinutes uint64  `json:"dailyLiveDurationMinutes" dc:"每天直播时长门槛(分钟)"`
	SalaryAmount             float64 `json:"salaryAmount" dc:"薪资金额"`
	Sort                     int     `json:"sort" dc:"排序值(越大越靠前)"`
}

type UpdateGuildSalaryCfgRes struct {
	Success bool `json:"success"`
}

type DeleteGuildSalaryCfgReq struct {
	g.Meta `path:"/deleteGuildSalaryCfg" method:"post" summary:"删除工会结算薪资分档" tags:"工会结算薪资配置"`
	ID     uint64 `json:"id" v:"required#配置ID不能为空" dc:"配置ID"`
}

type DeleteGuildSalaryCfgRes struct {
	Success bool `json:"success"`
}
