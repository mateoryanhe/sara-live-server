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
	miniGameWebViewLoadDailyLimit   = 30
	miniGameWebViewLoadWeeklyLimit  = 12
	miniGameWebViewLoadMonthlyLimit = 12

	// EventKeyMiniGameWebViewLoadSuccess 埋点 event_key(对齐埋点表 mini_game_webview_load_success)
	EventKeyMiniGameWebViewLoadSuccess = "mini_game_webview_load_success"
)

// ReportMiniGameWebViewLoadSuccess App 上报游戏 WebView 成功加载(登录即可,每次成功 +1)
func ReportMiniGameWebViewLoadSuccess(ctx context.Context, _ *trackingdto.ReportMiniGameWebViewLoadSuccessReq) (*trackingdto.ReportMiniGameWebViewLoadSuccessRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	event.Pub(gameevent.MiniGameWebViewLoadTrackingEvent, gameevent.NewMiniGameWebViewLoadTrackingEventData(userId, time.Now()))
	return &trackingdto.ReportMiniGameWebViewLoadSuccessRes{Success: true}, nil
}

func recordMiniGameWebViewLoad(userId uint64, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyMiniGameWebViewLoadStatKey(at)
	month := statentity.FormatMonthlyMiniGameWebViewLoadStatKey(at)

	daily := statdao.GetDailyMiniGameWebViewLoadStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyMiniGameWebViewLoadStat(daily)

	weekly := statdao.GetWeeklyMiniGameWebViewLoadStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyMiniGameWebViewLoadStat(weekly)

	monthly := statdao.GetMonthlyMiniGameWebViewLoadStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyMiniGameWebViewLoadStat(monthly)
}

func GetCMSMiniGameWebViewLoadTrend(_ context.Context, _ *trackingdto.CMSMiniGameWebViewLoadTrendReq) (*trackingdto.CMSMiniGameWebViewLoadTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyMiniGameWebViewLoadStatKey(now)
	month := statentity.FormatMonthlyMiniGameWebViewLoadStatKey(now)

	todayStat := statdao.GetDailyMiniGameWebViewLoadStatByDate(today)
	weekStat := statdao.GetWeeklyMiniGameWebViewLoadStatByWeek(week)
	monthStat := statdao.GetMonthlyMiniGameWebViewLoadStatByMonth(month)

	return &trackingdto.CMSMiniGameWebViewLoadTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toMiniGameWebViewLoadDailyPoints(overlayCurrentDailyMiniGameWebViewLoad(statdao.ListRecentDailyMiniGameWebViewLoadStats(miniGameWebViewLoadDailyLimit), now)),
		Weekly:     toMiniGameWebViewLoadWeeklyPoints(overlayCurrentWeeklyMiniGameWebViewLoad(statdao.ListRecentWeeklyMiniGameWebViewLoadStats(miniGameWebViewLoadWeeklyLimit), now)),
		Monthly:    toMiniGameWebViewLoadMonthlyPoints(overlayCurrentMonthlyMiniGameWebViewLoad(statdao.ListRecentMonthlyMiniGameWebViewLoadStats(miniGameWebViewLoadMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyMiniGameWebViewLoad(rows []*statentity.DailyMiniGameWebViewLoadStat, now time.Time) []*statentity.DailyMiniGameWebViewLoadStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyMiniGameWebViewLoadStatByDate(today)
	out := make([]*statentity.DailyMiniGameWebViewLoadStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyMiniGameWebViewLoad(rows []*statentity.WeeklyMiniGameWebViewLoadStat, now time.Time) []*statentity.WeeklyMiniGameWebViewLoadStat {
	week := statentity.FormatWeeklyMiniGameWebViewLoadStatKey(now)
	live := statdao.GetWeeklyMiniGameWebViewLoadStatByWeek(week)
	out := make([]*statentity.WeeklyMiniGameWebViewLoadStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyMiniGameWebViewLoad(rows []*statentity.MonthlyMiniGameWebViewLoadStat, now time.Time) []*statentity.MonthlyMiniGameWebViewLoadStat {
	month := statentity.FormatMonthlyMiniGameWebViewLoadStatKey(now)
	live := statdao.GetMonthlyMiniGameWebViewLoadStatByMonth(month)
	out := make([]*statentity.MonthlyMiniGameWebViewLoadStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toMiniGameWebViewLoadDailyPoints(rows []*statentity.DailyMiniGameWebViewLoadStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toMiniGameWebViewLoadWeeklyPoints(rows []*statentity.WeeklyMiniGameWebViewLoadStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toMiniGameWebViewLoadMonthlyPoints(rows []*statentity.MonthlyMiniGameWebViewLoadStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
