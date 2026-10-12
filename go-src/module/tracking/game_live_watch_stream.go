package tracking

import (
	"context"
	"time"

	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/statdao"
	"xr-game-server/dto/trackingdto"
	liveentity "xr-game-server/entity/live"
	statentity "xr-game-server/entity/stat"
)

const (
	gameLiveWatchStreamDailyLimit   = 30
	gameLiveWatchStreamWeeklyLimit  = 12
	gameLiveWatchStreamMonthlyLimit = 12

	// EventKeyGameLiveWatchStream 埋点 event_key(对齐埋点表 game_live_watch_stream)
	EventKeyGameLiveWatchStream = "game_live_watch_stream"
)

func recordGameLiveWatchStream(userId, roomId uint64, at time.Time) {
	if userId == 0 || roomId == 0 || userId == roomId {
		return
	}
	cfg := liveroomdao.GetLiveRoomCfg(roomId)
	if cfg == nil || cfg.Category != liveentity.LiveRoomCategoryGame {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyGameLiveWatchStreamStatKey(at)
	month := statentity.FormatMonthlyGameLiveWatchStreamStatKey(at)

	daily := statdao.GetDailyGameLiveWatchStreamStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyGameLiveWatchStreamStat(daily)

	weekly := statdao.GetWeeklyGameLiveWatchStreamStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyGameLiveWatchStreamStat(weekly)

	monthly := statdao.GetMonthlyGameLiveWatchStreamStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyGameLiveWatchStreamStat(monthly)
}

func GetCMSGameLiveWatchStreamTrend(_ context.Context, _ *trackingdto.CMSGameLiveWatchStreamTrendReq) (*trackingdto.CMSGameLiveWatchStreamTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyGameLiveWatchStreamStatKey(now)
	month := statentity.FormatMonthlyGameLiveWatchStreamStatKey(now)

	todayStat := statdao.GetDailyGameLiveWatchStreamStatByDate(today)
	weekStat := statdao.GetWeeklyGameLiveWatchStreamStatByWeek(week)
	monthStat := statdao.GetMonthlyGameLiveWatchStreamStatByMonth(month)

	return &trackingdto.CMSGameLiveWatchStreamTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toGameLiveWatchStreamDailyPoints(overlayCurrentDailyGameLiveWatchStream(statdao.ListRecentDailyGameLiveWatchStreamStats(gameLiveWatchStreamDailyLimit), now)),
		Weekly:     toGameLiveWatchStreamWeeklyPoints(overlayCurrentWeeklyGameLiveWatchStream(statdao.ListRecentWeeklyGameLiveWatchStreamStats(gameLiveWatchStreamWeeklyLimit), now)),
		Monthly:    toGameLiveWatchStreamMonthlyPoints(overlayCurrentMonthlyGameLiveWatchStream(statdao.ListRecentMonthlyGameLiveWatchStreamStats(gameLiveWatchStreamMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyGameLiveWatchStream(rows []*statentity.DailyGameLiveWatchStreamStat, now time.Time) []*statentity.DailyGameLiveWatchStreamStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyGameLiveWatchStreamStatByDate(today)
	out := make([]*statentity.DailyGameLiveWatchStreamStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyGameLiveWatchStream(rows []*statentity.WeeklyGameLiveWatchStreamStat, now time.Time) []*statentity.WeeklyGameLiveWatchStreamStat {
	week := statentity.FormatWeeklyGameLiveWatchStreamStatKey(now)
	live := statdao.GetWeeklyGameLiveWatchStreamStatByWeek(week)
	out := make([]*statentity.WeeklyGameLiveWatchStreamStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyGameLiveWatchStream(rows []*statentity.MonthlyGameLiveWatchStreamStat, now time.Time) []*statentity.MonthlyGameLiveWatchStreamStat {
	month := statentity.FormatMonthlyGameLiveWatchStreamStatKey(now)
	live := statdao.GetMonthlyGameLiveWatchStreamStatByMonth(month)
	out := make([]*statentity.MonthlyGameLiveWatchStreamStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toGameLiveWatchStreamDailyPoints(rows []*statentity.DailyGameLiveWatchStreamStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toGameLiveWatchStreamWeeklyPoints(rows []*statentity.WeeklyGameLiveWatchStreamStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toGameLiveWatchStreamMonthlyPoints(rows []*statentity.MonthlyGameLiveWatchStreamStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
