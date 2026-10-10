package statdto

import "github.com/gogf/gf/v2/frame/g"

// CMSResetBasicDashboardStatsReq CMS 重置仪表盘「基础数据」统计
type CMSResetBasicDashboardStatsReq struct {
	g.Meta `path:"/resetBasicDashboardStats" method:"post" summary:"CMS重置基础数据统计" tags:"仪表盘"`
}

// CMSResetBasicDashboardStatsRes 重置结果
type CMSResetBasicDashboardStatsRes struct {
	Success bool `json:"success"`
}
