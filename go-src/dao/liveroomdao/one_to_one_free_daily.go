package liveroomdao

import (
	"fmt"
	"time"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/snowflake"
	liveentity "xr-game-server/entity/live"
)

var oneToOneFreeDailyUsedCache = gmap.NewKVMap[string, struct{}](false)

func oneToOneFreeDailyCacheKey(payerId, anchorId uint64, statDate string) string {
	return fmt.Sprintf("%d:%d:%s", payerId, anchorId, statDate)
}

// OneToOneFreeStatDate 服务器本地自然日 YYYY-MM-DD。
func OneToOneFreeStatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func initOneToOneFreeDailyDao() {
	rows := make([]*liveentity.OneToOneCallFreeDailyUse, 0)
	today := OneToOneFreeStatDate(time.Now())
	_ = g.Model(string(liveentity.TbOneToOneCallFreeDailyUse)).Where("stat_date", today).Scan(&rows)
	for _, row := range rows {
		if row == nil || row.PayerId == 0 || row.AnchorId == 0 || row.StatDate == "" {
			continue
		}
		oneToOneFreeDailyUsedCache.Set(oneToOneFreeDailyCacheKey(row.PayerId, row.AnchorId, row.StatDate), struct{}{})
	}
}

// IsOneToOneDailyFreeAvailable 当日对该主播是否仍有 1v1 整段免费额度。
func IsOneToOneDailyFreeAvailable(payerId, anchorId uint64) bool {
	if payerId == 0 || anchorId == 0 {
		return false
	}
	statDate := OneToOneFreeStatDate(time.Now())
	key := oneToOneFreeDailyCacheKey(payerId, anchorId, statDate)
	if oneToOneFreeDailyUsedCache.Contains(key) {
		return false
	}
	count, err := g.Model(string(liveentity.TbOneToOneCallFreeDailyUse)).Ctx(gctx.New()).
		Where("payer_id", payerId).
		Where("anchor_id", anchorId).
		Where("stat_date", statDate).
		Count()
	if err != nil {
		return false
	}
	if count > 0 {
		oneToOneFreeDailyUsedCache.Set(key, struct{}{})
		return false
	}
	return true
}

// MarkOneToOneDailyFreeUsed 标记当日整段免费额度已使用(接通即作废,与通话时长无关)。
func MarkOneToOneDailyFreeUsed(payerId, anchorId uint64) {
	if payerId == 0 || anchorId == 0 {
		return
	}
	statDate := OneToOneFreeStatDate(time.Now())
	key := oneToOneFreeDailyCacheKey(payerId, anchorId, statDate)
	if oneToOneFreeDailyUsedCache.Contains(key) {
		return
	}
	row := &liveentity.OneToOneCallFreeDailyUse{
		PayerId:  payerId,
		AnchorId: anchorId,
		StatDate: statDate,
	}
	row.ID = snowflake.GetId()
	_, _ = g.Model(string(liveentity.TbOneToOneCallFreeDailyUse)).Ctx(gctx.New()).Data(row).Insert()
	oneToOneFreeDailyUsedCache.Set(key, struct{}{})
}
