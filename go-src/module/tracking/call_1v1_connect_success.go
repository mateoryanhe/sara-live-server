package tracking

import (
	"context"
	"time"

	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/statdao"
	"xr-game-server/dto/trackingdto"
	callentity "xr-game-server/entity/call"
	liveentity "xr-game-server/entity/live"
	statentity "xr-game-server/entity/stat"
)

const (
	call1v1ConnectSuccessDailyLimit   = 30
	call1v1ConnectSuccessWeeklyLimit  = 12
	call1v1ConnectSuccessMonthlyLimit = 12

	// EventKeyCall1v1ConnectSuccess 埋点 event_key(对齐埋点表 call_1v1_connect_success)
	EventKeyCall1v1ConnectSuccess = "call_1v1_connect_success"
)

func recordCall1v1ConnectSuccess(source, callType uint8, callerId, receiverId uint64, at time.Time) {
	if callType != callentity.CallOrderTypeVideo || callerId == 0 || receiverId == 0 {
		return
	}
	switch source {
	case callentity.CallOrderSourceOneToOneRoom:
	case callentity.CallOrderSourceLiveRoom:
		cfg := liveroomdao.GetLiveRoomCfg(receiverId)
		if cfg == nil || cfg.Category != liveentity.LiveRoomCategoryHot {
			return
		}
	default:
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyCall1v1ConnectSuccessStatKey(at)
	month := statentity.FormatMonthlyCall1v1ConnectSuccessStatKey(at)

	daily := statdao.GetDailyCall1v1ConnectSuccessStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyCall1v1ConnectSuccessStat(daily)

	weekly := statdao.GetWeeklyCall1v1ConnectSuccessStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyCall1v1ConnectSuccessStat(weekly)

	monthly := statdao.GetMonthlyCall1v1ConnectSuccessStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyCall1v1ConnectSuccessStat(monthly)
}

func GetCMSCall1v1ConnectSuccessTrend(_ context.Context, _ *trackingdto.CMSCall1v1ConnectSuccessTrendReq) (*trackingdto.CMSCall1v1ConnectSuccessTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyCall1v1ConnectSuccessStatKey(now)
	month := statentity.FormatMonthlyCall1v1ConnectSuccessStatKey(now)

	todayStat := statdao.GetDailyCall1v1ConnectSuccessStatByDate(today)
	weekStat := statdao.GetWeeklyCall1v1ConnectSuccessStatByWeek(week)
	monthStat := statdao.GetMonthlyCall1v1ConnectSuccessStatByMonth(month)

	return &trackingdto.CMSCall1v1ConnectSuccessTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toCall1v1ConnectSuccessDailyPoints(overlayCurrentDailyCall1v1ConnectSuccess(statdao.ListRecentDailyCall1v1ConnectSuccessStats(call1v1ConnectSuccessDailyLimit), now)),
		Weekly:     toCall1v1ConnectSuccessWeeklyPoints(overlayCurrentWeeklyCall1v1ConnectSuccess(statdao.ListRecentWeeklyCall1v1ConnectSuccessStats(call1v1ConnectSuccessWeeklyLimit), now)),
		Monthly:    toCall1v1ConnectSuccessMonthlyPoints(overlayCurrentMonthlyCall1v1ConnectSuccess(statdao.ListRecentMonthlyCall1v1ConnectSuccessStats(call1v1ConnectSuccessMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyCall1v1ConnectSuccess(rows []*statentity.DailyCall1v1ConnectSuccessStat, now time.Time) []*statentity.DailyCall1v1ConnectSuccessStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyCall1v1ConnectSuccessStatByDate(today)
	out := make([]*statentity.DailyCall1v1ConnectSuccessStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyCall1v1ConnectSuccess(rows []*statentity.WeeklyCall1v1ConnectSuccessStat, now time.Time) []*statentity.WeeklyCall1v1ConnectSuccessStat {
	week := statentity.FormatWeeklyCall1v1ConnectSuccessStatKey(now)
	live := statdao.GetWeeklyCall1v1ConnectSuccessStatByWeek(week)
	out := make([]*statentity.WeeklyCall1v1ConnectSuccessStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyCall1v1ConnectSuccess(rows []*statentity.MonthlyCall1v1ConnectSuccessStat, now time.Time) []*statentity.MonthlyCall1v1ConnectSuccessStat {
	month := statentity.FormatMonthlyCall1v1ConnectSuccessStatKey(now)
	live := statdao.GetMonthlyCall1v1ConnectSuccessStatByMonth(month)
	out := make([]*statentity.MonthlyCall1v1ConnectSuccessStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toCall1v1ConnectSuccessDailyPoints(rows []*statentity.DailyCall1v1ConnectSuccessStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toCall1v1ConnectSuccessWeeklyPoints(rows []*statentity.WeeklyCall1v1ConnectSuccessStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toCall1v1ConnectSuccessMonthlyPoints(rows []*statentity.MonthlyCall1v1ConnectSuccessStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
