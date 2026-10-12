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
	chargePanelEnterDailyLimit   = 30
	chargePanelEnterWeeklyLimit  = 12
	chargePanelEnterMonthlyLimit = 12

	// EventKeyChargePanelEnter 埋点 event_key(对齐埋点表 charge_panel_enter)
	EventKeyChargePanelEnter = "charge_panel_enter"
)

// ReportChargePanelEnter App 上报充值/礼物面板曝光(登录即可,每次成功 +1)
func ReportChargePanelEnter(ctx context.Context, _ *trackingdto.ReportChargePanelEnterReq) (*trackingdto.ReportChargePanelEnterRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	event.Pub(gameevent.ChargePanelEnterTrackingEvent, gameevent.NewChargePanelEnterTrackingEventData(userId, time.Now()))
	return &trackingdto.ReportChargePanelEnterRes{Success: true}, nil
}

func recordChargePanelEnter(userId uint64, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyChargePanelEnterStatKey(at)
	month := statentity.FormatMonthlyChargePanelEnterStatKey(at)

	daily := statdao.GetDailyChargePanelEnterStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyChargePanelEnterStat(daily)

	weekly := statdao.GetWeeklyChargePanelEnterStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyChargePanelEnterStat(weekly)

	monthly := statdao.GetMonthlyChargePanelEnterStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyChargePanelEnterStat(monthly)
}

func GetCMSChargePanelEnterTrend(_ context.Context, _ *trackingdto.CMSChargePanelEnterTrendReq) (*trackingdto.CMSChargePanelEnterTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyChargePanelEnterStatKey(now)
	month := statentity.FormatMonthlyChargePanelEnterStatKey(now)

	todayStat := statdao.GetDailyChargePanelEnterStatByDate(today)
	weekStat := statdao.GetWeeklyChargePanelEnterStatByWeek(week)
	monthStat := statdao.GetMonthlyChargePanelEnterStatByMonth(month)

	return &trackingdto.CMSChargePanelEnterTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toChargePanelEnterDailyPoints(overlayCurrentDailyChargePanelEnter(statdao.ListRecentDailyChargePanelEnterStats(chargePanelEnterDailyLimit), now)),
		Weekly:     toChargePanelEnterWeeklyPoints(overlayCurrentWeeklyChargePanelEnter(statdao.ListRecentWeeklyChargePanelEnterStats(chargePanelEnterWeeklyLimit), now)),
		Monthly:    toChargePanelEnterMonthlyPoints(overlayCurrentMonthlyChargePanelEnter(statdao.ListRecentMonthlyChargePanelEnterStats(chargePanelEnterMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyChargePanelEnter(rows []*statentity.DailyChargePanelEnterStat, now time.Time) []*statentity.DailyChargePanelEnterStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyChargePanelEnterStatByDate(today)
	out := make([]*statentity.DailyChargePanelEnterStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyChargePanelEnter(rows []*statentity.WeeklyChargePanelEnterStat, now time.Time) []*statentity.WeeklyChargePanelEnterStat {
	week := statentity.FormatWeeklyChargePanelEnterStatKey(now)
	live := statdao.GetWeeklyChargePanelEnterStatByWeek(week)
	out := make([]*statentity.WeeklyChargePanelEnterStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyChargePanelEnter(rows []*statentity.MonthlyChargePanelEnterStat, now time.Time) []*statentity.MonthlyChargePanelEnterStat {
	month := statentity.FormatMonthlyChargePanelEnterStatKey(now)
	live := statdao.GetMonthlyChargePanelEnterStatByMonth(month)
	out := make([]*statentity.MonthlyChargePanelEnterStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toChargePanelEnterDailyPoints(rows []*statentity.DailyChargePanelEnterStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toChargePanelEnterWeeklyPoints(rows []*statentity.WeeklyChargePanelEnterStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toChargePanelEnterMonthlyPoints(rows []*statentity.MonthlyChargePanelEnterStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
