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
	h5DownloadDailyLimit   = 30
	h5DownloadWeeklyLimit  = 12
	h5DownloadMonthlyLimit = 12

	// EventKeyH5Download 埋点 event_key(对齐埋点表 h5_download)
	EventKeyH5Download = "h5_download"
)

// ReportH5Download 上报 H5 下载点击(免登录,每次成功 +1;可带 channelCode 计入渠道日统计)
func ReportH5Download(_ context.Context, req *trackingdto.ReportH5DownloadReq) (*trackingdto.ReportH5DownloadRes, error) {
	channelCode := 0
	if req != nil {
		channelCode = req.ChannelCode
	}
	event.Pub(gameevent.H5DownloadTrackingEvent, gameevent.NewH5DownloadTrackingEventData(time.Now(), channelCode))
	return &trackingdto.ReportH5DownloadRes{Success: true}, nil
}

func recordH5Download(at time.Time) {
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyH5DownloadStatKey(at)
	month := statentity.FormatMonthlyH5DownloadStatKey(at)

	daily := statdao.GetDailyH5DownloadStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyH5DownloadStat(daily)

	weekly := statdao.GetWeeklyH5DownloadStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyH5DownloadStat(weekly)

	monthly := statdao.GetMonthlyH5DownloadStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyH5DownloadStat(monthly)
}

func GetCMSH5DownloadTrend(_ context.Context, _ *trackingdto.CMSH5DownloadTrendReq) (*trackingdto.CMSH5DownloadTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyH5DownloadStatKey(now)
	month := statentity.FormatMonthlyH5DownloadStatKey(now)

	todayStat := statdao.GetDailyH5DownloadStatByDate(today)
	weekStat := statdao.GetWeeklyH5DownloadStatByWeek(week)
	monthStat := statdao.GetMonthlyH5DownloadStatByMonth(month)

	return &trackingdto.CMSH5DownloadTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toH5DownloadDailyPoints(overlayCurrentDailyH5Download(statdao.ListRecentDailyH5DownloadStats(h5DownloadDailyLimit), now)),
		Weekly:     toH5DownloadWeeklyPoints(overlayCurrentWeeklyH5Download(statdao.ListRecentWeeklyH5DownloadStats(h5DownloadWeeklyLimit), now)),
		Monthly:    toH5DownloadMonthlyPoints(overlayCurrentMonthlyH5Download(statdao.ListRecentMonthlyH5DownloadStats(h5DownloadMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyH5Download(rows []*statentity.DailyH5DownloadStat, now time.Time) []*statentity.DailyH5DownloadStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyH5DownloadStatByDate(today)
	out := make([]*statentity.DailyH5DownloadStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyH5Download(rows []*statentity.WeeklyH5DownloadStat, now time.Time) []*statentity.WeeklyH5DownloadStat {
	week := statentity.FormatWeeklyH5DownloadStatKey(now)
	live := statdao.GetWeeklyH5DownloadStatByWeek(week)
	out := make([]*statentity.WeeklyH5DownloadStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyH5Download(rows []*statentity.MonthlyH5DownloadStat, now time.Time) []*statentity.MonthlyH5DownloadStat {
	month := statentity.FormatMonthlyH5DownloadStatKey(now)
	live := statdao.GetMonthlyH5DownloadStatByMonth(month)
	out := make([]*statentity.MonthlyH5DownloadStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toH5DownloadDailyPoints(rows []*statentity.DailyH5DownloadStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toH5DownloadWeeklyPoints(rows []*statentity.WeeklyH5DownloadStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toH5DownloadMonthlyPoints(rows []*statentity.MonthlyH5DownloadStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
