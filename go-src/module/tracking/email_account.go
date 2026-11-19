package tracking

import (
	"context"
	"time"

	"xr-game-server/dao/statdao"
	"xr-game-server/dto/trackingdto"
	statentity "xr-game-server/entity/stat"
)

const (
	emailAccountDailyLimit = 30

	// EventKeyEmailAccount 埋点 event_key：user_ext.email 按 user_id 首次写入 +1(绑定或邮箱新注册)
	EventKeyEmailAccount = "email_account"
)

func recordEmailAccount(at time.Time) {
	if at.IsZero() {
		at = time.Now()
	}
	date := statentity.FormatDailyLoginStatDate(at)
	daily := statdao.GetDailyEmailAccountStatByDate(date)
	daily.AddVal(1)
	statdao.PublishDailyEmailAccountStat(daily)
}

func GetCMSEmailAccountTrend(_ context.Context, _ *trackingdto.CMSEmailAccountTrendReq) (*trackingdto.CMSEmailAccountTrendRes, error) {
	now := time.Now()
	today := statentity.FormatDailyLoginStatDate(now)
	todayStat := statdao.GetDailyEmailAccountStatByDate(today)
	return &trackingdto.CMSEmailAccountTrendRes{
		TodayVal: todayStat.Val,
		Daily:    toEmailAccountDailyPoints(overlayCurrentDailyEmailAccount(statdao.ListRecentDailyEmailAccountStats(emailAccountDailyLimit), now)),
	}, nil
}

func overlayCurrentDailyEmailAccount(rows []*statentity.DailyEmailAccountStat, now time.Time) []*statentity.DailyEmailAccountStat {
	today := statentity.FormatDailyLoginStatDate(now)
	live := statdao.GetDailyEmailAccountStatByDate(today)
	out := make([]*statentity.DailyEmailAccountStat, len(rows))
	copy(out, rows)
	for i, row := range out {
		if row != nil && row.ID == today {
			out[i] = live
			return out
		}
	}
	return append(out, live)
}

func toEmailAccountDailyPoints(rows []*statentity.DailyEmailAccountStat) []*trackingdto.CMSTrackingEventTrendPoint {
	list := make([]*trackingdto.CMSTrackingEventTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trackingdto.CMSTrackingEventTrendPoint{Time: row.ID, Count: row.Val})
	}
	return list
}