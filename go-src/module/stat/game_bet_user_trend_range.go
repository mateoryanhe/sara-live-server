package stat

import (
	"strings"
	"time"

	"xr-game-server/dao/statdao"
	statentity "xr-game-server/entity/stat"
)

const (
	dashboardTrendDailyMaxPoints   = 366
	dashboardTrendWeeklyMaxPoints  = 104
	dashboardTrendMonthlyMaxPoints = 60
)

func parseDashboardTrendDate(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, false
	}
	t, err := time.ParseInLocation("2006-01-02", s, time.Local)
	if err != nil {
		return time.Time{}, false
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local), true
}

func normalizeDashboardTrendDateRange(startDate, endDate string, now time.Time) (start, end time.Time, ok bool) {
	end, ok = parseDashboardTrendDate(endDate)
	if !ok {
		end = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	}
	start, ok = parseDashboardTrendDate(startDate)
	if !ok {
		start = end.AddDate(0, 0, -29)
	}
	if start.After(end) {
		start, end = end, start
	}
	return start, end, true
}

func dashboardTrendUsesCustomRange(startDate, endDate string) bool {
	return strings.TrimSpace(startDate) != "" || strings.TrimSpace(endDate) != ""
}

func buildDailyBucketIDs(start, end time.Time, maxPoints int) ([]string, bool) {
	if maxPoints <= 0 {
		maxPoints = gameBetUserDailyLimit
	}
	if maxPoints > dashboardTrendDailyMaxPoints {
		maxPoints = dashboardTrendDailyMaxPoints
	}
	truncated := false
	days := int(end.Sub(start).Hours()/24) + 1
	if days > maxPoints {
		truncated = true
		start = end.AddDate(0, 0, -(maxPoints - 1))
		days = maxPoints
	}
	ids := make([]string, 0, days)
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		ids = append(ids, statentity.FormatDailyLoginStatDate(d))
	}
	return ids, truncated
}

func startOfISOWeek(t time.Time) time.Time {
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	wd := int(t.Weekday())
	if wd == 0 {
		wd = 7
	}
	return t.AddDate(0, 0, -(wd - 1))
}

func buildWeeklyBucketIDs(start, end time.Time, maxPoints int) ([]string, bool) {
	if maxPoints <= 0 {
		maxPoints = gameBetUserWeeklyLimit
	}
	if maxPoints > dashboardTrendWeeklyMaxPoints {
		maxPoints = dashboardTrendWeeklyMaxPoints
	}
	weekStart := startOfISOWeek(start)
	weekEnd := startOfISOWeek(end)
	all := make([]string, 0, 16)
	for cur := weekStart; !cur.After(weekEnd); cur = cur.AddDate(0, 0, 7) {
		all = append(all, statentity.FormatWeeklyLoginStatKey(cur))
	}
	truncated := false
	if len(all) > maxPoints {
		truncated = true
		all = all[len(all)-maxPoints:]
	}
	return all, truncated
}

func buildMonthlyBucketIDs(start, end time.Time, maxPoints int) ([]string, bool) {
	if maxPoints <= 0 {
		maxPoints = gameBetUserMonthlyLimit
	}
	if maxPoints > dashboardTrendMonthlyMaxPoints {
		maxPoints = dashboardTrendMonthlyMaxPoints
	}
	cur := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.Local)
	endMonth := time.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, time.Local)
	all := make([]string, 0, 16)
	for !cur.After(endMonth) {
		all = append(all, statentity.FormatMonthlyGameBetUserStatKey(cur))
		cur = cur.AddDate(0, 1, 0)
	}
	truncated := false
	if len(all) > maxPoints {
		truncated = true
		all = all[len(all)-maxPoints:]
	}
	return all, truncated
}

func loadDailyGameBetTrendByRange(startDate, endDate string, now time.Time, liveOnline uint64) ([]*statentity.DailyGameBetUserStat, []*statentity.DailyGameOnlinePeakStat, bool) {
	start, end, _ := normalizeDashboardTrendDateRange(startDate, endDate, now)
	ids, truncated := buildDailyBucketIDs(start, end, dashboardTrendDailyMaxPoints)
	betRows := statdao.ListDailyGameBetUserStatsByIDs(ids)
	onlineRows := statdao.ListDailyGameOnlinePeakStatsByIDs(ids)
	betRows = overlayCurrentDailyGameBetUser(betRows, now)
	onlineRows = overlayCurrentDailyGameOnlinePeak(onlineRows, now, liveOnline)
	return betRows, onlineRows, truncated
}

func loadWeeklyGameBetTrendByRange(startDate, endDate string, now time.Time, liveOnline uint64) ([]*statentity.WeeklyGameBetUserStat, []*statentity.WeeklyGameOnlinePeakStat, bool) {
	start, end, _ := normalizeDashboardTrendDateRange(startDate, endDate, now)
	ids, truncated := buildWeeklyBucketIDs(start, end, dashboardTrendWeeklyMaxPoints)
	betRows := statdao.ListWeeklyGameBetUserStatsByIDs(ids)
	onlineRows := statdao.ListWeeklyGameOnlinePeakStatsByIDs(ids)
	betRows = overlayCurrentWeeklyGameBetUser(betRows, now)
	onlineRows = overlayCurrentWeeklyGameOnlinePeak(onlineRows, now, liveOnline)
	return betRows, onlineRows, truncated
}

func loadMonthlyGameBetTrendByRange(startDate, endDate string, now time.Time, liveOnline uint64) ([]*statentity.MonthlyGameBetUserStat, []*statentity.MonthlyGameOnlinePeakStat, bool) {
	start, end, _ := normalizeDashboardTrendDateRange(startDate, endDate, now)
	ids, truncated := buildMonthlyBucketIDs(start, end, dashboardTrendMonthlyMaxPoints)
	betRows := statdao.ListMonthlyGameBetUserStatsByIDs(ids)
	onlineRows := statdao.ListMonthlyGameOnlinePeakStatsByIDs(ids)
	betRows = overlayCurrentMonthlyGameBetUser(betRows, now)
	onlineRows = overlayCurrentMonthlyGameOnlinePeak(onlineRows, now, liveOnline)
	return betRows, onlineRows, truncated
}

func dashboardTrendMaxPointsForGranularity(granularity string) int {
	switch strings.ToLower(strings.TrimSpace(granularity)) {
	case "weekly", "week":
		return dashboardTrendWeeklyMaxPoints
	case "monthly", "month":
		return dashboardTrendMonthlyMaxPoints
	default:
		return dashboardTrendDailyMaxPoints
	}
}
