package stat

import (
	"time"

	"xr-game-server/dao/statdao"
	statentity "xr-game-server/entity/stat"
)

func loadDailyUserStatTrendByRange(startDate, endDate string, now time.Time) ([]*statentity.DailyLoginStat, bool) {
	start, end, _ := normalizeDashboardTrendDateRange(startDate, endDate, now)
	ids, truncated := buildDailyBucketIDs(start, end, dashboardTrendDailyMaxPoints)
	rows := statdao.ListDailyLoginStatsByIDs(ids)
	return overlayCurrentDailyLoginStat(rows, now), truncated
}

func loadWeeklyUserStatTrendByRange(startDate, endDate string, now time.Time) ([]*statentity.WeeklyLoginStat, bool) {
	start, end, _ := normalizeDashboardTrendDateRange(startDate, endDate, now)
	ids, truncated := buildWeeklyBucketIDs(start, end, dashboardTrendWeeklyMaxPoints)
	rows := statdao.ListWeeklyLoginStatsByIDs(ids)
	return overlayCurrentWeeklyLoginStat(rows, now), truncated
}

func loadMonthlyUserStatTrendByRange(startDate, endDate string, now time.Time) ([]*statentity.MonthlyLoginStat, bool) {
	start, end, _ := normalizeDashboardTrendDateRange(startDate, endDate, now)
	ids, truncated := buildMonthlyBucketIDs(start, end, dashboardTrendMonthlyMaxPoints)
	rows := statdao.ListMonthlyLoginStatsByIDs(ids)
	return overlayCurrentMonthlyLoginStat(rows, now), truncated
}
