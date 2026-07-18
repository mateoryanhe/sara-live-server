package liveroomdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity"
)

const anchorRankTopLimit = 5500

// AnchorRevenueStatRow 主播收益聚合结果
type AnchorRevenueStatRow struct {
	ReceiverId  uint64 `json:"receiver_id"`
	TotalAmount uint64 `json:"total_amount"`
}

// SumRevenueByReceiver 统计指定时间范围内主播收益总额,按主播分组取前500名
func SumRevenueByReceiver(startTime, endTime time.Time) []*AnchorRevenueStatRow {
	list := make([]*AnchorRevenueStatRow, 0)
	if endTime.Before(startTime) {
		return list
	}
	ctx := gctx.New()
	now := time.Now()
	err := g.DB().Ctx(ctx).Raw(`
SELECT rl.receiver_id, SUM(rl.total_amount) AS total_amount
FROM `+string(entity.TbLiveRevenueLog)+` rl
INNER JOIN `+string(entity.TbAccount)+` a ON a.id = rl.receiver_id
LEFT JOIN `+string(entity.TbUserExt)+` ue ON ue.id = rl.receiver_id
WHERE rl.receiver_id > 0
  AND COALESCE(rl.`+string(entity.LiveRevenueLogStatus)+`, 0) = 0
  AND rl.created_at >= ?
  AND rl.created_at <= ?
  AND COALESCE(a.`+string(entity.AccountCancel)+`, false) = false
  AND (
    COALESCE(a.`+string(entity.AccountBan)+`, false) = false
    OR (a.`+string(entity.AccountBanApplyTime)+` IS NOT NULL AND a.`+string(entity.AccountBanApplyTime)+` <= ?)
  )
  AND (ue.id IS NULL OR COALESCE(ue.`+string(entity.UserExtCanRank)+`, true) = true)
GROUP BY rl.receiver_id
ORDER BY total_amount DESC
LIMIT ?
`, startTime, endTime, now, anchorRankTopLimit).Scan(&list)
	if err != nil {
		g.Log().Errorf(ctx, "SumRevenueByReceiver error: %v", err)
		return make([]*AnchorRevenueStatRow, 0)
	}
	return list
}
