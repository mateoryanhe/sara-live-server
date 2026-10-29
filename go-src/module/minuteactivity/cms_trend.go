package minuteactivity

import (
	"context"
	"time"

	"xr-game-server/dao/statdao"
	"xr-game-server/dto/trackingdto"
	statentity "xr-game-server/entity/stat"
)

const maxTrendPoints = 4320

func GetCMSMinuteLiveActivityTrend(_ context.Context, req *trackingdto.CMSMinuteLiveActivityTrendReq) (*trackingdto.CMSMinuteLiveActivityTrendRes, error) {
	bypass := req != nil && req.BypassDashboardCache
	if isDefaultMinuteLiveActivityTrendReq(req) {
		return getCachedDashboardUserEventsTrend(bypass, func() (*trackingdto.CMSMinuteLiveActivityTrendRes, error) {
			return buildCMSMinuteLiveActivityTrend(req)
		})
	}
	return buildCMSMinuteLiveActivityTrend(req)
}

func buildCMSMinuteLiveActivityTrend(req *trackingdto.CMSMinuteLiveActivityTrendReq) (*trackingdto.CMSMinuteLiveActivityTrendRes, error) {
	if req == nil {
		req = &trackingdto.CMSMinuteLiveActivityTrendReq{}
	}
	limit := req.Limit
	if limit <= 0 || limit > maxTrendPoints {
		limit = maxTrendPoints
	}
	now := time.Now().UTC()
	end := now
	start := end.AddDate(0, 0, -retentionDays)
	if req.EndTime > 0 {
		end = time.Unix(req.EndTime, 0).UTC()
	}
	if req.StartTime > 0 {
		start = time.Unix(req.StartTime, 0).UTC()
	}
	start = statentity.TruncateUTCMinute(start)
	end = statentity.TruncateUTCMinute(end)
	if end.Before(start) {
		start, end = end, start
	}

	rows := statdao.ListMinuteLiveActivityStats(start, end, limit)
	rows = overlayLiveMinuteBucket(rows, now)
	points := make([]*trackingdto.CMSMinuteLiveActivityTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		points = append(points, &trackingdto.CMSMinuteLiveActivityTrendPoint{
			Time:         row.BucketAt.UTC().Format("2006-01-02 15:04:05"),
			JoinCount:    row.JoinCount,
			LeaveCount:   row.LeaveCount,
			GiftCount:    row.GiftCount,
			ChatCount:    row.ChatCount,
			BetCount:           row.BetCount,
			Call1v1Count:       row.Call1v1Count,
			RechargeOrderCount: row.RechargeOrderCount,
		})
	}
	return &trackingdto.CMSMinuteLiveActivityTrendRes{
		Points:         points,
		RangeTruncated: len(rows) >= limit,
		TrendMaxPoints: limit,
	}, nil
}

func overlayLiveMinuteBucket(rows []*statentity.MinuteLiveActivityStat, now time.Time) []*statentity.MinuteLiveActivityStat {
	live := statdao.GetMinuteLiveActivityStatByTime(now)
	if live == nil {
		return rows
	}
	out := make([]*statentity.MinuteLiveActivityStat, len(rows))
	copy(out, rows)
	bucket := statentity.TruncateUTCMinute(now)
	for i, row := range out {
		if row != nil && row.BucketAt.Equal(bucket) {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}
