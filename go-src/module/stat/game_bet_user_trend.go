package stat

import (
	"context"
	"strings"
	"time"

	"xr-game-server/dao/statdao"
	"xr-game-server/dto/statdto"
	statentity "xr-game-server/entity/stat"
)

const (
	gameBetUserDailyLimit   = 30
	gameBetUserWeeklyLimit  = 12
	gameBetUserMonthlyLimit = 12
)

// GetCMSGameBetUserTrend CMS 游戏下注人数趋势(订阅 GameBetCreatedEvent 累计)
func GetCMSGameBetUserTrend(_ context.Context, req *statdto.CMSGameBetUserTrendReq) (*statdto.CMSGameBetUserTrendRes, error) {
	if req == nil {
		req = &statdto.CMSGameBetUserTrendReq{}
	}
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyLoginStatKey(now)
	month := statentity.FormatMonthlyGameBetUserStatKey(now)

	todayStat := statdao.GetDailyGameBetUserStatByDate(today)
	weekStat := statdao.GetWeeklyGameBetUserStatByWeek(week)
	monthStat := statdao.GetMonthlyGameBetUserStatByMonth(month)

	currentOnline := CurrentGameOnlineCount()
	todayOnlinePeak := statdao.GetDailyGameOnlinePeakStatByDate(today).PeakCount
	weekOnlinePeak := statdao.GetWeeklyGameOnlinePeakStatByWeek(week).PeakCount
	monthOnlinePeak := statdao.GetMonthlyGameOnlinePeakStatByMonth(month).PeakCount
	todayOnlinePeak = maxUint64(todayOnlinePeak, currentOnline)
	weekOnlinePeak = maxUint64(weekOnlinePeak, currentOnline)
	monthOnlinePeak = maxUint64(monthOnlinePeak, currentOnline)

	dailyBet := overlayCurrentDailyGameBetUser(statdao.ListRecentDailyGameBetUserStats(gameBetUserDailyLimit), now)
	weeklyBet := overlayCurrentWeeklyGameBetUser(statdao.ListRecentWeeklyGameBetUserStats(gameBetUserWeeklyLimit), now)
	monthlyBet := overlayCurrentMonthlyGameBetUser(statdao.ListRecentMonthlyGameBetUserStats(gameBetUserMonthlyLimit), now)

	dailyOnline := overlayCurrentDailyGameOnlinePeak(statdao.ListRecentDailyGameOnlinePeakStats(gameBetUserDailyLimit), now, currentOnline)
	weeklyOnline := overlayCurrentWeeklyGameOnlinePeak(statdao.ListRecentWeeklyGameOnlinePeakStats(gameBetUserWeeklyLimit), now, currentOnline)
	monthlyOnline := overlayCurrentMonthlyGameOnlinePeak(statdao.ListRecentMonthlyGameOnlinePeakStats(gameBetUserMonthlyLimit), now, currentOnline)

	rangeTruncated := false
	trendMaxPoints := 0
	granularity := normalizeDashboardTrendGranularity(req.Granularity)
	if dashboardTrendUsesCustomRange(req.StartDate, req.EndDate) {
		trendMaxPoints = dashboardTrendMaxPointsForGranularity(granularity)
		switch granularity {
		case "weekly":
			var truncated bool
			weeklyBet, weeklyOnline, truncated = loadWeeklyGameBetTrendByRange(req.StartDate, req.EndDate, now, currentOnline)
			rangeTruncated = truncated
		case "monthly":
			var truncated bool
			monthlyBet, monthlyOnline, truncated = loadMonthlyGameBetTrendByRange(req.StartDate, req.EndDate, now, currentOnline)
			rangeTruncated = truncated
		default:
			var truncated bool
			dailyBet, dailyOnline, truncated = loadDailyGameBetTrendByRange(req.StartDate, req.EndDate, now, currentOnline)
			rangeTruncated = truncated
		}
	}

	return &statdto.CMSGameBetUserTrendRes{
		LastHour:           buildLastHourGameTrendPoints(now),
		TodayCount:         todayStat.UserCount,
		WeekCount:          weekStat.UserCount,
		MonthCount:         monthStat.UserCount,
		CurrentOnlineCount: currentOnline,
		TodayOnlinePeak:    todayOnlinePeak,
		WeekOnlinePeak:     weekOnlinePeak,
		MonthOnlinePeak:    monthOnlinePeak,
		Daily: mergeGameBetTrendPoints(
			toGameBetUserDailyPoints(dailyBet),
			gameOnlinePeakMapDaily(dailyOnline),
			today,
			currentOnline,
		),
		Weekly: mergeGameBetTrendPoints(
			toGameBetUserWeeklyPoints(weeklyBet),
			gameOnlinePeakMapWeekly(weeklyOnline),
			week,
			currentOnline,
		),
		Monthly: mergeGameBetTrendPoints(
			toGameBetUserMonthlyPoints(monthlyBet),
			gameOnlinePeakMapMonthly(monthlyOnline),
			month,
			currentOnline,
		),
		RangeTruncated: rangeTruncated,
		TrendMaxPoints: trendMaxPoints,
	}, nil
}

func normalizeDashboardTrendGranularity(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "weekly", "week":
		return "weekly"
	case "monthly", "month":
		return "monthly"
	default:
		return "daily"
	}
}

func gameOnlinePeakMapDaily(rows []*statentity.DailyGameOnlinePeakStat) map[string]uint64 {
	m := make(map[string]uint64, len(rows))
	for _, row := range rows {
		if row != nil {
			m[row.ID] = row.PeakCount
		}
	}
	return m
}

func gameOnlinePeakMapWeekly(rows []*statentity.WeeklyGameOnlinePeakStat) map[string]uint64 {
	m := make(map[string]uint64, len(rows))
	for _, row := range rows {
		if row != nil {
			m[row.ID] = row.PeakCount
		}
	}
	return m
}

func gameOnlinePeakMapMonthly(rows []*statentity.MonthlyGameOnlinePeakStat) map[string]uint64 {
	m := make(map[string]uint64, len(rows))
	for _, row := range rows {
		if row != nil {
			m[row.ID] = row.PeakCount
		}
	}
	return m
}

func maxUint64(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func overlayCurrentDailyGameOnlinePeak(rows []*statentity.DailyGameOnlinePeakStat, now time.Time, liveOnline uint64) []*statentity.DailyGameOnlinePeakStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyGameOnlinePeakStatByDate(today)
	live.PeakCount = maxUint64(live.PeakCount, liveOnline)
	out := make([]*statentity.DailyGameOnlinePeakStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyGameOnlinePeak(rows []*statentity.WeeklyGameOnlinePeakStat, now time.Time, liveOnline uint64) []*statentity.WeeklyGameOnlinePeakStat {
	week := statentity.FormatWeeklyLoginStatKey(now)
	live := statdao.GetWeeklyGameOnlinePeakStatByWeek(week)
	live.PeakCount = maxUint64(live.PeakCount, liveOnline)
	out := make([]*statentity.WeeklyGameOnlinePeakStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyGameOnlinePeak(rows []*statentity.MonthlyGameOnlinePeakStat, now time.Time, liveOnline uint64) []*statentity.MonthlyGameOnlinePeakStat {
	month := statentity.FormatMonthlyGameBetUserStatKey(now)
	live := statdao.GetMonthlyGameOnlinePeakStatByMonth(month)
	live.PeakCount = maxUint64(live.PeakCount, liveOnline)
	out := make([]*statentity.MonthlyGameOnlinePeakStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func mergeGameBetTrendPoints(
	betPoints []*statdto.CMSGameBetUserTrendPoint,
	onlineByTime map[string]uint64,
	currentBucket string,
	liveOnline uint64,
) []*statdto.CMSGameBetUserTrendPoint {
	for _, p := range betPoints {
		if p == nil {
			continue
		}
		peak := onlineByTime[p.Time]
		if p.Time == currentBucket {
			peak = maxUint64(peak, liveOnline)
		}
		p.OnlineCount = peak
	}
	return betPoints
}

func overlayCurrentDailyGameBetUser(rows []*statentity.DailyGameBetUserStat, now time.Time) []*statentity.DailyGameBetUserStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyGameBetUserStatByDate(today)
	out := make([]*statentity.DailyGameBetUserStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyGameBetUser(rows []*statentity.WeeklyGameBetUserStat, now time.Time) []*statentity.WeeklyGameBetUserStat {
	week := statentity.FormatWeeklyLoginStatKey(now)
	live := statdao.GetWeeklyGameBetUserStatByWeek(week)
	out := make([]*statentity.WeeklyGameBetUserStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyGameBetUser(rows []*statentity.MonthlyGameBetUserStat, now time.Time) []*statentity.MonthlyGameBetUserStat {
	month := statentity.FormatMonthlyGameBetUserStatKey(now)
	live := statdao.GetMonthlyGameBetUserStatByMonth(month)
	out := make([]*statentity.MonthlyGameBetUserStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toGameBetUserDailyPoints(rows []*statentity.DailyGameBetUserStat) []*statdto.CMSGameBetUserTrendPoint {
	list := make([]*statdto.CMSGameBetUserTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &statdto.CMSGameBetUserTrendPoint{Time: row.ID, Count: row.UserCount})
	}
	return list
}

func toGameBetUserWeeklyPoints(rows []*statentity.WeeklyGameBetUserStat) []*statdto.CMSGameBetUserTrendPoint {
	list := make([]*statdto.CMSGameBetUserTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &statdto.CMSGameBetUserTrendPoint{Time: row.ID, Count: row.UserCount})
	}
	return list
}

func toGameBetUserMonthlyPoints(rows []*statentity.MonthlyGameBetUserStat) []*statdto.CMSGameBetUserTrendPoint {
	list := make([]*statdto.CMSGameBetUserTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &statdto.CMSGameBetUserTrendPoint{Time: row.ID, Count: row.UserCount})
	}
	return list
}
