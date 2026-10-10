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
	// EventKeyGameLiveRoomJoin 埋点 event_key(对齐埋点表 game_live_room_join)
	EventKeyGameLiveRoomJoin = "game_live_room_join"

	gameLiveRoomJoinDailyLimit   = 30
	gameLiveRoomJoinWeeklyLimit  = 12
	gameLiveRoomJoinMonthlyLimit = 12
)

// recordGameLiveRoomJoin 消费埋点事件:游戏直播类(category=2)进房次数 +1(含重复进入,不含主播本人).
func recordGameLiveRoomJoin(userId, roomId uint64, at time.Time) {
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
	week := statentity.FormatWeeklyGameLiveRoomJoinStatKey(at)
	month := statentity.FormatMonthlyGameLiveRoomJoinStatKey(at)

	daily := statdao.GetDailyGameLiveRoomJoinStatByDate(date)
	daily.AddJoinCount(1)
	statdao.PublishDailyGameLiveRoomJoinStat(daily)

	weekly := statdao.GetWeeklyGameLiveRoomJoinStatByWeek(week)
	weekly.AddJoinCount(1)
	statdao.PublishWeeklyGameLiveRoomJoinStat(weekly)

	monthly := statdao.GetMonthlyGameLiveRoomJoinStatByMonth(month)
	monthly.AddJoinCount(1)
	statdao.PublishMonthlyGameLiveRoomJoinStat(monthly)
}

// GetCMSGameLiveRoomJoinTrend CMS 埋点:进入游戏直播类直播间次数趋势
func GetCMSGameLiveRoomJoinTrend(_ context.Context, _ *trackingdto.CMSGameLiveRoomJoinTrendReq) (*trackingdto.CMSGameLiveRoomJoinTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyGameLiveRoomJoinStatKey(now)
	month := statentity.FormatMonthlyGameLiveRoomJoinStatKey(now)

	todayStat := statdao.GetDailyGameLiveRoomJoinStatByDate(today)
	weekStat := statdao.GetWeeklyGameLiveRoomJoinStatByWeek(week)
	monthStat := statdao.GetMonthlyGameLiveRoomJoinStatByMonth(month)

	return &trackingdto.CMSGameLiveRoomJoinTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toGameJoinDailyPoints(overlayCurrentDailyGameJoin(statdao.ListRecentDailyGameLiveRoomJoinStats(gameLiveRoomJoinDailyLimit), now)),
		Weekly:     toGameJoinWeeklyPoints(overlayCurrentWeeklyGameJoin(statdao.ListRecentWeeklyGameLiveRoomJoinStats(gameLiveRoomJoinWeeklyLimit), now)),
		Monthly:    toGameJoinMonthlyPoints(overlayCurrentMonthlyGameJoin(statdao.ListRecentMonthlyGameLiveRoomJoinStats(gameLiveRoomJoinMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyGameJoin(rows []*statentity.DailyGameLiveRoomJoinStat, now time.Time) []*statentity.DailyGameLiveRoomJoinStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyGameLiveRoomJoinStatByDate(today)
	out := make([]*statentity.DailyGameLiveRoomJoinStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyGameJoin(rows []*statentity.WeeklyGameLiveRoomJoinStat, now time.Time) []*statentity.WeeklyGameLiveRoomJoinStat {
	week := statentity.FormatWeeklyGameLiveRoomJoinStatKey(now)
	live := statdao.GetWeeklyGameLiveRoomJoinStatByWeek(week)
	out := make([]*statentity.WeeklyGameLiveRoomJoinStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyGameJoin(rows []*statentity.MonthlyGameLiveRoomJoinStat, now time.Time) []*statentity.MonthlyGameLiveRoomJoinStat {
	month := statentity.FormatMonthlyGameLiveRoomJoinStatKey(now)
	live := statdao.GetMonthlyGameLiveRoomJoinStatByMonth(month)
	out := make([]*statentity.MonthlyGameLiveRoomJoinStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toGameJoinDailyPoints(rows []*statentity.DailyGameLiveRoomJoinStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toGameJoinWeeklyPoints(rows []*statentity.WeeklyGameLiveRoomJoinStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toGameJoinMonthlyPoints(rows []*statentity.MonthlyGameLiveRoomJoinStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
