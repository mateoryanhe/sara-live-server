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
	hotLiveRoomGiftClickDailyLimit   = 30
	hotLiveRoomGiftClickWeeklyLimit  = 12
	hotLiveRoomGiftClickMonthlyLimit = 12

	// EventKeyHotLiveRoomGiftClick 埋点 event_key(对齐埋点表 hot_live_room_gift_click)
	EventKeyHotLiveRoomGiftClick = "hot_live_room_gift_click"
)

// ReportHotLiveRoomGiftClick App 上报秀场直播间礼物面板点击(登录即可,每次成功 +1)
func ReportHotLiveRoomGiftClick(ctx context.Context, _ *trackingdto.ReportHotLiveRoomGiftClickReq) (*trackingdto.ReportHotLiveRoomGiftClickRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	event.Pub(gameevent.HotLiveRoomGiftClickTrackingEvent, gameevent.NewHotLiveRoomGiftClickTrackingEventData(userId, time.Now()))
	return &trackingdto.ReportHotLiveRoomGiftClickRes{Success: true}, nil
}

func recordHotLiveRoomGiftClick(userId uint64, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyHotLiveRoomGiftClickStatKey(at)
	month := statentity.FormatMonthlyHotLiveRoomGiftClickStatKey(at)

	daily := statdao.GetDailyHotLiveRoomGiftClickStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyHotLiveRoomGiftClickStat(daily)

	weekly := statdao.GetWeeklyHotLiveRoomGiftClickStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyHotLiveRoomGiftClickStat(weekly)

	monthly := statdao.GetMonthlyHotLiveRoomGiftClickStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyHotLiveRoomGiftClickStat(monthly)
}

func GetCMSHotLiveRoomGiftClickTrend(_ context.Context, _ *trackingdto.CMSHotLiveRoomGiftClickTrendReq) (*trackingdto.CMSHotLiveRoomGiftClickTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyHotLiveRoomGiftClickStatKey(now)
	month := statentity.FormatMonthlyHotLiveRoomGiftClickStatKey(now)

	todayStat := statdao.GetDailyHotLiveRoomGiftClickStatByDate(today)
	weekStat := statdao.GetWeeklyHotLiveRoomGiftClickStatByWeek(week)
	monthStat := statdao.GetMonthlyHotLiveRoomGiftClickStatByMonth(month)

	return &trackingdto.CMSHotLiveRoomGiftClickTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toHotLiveRoomGiftClickDailyPoints(overlayCurrentDailyHotLiveRoomGiftClick(statdao.ListRecentDailyHotLiveRoomGiftClickStats(hotLiveRoomGiftClickDailyLimit), now)),
		Weekly:     toHotLiveRoomGiftClickWeeklyPoints(overlayCurrentWeeklyHotLiveRoomGiftClick(statdao.ListRecentWeeklyHotLiveRoomGiftClickStats(hotLiveRoomGiftClickWeeklyLimit), now)),
		Monthly:    toHotLiveRoomGiftClickMonthlyPoints(overlayCurrentMonthlyHotLiveRoomGiftClick(statdao.ListRecentMonthlyHotLiveRoomGiftClickStats(hotLiveRoomGiftClickMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyHotLiveRoomGiftClick(rows []*statentity.DailyHotLiveRoomGiftClickStat, now time.Time) []*statentity.DailyHotLiveRoomGiftClickStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyHotLiveRoomGiftClickStatByDate(today)
	out := make([]*statentity.DailyHotLiveRoomGiftClickStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyHotLiveRoomGiftClick(rows []*statentity.WeeklyHotLiveRoomGiftClickStat, now time.Time) []*statentity.WeeklyHotLiveRoomGiftClickStat {
	week := statentity.FormatWeeklyHotLiveRoomGiftClickStatKey(now)
	live := statdao.GetWeeklyHotLiveRoomGiftClickStatByWeek(week)
	out := make([]*statentity.WeeklyHotLiveRoomGiftClickStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyHotLiveRoomGiftClick(rows []*statentity.MonthlyHotLiveRoomGiftClickStat, now time.Time) []*statentity.MonthlyHotLiveRoomGiftClickStat {
	month := statentity.FormatMonthlyHotLiveRoomGiftClickStatKey(now)
	live := statdao.GetMonthlyHotLiveRoomGiftClickStatByMonth(month)
	out := make([]*statentity.MonthlyHotLiveRoomGiftClickStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toHotLiveRoomGiftClickDailyPoints(rows []*statentity.DailyHotLiveRoomGiftClickStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toHotLiveRoomGiftClickWeeklyPoints(rows []*statentity.WeeklyHotLiveRoomGiftClickStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toHotLiveRoomGiftClickMonthlyPoints(rows []*statentity.MonthlyHotLiveRoomGiftClickStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
