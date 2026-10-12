package tracking

import (
	"context"
	"time"

	"xr-game-server/dao/statdao"
	"xr-game-server/dto/trackingdto"
	statentity "xr-game-server/entity/stat"
)

const (
	ShortVideoPaySuccessDailyLimit   = 30
	ShortVideoPaySuccessWeeklyLimit  = 12
	ShortVideoPaySuccessMonthlyLimit = 12

	// EventKeyShortVideoPaySuccess 埋点 event_key(对齐埋点表 short_video_pay_success)
	EventKeyShortVideoPaySuccess = "short_video_pay_success"
)

func recordShortVideoPaySuccess(userId uint64, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyShortVideoPaySuccessStatKey(at)
	month := statentity.FormatMonthlyShortVideoPaySuccessStatKey(at)

	daily := statdao.GetDailyShortVideoPaySuccessStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyShortVideoPaySuccessStat(daily)

	weekly := statdao.GetWeeklyShortVideoPaySuccessStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyShortVideoPaySuccessStat(weekly)

	monthly := statdao.GetMonthlyShortVideoPaySuccessStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyShortVideoPaySuccessStat(monthly)
}

func GetCMSShortVideoPaySuccessTrend(_ context.Context, _ *trackingdto.CMSShortVideoPaySuccessTrendReq) (*trackingdto.CMSShortVideoPaySuccessTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyShortVideoPaySuccessStatKey(now)
	month := statentity.FormatMonthlyShortVideoPaySuccessStatKey(now)

	todayStat := statdao.GetDailyShortVideoPaySuccessStatByDate(today)
	weekStat := statdao.GetWeeklyShortVideoPaySuccessStatByWeek(week)
	monthStat := statdao.GetMonthlyShortVideoPaySuccessStatByMonth(month)

	return &trackingdto.CMSShortVideoPaySuccessTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toShortVideoPaySuccessDailyPoints(overlayCurrentDailyShortVideoPaySuccess(statdao.ListRecentDailyShortVideoPaySuccessStats(ShortVideoPaySuccessDailyLimit), now)),
		Weekly:     toShortVideoPaySuccessWeeklyPoints(overlayCurrentWeeklyShortVideoPaySuccess(statdao.ListRecentWeeklyShortVideoPaySuccessStats(ShortVideoPaySuccessWeeklyLimit), now)),
		Monthly:    toShortVideoPaySuccessMonthlyPoints(overlayCurrentMonthlyShortVideoPaySuccess(statdao.ListRecentMonthlyShortVideoPaySuccessStats(ShortVideoPaySuccessMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyShortVideoPaySuccess(rows []*statentity.DailyShortVideoPaySuccessStat, now time.Time) []*statentity.DailyShortVideoPaySuccessStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyShortVideoPaySuccessStatByDate(today)
	out := make([]*statentity.DailyShortVideoPaySuccessStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyShortVideoPaySuccess(rows []*statentity.WeeklyShortVideoPaySuccessStat, now time.Time) []*statentity.WeeklyShortVideoPaySuccessStat {
	week := statentity.FormatWeeklyShortVideoPaySuccessStatKey(now)
	live := statdao.GetWeeklyShortVideoPaySuccessStatByWeek(week)
	out := make([]*statentity.WeeklyShortVideoPaySuccessStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyShortVideoPaySuccess(rows []*statentity.MonthlyShortVideoPaySuccessStat, now time.Time) []*statentity.MonthlyShortVideoPaySuccessStat {
	month := statentity.FormatMonthlyShortVideoPaySuccessStatKey(now)
	live := statdao.GetMonthlyShortVideoPaySuccessStatByMonth(month)
	out := make([]*statentity.MonthlyShortVideoPaySuccessStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toShortVideoPaySuccessDailyPoints(rows []*statentity.DailyShortVideoPaySuccessStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toShortVideoPaySuccessWeeklyPoints(rows []*statentity.WeeklyShortVideoPaySuccessStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toShortVideoPaySuccessMonthlyPoints(rows []*statentity.MonthlyShortVideoPaySuccessStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
