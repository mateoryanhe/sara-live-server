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
	// EventKeyVoiceRoomEnter 埋点 event_key(对齐埋点表 voice_room_enter)
	EventKeyVoiceRoomEnter = "voice_room_enter"

	voiceRoomEnterDailyLimit   = 30
	voiceRoomEnterWeeklyLimit  = 12
	voiceRoomEnterMonthlyLimit = 12
)

// recordVoiceRoomEnter 消费埋点事件:语聊房(category=5)进房次数 +1(含重复进入,不含主播本人).
func recordVoiceRoomEnter(userId, roomId uint64, at time.Time) {
	if userId == 0 || roomId == 0 || userId == roomId {
		return
	}
	cfg := liveroomdao.GetLiveRoomCfg(roomId)
	if cfg == nil || cfg.Category != liveentity.LiveRoomCategoryVoiceChat {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyVoiceRoomEnterStatKey(at)
	month := statentity.FormatMonthlyVoiceRoomEnterStatKey(at)

	daily := statdao.GetDailyVoiceRoomEnterStatByDate(date)
	daily.AddJoinCount(1)
	statdao.PublishDailyVoiceRoomEnterStat(daily)

	weekly := statdao.GetWeeklyVoiceRoomEnterStatByWeek(week)
	weekly.AddJoinCount(1)
	statdao.PublishWeeklyVoiceRoomEnterStat(weekly)

	monthly := statdao.GetMonthlyVoiceRoomEnterStatByMonth(month)
	monthly.AddJoinCount(1)
	statdao.PublishMonthlyVoiceRoomEnterStat(monthly)
}

// GetCMSVoiceRoomEnterTrend CMS 埋点:进入语聊房次数趋势
func GetCMSVoiceRoomEnterTrend(_ context.Context, _ *trackingdto.CMSVoiceRoomEnterTrendReq) (*trackingdto.CMSVoiceRoomEnterTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyVoiceRoomEnterStatKey(now)
	month := statentity.FormatMonthlyVoiceRoomEnterStatKey(now)

	todayStat := statdao.GetDailyVoiceRoomEnterStatByDate(today)
	weekStat := statdao.GetWeeklyVoiceRoomEnterStatByWeek(week)
	monthStat := statdao.GetMonthlyVoiceRoomEnterStatByMonth(month)

	return &trackingdto.CMSVoiceRoomEnterTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toVoiceRoomEnterDailyPoints(overlayCurrentDailyVoiceRoomEnter(statdao.ListRecentDailyVoiceRoomEnterStats(voiceRoomEnterDailyLimit), now)),
		Weekly:     toVoiceRoomEnterWeeklyPoints(overlayCurrentWeeklyVoiceRoomEnter(statdao.ListRecentWeeklyVoiceRoomEnterStats(voiceRoomEnterWeeklyLimit), now)),
		Monthly:    toVoiceRoomEnterMonthlyPoints(overlayCurrentMonthlyVoiceRoomEnter(statdao.ListRecentMonthlyVoiceRoomEnterStats(voiceRoomEnterMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyVoiceRoomEnter(rows []*statentity.DailyVoiceRoomEnterStat, now time.Time) []*statentity.DailyVoiceRoomEnterStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyVoiceRoomEnterStatByDate(today)
	out := make([]*statentity.DailyVoiceRoomEnterStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyVoiceRoomEnter(rows []*statentity.WeeklyVoiceRoomEnterStat, now time.Time) []*statentity.WeeklyVoiceRoomEnterStat {
	week := statentity.FormatWeeklyVoiceRoomEnterStatKey(now)
	live := statdao.GetWeeklyVoiceRoomEnterStatByWeek(week)
	out := make([]*statentity.WeeklyVoiceRoomEnterStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyVoiceRoomEnter(rows []*statentity.MonthlyVoiceRoomEnterStat, now time.Time) []*statentity.MonthlyVoiceRoomEnterStat {
	month := statentity.FormatMonthlyVoiceRoomEnterStatKey(now)
	live := statdao.GetMonthlyVoiceRoomEnterStatByMonth(month)
	out := make([]*statentity.MonthlyVoiceRoomEnterStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toVoiceRoomEnterDailyPoints(rows []*statentity.DailyVoiceRoomEnterStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toVoiceRoomEnterWeeklyPoints(rows []*statentity.WeeklyVoiceRoomEnterStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toVoiceRoomEnterMonthlyPoints(rows []*statentity.MonthlyVoiceRoomEnterStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
