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
	hotLiveRoomLeaveDailyLimit   = 30
	hotLiveRoomLeaveWeeklyLimit  = 12
	hotLiveRoomLeaveMonthlyLimit = 12

	// EventKeyHotLiveRoomLeave 埋点 event_key(文档/CMS)
	EventKeyHotLiveRoomLeave = "hot_live_room_leave"
)

// recordShowcaseLiveRoomLeave 消费埋点事件:秀场(category=1)退房次数 +1(含重复,不含主播本人).
func recordShowcaseLiveRoomLeave(userId, roomId uint64, at time.Time) {
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
	week := statentity.FormatWeeklyHotLiveRoomLeaveStatKey(at)
	month := statentity.FormatMonthlyHotLiveRoomLeaveStatKey(at)

	daily := statdao.GetDailyHotLiveRoomLeaveStatByDate(date)
	daily.AddLeaveCount(1)
	statdao.PublishDailyHotLiveRoomLeaveStat(daily)

	weekly := statdao.GetWeeklyHotLiveRoomLeaveStatByWeek(week)
	weekly.AddLeaveCount(1)
	statdao.PublishWeeklyHotLiveRoomLeaveStat(weekly)

	monthly := statdao.GetMonthlyHotLiveRoomLeaveStatByMonth(month)
	monthly.AddLeaveCount(1)
	statdao.PublishMonthlyHotLiveRoomLeaveStat(monthly)
}

// GetCMSHotLiveRoomLeaveTrend CMS 埋点:退出 Hot 直播间次数趋势
func GetCMSHotLiveRoomLeaveTrend(_ context.Context, _ *trackingdto.CMSHotLiveRoomLeaveTrendReq) (*trackingdto.CMSHotLiveRoomLeaveTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyHotLiveRoomLeaveStatKey(now)
	month := statentity.FormatMonthlyHotLiveRoomLeaveStatKey(now)

	todayStat := statdao.GetDailyHotLiveRoomLeaveStatByDate(today)
	weekStat := statdao.GetWeeklyHotLiveRoomLeaveStatByWeek(week)
	monthStat := statdao.GetMonthlyHotLiveRoomLeaveStatByMonth(month)

	return &trackingdto.CMSHotLiveRoomLeaveTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toHotLeaveDailyPoints(overlayCurrentDailyHotLeave(statdao.ListRecentDailyHotLiveRoomLeaveStats(hotLiveRoomLeaveDailyLimit), now)),
		Weekly:     toHotLeaveWeeklyPoints(overlayCurrentWeeklyHotLeave(statdao.ListRecentWeeklyHotLiveRoomLeaveStats(hotLiveRoomLeaveWeeklyLimit), now)),
		Monthly:    toHotLeaveMonthlyPoints(overlayCurrentMonthlyHotLeave(statdao.ListRecentMonthlyHotLiveRoomLeaveStats(hotLiveRoomLeaveMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyHotLeave(rows []*statentity.DailyHotLiveRoomLeaveStat, now time.Time) []*statentity.DailyHotLiveRoomLeaveStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyHotLiveRoomLeaveStatByDate(today)
	out := make([]*statentity.DailyHotLiveRoomLeaveStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyHotLeave(rows []*statentity.WeeklyHotLiveRoomLeaveStat, now time.Time) []*statentity.WeeklyHotLiveRoomLeaveStat {
	week := statentity.FormatWeeklyHotLiveRoomLeaveStatKey(now)
	live := statdao.GetWeeklyHotLiveRoomLeaveStatByWeek(week)
	out := make([]*statentity.WeeklyHotLiveRoomLeaveStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyHotLeave(rows []*statentity.MonthlyHotLiveRoomLeaveStat, now time.Time) []*statentity.MonthlyHotLiveRoomLeaveStat {
	month := statentity.FormatMonthlyHotLiveRoomLeaveStatKey(now)
	live := statdao.GetMonthlyHotLiveRoomLeaveStatByMonth(month)
	out := make([]*statentity.MonthlyHotLiveRoomLeaveStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toHotLeaveDailyPoints(rows []*statentity.DailyHotLiveRoomLeaveStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toHotLeaveWeeklyPoints(rows []*statentity.WeeklyHotLiveRoomLeaveStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toHotLeaveMonthlyPoints(rows []*statentity.MonthlyHotLiveRoomLeaveStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
