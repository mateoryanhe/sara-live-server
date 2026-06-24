package liveroomdao

import (
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	liverevenueconst "xr-game-server/constants/liverevenue"
	"xr-game-server/entity"
)

// AudienceContributionStatRow 观众贡献聚合结果
type AudienceContributionStatRow struct {
	SenderId    uint64  `json:"sender_id"`
	TotalAmount float64 `json:"total_amount"`
}

const contributionRankTopLimit = 500

// SumAudienceContributionByRoom 统计指定直播间在时间范围内、且当前在线观众的贡献(礼物+付费弹幕)
func SumAudienceContributionByRoom(roomId uint64, startTime, endTime time.Time, onlineUserIds []uint64) []*AudienceContributionStatRow {
	list := make([]*AudienceContributionStatRow, 0)
	if roomId == 0 || endTime.Before(startTime) || len(onlineUserIds) == 0 {
		return list
	}
	ctx := gctx.New()
	now := time.Now()

	inPlaceholders := strings.Repeat("?,", len(onlineUserIds))
	inPlaceholders = inPlaceholders[:len(inPlaceholders)-1]

	args := make([]any, 0, 6+len(onlineUserIds))
	args = append(args, roomId)
	for _, userId := range onlineUserIds {
		args = append(args, userId)
	}
	args = append(args,
		uint8(liverevenueconst.Gift),
		uint8(liverevenueconst.PaidDanmaku),
		startTime,
		endTime,
		now,
		contributionRankTopLimit,
	)

	sql := `
SELECT rl.sender_id, SUM(rl.total_amount) AS total_amount
FROM ` + string(entity.TbLiveRevenueLog) + ` rl
INNER JOIN ` + string(entity.TbAccount) + ` a ON a.id = rl.sender_id
WHERE rl.room_id = ?
  AND rl.sender_id IN (` + inPlaceholders + `)
  AND rl.revenue_type IN (?, ?)
  AND rl.created_at >= ?
  AND rl.created_at <= ?
  AND IFNULL(a.` + string(entity.AccountCancel) + `, 0) = 0
  AND (
    IFNULL(a.` + string(entity.AccountBan) + `, 0) = 0
    OR (a.` + string(entity.AccountBanApplyTime) + ` IS NOT NULL AND a.` + string(entity.AccountBanApplyTime) + ` <= ?)
  )
GROUP BY rl.sender_id
HAVING SUM(rl.total_amount) > 0
ORDER BY total_amount DESC
LIMIT ?
`
	err := g.DB().Ctx(ctx).Raw(sql, args...).Scan(&list)
	if err != nil {
		g.Log().Errorf(ctx, "SumAudienceContributionByRoom roomId=%d error: %v", roomId, err)
		return make([]*AudienceContributionStatRow, 0)
	}
	return list
}
