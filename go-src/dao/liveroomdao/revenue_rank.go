package liveroomdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity"
)

const anchorRankTopLimit = 500

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
	err := g.DB().Ctx(ctx).Raw(`
SELECT receiver_id, SUM(total_amount) AS total_amount
FROM `+string(entity.TbLiveRevenueLog)+`
WHERE receiver_id > 0
  AND created_at >= ?
  AND created_at <= ?
GROUP BY receiver_id
ORDER BY total_amount DESC
LIMIT ?
`, startTime, endTime, anchorRankTopLimit).Scan(&list)
	if err != nil {
		g.Log().Errorf(ctx, "SumRevenueByReceiver error: %v", err)
		return make([]*AnchorRevenueStatRow, 0)
	}
	return list
}
