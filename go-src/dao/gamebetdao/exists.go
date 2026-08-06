package gamebetdao

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity"
)

// ExistsGameBetLogByOrderId 是否已存在相同订单/交易 ID 的下注记录(幂等).
func ExistsGameBetLogByOrderId(orderId string) bool {
	orderId = strings.TrimSpace(orderId)
	if orderId == "" {
		return false
	}
	count, err := g.Model(string(entity.TbGameBetLog)).Ctx(gctx.New()).
		Where(string(entity.GameBetLogOrderId)+" = ?", orderId).
		Count()
	return err == nil && count > 0
}
