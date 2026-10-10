package stat

import (
	"context"
	"time"

	"xr-game-server/dao/statdao"
	"xr-game-server/dto/statdto"
	"xr-game-server/entity/stat"
)

// ResetCMSBasicDashboardStats 将仪表盘「基础数据」Tab 展示的全部累计/今日统计归零(不改用户钱包、不动主播代付 Tab).
func ResetCMSBasicDashboardStats(_ context.Context, _ *statdto.CMSResetBasicDashboardStatsReq) (*statdto.CMSResetBasicDashboardStatsRes, error) {
	sys := statdao.GetSysStat()
	if sys == nil {
		sys = entity.NewSystemTotalStat(entity.SystemTotalStatDefaultID)
	}
	sys.ResetBasicDashboardCounters()

	now := time.Now()
	today := entity.FormatDailyLoginStatDate(now)
	todayStat := statdao.GetDailyLoginStatByDate(today)
	todayStat.ResetBasicDashboardTodayMetrics()
	statdao.PublishDailyLoginStat(todayStat)

	return &statdto.CMSResetBasicDashboardStatsRes{Success: true}, nil
}
