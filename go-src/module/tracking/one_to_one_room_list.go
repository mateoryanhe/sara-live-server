package tracking

import (
	"context"
	"time"

	"xr-game-server/dao/statdao"
	"xr-game-server/dto/trackingdto"
	statentity "xr-game-server/entity/stat"
)

const (
	oneToOneRoomListDailyLimit   = 30
	oneToOneRoomListWeeklyLimit  = 12
	oneToOneRoomListMonthlyLimit = 12

	// EventKeyOneToOneRoomList 埋点 event_key(对齐埋点表 one_to_one_room_list)
	EventKeyOneToOneRoomList = "one_to_one_room_list"
)

func recordOneToOneRoomList(userId uint64, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyOneToOneRoomListStatKey(at)
	month := statentity.FormatMonthlyOneToOneRoomListStatKey(at)

	daily := statdao.GetDailyOneToOneRoomListStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyOneToOneRoomListStat(daily)

	weekly := statdao.GetWeeklyOneToOneRoomListStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyOneToOneRoomListStat(weekly)

	monthly := statdao.GetMonthlyOneToOneRoomListStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyOneToOneRoomListStat(monthly)
}

func GetCMSOneToOneRoomListTrend(_ context.Context, _ *trackingdto.CMSOneToOneRoomListTrendReq) (*trackingdto.CMSOneToOneRoomListTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyOneToOneRoomListStatKey(now)
	month := statentity.FormatMonthlyOneToOneRoomListStatKey(now)

	todayStat := statdao.GetDailyOneToOneRoomListStatByDate(today)
	weekStat := statdao.GetWeeklyOneToOneRoomListStatByWeek(week)
	monthStat := statdao.GetMonthlyOneToOneRoomListStatByMonth(month)

	return &trackingdto.CMSOneToOneRoomListTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toOneToOneRoomListDailyPoints(overlayCurrentDailyOneToOneRoomList(statdao.ListRecentDailyOneToOneRoomListStats(oneToOneRoomListDailyLimit), now)),
		Weekly:     toOneToOneRoomListWeeklyPoints(overlayCurrentWeeklyOneToOneRoomList(statdao.ListRecentWeeklyOneToOneRoomListStats(oneToOneRoomListWeeklyLimit), now)),
		Monthly:    toOneToOneRoomListMonthlyPoints(overlayCurrentMonthlyOneToOneRoomList(statdao.ListRecentMonthlyOneToOneRoomListStats(oneToOneRoomListMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyOneToOneRoomList(rows []*statentity.DailyOneToOneRoomListStat, now time.Time) []*statentity.DailyOneToOneRoomListStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyOneToOneRoomListStatByDate(today)
	out := make([]*statentity.DailyOneToOneRoomListStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyOneToOneRoomList(rows []*statentity.WeeklyOneToOneRoomListStat, now time.Time) []*statentity.WeeklyOneToOneRoomListStat {
	week := statentity.FormatWeeklyOneToOneRoomListStatKey(now)
	live := statdao.GetWeeklyOneToOneRoomListStatByWeek(week)
	out := make([]*statentity.WeeklyOneToOneRoomListStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyOneToOneRoomList(rows []*statentity.MonthlyOneToOneRoomListStat, now time.Time) []*statentity.MonthlyOneToOneRoomListStat {
	month := statentity.FormatMonthlyOneToOneRoomListStatKey(now)
	live := statdao.GetMonthlyOneToOneRoomListStatByMonth(month)
	out := make([]*statentity.MonthlyOneToOneRoomListStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toOneToOneRoomListDailyPoints(rows []*statentity.DailyOneToOneRoomListStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toOneToOneRoomListWeeklyPoints(rows []*statentity.WeeklyOneToOneRoomListStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toOneToOneRoomListMonthlyPoints(rows []*statentity.MonthlyOneToOneRoomListStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
