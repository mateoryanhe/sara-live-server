package tracking

import (
	"context"
	"time"

	"xr-game-server/dao/statdao"
	"xr-game-server/dto/trackingdto"
	statentity "xr-game-server/entity/stat"
)

const (
	chargeResultServerDailyLimit   = 30
	chargeResultServerWeeklyLimit  = 12
	chargeResultServerMonthlyLimit = 12

	// EventKeyChargeResultServer 埋点 event_key(对齐埋点表 charge_result_server)
	EventKeyChargeResultServer = "charge_result_server"
)

func recordChargeResultServer(userId uint64, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyChargeResultServerStatKey(at)
	month := statentity.FormatMonthlyChargeResultServerStatKey(at)

	daily := statdao.GetDailyChargeResultServerStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyChargeResultServerStat(daily)

	weekly := statdao.GetWeeklyChargeResultServerStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyChargeResultServerStat(weekly)

	monthly := statdao.GetMonthlyChargeResultServerStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyChargeResultServerStat(monthly)
}

func GetCMSChargeResultServerTrend(_ context.Context, _ *trackingdto.CMSChargeResultServerTrendReq) (*trackingdto.CMSChargeResultServerTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyChargeResultServerStatKey(now)
	month := statentity.FormatMonthlyChargeResultServerStatKey(now)

	todayStat := statdao.GetDailyChargeResultServerStatByDate(today)
	weekStat := statdao.GetWeeklyChargeResultServerStatByWeek(week)
	monthStat := statdao.GetMonthlyChargeResultServerStatByMonth(month)

	return &trackingdto.CMSChargeResultServerTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toChargeResultServerDailyPoints(overlayCurrentDailyChargeResultServer(statdao.ListRecentDailyChargeResultServerStats(chargeResultServerDailyLimit), now)),
		Weekly:     toChargeResultServerWeeklyPoints(overlayCurrentWeeklyChargeResultServer(statdao.ListRecentWeeklyChargeResultServerStats(chargeResultServerWeeklyLimit), now)),
		Monthly:    toChargeResultServerMonthlyPoints(overlayCurrentMonthlyChargeResultServer(statdao.ListRecentMonthlyChargeResultServerStats(chargeResultServerMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyChargeResultServer(rows []*statentity.DailyChargeResultServerStat, now time.Time) []*statentity.DailyChargeResultServerStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyChargeResultServerStatByDate(today)
	out := make([]*statentity.DailyChargeResultServerStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyChargeResultServer(rows []*statentity.WeeklyChargeResultServerStat, now time.Time) []*statentity.WeeklyChargeResultServerStat {
	week := statentity.FormatWeeklyChargeResultServerStatKey(now)
	live := statdao.GetWeeklyChargeResultServerStatByWeek(week)
	out := make([]*statentity.WeeklyChargeResultServerStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyChargeResultServer(rows []*statentity.MonthlyChargeResultServerStat, now time.Time) []*statentity.MonthlyChargeResultServerStat {
	month := statentity.FormatMonthlyChargeResultServerStatKey(now)
	live := statdao.GetMonthlyChargeResultServerStatByMonth(month)
	out := make([]*statentity.MonthlyChargeResultServerStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toChargeResultServerDailyPoints(rows []*statentity.DailyChargeResultServerStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toChargeResultServerWeeklyPoints(rows []*statentity.WeeklyChargeResultServerStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toChargeResultServerMonthlyPoints(rows []*statentity.MonthlyChargeResultServerStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
