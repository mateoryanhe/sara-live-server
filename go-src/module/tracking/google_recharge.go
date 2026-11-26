package tracking

import (
	"context"
	"time"

	"xr-game-server/dao/statdao"
	"xr-game-server/dto/trackingdto"
	statentity "xr-game-server/entity/stat"
)

const (
	googleRechargeDailyLimit = 30

	EventKeyGoogleRechargeList   = "google_recharge_list"
	EventKeyGoogleRechargeOrder  = "google_recharge_order"
	EventKeyGoogleRechargeResult = "google_recharge_result"
)

func recordGoogleRechargeList(at time.Time) {
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	daily := statdao.GetDailyGoogleRechargeListStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyGoogleRechargeListStat(daily)
}

func recordGoogleRechargeOrder(_ uint64, at time.Time) {
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	daily := statdao.GetDailyGoogleRechargeOrderStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyGoogleRechargeOrderStat(daily)
}

func recordGoogleRechargeResult(_ uint64, at time.Time) {
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	daily := statdao.GetDailyGoogleRechargeResultStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyGoogleRechargeResultStat(daily)
}

func GetCMSGoogleRechargeListTrend(_ context.Context, _ *trackingdto.CMSGoogleRechargeListTrendReq) (*trackingdto.CMSGoogleRechargeDailyTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	todayStat := statdao.GetDailyGoogleRechargeListStatByDate(today)
	return &trackingdto.CMSGoogleRechargeDailyTrendRes{
		TodayCount: todayStat.Count,
		Daily:      toCountDailyPoints(overlayListDaily(statdao.ListRecentDailyGoogleRechargeListStats(googleRechargeDailyLimit), now)),
	}, nil
}

func GetCMSGoogleRechargeOrderTrend(_ context.Context, _ *trackingdto.CMSGoogleRechargeOrderTrendReq) (*trackingdto.CMSGoogleRechargeDailyTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	todayStat := statdao.GetDailyGoogleRechargeOrderStatByDate(today)
	return &trackingdto.CMSGoogleRechargeDailyTrendRes{
		TodayCount: todayStat.Count,
		Daily:      toCountDailyPointsFromOrder(overlayOrderDaily(statdao.ListRecentDailyGoogleRechargeOrderStats(googleRechargeDailyLimit), now)),
	}, nil
}

func GetCMSGoogleRechargeResultTrend(_ context.Context, _ *trackingdto.CMSGoogleRechargeResultTrendReq) (*trackingdto.CMSGoogleRechargeDailyTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	todayStat := statdao.GetDailyGoogleRechargeResultStatByDate(today)
	return &trackingdto.CMSGoogleRechargeDailyTrendRes{
		TodayCount: todayStat.Count,
		Daily:      toCountDailyPointsFromResult(overlayResultDaily(statdao.ListRecentDailyGoogleRechargeResultStats(googleRechargeDailyLimit), now)),
	}, nil
}

func overlayListDaily(rows []*statentity.DailyGoogleRechargeListStat, now time.Time) []*statentity.DailyGoogleRechargeListStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyGoogleRechargeListStatByDate(today)
	out := make([]*statentity.DailyGoogleRechargeListStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayOrderDaily(rows []*statentity.DailyGoogleRechargeOrderStat, now time.Time) []*statentity.DailyGoogleRechargeOrderStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyGoogleRechargeOrderStatByDate(today)
	out := make([]*statentity.DailyGoogleRechargeOrderStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayResultDaily(rows []*statentity.DailyGoogleRechargeResultStat, now time.Time) []*statentity.DailyGoogleRechargeResultStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyGoogleRechargeResultStatByDate(today)
	out := make([]*statentity.DailyGoogleRechargeResultStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toCountDailyPoints(rows []*statentity.DailyGoogleRechargeListStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toCountDailyPointsFromOrder(rows []*statentity.DailyGoogleRechargeOrderStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toCountDailyPointsFromResult(rows []*statentity.DailyGoogleRechargeResultStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}