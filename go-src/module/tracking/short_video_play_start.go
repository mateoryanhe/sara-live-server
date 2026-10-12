package tracking

import (
	"context"
	"time"

	"xr-game-server/dao/statdao"
	"xr-game-server/dto/trackingdto"
	statentity "xr-game-server/entity/stat"
)

const (
	ShortVideoPlayStartDailyLimit   = 30
	ShortVideoPlayStartWeeklyLimit  = 12
	ShortVideoPlayStartMonthlyLimit = 12

	// EventKeyShortVideoPlayStart 埋点 event_key(对齐埋点表 short_video_play_start)
	EventKeyShortVideoPlayStart = "short_video_play_start"
)

func recordShortVideoPlayStart(userId uint64, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyShortVideoPlayStartStatKey(at)
	month := statentity.FormatMonthlyShortVideoPlayStartStatKey(at)

	daily := statdao.GetDailyShortVideoPlayStartStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyShortVideoPlayStartStat(daily)

	weekly := statdao.GetWeeklyShortVideoPlayStartStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyShortVideoPlayStartStat(weekly)

	monthly := statdao.GetMonthlyShortVideoPlayStartStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyShortVideoPlayStartStat(monthly)
}

func GetCMSShortVideoPlayStartTrend(_ context.Context, _ *trackingdto.CMSShortVideoPlayStartTrendReq) (*trackingdto.CMSShortVideoPlayStartTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyShortVideoPlayStartStatKey(now)
	month := statentity.FormatMonthlyShortVideoPlayStartStatKey(now)

	todayStat := statdao.GetDailyShortVideoPlayStartStatByDate(today)
	weekStat := statdao.GetWeeklyShortVideoPlayStartStatByWeek(week)
	monthStat := statdao.GetMonthlyShortVideoPlayStartStatByMonth(month)

	return &trackingdto.CMSShortVideoPlayStartTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toShortVideoPlayStartDailyPoints(overlayCurrentDailyShortVideoPlayStart(statdao.ListRecentDailyShortVideoPlayStartStats(ShortVideoPlayStartDailyLimit), now)),
		Weekly:     toShortVideoPlayStartWeeklyPoints(overlayCurrentWeeklyShortVideoPlayStart(statdao.ListRecentWeeklyShortVideoPlayStartStats(ShortVideoPlayStartWeeklyLimit), now)),
		Monthly:    toShortVideoPlayStartMonthlyPoints(overlayCurrentMonthlyShortVideoPlayStart(statdao.ListRecentMonthlyShortVideoPlayStartStats(ShortVideoPlayStartMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyShortVideoPlayStart(rows []*statentity.DailyShortVideoPlayStartStat, now time.Time) []*statentity.DailyShortVideoPlayStartStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyShortVideoPlayStartStatByDate(today)
	out := make([]*statentity.DailyShortVideoPlayStartStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyShortVideoPlayStart(rows []*statentity.WeeklyShortVideoPlayStartStat, now time.Time) []*statentity.WeeklyShortVideoPlayStartStat {
	week := statentity.FormatWeeklyShortVideoPlayStartStatKey(now)
	live := statdao.GetWeeklyShortVideoPlayStartStatByWeek(week)
	out := make([]*statentity.WeeklyShortVideoPlayStartStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyShortVideoPlayStart(rows []*statentity.MonthlyShortVideoPlayStartStat, now time.Time) []*statentity.MonthlyShortVideoPlayStartStat {
	month := statentity.FormatMonthlyShortVideoPlayStartStatKey(now)
	live := statdao.GetMonthlyShortVideoPlayStartStatByMonth(month)
	out := make([]*statentity.MonthlyShortVideoPlayStartStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toShortVideoPlayStartDailyPoints(rows []*statentity.DailyShortVideoPlayStartStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toShortVideoPlayStartWeeklyPoints(rows []*statentity.WeeklyShortVideoPlayStartStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toShortVideoPlayStartMonthlyPoints(rows []*statentity.MonthlyShortVideoPlayStartStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
