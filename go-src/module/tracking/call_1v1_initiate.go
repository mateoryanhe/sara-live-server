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
	call1v1InitiateDailyLimit   = 30
	call1v1InitiateWeeklyLimit  = 12
	call1v1InitiateMonthlyLimit = 12

	// EventKeyCall1v1Initiate 埋点 event_key(对齐埋点表 call_1v1_initiate)
	EventKeyCall1v1Initiate = "call_1v1_initiate"
)

// recordCall1v1Initiate 消费埋点:秀场直播间 source=1 视频通话发起 +1.
func recordCall1v1Initiate(callerId, anchorId uint64, at time.Time) {
	if callerId == 0 || anchorId == 0 || callerId == anchorId {
		return
	}
	cfg := liveroomdao.GetLiveRoomCfg(anchorId)
	if cfg == nil || cfg.Category != liveentity.LiveRoomCategoryHot {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyCall1v1InitiateStatKey(at)
	month := statentity.FormatMonthlyCall1v1InitiateStatKey(at)

	daily := statdao.GetDailyCall1v1InitiateStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyCall1v1InitiateStat(daily)

	weekly := statdao.GetWeeklyCall1v1InitiateStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyCall1v1InitiateStat(weekly)

	monthly := statdao.GetMonthlyCall1v1InitiateStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyCall1v1InitiateStat(monthly)
}

func GetCMSCall1v1InitiateTrend(_ context.Context, _ *trackingdto.CMSCall1v1InitiateTrendReq) (*trackingdto.CMSCall1v1InitiateTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyCall1v1InitiateStatKey(now)
	month := statentity.FormatMonthlyCall1v1InitiateStatKey(now)

	todayStat := statdao.GetDailyCall1v1InitiateStatByDate(today)
	weekStat := statdao.GetWeeklyCall1v1InitiateStatByWeek(week)
	monthStat := statdao.GetMonthlyCall1v1InitiateStatByMonth(month)

	return &trackingdto.CMSCall1v1InitiateTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toCall1v1InitiateDailyPoints(overlayCurrentDailyCall1v1Initiate(statdao.ListRecentDailyCall1v1InitiateStats(call1v1InitiateDailyLimit), now)),
		Weekly:     toCall1v1InitiateWeeklyPoints(overlayCurrentWeeklyCall1v1Initiate(statdao.ListRecentWeeklyCall1v1InitiateStats(call1v1InitiateWeeklyLimit), now)),
		Monthly:    toCall1v1InitiateMonthlyPoints(overlayCurrentMonthlyCall1v1Initiate(statdao.ListRecentMonthlyCall1v1InitiateStats(call1v1InitiateMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyCall1v1Initiate(rows []*statentity.DailyCall1v1InitiateStat, now time.Time) []*statentity.DailyCall1v1InitiateStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyCall1v1InitiateStatByDate(today)
	out := make([]*statentity.DailyCall1v1InitiateStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyCall1v1Initiate(rows []*statentity.WeeklyCall1v1InitiateStat, now time.Time) []*statentity.WeeklyCall1v1InitiateStat {
	week := statentity.FormatWeeklyCall1v1InitiateStatKey(now)
	live := statdao.GetWeeklyCall1v1InitiateStatByWeek(week)
	out := make([]*statentity.WeeklyCall1v1InitiateStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyCall1v1Initiate(rows []*statentity.MonthlyCall1v1InitiateStat, now time.Time) []*statentity.MonthlyCall1v1InitiateStat {
	month := statentity.FormatMonthlyCall1v1InitiateStatKey(now)
	live := statdao.GetMonthlyCall1v1InitiateStatByMonth(month)
	out := make([]*statentity.MonthlyCall1v1InitiateStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toCall1v1InitiateDailyPoints(rows []*statentity.DailyCall1v1InitiateStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toCall1v1InitiateWeeklyPoints(rows []*statentity.WeeklyCall1v1InitiateStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toCall1v1InitiateMonthlyPoints(rows []*statentity.MonthlyCall1v1InitiateStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
