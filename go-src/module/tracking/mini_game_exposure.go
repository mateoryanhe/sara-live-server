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
	miniGameExposureDailyLimit   = 30
	miniGameExposureWeeklyLimit  = 12
	miniGameExposureMonthlyLimit = 12

	// EventKeyMiniGameExposure 埋点 event_key(对齐埋点表 mini_game_exposure)
	EventKeyMiniGameExposure = "mini_game_exposure"
)

// ReportMiniGameExposure App 上报半屏游戏窗口曝光(登录即可,每次成功 +1)
func ReportMiniGameExposure(ctx context.Context, _ *trackingdto.ReportMiniGameExposureReq) (*trackingdto.ReportMiniGameExposureRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	event.Pub(gameevent.MiniGameExposureTrackingEvent, gameevent.NewMiniGameExposureTrackingEventData(userId, time.Now()))
	return &trackingdto.ReportMiniGameExposureRes{Success: true}, nil
}

func recordMiniGameExposure(userId uint64, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyMiniGameExposureStatKey(at)
	month := statentity.FormatMonthlyMiniGameExposureStatKey(at)

	daily := statdao.GetDailyMiniGameExposureStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyMiniGameExposureStat(daily)

	weekly := statdao.GetWeeklyMiniGameExposureStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyMiniGameExposureStat(weekly)

	monthly := statdao.GetMonthlyMiniGameExposureStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyMiniGameExposureStat(monthly)
}

func GetCMSMiniGameExposureTrend(_ context.Context, _ *trackingdto.CMSMiniGameExposureTrendReq) (*trackingdto.CMSMiniGameExposureTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyMiniGameExposureStatKey(now)
	month := statentity.FormatMonthlyMiniGameExposureStatKey(now)

	todayStat := statdao.GetDailyMiniGameExposureStatByDate(today)
	weekStat := statdao.GetWeeklyMiniGameExposureStatByWeek(week)
	monthStat := statdao.GetMonthlyMiniGameExposureStatByMonth(month)

	return &trackingdto.CMSMiniGameExposureTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toMiniGameExposureDailyPoints(overlayCurrentDailyMiniGameExposure(statdao.ListRecentDailyMiniGameExposureStats(miniGameExposureDailyLimit), now)),
		Weekly:     toMiniGameExposureWeeklyPoints(overlayCurrentWeeklyMiniGameExposure(statdao.ListRecentWeeklyMiniGameExposureStats(miniGameExposureWeeklyLimit), now)),
		Monthly:    toMiniGameExposureMonthlyPoints(overlayCurrentMonthlyMiniGameExposure(statdao.ListRecentMonthlyMiniGameExposureStats(miniGameExposureMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyMiniGameExposure(rows []*statentity.DailyMiniGameExposureStat, now time.Time) []*statentity.DailyMiniGameExposureStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyMiniGameExposureStatByDate(today)
	out := make([]*statentity.DailyMiniGameExposureStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyMiniGameExposure(rows []*statentity.WeeklyMiniGameExposureStat, now time.Time) []*statentity.WeeklyMiniGameExposureStat {
	week := statentity.FormatWeeklyMiniGameExposureStatKey(now)
	live := statdao.GetWeeklyMiniGameExposureStatByWeek(week)
	out := make([]*statentity.WeeklyMiniGameExposureStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyMiniGameExposure(rows []*statentity.MonthlyMiniGameExposureStat, now time.Time) []*statentity.MonthlyMiniGameExposureStat {
	month := statentity.FormatMonthlyMiniGameExposureStatKey(now)
	live := statdao.GetMonthlyMiniGameExposureStatByMonth(month)
	out := make([]*statentity.MonthlyMiniGameExposureStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toMiniGameExposureDailyPoints(rows []*statentity.DailyMiniGameExposureStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toMiniGameExposureWeeklyPoints(rows []*statentity.WeeklyMiniGameExposureStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toMiniGameExposureMonthlyPoints(rows []*statentity.MonthlyMiniGameExposureStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
