package currencylogdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity"
	"xr-game-server/gameevent"
)

const richRankTopLimit = 500

// DiamondConsumeStatRow 钻石消费聚合结果
type DiamondConsumeStatRow struct {
	UserId uint64  `json:"user_id"`
	Total  float64 `json:"total"`
}

// SumDiamondConsumeByUser 统计指定时间范围内钻石消费总额,按用户分组取前500名
func SumDiamondConsumeByUser(startTime, endTime time.Time) []*DiamondConsumeStatRow {
	list := make([]*DiamondConsumeStatRow, 0)
	if endTime.Before(startTime) {
		return list
	}
	ctx := gctx.New()
	err := g.DB().Ctx(ctx).Raw(`
SELECT user_id, SUM(amount) AS total
FROM `+string(entity.TbCurrencyLog)+`
WHERE `+string(entity.CurrencyLogType)+` = ?
  AND `+string(entity.CurrencyLogAction)+` = ?
  AND created_at >= ?
  AND created_at <= ?
GROUP BY user_id
ORDER BY total DESC
LIMIT ?
`, gameevent.CurrencyTypeDiamond, gameevent.CurrencyActionSub, startTime, endTime, richRankTopLimit).Scan(&list)
	if err != nil {
		g.Log().Errorf(ctx, "SumDiamondConsumeByUser error: %v", err)
		return make([]*DiamondConsumeStatRow, 0)
	}
	return list
}
