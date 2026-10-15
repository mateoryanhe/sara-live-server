package userwithdraw

import userentity "xr-game-server/entity/user"

func withdrawLogDiamondDisplay(row *userentity.UserWithdrawLog) float64 {
	if row == nil {
		return 0
	}
	if row.DiamondAmount > 0 {
		return row.DiamondAmount
	}
	return 0
}
