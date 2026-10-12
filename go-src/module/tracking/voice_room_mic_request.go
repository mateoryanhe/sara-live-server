package tracking

import (
	"context"
	"time"

	"xr-game-server/dao/statdao"
	"xr-game-server/dto/trackingdto"
	statentity "xr-game-server/entity/stat"
)

const (
	VoiceRoomMicRequestDailyLimit   = 30
	VoiceRoomMicRequestWeeklyLimit  = 12
	VoiceRoomMicRequestMonthlyLimit = 12

	// EventKeyVoiceRoomMicRequest 埋点 event_key(对齐埋点表 voice_room_mic_request)
	EventKeyVoiceRoomMicRequest = "voice_room_mic_request"
)

func recordVoiceRoomMicRequest(userId uint64, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyVoiceRoomMicRequestStatKey(at)
	month := statentity.FormatMonthlyVoiceRoomMicRequestStatKey(at)

	daily := statdao.GetDailyVoiceRoomMicRequestStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyVoiceRoomMicRequestStat(daily)

	weekly := statdao.GetWeeklyVoiceRoomMicRequestStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyVoiceRoomMicRequestStat(weekly)

	monthly := statdao.GetMonthlyVoiceRoomMicRequestStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyVoiceRoomMicRequestStat(monthly)
}

func GetCMSVoiceRoomMicRequestTrend(_ context.Context, _ *trackingdto.CMSVoiceRoomMicRequestTrendReq) (*trackingdto.CMSVoiceRoomMicRequestTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyVoiceRoomMicRequestStatKey(now)
	month := statentity.FormatMonthlyVoiceRoomMicRequestStatKey(now)

	todayStat := statdao.GetDailyVoiceRoomMicRequestStatByDate(today)
	weekStat := statdao.GetWeeklyVoiceRoomMicRequestStatByWeek(week)
	monthStat := statdao.GetMonthlyVoiceRoomMicRequestStatByMonth(month)

	return &trackingdto.CMSVoiceRoomMicRequestTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toVoiceRoomMicRequestDailyPoints(overlayCurrentDailyVoiceRoomMicRequest(statdao.ListRecentDailyVoiceRoomMicRequestStats(VoiceRoomMicRequestDailyLimit), now)),
		Weekly:     toVoiceRoomMicRequestWeeklyPoints(overlayCurrentWeeklyVoiceRoomMicRequest(statdao.ListRecentWeeklyVoiceRoomMicRequestStats(VoiceRoomMicRequestWeeklyLimit), now)),
		Monthly:    toVoiceRoomMicRequestMonthlyPoints(overlayCurrentMonthlyVoiceRoomMicRequest(statdao.ListRecentMonthlyVoiceRoomMicRequestStats(VoiceRoomMicRequestMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyVoiceRoomMicRequest(rows []*statentity.DailyVoiceRoomMicRequestStat, now time.Time) []*statentity.DailyVoiceRoomMicRequestStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyVoiceRoomMicRequestStatByDate(today)
	out := make([]*statentity.DailyVoiceRoomMicRequestStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyVoiceRoomMicRequest(rows []*statentity.WeeklyVoiceRoomMicRequestStat, now time.Time) []*statentity.WeeklyVoiceRoomMicRequestStat {
	week := statentity.FormatWeeklyVoiceRoomMicRequestStatKey(now)
	live := statdao.GetWeeklyVoiceRoomMicRequestStatByWeek(week)
	out := make([]*statentity.WeeklyVoiceRoomMicRequestStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyVoiceRoomMicRequest(rows []*statentity.MonthlyVoiceRoomMicRequestStat, now time.Time) []*statentity.MonthlyVoiceRoomMicRequestStat {
	month := statentity.FormatMonthlyVoiceRoomMicRequestStatKey(now)
	live := statdao.GetMonthlyVoiceRoomMicRequestStatByMonth(month)
	out := make([]*statentity.MonthlyVoiceRoomMicRequestStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toVoiceRoomMicRequestDailyPoints(rows []*statentity.DailyVoiceRoomMicRequestStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toVoiceRoomMicRequestWeeklyPoints(rows []*statentity.WeeklyVoiceRoomMicRequestStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toVoiceRoomMicRequestMonthlyPoints(rows []*statentity.MonthlyVoiceRoomMicRequestStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
