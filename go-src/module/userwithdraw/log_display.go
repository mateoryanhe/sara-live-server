package userwithdraw

import (
	"strings"

	"xr-game-server/dao/userwithdrawdao"
	userentity "xr-game-server/entity/user"
)

func withdrawLogDiamondDisplay(row *userentity.UserWithdrawLog) float64 {
	if row == nil {
		return 0
	}
	if row.DiamondAmount > 0 {
		return row.DiamondAmount
	}
	return 0
}

// cmsWithdrawLogTransferCurrency 列表展示：优先提现单字段；待审核旧单未落库时回退用户收款信息币种。
func cmsWithdrawLogTransferCurrency(row *userentity.UserWithdrawLog) string {
	if row == nil {
		return ""
	}
	if c := strings.ToUpper(strings.TrimSpace(row.TransferCurrency)); c != "" {
		return c
	}
	info := userwithdrawdao.GetUserPayoutInfo(row.UserId)
	if info == nil {
		return ""
	}
	return strings.ToUpper(strings.TrimSpace(info.Currency))
}
