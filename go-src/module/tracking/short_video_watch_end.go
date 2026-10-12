package tracking

import (
	"context"
	"time"

	"xr-game-server/dao/statdao"
	"xr-game-server/dto/trackingdto"
	statentity "xr-game-server/entity/stat"
)

const (
	ShortVideoWatchEndDailyLimit   = 30
	ShortVideoWatchEndWeeklyLimit  = 12
	ShortVideoWatchEndMonthlyLimit = 12

	// EventKeyShortVideoWatchEnd 埋点 event_key(对齐埋点表 short_video_watch_end)
	EventKeyShortVideoWatchEnd = "short_video_watch_end"
)

func recordShortVideoWatchEnd(userId uint64, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyShortVideoWatchEndStatKey(at)
	month := statentity.FormatMonthlyShortVideoWatchEndStatKey(at)

	daily := statdao.GetDailyShortVideoWatchEndStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyShortVideoWatchEndStat(daily)

	weekly := statdao.GetWeeklyShortVideoWatchEndStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyShortVideoWatchEndStat(weekly)

	monthly := statdao.GetMonthlyShortVideoWatchEndStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyShortVideoWatchEndStat(monthly)
}

func GetCMSShortVideoWatchEndTrend(_ context.Context, _ *trackingdto.CMSShortVideoWatchEndTrendReq) (*trackingdto.CMSShortVideoWatchEndTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyShortVideoWatchEndStatKey(now)
	month := statentity.FormatMonthlyShortVideoWatchEndStatKey(now)

	todayStat := statdao.GetDailyShortVideoWatchEndStatByDate(today)
	weekStat := statdao.GetWeeklyShortVideoWatchEndStatByWeek(week)
	monthStat := statdao.GetMonthlyShortVideoWatchEndStatByMonth(month)

	return &trackingdto.CMSShortVideoWatchEndTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toShortVideoWatchEndDailyPoints(overlayCurrentDailyShortVideoWatchEnd(statdao.ListRecentDailyShortVideoWatchEndStats(ShortVideoWatchEndDailyLimit), now)),
		Weekly:     toShortVideoWatchEndWeeklyPoints(overlayCurrentWeeklyShortVideoWatchEnd(statdao.ListRecentWeeklyShortVideoWatchEndStats(ShortVideoWatchEndWeeklyLimit), now)),
		Monthly:    toShortVideoWatchEndMonthlyPoints(overlayCurrentMonthlyShortVideoWatchEnd(statdao.ListRecentMonthlyShortVideoWatchEndStats(ShortVideoWatchEndMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyShortVideoWatchEnd(rows []*statentity.DailyShortVideoWatchEndStat, now time.Time) []*statentity.DailyShortVideoWatchEndStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyShortVideoWatchEndStatByDate(today)
	out := make([]*statentity.DailyShortVideoWatchEndStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyShortVideoWatchEnd(rows []*statentity.WeeklyShortVideoWatchEndStat, now time.Time) []*statentity.WeeklyShortVideoWatchEndStat {
	week := statentity.FormatWeeklyShortVideoWatchEndStatKey(now)
	live := statdao.GetWeeklyShortVideoWatchEndStatByWeek(week)
	out := make([]*statentity.WeeklyShortVideoWatchEndStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyShortVideoWatchEnd(rows []*statentity.MonthlyShortVideoWatchEndStat, now time.Time) []*statentity.MonthlyShortVideoWatchEndStat {
	month := statentity.FormatMonthlyShortVideoWatchEndStatKey(now)
	live := statdao.GetMonthlyShortVideoWatchEndStatByMonth(month)
	out := make([]*statentity.MonthlyShortVideoWatchEndStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toShortVideoWatchEndDailyPoints(rows []*statentity.DailyShortVideoWatchEndStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toShortVideoWatchEndWeeklyPoints(rows []*statentity.WeeklyShortVideoWatchEndStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toShortVideoWatchEndMonthlyPoints(rows []*statentity.MonthlyShortVideoWatchEndStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
