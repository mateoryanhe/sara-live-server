package liveroomdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity/live"
	userentity "xr-game-server/entity/user"
)

const anchorRankTopLimit = 5500

// AnchorRevenueStatRow 主播社交流水聚合结果
type AnchorRevenueStatRow struct {
	RoomId      uint64 `json:"room_id"`
	TotalAmount uint64 `json:"total_amount"`
}

// SumRevenueByReceiver 统计指定时间范围内主播社交流水总额,按直播间分组取前5500名.
func SumRevenueByReceiver(startTime, endTime time.Time) []*AnchorRevenueStatRow {
	list := make([]*AnchorRevenueStatRow, 0)
	if endTime.Before(startTime) {
		return list
	}
	ctx := gctx.New()
	now := time.Now()
	err := g.DB().Ctx(ctx).Raw(`
SELECT rl.room_id, SUM(rl.total_amount) AS total_amount
FROM `+string(entity.TbLiveRevenueLog)+` rl
INNER JOIN `+string(userentity.TbAccount)+` a ON a.id = rl.room_id
LEFT JOIN `+string(userentity.TbUserExt)+` ue ON ue.id = rl.room_id
WHERE rl.room_id > 0
  AND IFNULL(rl.`+string(entity.LiveRevenueLogStatus)+`, 0) = 0
  AND rl.created_at >= ?
  AND rl.created_at <= ?
  AND IFNULL(a.`+string(userentity.AccountCancel)+`, 0) = 0
  AND (
    IFNULL(a.`+string(userentity.AccountBan)+`, 0) = 0
    OR (a.`+string(userentity.AccountBanApplyTime)+` IS NOT NULL AND a.`+string(userentity.AccountBanApplyTime)+` <= ?)
  )
  AND (ue.id IS NULL OR IFNULL(ue.`+string(userentity.UserExtCanRank)+`, 1) = 1)
GROUP BY rl.room_id
ORDER BY total_amount DESC
LIMIT ?
`, startTime, endTime, now, anchorRankTopLimit).Scan(&list)
	if err != nil {
		g.Log().Errorf(ctx, "SumRevenueByReceiver error: %v", err)
		return make([]*AnchorRevenueStatRow, 0)
	}
	return list
}
