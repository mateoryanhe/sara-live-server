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
	now := time.Now()
	err := g.DB().Ctx(ctx).Raw(`
SELECT cl.user_id, SUM(cl.amount) AS total
FROM `+string(entity.TbCurrencyLog)+` cl
INNER JOIN `+string(entity.TbAccount)+` a ON a.id = cl.user_id
WHERE cl.`+string(entity.CurrencyLogType)+` = ?
  AND cl.`+string(entity.CurrencyLogAction)+` = ?
  AND cl.created_at >= ?
  AND cl.created_at <= ?
  AND IFNULL(a.`+string(entity.AccountCancel)+`, 0) = 0
  AND (
    IFNULL(a.`+string(entity.AccountBan)+`, 0) = 0
    OR (a.`+string(entity.AccountBanApplyTime)+` IS NOT NULL AND a.`+string(entity.AccountBanApplyTime)+` <= ?)
  )
GROUP BY cl.user_id
ORDER BY total DESC
LIMIT ?
`, gameevent.CurrencyTypeDiamond, gameevent.CurrencyActionSub, startTime, endTime, now, richRankTopLimit).Scan(&list)
	if err != nil {
		g.Log().Errorf(ctx, "SumDiamondConsumeByUser error: %v", err)
		return make([]*DiamondConsumeStatRow, 0)
	}
	return list
}
