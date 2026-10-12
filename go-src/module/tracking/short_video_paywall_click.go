package tracking

import (
	"context"
	"time"

	"xr-game-server/dao/statdao"
	"xr-game-server/dto/trackingdto"
	statentity "xr-game-server/entity/stat"
)

const (
	ShortVideoPaywallClickDailyLimit   = 30
	ShortVideoPaywallClickWeeklyLimit  = 12
	ShortVideoPaywallClickMonthlyLimit = 12

	// EventKeyShortVideoPaywallClick 埋点 event_key(对齐埋点表 short_video_paywall_click)
	EventKeyShortVideoPaywallClick = "short_video_paywall_click"
)

func recordShortVideoPaywallClick(userId uint64, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyShortVideoPaywallClickStatKey(at)
	month := statentity.FormatMonthlyShortVideoPaywallClickStatKey(at)

	daily := statdao.GetDailyShortVideoPaywallClickStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyShortVideoPaywallClickStat(daily)

	weekly := statdao.GetWeeklyShortVideoPaywallClickStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyShortVideoPaywallClickStat(weekly)

	monthly := statdao.GetMonthlyShortVideoPaywallClickStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyShortVideoPaywallClickStat(monthly)
}

func GetCMSShortVideoPaywallClickTrend(_ context.Context, _ *trackingdto.CMSShortVideoPaywallClickTrendReq) (*trackingdto.CMSShortVideoPaywallClickTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyShortVideoPaywallClickStatKey(now)
	month := statentity.FormatMonthlyShortVideoPaywallClickStatKey(now)

	todayStat := statdao.GetDailyShortVideoPaywallClickStatByDate(today)
	weekStat := statdao.GetWeeklyShortVideoPaywallClickStatByWeek(week)
	monthStat := statdao.GetMonthlyShortVideoPaywallClickStatByMonth(month)

	return &trackingdto.CMSShortVideoPaywallClickTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toShortVideoPaywallClickDailyPoints(overlayCurrentDailyShortVideoPaywallClick(statdao.ListRecentDailyShortVideoPaywallClickStats(ShortVideoPaywallClickDailyLimit), now)),
		Weekly:     toShortVideoPaywallClickWeeklyPoints(overlayCurrentWeeklyShortVideoPaywallClick(statdao.ListRecentWeeklyShortVideoPaywallClickStats(ShortVideoPaywallClickWeeklyLimit), now)),
		Monthly:    toShortVideoPaywallClickMonthlyPoints(overlayCurrentMonthlyShortVideoPaywallClick(statdao.ListRecentMonthlyShortVideoPaywallClickStats(ShortVideoPaywallClickMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyShortVideoPaywallClick(rows []*statentity.DailyShortVideoPaywallClickStat, now time.Time) []*statentity.DailyShortVideoPaywallClickStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyShortVideoPaywallClickStatByDate(today)
	out := make([]*statentity.DailyShortVideoPaywallClickStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyShortVideoPaywallClick(rows []*statentity.WeeklyShortVideoPaywallClickStat, now time.Time) []*statentity.WeeklyShortVideoPaywallClickStat {
	week := statentity.FormatWeeklyShortVideoPaywallClickStatKey(now)
	live := statdao.GetWeeklyShortVideoPaywallClickStatByWeek(week)
	out := make([]*statentity.WeeklyShortVideoPaywallClickStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyShortVideoPaywallClick(rows []*statentity.MonthlyShortVideoPaywallClickStat, now time.Time) []*statentity.MonthlyShortVideoPaywallClickStat {
	month := statentity.FormatMonthlyShortVideoPaywallClickStatKey(now)
	live := statdao.GetMonthlyShortVideoPaywallClickStatByMonth(month)
	out := make([]*statentity.MonthlyShortVideoPaywallClickStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toShortVideoPaywallClickDailyPoints(rows []*statentity.DailyShortVideoPaywallClickStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toShortVideoPaywallClickWeeklyPoints(rows []*statentity.WeeklyShortVideoPaywallClickStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toShortVideoPaywallClickMonthlyPoints(rows []*statentity.MonthlyShortVideoPaywallClickStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
