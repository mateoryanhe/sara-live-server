package userwithdraw

import (
	"xr-game-server/dao/gamebetdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/userwithdrawtiercfgdto"
	"xr-game-server/errercode"
)

const withdrawGameConsumeEpsilon = 0.0001

func checkWithdrawTierGameConsumeGate(userId uint64, tier *userwithdrawtiercfgdto.AppUserWithdrawTierItem) error {
	if tier == nil {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	ext := userinfodao.GetUserExtByUserId(userId)
	if ext == nil || ext.GameWinCount == 0 {
		return nil
	}
	required := tier.MinGameConsumeGold
	if required <= 0 {
		return nil
	}
	total := gamebetdao.SumGameBetAmountByUserId(userId)
	if total+withdrawGameConsumeEpsilon < required {
		return errercode.CreateCode(errercode.UserWithdrawGameConsumeNotEnough)
	}
	return nil
}
