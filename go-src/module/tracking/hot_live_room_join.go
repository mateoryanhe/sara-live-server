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
	hotLiveRoomJoinDailyLimit   = 30
	hotLiveRoomJoinWeeklyLimit  = 12
	hotLiveRoomJoinMonthlyLimit = 12
)

// recordShowcaseLiveRoomJoin 消费埋点事件:秀场(category=1)进房次数 +1(含重复进入,不含主播本人).
func recordShowcaseLiveRoomJoin(userId, roomId uint64, at time.Time) {
	if userId == 0 || roomId == 0 || userId == roomId {
		return
	}
	cfg := liveroomdao.GetLiveRoomCfg(roomId)
	if cfg == nil || cfg.Category != liveentity.LiveRoomCategoryHot {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyHotLiveRoomJoinStatKey(at)
	month := statentity.FormatMonthlyHotLiveRoomJoinStatKey(at)

	daily := statdao.GetDailyHotLiveRoomJoinStatByDate(date)
	daily.AddJoinCount(1)
	statdao.PublishDailyHotLiveRoomJoinStat(daily)

	weekly := statdao.GetWeeklyHotLiveRoomJoinStatByWeek(week)
	weekly.AddJoinCount(1)
	statdao.PublishWeeklyHotLiveRoomJoinStat(weekly)

	monthly := statdao.GetMonthlyHotLiveRoomJoinStatByMonth(month)
	monthly.AddJoinCount(1)
	statdao.PublishMonthlyHotLiveRoomJoinStat(monthly)
}

// GetCMSHotLiveRoomJoinTrend CMS 埋点:进入 Hot 直播间次数趋势
func GetCMSHotLiveRoomJoinTrend(_ context.Context, _ *trackingdto.CMSHotLiveRoomJoinTrendReq) (*trackingdto.CMSHotLiveRoomJoinTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyHotLiveRoomJoinStatKey(now)
	month := statentity.FormatMonthlyHotLiveRoomJoinStatKey(now)

	todayStat := statdao.GetDailyHotLiveRoomJoinStatByDate(today)
	weekStat := statdao.GetWeeklyHotLiveRoomJoinStatByWeek(week)
	monthStat := statdao.GetMonthlyHotLiveRoomJoinStatByMonth(month)

	return &trackingdto.CMSHotLiveRoomJoinTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toHotJoinDailyPoints(overlayCurrentDailyHotJoin(statdao.ListRecentDailyHotLiveRoomJoinStats(hotLiveRoomJoinDailyLimit), now)),
		Weekly:     toHotJoinWeeklyPoints(overlayCurrentWeeklyHotJoin(statdao.ListRecentWeeklyHotLiveRoomJoinStats(hotLiveRoomJoinWeeklyLimit), now)),
		Monthly:    toHotJoinMonthlyPoints(overlayCurrentMonthlyHotJoin(statdao.ListRecentMonthlyHotLiveRoomJoinStats(hotLiveRoomJoinMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyHotJoin(rows []*statentity.DailyHotLiveRoomJoinStat, now time.Time) []*statentity.DailyHotLiveRoomJoinStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyHotLiveRoomJoinStatByDate(today)
	out := make([]*statentity.DailyHotLiveRoomJoinStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyHotJoin(rows []*statentity.WeeklyHotLiveRoomJoinStat, now time.Time) []*statentity.WeeklyHotLiveRoomJoinStat {
	week := statentity.FormatWeeklyHotLiveRoomJoinStatKey(now)
	live := statdao.GetWeeklyHotLiveRoomJoinStatByWeek(week)
	out := make([]*statentity.WeeklyHotLiveRoomJoinStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyHotJoin(rows []*statentity.MonthlyHotLiveRoomJoinStat, now time.Time) []*statentity.MonthlyHotLiveRoomJoinStat {
	month := statentity.FormatMonthlyHotLiveRoomJoinStatKey(now)
	live := statdao.GetMonthlyHotLiveRoomJoinStatByMonth(month)
	out := make([]*statentity.MonthlyHotLiveRoomJoinStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toHotJoinDailyPoints(rows []*statentity.DailyHotLiveRoomJoinStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toHotJoinWeeklyPoints(rows []*statentity.WeeklyHotLiveRoomJoinStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toHotJoinMonthlyPoints(rows []*statentity.MonthlyHotLiveRoomJoinStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
