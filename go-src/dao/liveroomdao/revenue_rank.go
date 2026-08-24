package liveroomdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/liverevenue"
	liveentity "xr-game-server/entity/live"
	userentity "xr-game-server/entity/user"
)

const anchorRankTopLimit = 5500

// AnchorRevenueStatRow 主播收益聚合结果
type AnchorRevenueStatRow struct {
	ReceiverId  uint64 `json:"receiver_id"`
	TotalAmount uint64 `json:"total_amount"`
}

// SumRevenueByReceiver 统计指定时间范围内主播社交类收益总额,按主播分组取前500名(不含游戏收益).
func SumRevenueByReceiver(startTime, endTime time.Time) []*AnchorRevenueStatRow {
	list := make([]*AnchorRevenueStatRow, 0)
	if endTime.Before(startTime) {
		return list
	}
	ctx := gctx.New()
	now := time.Now()
	err := g.DB().Ctx(ctx).Raw(`
SELECT rl.receiver_id, SUM(rl.total_amount) AS total_amount
FROM `+string(liveentity.TbLiveRevenueLog)+` rl
INNER JOIN `+string(userentity.TbAccount)+` a ON a.id = rl.receiver_id
LEFT JOIN `+string(userentity.TbUserExt)+` ue ON ue.id = rl.receiver_id
WHERE rl.receiver_id > 0
  AND IFNULL(rl.`+string(liveentity.LiveRevenueLogStatus)+`, 0) = 0
  AND rl.created_at >= ?
  AND rl.created_at <= ?
  AND IFNULL(a.`+string(userentity.AccountCancel)+`, 0) = 0
  AND (
    IFNULL(a.`+string(userentity.AccountBan)+`, 0) = 0
    OR (a.`+string(userentity.AccountBanApplyTime)+` IS NOT NULL AND a.`+string(userentity.AccountBanApplyTime)+` <= ?)
  )
  AND (ue.id IS NULL OR IFNULL(ue.`+string(userentity.UserExtCanRank)+`, 1) = 1)
  AND IFNULL(rl.`+string(liveentity.LiveRevenueLogRevenueType)+`, 0) <> ?
GROUP BY rl.receiver_id
ORDER BY total_amount DESC
LIMIT ?
`, startTime, endTime, now, liverevenue.GameBet, anchorRankTopLimit).Scan(&list)
	if err != nil {
		g.Log().Errorf(ctx, "SumRevenueByReceiver error: %v", err)
		return make([]*AnchorRevenueStatRow, 0)
	}
	return list
}
