package userwithdraw

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/constants/currency"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/userwithdrawdao"
	"xr-game-server/dto/userwithdrawdto"
	userentity "xr-game-server/entity/user"
	"xr-game-server/errercode"
	"xr-game-server/module/accountcfg"
	"xr-game-server/module/wallet"
)

func withdrawLockKey(userId uint64) string {
	return fmt.Sprintf("user_withdraw_%d", userId)
}

func GetWithdrawCfg(_ context.Context, _ *userwithdrawdto.AppGetWithdrawCfgReq) (*userwithdrawdto.AppGetWithdrawCfgRes, error) {
	return &userwithdrawdto.AppGetWithdrawCfgRes{
		Enabled: accountcfg.IsUserWithdrawEnabled(),
		MinUsd:  accountcfg.GetUserWithdrawMinUsd(),
	}, nil
}

func AppGetPayoutInfo(ctx context.Context, _ *userwithdrawdto.AppGetPayoutInfoReq) (*userwithdrawdto.AppGetPayoutInfoRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	return GetPayoutInfo(ctx, userId)
}

func AppSavePayoutInfo(ctx context.Context, req *userwithdrawdto.AppSavePayoutInfoReq) (*userwithdrawdto.AppSavePayoutInfoRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if !accountcfg.IsUserWithdrawEnabled() {
		return nil, errercode.CreateCode(errercode.UserWithdrawDisabled)
	}
	return SavePayoutInfo(ctx, userId, req)
}

func AppApplyWithdraw(ctx context.Context, req *userwithdrawdto.AppApplyWithdrawReq) (*userwithdrawdto.AppApplyWithdrawRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if !accountcfg.IsUserWithdrawEnabled() {
		return nil, errercode.CreateCode(errercode.UserWithdrawDisabled)
	}
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	goldAmount := math.Round(req.GoldAmount*100) / 100
	if goldAmount <= 0 {
		return nil, errercode.CreateCode(errercode.GoldAmountInvalid)
	}
	if userwithdrawdao.GetUserPayoutInfo(userId) == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	gmlock.Lock(withdrawLockKey(userId))
	defer gmlock.Unlock(withdrawLockKey(userId))

	snap := wallet.GetExchangeCfgSnapshot()
	withdrawUsd := wallet.CalcGoldToUsdWithSnapshot(goldAmount, snap)
	minUsd := accountcfg.GetUserWithdrawMinUsd()
	if withdrawUsd < minUsd {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	if _, err := wallet.GoldSub(userId, goldAmount, currency.ReasonUserWithdraw); err != nil {
		return nil, err
	}

	row := userentity.NewUserWithdrawLog(userId, goldAmount, withdrawUsd, snap.UsdToGoldRate)
	if err := userwithdrawdao.InsertUserWithdrawLog(row); err != nil {
		if _, refundErr := wallet.GoldAdd(userId, goldAmount, currency.ReasonUserWithdrawRefund); refundErr != nil {
			return nil, err
		}
		return nil, err
	}

	return &userwithdrawdto.AppApplyWithdrawRes{
		Id:          strconv.FormatUint(row.ID, 10),
		GoldAmount:  goldAmount,
		WithdrawUsd: withdrawUsd,
		Status:      row.Status,
	}, nil
}

func AppListWithdrawLogs(ctx context.Context, req *userwithdrawdto.AppWithdrawLogListReq) (*userwithdrawdto.AppWithdrawLogListRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if req == nil {
		req = &userwithdrawdto.AppWithdrawLogListReq{}
	}
	total, rows, err := userwithdrawdao.ListUserWithdrawLogs(&userwithdrawdao.UserWithdrawListFilter{
		UserId:   userId,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	list := make([]*userwithdrawdto.AppWithdrawLogItem, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &userwithdrawdto.AppWithdrawLogItem{
			Id:                  strconv.FormatUint(row.ID, 10),
			GoldAmount:          row.GoldAmount,
			WithdrawUsd:         row.WithdrawUsd,
			Status:              row.Status,
			TransferCurrency:    row.TransferCurrency,
			TransferLocalAmount: row.TransferLocalAmount,
			TransferFailMsg:     row.TransferFailMsg,
			RejectReason:        row.RejectReason,
			CreatedAt:           formatLogTime(row.CreatedAt),
			TransferAt:          formatLogTimePtr(row.TransferAt),
		})
	}
	return &userwithdrawdto.AppWithdrawLogListRes{Total: total, List: list}, nil
}
