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
	voiceRoomGiftClickDailyLimit   = 30
	voiceRoomGiftClickWeeklyLimit  = 12
	voiceRoomGiftClickMonthlyLimit = 12

	// EventKeyVoiceRoomGiftClick 埋点 event_key(对齐埋点表 voice_room_gift_click)
	EventKeyVoiceRoomGiftClick = "voice_room_gift_click"
)

// ReportVoiceRoomGiftClick App 上报语聊房礼物面板点击(登录即可,每次成功 +1)
func ReportVoiceRoomGiftClick(ctx context.Context, _ *trackingdto.ReportVoiceRoomGiftClickReq) (*trackingdto.ReportVoiceRoomGiftClickRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	event.Pub(gameevent.VoiceRoomGiftClickTrackingEvent, gameevent.NewVoiceRoomGiftClickTrackingEventData(userId, time.Now()))
	return &trackingdto.ReportVoiceRoomGiftClickRes{Success: true}, nil
}

func recordVoiceRoomGiftClick(userId uint64, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyVoiceRoomGiftClickStatKey(at)
	month := statentity.FormatMonthlyVoiceRoomGiftClickStatKey(at)

	daily := statdao.GetDailyVoiceRoomGiftClickStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyVoiceRoomGiftClickStat(daily)

	weekly := statdao.GetWeeklyVoiceRoomGiftClickStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyVoiceRoomGiftClickStat(weekly)

	monthly := statdao.GetMonthlyVoiceRoomGiftClickStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyVoiceRoomGiftClickStat(monthly)
}

func GetCMSVoiceRoomGiftClickTrend(_ context.Context, _ *trackingdto.CMSVoiceRoomGiftClickTrendReq) (*trackingdto.CMSVoiceRoomGiftClickTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyVoiceRoomGiftClickStatKey(now)
	month := statentity.FormatMonthlyVoiceRoomGiftClickStatKey(now)

	todayStat := statdao.GetDailyVoiceRoomGiftClickStatByDate(today)
	weekStat := statdao.GetWeeklyVoiceRoomGiftClickStatByWeek(week)
	monthStat := statdao.GetMonthlyVoiceRoomGiftClickStatByMonth(month)

	return &trackingdto.CMSVoiceRoomGiftClickTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toVoiceRoomGiftClickDailyPoints(overlayCurrentDailyVoiceRoomGiftClick(statdao.ListRecentDailyVoiceRoomGiftClickStats(voiceRoomGiftClickDailyLimit), now)),
		Weekly:     toVoiceRoomGiftClickWeeklyPoints(overlayCurrentWeeklyVoiceRoomGiftClick(statdao.ListRecentWeeklyVoiceRoomGiftClickStats(voiceRoomGiftClickWeeklyLimit), now)),
		Monthly:    toVoiceRoomGiftClickMonthlyPoints(overlayCurrentMonthlyVoiceRoomGiftClick(statdao.ListRecentMonthlyVoiceRoomGiftClickStats(voiceRoomGiftClickMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyVoiceRoomGiftClick(rows []*statentity.DailyVoiceRoomGiftClickStat, now time.Time) []*statentity.DailyVoiceRoomGiftClickStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyVoiceRoomGiftClickStatByDate(today)
	out := make([]*statentity.DailyVoiceRoomGiftClickStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyVoiceRoomGiftClick(rows []*statentity.WeeklyVoiceRoomGiftClickStat, now time.Time) []*statentity.WeeklyVoiceRoomGiftClickStat {
	week := statentity.FormatWeeklyVoiceRoomGiftClickStatKey(now)
	live := statdao.GetWeeklyVoiceRoomGiftClickStatByWeek(week)
	out := make([]*statentity.WeeklyVoiceRoomGiftClickStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyVoiceRoomGiftClick(rows []*statentity.MonthlyVoiceRoomGiftClickStat, now time.Time) []*statentity.MonthlyVoiceRoomGiftClickStat {
	month := statentity.FormatMonthlyVoiceRoomGiftClickStatKey(now)
	live := statdao.GetMonthlyVoiceRoomGiftClickStatByMonth(month)
	out := make([]*statentity.MonthlyVoiceRoomGiftClickStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toVoiceRoomGiftClickDailyPoints(rows []*statentity.DailyVoiceRoomGiftClickStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toVoiceRoomGiftClickWeeklyPoints(rows []*statentity.WeeklyVoiceRoomGiftClickStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toVoiceRoomGiftClickMonthlyPoints(rows []*statentity.MonthlyVoiceRoomGiftClickStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
