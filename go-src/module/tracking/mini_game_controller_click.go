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
	miniGameControllerClickDailyLimit   = 30
	miniGameControllerClickWeeklyLimit  = 12
	miniGameControllerClickMonthlyLimit = 12

	// EventKeyMiniGameControllerClick 埋点 event_key(对齐埋点表 mini_game_controller_click)
	EventKeyMiniGameControllerClick = "mini_game_controller_click"
)

// ReportMiniGameControllerClick App 上报点击游戏手柄(登录即可,每次成功 +1)
func ReportMiniGameControllerClick(ctx context.Context, _ *trackingdto.ReportMiniGameControllerClickReq) (*trackingdto.ReportMiniGameControllerClickRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if gameLiveRoomTrackingPersistDisabled() {
		return &trackingdto.ReportMiniGameControllerClickRes{Success: true}, nil
	}
	event.Pub(gameevent.MiniGameControllerClickTrackingEvent, gameevent.NewMiniGameControllerClickTrackingEventData(userId, time.Now()))
	return &trackingdto.ReportMiniGameControllerClickRes{Success: true}, nil
}

func recordMiniGameControllerClick(userId uint64, at time.Time) {
	if gameLiveRoomTrackingPersistDisabled() {
		return
	}
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyMiniGameControllerClickStatKey(at)
	month := statentity.FormatMonthlyMiniGameControllerClickStatKey(at)

	daily := statdao.GetDailyMiniGameControllerClickStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyMiniGameControllerClickStat(daily)

	weekly := statdao.GetWeeklyMiniGameControllerClickStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyMiniGameControllerClickStat(weekly)

	monthly := statdao.GetMonthlyMiniGameControllerClickStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyMiniGameControllerClickStat(monthly)
}

func GetCMSMiniGameControllerClickTrend(_ context.Context, _ *trackingdto.CMSMiniGameControllerClickTrendReq) (*trackingdto.CMSMiniGameControllerClickTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyMiniGameControllerClickStatKey(now)
	month := statentity.FormatMonthlyMiniGameControllerClickStatKey(now)

	todayStat := statdao.GetDailyMiniGameControllerClickStatByDate(today)
	weekStat := statdao.GetWeeklyMiniGameControllerClickStatByWeek(week)
	monthStat := statdao.GetMonthlyMiniGameControllerClickStatByMonth(month)

	return &trackingdto.CMSMiniGameControllerClickTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toMiniGameControllerClickDailyPoints(overlayCurrentDailyMiniGameControllerClick(statdao.ListRecentDailyMiniGameControllerClickStats(miniGameControllerClickDailyLimit), now)),
		Weekly:     toMiniGameControllerClickWeeklyPoints(overlayCurrentWeeklyMiniGameControllerClick(statdao.ListRecentWeeklyMiniGameControllerClickStats(miniGameControllerClickWeeklyLimit), now)),
		Monthly:    toMiniGameControllerClickMonthlyPoints(overlayCurrentMonthlyMiniGameControllerClick(statdao.ListRecentMonthlyMiniGameControllerClickStats(miniGameControllerClickMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyMiniGameControllerClick(rows []*statentity.DailyMiniGameControllerClickStat, now time.Time) []*statentity.DailyMiniGameControllerClickStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyMiniGameControllerClickStatByDate(today)
	out := make([]*statentity.DailyMiniGameControllerClickStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyMiniGameControllerClick(rows []*statentity.WeeklyMiniGameControllerClickStat, now time.Time) []*statentity.WeeklyMiniGameControllerClickStat {
	week := statentity.FormatWeeklyMiniGameControllerClickStatKey(now)
	live := statdao.GetWeeklyMiniGameControllerClickStatByWeek(week)
	out := make([]*statentity.WeeklyMiniGameControllerClickStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyMiniGameControllerClick(rows []*statentity.MonthlyMiniGameControllerClickStat, now time.Time) []*statentity.MonthlyMiniGameControllerClickStat {
	month := statentity.FormatMonthlyMiniGameControllerClickStatKey(now)
	live := statdao.GetMonthlyMiniGameControllerClickStatByMonth(month)
	out := make([]*statentity.MonthlyMiniGameControllerClickStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toMiniGameControllerClickDailyPoints(rows []*statentity.DailyMiniGameControllerClickStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toMiniGameControllerClickWeeklyPoints(rows []*statentity.WeeklyMiniGameControllerClickStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toMiniGameControllerClickMonthlyPoints(rows []*statentity.MonthlyMiniGameControllerClickStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
