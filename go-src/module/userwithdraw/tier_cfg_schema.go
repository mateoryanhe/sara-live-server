package userwithdraw

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	userentity "xr-game-server/entity/user"
)

func ensureUserWithdrawTierCfgGameConsumeGoldColumn() {
	ctx := gctx.New()
	table := string(userentity.TbUserWithdrawTierCfg)
	hasOld, _ := columnExists(ctx, table, "min_bet_total_usd")
	hasNew, _ := columnExists(ctx, table, "min_game_consume_gold")
	if hasOld && !hasNew {
		_, _ = g.DB().Exec(ctx, `ALTER TABLE user_withdraw_tier_cfgs CHANGE COLUMN min_bet_total_usd min_game_consume_gold DECIMAL(20,2) NOT NULL DEFAULT 0 COMMENT '游戏消费累计金币门槛(配置,逻辑待接)'`)
		return
	}
	if hasOld && hasNew {
		_, _ = g.DB().Exec(ctx, `UPDATE user_withdraw_tier_cfgs SET min_game_consume_gold = min_bet_total_usd WHERE min_game_consume_gold = 0 AND min_bet_total_usd <> 0`)
		_, _ = g.DB().Exec(ctx, `ALTER TABLE user_withdraw_tier_cfgs DROP COLUMN min_bet_total_usd`)
	}
}

func columnExists(ctx context.Context, table, column string) (bool, error) {
	var count int
	err := g.DB().GetScan(ctx, &count, `
		SELECT COUNT(*) FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = ?`, table, column)
	return count > 0, err
}
