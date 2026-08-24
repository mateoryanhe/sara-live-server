package currencylogdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/currency"
	"xr-game-server/entity/user"
	"xr-game-server/gameevent"
)

const gameConsumeRankTopLimit = 500

// SumGameGoldConsumeByUser 统计指定时间范围内游戏类金币消费总额,按用户分组取前500名.
func SumGameGoldConsumeByUser(startTime, endTime time.Time) []*DiamondConsumeStatRow {
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
LEFT JOIN `+string(entity.TbUserExt)+` ue ON ue.id = cl.user_id
WHERE cl.`+string(entity.CurrencyLogType)+` = ?
  AND cl.`+string(entity.CurrencyLogAction)+` = ?
  AND cl.created_at >= ?
  AND cl.created_at <= ?
  AND IFNULL(a.`+string(entity.AccountCancel)+`, 0) = 0
  AND (
    IFNULL(a.`+string(entity.AccountBan)+`, 0) = 0
    OR (a.`+string(entity.AccountBanApplyTime)+` IS NOT NULL AND a.`+string(entity.AccountBanApplyTime)+` <= ?)
  )
  AND (ue.id IS NULL OR IFNULL(ue.`+string(entity.UserExtCanRank)+`, 1) = 1)
  AND IFNULL(cl.`+string(entity.CurrencyLogBusinessType)+`, ?) = ?
GROUP BY cl.user_id
ORDER BY total DESC
LIMIT ?
`, gameevent.CurrencyTypeGold, gameevent.CurrencyActionSub, startTime, endTime, now, currency.BusinessTypeSocial, currency.BusinessTypeGame, gameConsumeRankTopLimit).Scan(&list)
	if err != nil {
		g.Log().Errorf(ctx, "SumGameGoldConsumeByUser error: %v", err)
		return make([]*DiamondConsumeStatRow, 0)
	}
	return list
}
