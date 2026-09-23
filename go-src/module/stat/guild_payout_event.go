package stat

import (
	"math"
	"time"

	"xr-game-server/core/event"
	"xr-game-server/dao/statdao"
	statentity "xr-game-server/entity/stat"
	"xr-game-server/gameevent"
)

func initGuildPayoutEvent() {
	event.Sub(gameevent.GuildPayoutSucceededEvent, onGuildPayoutSucceededEvent)
}

func onGuildPayoutSucceededEvent(data any) {
	payload, ok := data.(*gameevent.GuildPayoutSucceededEventData)
	if !ok || payload == nil || payload.SettlementID == 0 || payload.UsdAmount <= 0 ||
		math.IsNaN(payload.UsdAmount) || math.IsInf(payload.UsdAmount, 0) {
		return
	}
	enqueue(&statJob{Kind: jobGuildPayout, At: payload.PaidAt, Payload: payload})
}

func consumeGuildPayoutJob(job *statJob) {
	data, ok := job.Payload.(*gameevent.GuildPayoutSucceededEventData)
	if !ok || data == nil || data.UsdAmount <= 0 || math.IsNaN(data.UsdAmount) || math.IsInf(data.UsdAmount, 0) {
		return
	}

	if total := statdao.GetSysStat(); total != nil {
		total.AddTotalAnchorPayout(data.UsdAmount)
	}

	statAt := data.PaidAt
	if statAt.IsZero() {
		statAt = job.At
	}
	if statAt.IsZero() {
		statAt = time.Now()
	}
	weekly := statdao.GetWeeklyLoginStatByWeek(statentity.FormatWeeklyLoginStatKey(statAt))
	weekly.AddAnchorPayoutAmount(data.UsdAmount)
}
