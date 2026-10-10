package tracking

import (
	"context"
	"time"

	"xr-game-server/dao/statdao"
	"xr-game-server/dto/trackingdto"
	statentity "xr-game-server/entity/stat"
)

const (
	call1v1RoomDailyLimit   = 30
	call1v1RoomWeeklyLimit  = 12
	call1v1RoomMonthlyLimit = 12

	// EventKeyCall1v1Room 埋点 event_key(对齐埋点表 call_1v1_room)
	EventKeyCall1v1Room = "call_1v1_room"
)

// recordCall1v1Room 消费埋点:1v1 房间 source=3 视频通话发起 +1.
func recordCall1v1Room(callerId, anchorId uint64, at time.Time) {
	if callerId == 0 || anchorId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyCall1v1RoomStatKey(at)
	month := statentity.FormatMonthlyCall1v1RoomStatKey(at)

	daily := statdao.GetDailyCall1v1RoomStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyCall1v1RoomStat(daily)

	weekly := statdao.GetWeeklyCall1v1RoomStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyCall1v1RoomStat(weekly)

	monthly := statdao.GetMonthlyCall1v1RoomStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyCall1v1RoomStat(monthly)
}

func GetCMSCall1v1RoomCallTrend(_ context.Context, _ *trackingdto.CMSCall1v1RoomCallTrendReq) (*trackingdto.CMSCall1v1RoomCallTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyCall1v1RoomStatKey(now)
	month := statentity.FormatMonthlyCall1v1RoomStatKey(now)

	todayStat := statdao.GetDailyCall1v1RoomStatByDate(today)
	weekStat := statdao.GetWeeklyCall1v1RoomStatByWeek(week)
	monthStat := statdao.GetMonthlyCall1v1RoomStatByMonth(month)

	return &trackingdto.CMSCall1v1RoomCallTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toCall1v1RoomDailyPoints(overlayCurrentDailyCall1v1Room(statdao.ListRecentDailyCall1v1RoomStats(call1v1RoomDailyLimit), now)),
		Weekly:     toCall1v1RoomWeeklyPoints(overlayCurrentWeeklyCall1v1Room(statdao.ListRecentWeeklyCall1v1RoomStats(call1v1RoomWeeklyLimit), now)),
		Monthly:    toCall1v1RoomMonthlyPoints(overlayCurrentMonthlyCall1v1Room(statdao.ListRecentMonthlyCall1v1RoomStats(call1v1RoomMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyCall1v1Room(rows []*statentity.DailyCall1v1RoomStat, now time.Time) []*statentity.DailyCall1v1RoomStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyCall1v1RoomStatByDate(today)
	out := make([]*statentity.DailyCall1v1RoomStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyCall1v1Room(rows []*statentity.WeeklyCall1v1RoomStat, now time.Time) []*statentity.WeeklyCall1v1RoomStat {
	week := statentity.FormatWeeklyCall1v1RoomStatKey(now)
	live := statdao.GetWeeklyCall1v1RoomStatByWeek(week)
	out := make([]*statentity.WeeklyCall1v1RoomStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyCall1v1Room(rows []*statentity.MonthlyCall1v1RoomStat, now time.Time) []*statentity.MonthlyCall1v1RoomStat {
	month := statentity.FormatMonthlyCall1v1RoomStatKey(now)
	live := statdao.GetMonthlyCall1v1RoomStatByMonth(month)
	out := make([]*statentity.MonthlyCall1v1RoomStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toCall1v1RoomDailyPoints(rows []*statentity.DailyCall1v1RoomStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toCall1v1RoomWeeklyPoints(rows []*statentity.WeeklyCall1v1RoomStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toCall1v1RoomMonthlyPoints(rows []*statentity.MonthlyCall1v1RoomStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
