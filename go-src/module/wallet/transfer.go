package wallet

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/constants/currency"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/golddto"
	userentity "xr-game-server/entity/user"
	"xr-game-server/errercode"
)

func goldTransferLockKey(userId uint64) string {
	return fmt.Sprintf("gold_transfer_%d", userId)
}

func normalizeGoldAmount2(amount float64) float64 {
	return math.Round(amount*100) / 100
}

// AppTransferGold 当前登录用户(币商)向指定用户转赠金币
func AppTransferGold(ctx context.Context, req *golddto.AppTransferGoldReq) (*golddto.AppTransferGoldRes, error) {
	fromUserId := httpserver.GetAuthId(ctx)
	if fromUserId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	toUserId, err := strconv.ParseUint(strings.TrimSpace(req.TargetUserId), 10, 64)
	if err != nil || toUserId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if toUserId == fromUserId {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	amount := normalizeGoldAmount2(req.Amount)
	if amount < 0.01 {
		return nil, errercode.CreateCode(errercode.GoldAmountInvalid)
	}

	fromUser := userinfodao.GetUserInfoByUserId(fromUserId)
	if fromUser == nil || fromUser.UserType != userentity.UserTypeCoinMerchant {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}

	toAccount := accountdao.GetAccountById(toUserId)
	if toAccount == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if toAccount.Cancel {
		return nil, errercode.CreateCode(errercode.AccountCanceled)
	}

	first, second := fromUserId, toUserId
	if first > second {
		first, second = second, first
	}
	gmlock.Lock(goldTransferLockKey(first))
	defer gmlock.Unlock(goldTransferLockKey(first))
	gmlock.Lock(goldTransferLockKey(second))
	defer gmlock.Unlock(goldTransferLockKey(second))

	senderGold, err := GoldSub(fromUserId, amount, currency.ReasonGoldTransferOut)
	if err != nil {
		return nil, err
	}
	targetGold, err := GoldAdd(toUserId, amount, currency.ReasonGoldTransferIn)
	if err != nil {
		if _, refundErr := GoldAdd(fromUserId, amount, currency.ReasonRefund); refundErr != nil {
			return nil, err
		}
		return nil, err
	}

	return &golddto.AppTransferGoldRes{
		Amount:         amount,
		SenderGold:     senderGold,
		TargetUserId:   strconv.FormatUint(toUserId, 10),
		TargetUserGold: targetGold,
	}, nil
}
