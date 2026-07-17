package logquery

import (
	"context"
	"sort"
	"time"

	"xr-game-server/dto/logquerydto"
)

const maxTrendPoints = 2000

var trendIntervalOptions = []int{1, 5, 15, 60}

func GetAccessTrend(_ context.Context, req *logquerydto.CMSGetAccessTrendReq) (*logquerydto.CMSGetAccessTrendRes, error) {
	rangeStart, rangeEnd, err := buildTimeRange(req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	intervalMinutes := resolveTrendIntervalMinutes(req.StartDate, req.EndDate, req.IntervalMinutes)
	var points []logquerydto.AccessTrendPoint
	var totalCount int64

	for {
		counter := make(map[string]int64)
		totalCount = 0
		for _, filePath := range listAccessLogFiles() {
			_ = scanAccessLogFile(filePath, func(entry *AccessLogEntry) bool {
				if !logTimeInRange(entry.Time, rangeStart, rangeEnd) {
					return true
				}
				if !matchAccessEntry(entry, req.TraceId, req.Url, req.Ip, req.StatusCode, req.MinHandlerMs, req.MaxHandlerMs) {
					return true
				}
				t, ok := parseAccessLogTime(entry.Time)
				if !ok {
					return true
				}
				bucket := bucketAccessTime(t, intervalMinutes)
				counter[bucket]++
				totalCount++
				return true
			})
		}

		points = buildTrendPoints(counter)
		if len(points) <= maxTrendPoints || intervalMinutes >= 60 {
			break
		}
		intervalMinutes = bumpTrendInterval(intervalMinutes)
	}

	peakTime, peakCount := findTrendPeak(points)
	return &logquerydto.CMSGetAccessTrendRes{
		IntervalMinutes: intervalMinutes,
		Points:          points,
		TotalCount:      totalCount,
		PeakTime:        peakTime,
		PeakCount:       peakCount,
	}, nil
}

func resolveTrendIntervalMinutes(startDate, endDate string, requested int) int {
	if requested > 0 {
		return normalizeTrendInterval(requested)
	}
	start, _ := parseDate(startDate)
	end, _ := parseDate(endDate)
	days := int(end.Sub(start).Hours()/24) + 1
	switch {
	case days <= 1:
		return 1
	case days <= 3:
		return 5
	default:
		return 15
	}
}

func normalizeTrendInterval(minutes int) int {
	for _, item := range trendIntervalOptions {
		if minutes <= item {
			return item
		}
	}
	return 60
}

func bumpTrendInterval(current int) int {
	for _, item := range trendIntervalOptions {
		if item > current {
			return item
		}
	}
	return 60
}

func bucketAccessTime(t time.Time, intervalMinutes int) string {
	if intervalMinutes <= 1 {
		return t.Format("2006-01-02 15:04")
	}
	minute := (t.Minute() / intervalMinutes) * intervalMinutes
	bucket := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), minute, 0, 0, t.Location())
	return bucket.Format("2006-01-02 15:04")
}

func buildTrendPoints(counter map[string]int64) []logquerydto.AccessTrendPoint {
	points := make([]logquerydto.AccessTrendPoint, 0, len(counter))
	for bucket, count := range counter {
		points = append(points, logquerydto.AccessTrendPoint{
			Time:  bucket,
			Count: count,
		})
	}
	sort.Slice(points, func(i, j int) bool {
		return points[i].Time < points[j].Time
	})
	return points
}

func findTrendPeak(points []logquerydto.AccessTrendPoint) (string, int64) {
	var peakTime string
	var peakCount int64
	for _, point := range points {
		if point.Count > peakCount {
			peakCount = point.Count
			peakTime = point.Time
		}
	}
	return peakTime, peakCount
}
