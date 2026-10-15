package gamebetdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	gameentity "xr-game-server/entity/game"
)

// SumGameBetAmountByUserId 用户历史游戏下注(消费)金币累计.
func SumGameBetAmountByUserId(userId uint64) float64 {
	if userId == 0 {
		return 0
	}
	ctx := gctx.New()
	val, err := g.Model(string(gameentity.TbGameBetLog)).Ctx(ctx).
		Where(string(gameentity.GameBetLogUserId)+" = ?", userId).
		Sum(string(gameentity.GameBetLogAmount))
	if err != nil {
		return 0
	}
	return val
}
