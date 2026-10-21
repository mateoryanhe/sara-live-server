package tracking

import (
	"context"
	"time"

	"xr-game-server/core/event"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/statdao"
	"xr-game-server/dto/trackingdto"
	statentity "xr-game-server/entity/stat"
	"xr-game-server/errercode"
	"xr-game-server/gameevent"
)

const (
	miniGamePlusClickDailyLimit   = 30
	miniGamePlusClickWeeklyLimit  = 12
	miniGamePlusClickMonthlyLimit = 12

	// EventKeyMiniGamePlusClick 埋点 event_key(对齐埋点表 mini_game_plus_click)
	EventKeyMiniGamePlusClick = "mini_game_plus_click"
)

// ReportMiniGamePlusClick App 上报游戏+(登录即可,每次成功 +1)
func ReportMiniGamePlusClick(ctx context.Context, _ *trackingdto.ReportMiniGamePlusClickReq) (*trackingdto.ReportMiniGamePlusClickRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	event.Pub(gameevent.MiniGamePlusClickTrackingEvent, gameevent.NewMiniGamePlusClickTrackingEventData(userId, time.Now()))
	return &trackingdto.ReportMiniGamePlusClickRes{Success: true}, nil
}

func recordMiniGamePlusClick(userId uint64, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyMiniGamePlusClickStatKey(at)
	month := statentity.FormatMonthlyMiniGamePlusClickStatKey(at)

	daily := statdao.GetDailyMiniGamePlusClickStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyMiniGamePlusClickStat(daily)

	weekly := statdao.GetWeeklyMiniGamePlusClickStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyMiniGamePlusClickStat(weekly)

	monthly := statdao.GetMonthlyMiniGamePlusClickStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyMiniGamePlusClickStat(monthly)
}

func GetCMSMiniGamePlusClickTrend(_ context.Context, _ *trackingdto.CMSMiniGamePlusClickTrendReq) (*trackingdto.CMSMiniGamePlusClickTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyMiniGamePlusClickStatKey(now)
	month := statentity.FormatMonthlyMiniGamePlusClickStatKey(now)

	todayStat := statdao.GetDailyMiniGamePlusClickStatByDate(today)
	weekStat := statdao.GetWeeklyMiniGamePlusClickStatByWeek(week)
	monthStat := statdao.GetMonthlyMiniGamePlusClickStatByMonth(month)

	return &trackingdto.CMSMiniGamePlusClickTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toMiniGamePlusClickDailyPoints(overlayCurrentDailyMiniGamePlusClick(statdao.ListRecentDailyMiniGamePlusClickStats(miniGamePlusClickDailyLimit), now)),
		Weekly:     toMiniGamePlusClickWeeklyPoints(overlayCurrentWeeklyMiniGamePlusClick(statdao.ListRecentWeeklyMiniGamePlusClickStats(miniGamePlusClickWeeklyLimit), now)),
		Monthly:    toMiniGamePlusClickMonthlyPoints(overlayCurrentMonthlyMiniGamePlusClick(statdao.ListRecentMonthlyMiniGamePlusClickStats(miniGamePlusClickMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyMiniGamePlusClick(rows []*statentity.DailyMiniGamePlusClickStat, now time.Time) []*statentity.DailyMiniGamePlusClickStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyMiniGamePlusClickStatByDate(today)
	out := make([]*statentity.DailyMiniGamePlusClickStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyMiniGamePlusClick(rows []*statentity.WeeklyMiniGamePlusClickStat, now time.Time) []*statentity.WeeklyMiniGamePlusClickStat {
	week := statentity.FormatWeeklyMiniGamePlusClickStatKey(now)
	live := statdao.GetWeeklyMiniGamePlusClickStatByWeek(week)
	out := make([]*statentity.WeeklyMiniGamePlusClickStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyMiniGamePlusClick(rows []*statentity.MonthlyMiniGamePlusClickStat, now time.Time) []*statentity.MonthlyMiniGamePlusClickStat {
	month := statentity.FormatMonthlyMiniGamePlusClickStatKey(now)
	live := statdao.GetMonthlyMiniGamePlusClickStatByMonth(month)
	out := make([]*statentity.MonthlyMiniGamePlusClickStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toMiniGamePlusClickDailyPoints(rows []*statentity.DailyMiniGamePlusClickStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toMiniGamePlusClickWeeklyPoints(rows []*statentity.WeeklyMiniGamePlusClickStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toMiniGamePlusClickMonthlyPoints(rows []*statentity.MonthlyMiniGamePlusClickStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
