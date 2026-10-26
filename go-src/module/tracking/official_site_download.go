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
	officialSiteDownloadDailyLimit   = 30
	officialSiteDownloadWeeklyLimit  = 12
	officialSiteDownloadMonthlyLimit = 12

	// EventKeyOfficialSiteDownload 埋点 event_key(对齐埋点表 official_site_download)
	EventKeyOfficialSiteDownload = "official_site_download"
)

// ReportOfficialSiteDownload 上报官网下载点击(免登录,每次成功 +1)
func ReportOfficialSiteDownload(_ context.Context, _ *trackingdto.ReportOfficialSiteDownloadReq) (*trackingdto.ReportOfficialSiteDownloadRes, error) {
	event.Pub(gameevent.OfficialSiteDownloadTrackingEvent, gameevent.NewOfficialSiteDownloadTrackingEventData(time.Now()))
	return &trackingdto.ReportOfficialSiteDownloadRes{Success: true}, nil
}

func recordOfficialSiteDownload(at time.Time) {
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyOfficialSiteDownloadStatKey(at)
	month := statentity.FormatMonthlyOfficialSiteDownloadStatKey(at)

	daily := statdao.GetDailyOfficialSiteDownloadStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyOfficialSiteDownloadStat(daily)

	weekly := statdao.GetWeeklyOfficialSiteDownloadStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyOfficialSiteDownloadStat(weekly)

	monthly := statdao.GetMonthlyOfficialSiteDownloadStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyOfficialSiteDownloadStat(monthly)
}

func GetCMSOfficialSiteDownloadTrend(_ context.Context, _ *trackingdto.CMSOfficialSiteDownloadTrendReq) (*trackingdto.CMSOfficialSiteDownloadTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyOfficialSiteDownloadStatKey(now)
	month := statentity.FormatMonthlyOfficialSiteDownloadStatKey(now)

	todayStat := statdao.GetDailyOfficialSiteDownloadStatByDate(today)
	weekStat := statdao.GetWeeklyOfficialSiteDownloadStatByWeek(week)
	monthStat := statdao.GetMonthlyOfficialSiteDownloadStatByMonth(month)

	return &trackingdto.CMSOfficialSiteDownloadTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toOfficialSiteDownloadDailyPoints(overlayCurrentDailyOfficialSiteDownload(statdao.ListRecentDailyOfficialSiteDownloadStats(officialSiteDownloadDailyLimit), now)),
		Weekly:     toOfficialSiteDownloadWeeklyPoints(overlayCurrentWeeklyOfficialSiteDownload(statdao.ListRecentWeeklyOfficialSiteDownloadStats(officialSiteDownloadWeeklyLimit), now)),
		Monthly:    toOfficialSiteDownloadMonthlyPoints(overlayCurrentMonthlyOfficialSiteDownload(statdao.ListRecentMonthlyOfficialSiteDownloadStats(officialSiteDownloadMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyOfficialSiteDownload(rows []*statentity.DailyOfficialSiteDownloadStat, now time.Time) []*statentity.DailyOfficialSiteDownloadStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyOfficialSiteDownloadStatByDate(today)
	out := make([]*statentity.DailyOfficialSiteDownloadStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyOfficialSiteDownload(rows []*statentity.WeeklyOfficialSiteDownloadStat, now time.Time) []*statentity.WeeklyOfficialSiteDownloadStat {
	week := statentity.FormatWeeklyOfficialSiteDownloadStatKey(now)
	live := statdao.GetWeeklyOfficialSiteDownloadStatByWeek(week)
	out := make([]*statentity.WeeklyOfficialSiteDownloadStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyOfficialSiteDownload(rows []*statentity.MonthlyOfficialSiteDownloadStat, now time.Time) []*statentity.MonthlyOfficialSiteDownloadStat {
	month := statentity.FormatMonthlyOfficialSiteDownloadStatKey(now)
	live := statdao.GetMonthlyOfficialSiteDownloadStatByMonth(month)
	out := make([]*statentity.MonthlyOfficialSiteDownloadStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toOfficialSiteDownloadDailyPoints(rows []*statentity.DailyOfficialSiteDownloadStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toOfficialSiteDownloadWeeklyPoints(rows []*statentity.WeeklyOfficialSiteDownloadStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toOfficialSiteDownloadMonthlyPoints(rows []*statentity.MonthlyOfficialSiteDownloadStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
