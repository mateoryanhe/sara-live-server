package tracking

import (
	"context"
	"time"

	"xr-game-server/core/event"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/statdao"
	"xr-game-server/dto/trackingdto"
	liveentity "xr-game-server/entity/live"
	statentity "xr-game-server/entity/stat"
	"xr-game-server/errercode"
	"xr-game-server/gameevent"
	"xr-game-server/module/liveroom"
)

const (
	liveFirstFrameRenderDailyLimit   = 30
	liveFirstFrameRenderWeeklyLimit  = 12
	liveFirstFrameRenderMonthlyLimit = 12

	// EventKeyLiveFirstFrameRendered 埋点 event_key(文档/CMS)
	EventKeyLiveFirstFrameRendered = "live_first_frame_rendered"
)

// ReportLiveFirstFrameRendered App 上报秀场直播首帧画面渲染完成(每次成功请求 +1,异步落库)
func ReportLiveFirstFrameRendered(ctx context.Context, req *trackingdto.ReportLiveFirstFrameRenderedReq) (*trackingdto.ReportLiveFirstFrameRenderedRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 || req == nil || req.RoomId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	room := liveroomdao.GetRoomById(req.RoomId)
	if room == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}
	if userId != room.ID && liveroom.IsRoomOffShelf(room) {
		return nil, errercode.CreateCode(errercode.LiveRoomOffShelf)
	}
	if userId != room.ID && room.LiveRecordId == 0 {
		return nil, errercode.CreateCode(errercode.LiveRoomNotLive)
	}
	cfg := liveroomdao.GetLiveRoomCfg(req.RoomId)
	if cfg == nil || cfg.Category != liveentity.LiveRoomCategoryHot {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if userId == room.ID {
		return &trackingdto.ReportLiveFirstFrameRenderedRes{Success: true}, nil
	}
	event.Pub(gameevent.LiveFirstFrameRenderedTrackingEvent, gameevent.NewLiveFirstFrameRenderedTrackingEventData(userId, req.RoomId, time.Now()))
	return &trackingdto.ReportLiveFirstFrameRenderedRes{Success: true}, nil
}

func recordLiveFirstFrameRendered(userId, roomId uint64, at time.Time) {
	if userId == 0 || roomId == 0 || userId == roomId {
		return
	}
	cfg := liveroomdao.GetLiveRoomCfg(roomId)
	if cfg == nil || cfg.Category != liveentity.LiveRoomCategoryHot {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	week := statentity.FormatWeeklyLiveFirstFrameRenderStatKey(at)
	month := statentity.FormatMonthlyLiveFirstFrameRenderStatKey(at)

	daily := statdao.GetDailyLiveFirstFrameRenderStatByDate(date)
	daily.AddCount(1)
	statdao.PublishDailyLiveFirstFrameRenderStat(daily)

	weekly := statdao.GetWeeklyLiveFirstFrameRenderStatByWeek(week)
	weekly.AddCount(1)
	statdao.PublishWeeklyLiveFirstFrameRenderStat(weekly)

	monthly := statdao.GetMonthlyLiveFirstFrameRenderStatByMonth(month)
	monthly.AddCount(1)
	statdao.PublishMonthlyLiveFirstFrameRenderStat(monthly)
}

// GetCMSLiveFirstFrameRenderTrend CMS 埋点:直播首帧渲染完成次数趋势
func GetCMSLiveFirstFrameRenderTrend(_ context.Context, _ *trackingdto.CMSLiveFirstFrameRenderTrendReq) (*trackingdto.CMSLiveFirstFrameRenderTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	week := statentity.FormatWeeklyLiveFirstFrameRenderStatKey(now)
	month := statentity.FormatMonthlyLiveFirstFrameRenderStatKey(now)

	todayStat := statdao.GetDailyLiveFirstFrameRenderStatByDate(today)
	weekStat := statdao.GetWeeklyLiveFirstFrameRenderStatByWeek(week)
	monthStat := statdao.GetMonthlyLiveFirstFrameRenderStatByMonth(month)

	return &trackingdto.CMSLiveFirstFrameRenderTrendRes{
		TodayCount: todayStat.Count,
		WeekCount:  weekStat.Count,
		MonthCount: monthStat.Count,
		Daily:      toFirstFrameDailyPoints(overlayCurrentDailyFirstFrame(statdao.ListRecentDailyLiveFirstFrameRenderStats(liveFirstFrameRenderDailyLimit), now)),
		Weekly:     toFirstFrameWeeklyPoints(overlayCurrentWeeklyFirstFrame(statdao.ListRecentWeeklyLiveFirstFrameRenderStats(liveFirstFrameRenderWeeklyLimit), now)),
		Monthly:    toFirstFrameMonthlyPoints(overlayCurrentMonthlyFirstFrame(statdao.ListRecentMonthlyLiveFirstFrameRenderStats(liveFirstFrameRenderMonthlyLimit), now)),
	}, nil
}

func overlayCurrentDailyFirstFrame(rows []*statentity.DailyLiveFirstFrameRenderStat, now time.Time) []*statentity.DailyLiveFirstFrameRenderStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyLiveFirstFrameRenderStatByDate(today)
	out := make([]*statentity.DailyLiveFirstFrameRenderStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentWeeklyFirstFrame(rows []*statentity.WeeklyLiveFirstFrameRenderStat, now time.Time) []*statentity.WeeklyLiveFirstFrameRenderStat {
	week := statentity.FormatWeeklyLiveFirstFrameRenderStatKey(now)
	live := statdao.GetWeeklyLiveFirstFrameRenderStatByWeek(week)
	out := make([]*statentity.WeeklyLiveFirstFrameRenderStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == week {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func overlayCurrentMonthlyFirstFrame(rows []*statentity.MonthlyLiveFirstFrameRenderStat, now time.Time) []*statentity.MonthlyLiveFirstFrameRenderStat {
	month := statentity.FormatMonthlyLiveFirstFrameRenderStatKey(now)
	live := statdao.GetMonthlyLiveFirstFrameRenderStatByMonth(month)
	out := make([]*statentity.MonthlyLiveFirstFrameRenderStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == month {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toFirstFrameDailyPoints(rows []*statentity.DailyLiveFirstFrameRenderStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toFirstFrameWeeklyPoints(rows []*statentity.WeeklyLiveFirstFrameRenderStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}

func toFirstFrameMonthlyPoints(rows []*statentity.MonthlyLiveFirstFrameRenderStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Count})
	}
	return list
}
