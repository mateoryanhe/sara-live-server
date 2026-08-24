package currencylogdao

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity/user"
)

// ExistsByTransactionId 是否已存在相同第三方交易 ID 的货币流水.
func ExistsByTransactionId(transactionId string) bool {
	transactionId = strings.TrimSpace(transactionId)
	if transactionId == "" {
		return false
	}
	count, err := g.Model(string(entity.TbCurrencyLog)).Ctx(gctx.New()).
		Where(string(entity.CurrencyLogTransactionId)+" = ?", transactionId).
		Count()
	return err == nil && count > 0
}
