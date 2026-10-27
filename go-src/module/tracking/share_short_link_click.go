package tracking

import (
	"context"
	"time"

	"xr-game-server/core/event"
	"xr-game-server/dao/statdao"
	"xr-game-server/dto/trackingdto"
	statentity "xr-game-server/entity/stat"
	"xr-game-server/gameevent"
)

const (
	shareShortLinkClickDailyLimit   = 30
	shareShortLinkClickWeeklyLimit  = 12
	shareShortLinkClickMonthlyLimit = 12

	// EventKeyShareShortLinkClick 埋点 event_key(对齐埋点表 share_short_link_click)
	EventKeyShareShortLinkClick = "share_short_link_click"
)

// TrackShareShortLinkClick 分享短链 GET /dl 成功 302 时 +1(免登录,服务端计数).
func TrackShareShortLinkClick() {
	event.Pub(gameevent.ShareShortLinkClickTrackingEvent, gameevent.NewShareShortLinkClickTrackingEventData(time.Now()))
}

func recordShareShortLinkClick(at time.Time) {
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyShareShortLinkClickStatKey(at)
	month := statentity.FormatMonthlyShareShortLinkClickStatKey(at)

	daily := statdao.GetDailyShareShortLinkClickStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyShareShortLinkClickStat(daily)

	weekly := statdao.GetWeeklyShareShortLinkClickStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyShareShortLinkClickStat(weekly)

	monthly := statdao.GetMonthlyShareShortLinkClickStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyShareShortLinkClickStat(monthly)
}

func GetCMSShareShortLinkClickTrend(_ context.Context, _ *trackingdto.CMSShareShortLinkClickTrendReq) (*trackingdto.CMSShareShortLinkClickTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyShareShortLinkClickStatKey(now)
	month := statentity.FormatMonthlyShareShortLinkClickStatKey(now)

	todayStat := statdao.GetDailyShareShortLinkClickStatByDate(today)
	weekStat := statdao.GetWeeklyShareShortLinkClickStatByWeek(week)
	monthStat := statdao.GetMonthlyShareShortLinkClickStatByMonth(month)

	return &trackingdto.CMSShareShortLinkClickTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toShareShortLinkClickDailyPoints(overlayCurrentDailyShareShortLinkClick(statdao.ListRecentDailyShareShortLinkClickStats(shareShortLinkClickDailyLimit), now)),
		Weekly:     toShareShortLinkClickWeeklyPoints(overlayCurrentWeeklyShareShortLinkClick(statdao.ListRecentWeeklyShareShortLinkClickStats(shareShortLinkClickWeeklyLimit), now)),
		Monthly:    toShareShortLinkClickMonthlyPoints(overlayCurrentMonthlyShareShortLinkClick(statdao.ListRecentMonthlyShareShortLinkClickStats(shareShortLinkClickMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyShareShortLinkClick(rows []*statentity.DailyShareShortLinkClickStat, now time.Time) []*statentity.DailyShareShortLinkClickStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyShareShortLinkClickStatByDate(today)
	out := make([]*statentity.DailyShareShortLinkClickStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyShareShortLinkClick(rows []*statentity.WeeklyShareShortLinkClickStat, now time.Time) []*statentity.WeeklyShareShortLinkClickStat {
	week := statentity.FormatWeeklyShareShortLinkClickStatKey(now)
	live := statdao.GetWeeklyShareShortLinkClickStatByWeek(week)
	out := make([]*statentity.WeeklyShareShortLinkClickStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyShareShortLinkClick(rows []*statentity.MonthlyShareShortLinkClickStat, now time.Time) []*statentity.MonthlyShareShortLinkClickStat {
	month := statentity.FormatMonthlyShareShortLinkClickStatKey(now)
	live := statdao.GetMonthlyShareShortLinkClickStatByMonth(month)
	out := make([]*statentity.MonthlyShareShortLinkClickStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toShareShortLinkClickDailyPoints(rows []*statentity.DailyShareShortLinkClickStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toShareShortLinkClickWeeklyPoints(rows []*statentity.WeeklyShareShortLinkClickStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toShareShortLinkClickMonthlyPoints(rows []*statentity.MonthlyShareShortLinkClickStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
