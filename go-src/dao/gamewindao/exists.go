package gamewindao

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity"
)

// ExistsGameWinLogByOrderId 是否已存在相同订单/交易 ID 的派彩记录(幂等).
func ExistsGameWinLogByOrderId(orderId string) bool {
	orderId = strings.TrimSpace(orderId)
	if orderId == "" {
		return false
	}
	count, err := g.Model(string(entity.TbGameWinLog)).Ctx(gctx.New()).
		Where(string(entity.GameWinLogOrderId)+" = ?", orderId).
		Count()
	return err == nil && count > 0
}
