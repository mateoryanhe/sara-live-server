package tracking

import (
	"context"
	"time"

	"xr-game-server/dao/statdao"
	"xr-game-server/dto/trackingdto"
	statentity "xr-game-server/entity/stat"
)

const (
	miniGameRoundResultDailyLimit   = 30
	miniGameRoundResultWeeklyLimit  = 12
	miniGameRoundResultMonthlyLimit = 12

	// EventKeyMiniGameRoundResult 埋点 event_key(对齐埋点表 mini_game_round_result)
	EventKeyMiniGameRoundResult = "mini_game_round_result"
)

func recordMiniGameRoundResult(userId uint64, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyMiniGameRoundResultStatKey(at)
	month := statentity.FormatMonthlyMiniGameRoundResultStatKey(at)

	daily := statdao.GetDailyMiniGameRoundResultStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyMiniGameRoundResultStat(daily)

	weekly := statdao.GetWeeklyMiniGameRoundResultStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyMiniGameRoundResultStat(weekly)

	monthly := statdao.GetMonthlyMiniGameRoundResultStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyMiniGameRoundResultStat(monthly)
}

func GetCMSMiniGameRoundResultTrend(_ context.Context, _ *trackingdto.CMSMiniGameRoundResultTrendReq) (*trackingdto.CMSMiniGameRoundResultTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyMiniGameRoundResultStatKey(now)
	month := statentity.FormatMonthlyMiniGameRoundResultStatKey(now)

	todayStat := statdao.GetDailyMiniGameRoundResultStatByDate(today)
	weekStat := statdao.GetWeeklyMiniGameRoundResultStatByWeek(week)
	monthStat := statdao.GetMonthlyMiniGameRoundResultStatByMonth(month)

	return &trackingdto.CMSMiniGameRoundResultTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toMiniGameRoundResultDailyPoints(overlayCurrentDailyMiniGameRoundResult(statdao.ListRecentDailyMiniGameRoundResultStats(miniGameRoundResultDailyLimit), now)),
		Weekly:     toMiniGameRoundResultWeeklyPoints(overlayCurrentWeeklyMiniGameRoundResult(statdao.ListRecentWeeklyMiniGameRoundResultStats(miniGameRoundResultWeeklyLimit), now)),
		Monthly:    toMiniGameRoundResultMonthlyPoints(overlayCurrentMonthlyMiniGameRoundResult(statdao.ListRecentMonthlyMiniGameRoundResultStats(miniGameRoundResultMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyMiniGameRoundResult(rows []*statentity.DailyMiniGameRoundResultStat, now time.Time) []*statentity.DailyMiniGameRoundResultStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyMiniGameRoundResultStatByDate(today)
	out := make([]*statentity.DailyMiniGameRoundResultStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyMiniGameRoundResult(rows []*statentity.WeeklyMiniGameRoundResultStat, now time.Time) []*statentity.WeeklyMiniGameRoundResultStat {
	week := statentity.FormatWeeklyMiniGameRoundResultStatKey(now)
	live := statdao.GetWeeklyMiniGameRoundResultStatByWeek(week)
	out := make([]*statentity.WeeklyMiniGameRoundResultStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyMiniGameRoundResult(rows []*statentity.MonthlyMiniGameRoundResultStat, now time.Time) []*statentity.MonthlyMiniGameRoundResultStat {
	month := statentity.FormatMonthlyMiniGameRoundResultStatKey(now)
	live := statdao.GetMonthlyMiniGameRoundResultStatByMonth(month)
	out := make([]*statentity.MonthlyMiniGameRoundResultStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toMiniGameRoundResultDailyPoints(rows []*statentity.DailyMiniGameRoundResultStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toMiniGameRoundResultWeeklyPoints(rows []*statentity.WeeklyMiniGameRoundResultStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toMiniGameRoundResultMonthlyPoints(rows []*statentity.MonthlyMiniGameRoundResultStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
