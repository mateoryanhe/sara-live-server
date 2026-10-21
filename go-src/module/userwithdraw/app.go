package userwithdraw

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

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

func startOfLocalDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func GetWithdrawCfg(_ context.Context, _ *userwithdrawdto.AppGetWithdrawCfgReq) (*userwithdrawdto.AppGetWithdrawCfgRes, error) {
	return &userwithdrawdto.AppGetWithdrawCfgRes{
		Enabled:           accountcfg.IsUserWithdrawEnabled(),
		DailyMaxCount:     accountcfg.GetUserWithdrawDailyMaxCount(),
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
	if req == nil || req.TierId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	tier := GetUserWithdrawTierFromCacheById(req.TierId)
	if tier == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	currencyCode := strings.ToUpper(strings.TrimSpace(req.Currency))
	if err := validateAppWithdrawCurrency(currencyCode); err != nil {
		return nil, err
	}
	payoutInfo := userwithdrawdao.GetUserPayoutInfo(userId)
	if err := validateUserPayoutInfoForWithdrawCurrency(payoutInfo, currencyCode); err != nil {
		return nil, err
	}
	if dailyMax := accountcfg.GetUserWithdrawDailyMaxCount(); dailyMax > 0 {
		since := startOfLocalDay(time.Now())
		cnt, err := userwithdrawdao.CountUserWithdrawLogsSince(userId, since)
		if err != nil {
			return nil, err
		}
		if cnt >= dailyMax {
			return nil, errercode.CreateCode(errercode.DailyLimitExceeded)
		}
	}
	if err := checkWithdrawTierGameConsumeGate(userId, tier); err != nil {
		return nil, err
	}

	gmlock.Lock(withdrawLockKey(userId))
	defer gmlock.Unlock(withdrawLockKey(userId))

	snap := wallet.GetExchangeCfgSnapshot()
	withdrawUsd := math.Round(tier.WithdrawUsd*10000) / 10000
	diamondAmount := wallet.CalcUsdToDiamondWithSnapshot(withdrawUsd, snap)
	if withdrawUsd <= 0 || diamondAmount <= 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	if _, err := wallet.DiamondSubWithGoldExchange(userId, diamondAmount, currency.ReasonUserWithdraw); err != nil {
		return nil, err
	}

	row := userentity.NewUserWithdrawLog(userId, req.TierId, diamondAmount, 0, withdrawUsd, snap.UsdToGoldRate)
	row.SetTransferCurrency(currencyCode)
	userwithdrawdao.PublishNewUserWithdrawLog(row)

	return &userwithdrawdto.AppApplyWithdrawRes{
		Id:            strconv.FormatUint(row.ID, 10),
		DiamondAmount: diamondAmount,
		WithdrawUsd:   withdrawUsd,
		Status:        row.Status,
	}, nil
}

func AppListWithdrawLogs(ctx context.Context, req *userwithdrawdto.AppWithdrawLogListReq) (*userwithdrawdto.AppWithdrawLogListRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	page := 1
	if req != nil && req.Page > 0 {
		page = req.Page
	}
	total, err := userwithdrawdao.CountUserWithdrawLogsByUser(userId)
	if err != nil {
		return nil, err
	}
	rows := userwithdrawdao.GetAppUserWithdrawLogsByUser(userId, page, userwithdrawdao.AppUserWithdrawLogListCachePageSize)
	list := make([]*userwithdrawdto.AppWithdrawLogItem, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &userwithdrawdto.AppWithdrawLogItem{
			Id:                  strconv.FormatUint(row.ID, 10),
			DiamondAmount:       withdrawLogDiamondDisplay(row),
			WithdrawUsd:         row.WithdrawUsd,
			Status:              row.Status,
			TransferCurrency:    cmsWithdrawLogTransferCurrency(row),
			TransferLocalAmount: row.TransferLocalAmount,
			TransferFailMsg:     row.TransferFailMsg,
			RejectReason:        row.RejectReason,
			CreatedAt:           formatLogTime(row.CreatedAt),
			TransferAt:          formatLogTimePtr(row.TransferAt),
		})
	}
	return &userwithdrawdto.AppWithdrawLogListRes{Total: total, List: list}, nil
}
