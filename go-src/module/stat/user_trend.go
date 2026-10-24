package stat

import (
	"context"
	"time"

	"xr-game-server/dao/statdao"
	"xr-game-server/dto/statdto"
	"xr-game-server/entity/stat"
)

const (
	userStatDailyLimit   = 30
	userStatWeeklyLimit  = 12
	userStatMonthlyLimit = 12
)

// GetCMSUserStatTrend CMS获取用户数据趋势
func GetCMSUserStatTrend(_ context.Context, req *statdto.CMSUserStatTrendReq) (*statdto.CMSUserStatTrendRes, error) {
	if req == nil {
		req = &statdto.CMSUserStatTrendReq{}
	}
	now := time.Now()

	dailyRows := overlayCurrentDailyLoginStat(statdao.ListRecentDailyLoginStats(userStatDailyLimit), now)
	weeklyRows := overlayCurrentWeeklyLoginStat(statdao.ListRecentWeeklyLoginStats(userStatWeeklyLimit), now)
	monthlyRows := overlayCurrentMonthlyLoginStat(statdao.ListRecentMonthlyLoginStats(userStatMonthlyLimit), now)

	rangeTruncated := false
	trendMaxPoints := 0
	granularity := normalizeDashboardTrendGranularity(req.Granularity)
	if dashboardTrendUsesCustomRange(req.StartDate, req.EndDate) {
		trendMaxPoints = dashboardTrendMaxPointsForGranularity(granularity)
		switch granularity {
		case "weekly":
			var truncated bool
			weeklyRows, truncated = loadWeeklyUserStatTrendByRange(req.StartDate, req.EndDate, now)
			rangeTruncated = truncated
		case "monthly":
			var truncated bool
			monthlyRows, truncated = loadMonthlyUserStatTrendByRange(req.StartDate, req.EndDate, now)
			rangeTruncated = truncated
		default:
			var truncated bool
			dailyRows, truncated = loadDailyUserStatTrendByRange(req.StartDate, req.EndDate, now)
			rangeTruncated = truncated
		}
	}

	return &statdto.CMSUserStatTrendRes{
		Daily:          toDailyTrendPoints(dailyRows),
		Weekly:         toWeeklyTrendPoints(weeklyRows),
		Monthly:        toMonthlyTrendPoints(monthlyRows),
		RangeTruncated: rangeTruncated,
		TrendMaxPoints: trendMaxPoints,
	}, nil
}

// 历史列表走 ListCache; 当前日/周/月数据点用 RowCache 覆盖或追加, 列表缓存不含当天时曲线仍显示当日.
func overlayCurrentDailyLoginStat(rows []*entity.DailyLoginStat, now time.Time) []*entity.DailyLoginStat {
	today := entity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyLoginStatByDate(today)
	out := make([]*entity.DailyLoginStat, len(rows))
	copy(out, rows)
	found := false
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			found = true
			break
		}
	}
	if !found && live != nil {
		out = append(out, live)
	}
	return out
}

func overlayCurrentWeeklyLoginStat(rows []*entity.WeeklyLoginStat, now time.Time) []*entity.WeeklyLoginStat {
	week := entity.FormatWeeklyLoginStatKey(now)
	live := statdao.GetWeeklyLoginStatByWeek(week)
	out := make([]*entity.WeeklyLoginStat, len(rows))
	copy(out, rows)
	found := false
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			found = true
			break
		}
	}
	if !found && live != nil {
		out = append(out, live)
	}
	return out
}

func overlayCurrentMonthlyLoginStat(rows []*entity.MonthlyLoginStat, now time.Time) []*entity.MonthlyLoginStat {
	month := entity.FormatMonthlyLoginStatKey(now)
	live := statdao.GetMonthlyLoginStatByMonth(month)
	out := make([]*entity.MonthlyLoginStat, len(rows))
	copy(out, rows)
	found := false
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			found = true
			break
		}
	}
	if !found && live != nil {
		out = append(out, live)
	}
	return out
}

func toDailyTrendPoints(rows []*entity.DailyLoginStat) []*statdto.CMSUserStatTrendPoint {
	list := make([]*statdto.CMSUserStatTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &statdto.CMSUserStatTrendPoint{
			Time:                    row.ID,
			ActiveUserCount:         row.Count,
			RegisterUserCount:       row.RegisterCount,
			AppDeviceRechargeAmount: row.AppDeviceRechargeAmount,
			H5DeviceRechargeAmount:  row.H5DeviceRechargeAmount,
			AppDeviceLoginCount:     row.AppDeviceLoginCount,
			H5DeviceLoginCount:      row.H5DeviceLoginCount,
			BarMetrics:              statdto.BuildUserStatBarMetrics(row.RechargeUserCount),
		})
	}
	return list
}

func toWeeklyTrendPoints(rows []*entity.WeeklyLoginStat) []*statdto.CMSUserStatTrendPoint {
	list := make([]*statdto.CMSUserStatTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &statdto.CMSUserStatTrendPoint{
			Time:                    row.ID,
			ActiveUserCount:         row.Count,
			RegisterUserCount:       row.RegisterCount,
			AppDeviceRechargeAmount: row.AppDeviceRechargeAmount,
			H5DeviceRechargeAmount:  row.H5DeviceRechargeAmount,
			AppDeviceLoginCount:     row.AppDeviceLoginCount,
			H5DeviceLoginCount:      row.H5DeviceLoginCount,
			BarMetrics:              statdto.BuildUserStatBarMetrics(row.RechargeUserCount),
		})
	}
	return list
}

func toMonthlyTrendPoints(rows []*entity.MonthlyLoginStat) []*statdto.CMSUserStatTrendPoint {
	list := make([]*statdto.CMSUserStatTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &statdto.CMSUserStatTrendPoint{
			Time:                    row.ID,
			ActiveUserCount:         row.Count,
			RegisterUserCount:       row.RegisterCount,
			AppDeviceRechargeAmount: row.AppDeviceRechargeAmount,
			H5DeviceRechargeAmount:  row.H5DeviceRechargeAmount,
			AppDeviceLoginCount:     row.AppDeviceLoginCount,
			H5DeviceLoginCount:      row.H5DeviceLoginCount,
			BarMetrics:              statdto.BuildUserStatBarMetrics(row.RechargeUserCount),
		})
	}
	return list
}
