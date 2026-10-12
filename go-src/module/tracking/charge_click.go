package tracking

import (
	"context"
	"time"

	"xr-game-server/dao/statdao"
	"xr-game-server/dto/trackingdto"
	statentity "xr-game-server/entity/stat"
)

const (
	chargeClickDailyLimit   = 30
	chargeClickWeeklyLimit  = 12
	chargeClickMonthlyLimit = 12

	// EventKeyChargeClick 埋点 event_key(对齐埋点表 charge_click)
	EventKeyChargeClick = "charge_click"
)

func recordChargeClick(userId uint64, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyChargeClickStatKey(at)
	month := statentity.FormatMonthlyChargeClickStatKey(at)

	daily := statdao.GetDailyChargeClickStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyChargeClickStat(daily)

	weekly := statdao.GetWeeklyChargeClickStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyChargeClickStat(weekly)

	monthly := statdao.GetMonthlyChargeClickStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyChargeClickStat(monthly)
}

func GetCMSChargeClickTrend(_ context.Context, _ *trackingdto.CMSChargeClickTrendReq) (*trackingdto.CMSChargeClickTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyChargeClickStatKey(now)
	month := statentity.FormatMonthlyChargeClickStatKey(now)

	todayStat := statdao.GetDailyChargeClickStatByDate(today)
	weekStat := statdao.GetWeeklyChargeClickStatByWeek(week)
	monthStat := statdao.GetMonthlyChargeClickStatByMonth(month)

	return &trackingdto.CMSChargeClickTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toChargeClickDailyPoints(overlayCurrentDailyChargeClick(statdao.ListRecentDailyChargeClickStats(chargeClickDailyLimit), now)),
		Weekly:     toChargeClickWeeklyPoints(overlayCurrentWeeklyChargeClick(statdao.ListRecentWeeklyChargeClickStats(chargeClickWeeklyLimit), now)),
		Monthly:    toChargeClickMonthlyPoints(overlayCurrentMonthlyChargeClick(statdao.ListRecentMonthlyChargeClickStats(chargeClickMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyChargeClick(rows []*statentity.DailyChargeClickStat, now time.Time) []*statentity.DailyChargeClickStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyChargeClickStatByDate(today)
	out := make([]*statentity.DailyChargeClickStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyChargeClick(rows []*statentity.WeeklyChargeClickStat, now time.Time) []*statentity.WeeklyChargeClickStat {
	week := statentity.FormatWeeklyChargeClickStatKey(now)
	live := statdao.GetWeeklyChargeClickStatByWeek(week)
	out := make([]*statentity.WeeklyChargeClickStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyChargeClick(rows []*statentity.MonthlyChargeClickStat, now time.Time) []*statentity.MonthlyChargeClickStat {
	month := statentity.FormatMonthlyChargeClickStatKey(now)
	live := statdao.GetMonthlyChargeClickStatByMonth(month)
	out := make([]*statentity.MonthlyChargeClickStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toChargeClickDailyPoints(rows []*statentity.DailyChargeClickStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toChargeClickWeeklyPoints(rows []*statentity.WeeklyChargeClickStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toChargeClickMonthlyPoints(rows []*statentity.MonthlyChargeClickStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
