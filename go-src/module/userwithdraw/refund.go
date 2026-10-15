package userwithdraw

import (
	"xr-game-server/constants/currency"
	userentity "xr-game-server/entity/user"
	"xr-game-server/module/wallet"
)

func refundUserWithdraw(userId uint64, row *userentity.UserWithdrawLog) error {
	if row == nil {
		return nil
	}
	if row.DiamondAmount > 0 {
		_, err := wallet.DiamondAdd(userId, row.DiamondAmount, currency.ReasonUserWithdrawRefund)
		return err
	}
	if row.GoldAmount > 0 {
		_, err := wallet.GoldAdd(userId, row.GoldAmount, currency.ReasonUserWithdrawRefund)
		return err
	}
	return nil
}
