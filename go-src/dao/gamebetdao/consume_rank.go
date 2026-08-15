package gamebetdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity/game"
	userentity "xr-game-server/entity/user"
)

// GameConsumeRankRow 单场直播游戏消费聚合结果
type GameConsumeRankRow struct {
	UserId      uint64  `json:"user_id"`
	TotalAmount float64 `json:"total_amount"`
}

const gameConsumeRankTopLimit = 500

// SumGameBetByLiveRecord 统计指定直播记录下各用户游戏下注总额(按金额降序,最多500条)
func SumGameBetByLiveRecord(liveRecordId uint64) []*GameConsumeRankRow {
	list := make([]*GameConsumeRankRow, 0)
	if liveRecordId == 0 {
		return list
	}
	ctx := gctx.New()
	now := time.Now()

	sql := `
SELECT bl.` + string(entity.GameBetLogUserId) + ` AS user_id, SUM(bl.` + string(entity.GameBetLogAmount) + `) AS total_amount
FROM ` + string(entity.TbGameBetLog) + ` bl
INNER JOIN ` + string(userentity.TbAccount) + ` a ON a.id = bl.` + string(entity.GameBetLogUserId) + `
LEFT JOIN ` + string(userentity.TbUserExt) + ` ue ON ue.id = bl.` + string(entity.GameBetLogUserId) + `
WHERE bl.` + string(entity.GameBetLogLiveRecordId) + ` = ?
  AND bl.` + string(entity.GameBetLogAmount) + ` > 0
  AND IFNULL(a.` + string(userentity.AccountCancel) + `, 0) = 0
  AND (
    IFNULL(a.` + string(userentity.AccountBan) + `, 0) = 0
    OR (a.` + string(userentity.AccountBanApplyTime) + ` IS NOT NULL AND a.` + string(userentity.AccountBanApplyTime) + ` <= ?)
  )
  AND (ue.id IS NULL OR IFNULL(ue.` + string(userentity.UserExtCanRank) + `, 1) = 1)
GROUP BY bl.` + string(entity.GameBetLogUserId) + `
HAVING SUM(bl.` + string(entity.GameBetLogAmount) + `) > 0
ORDER BY total_amount DESC
LIMIT ?
`
	err := g.DB().Ctx(ctx).Raw(sql, liveRecordId, now, gameConsumeRankTopLimit).Scan(&list)
	if err != nil {
		g.Log().Errorf(ctx, "SumGameBetByLiveRecord liveRecordId=%d error: %v", liveRecordId, err)
		return make([]*GameConsumeRankRow, 0)
	}
	return list
}
