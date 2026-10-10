package tracking

import (
	"context"
	"time"

	"xr-game-server/dao/statdao"
	"xr-game-server/dto/trackingdto"
	statentity "xr-game-server/entity/stat"
)

const (
	miniGameRoundStartDailyLimit   = 30
	miniGameRoundStartWeeklyLimit  = 12
	miniGameRoundStartMonthlyLimit = 12

	// EventKeyMiniGameRoundStart 埋点 event_key(对齐埋点表 mini_game_round_start)
	EventKeyMiniGameRoundStart = "mini_game_round_start"
)

func recordMiniGameRoundStart(userId uint64, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyMiniGameRoundStartStatKey(at)
	month := statentity.FormatMonthlyMiniGameRoundStartStatKey(at)

	daily := statdao.GetDailyMiniGameRoundStartStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyMiniGameRoundStartStat(daily)

	weekly := statdao.GetWeeklyMiniGameRoundStartStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyMiniGameRoundStartStat(weekly)

	monthly := statdao.GetMonthlyMiniGameRoundStartStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyMiniGameRoundStartStat(monthly)
}

func GetCMSMiniGameRoundStartTrend(_ context.Context, _ *trackingdto.CMSMiniGameRoundStartTrendReq) (*trackingdto.CMSMiniGameRoundStartTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyMiniGameRoundStartStatKey(now)
	month := statentity.FormatMonthlyMiniGameRoundStartStatKey(now)

	todayStat := statdao.GetDailyMiniGameRoundStartStatByDate(today)
	weekStat := statdao.GetWeeklyMiniGameRoundStartStatByWeek(week)
	monthStat := statdao.GetMonthlyMiniGameRoundStartStatByMonth(month)

	return &trackingdto.CMSMiniGameRoundStartTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toMiniGameRoundStartDailyPoints(overlayCurrentDailyMiniGameRoundStart(statdao.ListRecentDailyMiniGameRoundStartStats(miniGameRoundStartDailyLimit), now)),
		Weekly:     toMiniGameRoundStartWeeklyPoints(overlayCurrentWeeklyMiniGameRoundStart(statdao.ListRecentWeeklyMiniGameRoundStartStats(miniGameRoundStartWeeklyLimit), now)),
		Monthly:    toMiniGameRoundStartMonthlyPoints(overlayCurrentMonthlyMiniGameRoundStart(statdao.ListRecentMonthlyMiniGameRoundStartStats(miniGameRoundStartMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyMiniGameRoundStart(rows []*statentity.DailyMiniGameRoundStartStat, now time.Time) []*statentity.DailyMiniGameRoundStartStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyMiniGameRoundStartStatByDate(today)
	out := make([]*statentity.DailyMiniGameRoundStartStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyMiniGameRoundStart(rows []*statentity.WeeklyMiniGameRoundStartStat, now time.Time) []*statentity.WeeklyMiniGameRoundStartStat {
	week := statentity.FormatWeeklyMiniGameRoundStartStatKey(now)
	live := statdao.GetWeeklyMiniGameRoundStartStatByWeek(week)
	out := make([]*statentity.WeeklyMiniGameRoundStartStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyMiniGameRoundStart(rows []*statentity.MonthlyMiniGameRoundStartStat, now time.Time) []*statentity.MonthlyMiniGameRoundStartStat {
	month := statentity.FormatMonthlyMiniGameRoundStartStatKey(now)
	live := statdao.GetMonthlyMiniGameRoundStartStatByMonth(month)
	out := make([]*statentity.MonthlyMiniGameRoundStartStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toMiniGameRoundStartDailyPoints(rows []*statentity.DailyMiniGameRoundStartStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toMiniGameRoundStartWeeklyPoints(rows []*statentity.WeeklyMiniGameRoundStartStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toMiniGameRoundStartMonthlyPoints(rows []*statentity.MonthlyMiniGameRoundStartStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
