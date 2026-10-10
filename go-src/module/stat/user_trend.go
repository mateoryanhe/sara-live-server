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
func GetCMSUserStatTrend(_ context.Context, _ *statdto.CMSUserStatTrendReq) (*statdto.CMSUserStatTrendRes, error) {
	now := time.Now()
	return &statdto.CMSUserStatTrendRes{
		Daily: toDailyTrendPoints(
			overlayCurrentDailyLoginStat(statdao.ListRecentDailyLoginStats(userStatDailyLimit), now),
		),
		Weekly: toWeeklyTrendPoints(
			overlayCurrentWeeklyLoginStat(statdao.ListRecentWeeklyLoginStats(userStatWeeklyLimit), now),
		),
		Monthly: toMonthlyTrendPoints(
			overlayCurrentMonthlyLoginStat(statdao.ListRecentMonthlyLoginStats(userStatMonthlyLimit), now),
		),
	}, nil
}

// 历史列表走 ListCache; 当前日/周/月数据点用 RowCache 覆盖, 自动刷新不必重复查库也能更新当日曲线.
func overlayCurrentDailyLoginStat(rows []*entity.DailyLoginStat, now time.Time) []*entity.DailyLoginStat {
	today := entity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyLoginStatByDate(today)
	out := make([]*entity.DailyLoginStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			break
		}
	}
	return out
}

func overlayCurrentWeeklyLoginStat(rows []*entity.WeeklyLoginStat, now time.Time) []*entity.WeeklyLoginStat {
	week := entity.FormatWeeklyLoginStatKey(now)
	live := statdao.GetWeeklyLoginStatByWeek(week)
	out := make([]*entity.WeeklyLoginStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			break
		}
	}
	return out
}

func overlayCurrentMonthlyLoginStat(rows []*entity.MonthlyLoginStat, now time.Time) []*entity.MonthlyLoginStat {
	month := entity.FormatMonthlyLoginStatKey(now)
	live := statdao.GetMonthlyLoginStatByMonth(month)
	out := make([]*entity.MonthlyLoginStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			break
		}
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
			Time:              row.ID,
			ActiveUserCount:   row.Count,
			RegisterUserCount: row.RegisterCount,
			BarMetrics: statdto.BuildUserStatBarMetrics(
				row.RechargeUserCount, row.GoldConsumeUserCount, row.DiamondConsumeUserCount,
			),
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
			Time:              row.ID,
			ActiveUserCount:   row.Count,
			RegisterUserCount: row.RegisterCount,
			BarMetrics: statdto.BuildUserStatBarMetrics(
				row.RechargeUserCount, row.GoldConsumeUserCount, row.DiamondConsumeUserCount,
			),
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
			Time:              row.ID,
			ActiveUserCount:   row.Count,
			RegisterUserCount: row.RegisterCount,
			BarMetrics: statdto.BuildUserStatBarMetrics(
				row.RechargeUserCount, row.GoldConsumeUserCount, row.DiamondConsumeUserCount,
			),
		})
	}
	return list
}
