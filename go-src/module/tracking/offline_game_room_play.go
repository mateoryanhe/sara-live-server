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
	offlineGameRoomPlayDailyLimit   = 30
	offlineGameRoomPlayWeeklyLimit  = 12
	offlineGameRoomPlayMonthlyLimit = 12

	// EventKeyOfflineGameRoomPlay 埋点 event_key(对齐埋点表 offline_game_room_play)
	EventKeyOfflineGameRoomPlay = "offline_game_room_play"
)

func recordOfflineGameRoomPlay(userId, roomId uint64, at time.Time) {
	if userId == 0 || roomId == 0 || userId == roomId {
		return
	}
	cfg := liveroomdao.GetLiveRoomCfg(roomId)
	if cfg == nil || cfg.Category != liveentity.LiveRoomCategoryGame {
		return
	}
	room := liveroomdao.GetRoomById(roomId)
	if room == nil || room.LiveRecordId != 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyOfflineGameRoomPlayStatKey(at)
	month := statentity.FormatMonthlyOfflineGameRoomPlayStatKey(at)

	daily := statdao.GetDailyOfflineGameRoomPlayStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyOfflineGameRoomPlayStat(daily)

	weekly := statdao.GetWeeklyOfflineGameRoomPlayStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyOfflineGameRoomPlayStat(weekly)

	monthly := statdao.GetMonthlyOfflineGameRoomPlayStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyOfflineGameRoomPlayStat(monthly)
}

func GetCMSOfflineGameRoomPlayTrend(_ context.Context, _ *trackingdto.CMSOfflineGameRoomPlayTrendReq) (*trackingdto.CMSOfflineGameRoomPlayTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyOfflineGameRoomPlayStatKey(now)
	month := statentity.FormatMonthlyOfflineGameRoomPlayStatKey(now)

	todayStat := statdao.GetDailyOfflineGameRoomPlayStatByDate(today)
	weekStat := statdao.GetWeeklyOfflineGameRoomPlayStatByWeek(week)
	monthStat := statdao.GetMonthlyOfflineGameRoomPlayStatByMonth(month)

	return &trackingdto.CMSOfflineGameRoomPlayTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toOfflineGameRoomPlayDailyPoints(overlayCurrentDailyOfflineGameRoomPlay(statdao.ListRecentDailyOfflineGameRoomPlayStats(offlineGameRoomPlayDailyLimit), now)),
		Weekly:     toOfflineGameRoomPlayWeeklyPoints(overlayCurrentWeeklyOfflineGameRoomPlay(statdao.ListRecentWeeklyOfflineGameRoomPlayStats(offlineGameRoomPlayWeeklyLimit), now)),
		Monthly:    toOfflineGameRoomPlayMonthlyPoints(overlayCurrentMonthlyOfflineGameRoomPlay(statdao.ListRecentMonthlyOfflineGameRoomPlayStats(offlineGameRoomPlayMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyOfflineGameRoomPlay(rows []*statentity.DailyOfflineGameRoomPlayStat, now time.Time) []*statentity.DailyOfflineGameRoomPlayStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyOfflineGameRoomPlayStatByDate(today)
	out := make([]*statentity.DailyOfflineGameRoomPlayStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyOfflineGameRoomPlay(rows []*statentity.WeeklyOfflineGameRoomPlayStat, now time.Time) []*statentity.WeeklyOfflineGameRoomPlayStat {
	week := statentity.FormatWeeklyOfflineGameRoomPlayStatKey(now)
	live := statdao.GetWeeklyOfflineGameRoomPlayStatByWeek(week)
	out := make([]*statentity.WeeklyOfflineGameRoomPlayStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyOfflineGameRoomPlay(rows []*statentity.MonthlyOfflineGameRoomPlayStat, now time.Time) []*statentity.MonthlyOfflineGameRoomPlayStat {
	month := statentity.FormatMonthlyOfflineGameRoomPlayStatKey(now)
	live := statdao.GetMonthlyOfflineGameRoomPlayStatByMonth(month)
	out := make([]*statentity.MonthlyOfflineGameRoomPlayStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toOfflineGameRoomPlayDailyPoints(rows []*statentity.DailyOfflineGameRoomPlayStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toOfflineGameRoomPlayWeeklyPoints(rows []*statentity.WeeklyOfflineGameRoomPlayStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toOfflineGameRoomPlayMonthlyPoints(rows []*statentity.MonthlyOfflineGameRoomPlayStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
