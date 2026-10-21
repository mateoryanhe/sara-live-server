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
	miniGameSpinClickDailyLimit   = 30
	miniGameSpinClickWeeklyLimit  = 12
	miniGameSpinClickMonthlyLimit = 12

	// EventKeyMiniGameSpinClick 埋点 event_key(对齐埋点表 mini_game_spin_click)
	EventKeyMiniGameSpinClick = "mini_game_spin_click"
)

// ReportMiniGameSpinClick App 上报 Spin 点击(登录即可,每次成功 +1)
func ReportMiniGameSpinClick(ctx context.Context, _ *trackingdto.ReportMiniGameSpinClickReq) (*trackingdto.ReportMiniGameSpinClickRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	event.Pub(gameevent.MiniGameSpinClickTrackingEvent, gameevent.NewMiniGameSpinClickTrackingEventData(userId, time.Now()))
	return &trackingdto.ReportMiniGameSpinClickRes{Success: true}, nil
}

func recordMiniGameSpinClick(userId uint64, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyMiniGameSpinClickStatKey(at)
	month := statentity.FormatMonthlyMiniGameSpinClickStatKey(at)

	daily := statdao.GetDailyMiniGameSpinClickStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyMiniGameSpinClickStat(daily)

	weekly := statdao.GetWeeklyMiniGameSpinClickStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyMiniGameSpinClickStat(weekly)

	monthly := statdao.GetMonthlyMiniGameSpinClickStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyMiniGameSpinClickStat(monthly)
}

func GetCMSMiniGameSpinClickTrend(_ context.Context, _ *trackingdto.CMSMiniGameSpinClickTrendReq) (*trackingdto.CMSMiniGameSpinClickTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyMiniGameSpinClickStatKey(now)
	month := statentity.FormatMonthlyMiniGameSpinClickStatKey(now)

	todayStat := statdao.GetDailyMiniGameSpinClickStatByDate(today)
	weekStat := statdao.GetWeeklyMiniGameSpinClickStatByWeek(week)
	monthStat := statdao.GetMonthlyMiniGameSpinClickStatByMonth(month)

	return &trackingdto.CMSMiniGameSpinClickTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toMiniGameSpinClickDailyPoints(overlayCurrentDailyMiniGameSpinClick(statdao.ListRecentDailyMiniGameSpinClickStats(miniGameSpinClickDailyLimit), now)),
		Weekly:     toMiniGameSpinClickWeeklyPoints(overlayCurrentWeeklyMiniGameSpinClick(statdao.ListRecentWeeklyMiniGameSpinClickStats(miniGameSpinClickWeeklyLimit), now)),
		Monthly:    toMiniGameSpinClickMonthlyPoints(overlayCurrentMonthlyMiniGameSpinClick(statdao.ListRecentMonthlyMiniGameSpinClickStats(miniGameSpinClickMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyMiniGameSpinClick(rows []*statentity.DailyMiniGameSpinClickStat, now time.Time) []*statentity.DailyMiniGameSpinClickStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyMiniGameSpinClickStatByDate(today)
	out := make([]*statentity.DailyMiniGameSpinClickStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyMiniGameSpinClick(rows []*statentity.WeeklyMiniGameSpinClickStat, now time.Time) []*statentity.WeeklyMiniGameSpinClickStat {
	week := statentity.FormatWeeklyMiniGameSpinClickStatKey(now)
	live := statdao.GetWeeklyMiniGameSpinClickStatByWeek(week)
	out := make([]*statentity.WeeklyMiniGameSpinClickStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyMiniGameSpinClick(rows []*statentity.MonthlyMiniGameSpinClickStat, now time.Time) []*statentity.MonthlyMiniGameSpinClickStat {
	month := statentity.FormatMonthlyMiniGameSpinClickStatKey(now)
	live := statdao.GetMonthlyMiniGameSpinClickStatByMonth(month)
	out := make([]*statentity.MonthlyMiniGameSpinClickStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toMiniGameSpinClickDailyPoints(rows []*statentity.DailyMiniGameSpinClickStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toMiniGameSpinClickWeeklyPoints(rows []*statentity.WeeklyMiniGameSpinClickStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toMiniGameSpinClickMonthlyPoints(rows []*statentity.MonthlyMiniGameSpinClickStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
