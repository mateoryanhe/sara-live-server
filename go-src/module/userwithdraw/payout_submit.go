package userwithdraw

import (
	"context"
	"fmt"
	"strings"

	"xr-game-server/constants/country"
	"xr-game-server/core/xrlog"
	"xr-game-server/dao/userwithdrawdao"
	userentity "xr-game-server/entity/user"
	"xr-game-server/module/fxrate"
	"xr-game-server/module/recharge"
)

func submitUserHaiPayPayout(ctx context.Context, row *userentity.UserWithdrawLog) error {
	if row == nil || row.ID == 0 || row.UserId == 0 {
		return fmt.Errorf("empty withdraw log")
	}
	if row.WithdrawUsd <= 0 {
		return fmt.Errorf("withdraw usd <= 0")
	}
	info := userwithdrawdao.GetUserPayoutInfo(row.UserId)
	if info == nil {
		return fmt.Errorf("user payout info missing")
	}
	countryCode := strings.ToUpper(strings.TrimSpace(info.CountryCode))
	if countryCode == "" {
		countryCode = country.HaiPayPayoutCountryFromCurrency(strings.ToUpper(strings.TrimSpace(info.Currency)))
	}
	currencyCode, accountType, bankCode, err := payoutMethodFromUserWithdrawCountryCfg(countryCode)
	if err != nil {
		return fmt.Errorf("payout country cfg invalid")
	}
	if err := validateUserPayoutInfoForWithdrawCountry(info, countryCode); err != nil {
		return fmt.Errorf("user payout info mismatch country cfg")
	}
	conversion, err := fxrate.ConvertUSD(ctx, row.WithdrawUsd, currencyCode)
	if err != nil {
		return err
	}
	localAmount := recharge.HaiPayNormalizePayoutAmount(currencyCode, conversion.TargetAmount)
	if localAmount <= 0 {
		return fmt.Errorf("local amount <= 0")
	}
	xrlog.DetailLog.Infof(ctx,
		"user withdraw payout fx converted logId=%d userId=%d usdAmount=%.4f currency=%s rate=%v localAmount=%v source=%s cached=%t rateDate=%s",
		row.ID, row.UserId, row.WithdrawUsd, currencyCode, conversion.Rate, localAmount,
		conversion.Source, conversion.Cached, conversion.RateDate)

	orderId := FormatPayoutOrderID(row.ID)
	payRes, err := recharge.HaiPayApplyPayout(ctx, &recharge.HaiPayPayoutApplyReq{
		OrderID:       orderId,
		Currency:      currencyCode,
		Amount:        localAmount,
		AccountType:   accountType,
		BankCode:      bankCode,
		AccountNo:     info.AccountNo,
		Name:          info.PayeeName,
		Phone:         info.Phone,
		Email:         info.Email,
		Country:       countryCode,
		IdentifyType:  info.IdentifyType,
		Address1:      info.Address1,
		Address2:      info.Address2,
		Address3:      info.Address3,
		PostalCode:    info.PostalCode,
		PartnerUserID: fmt.Sprintf("%d", row.UserId),
		Body:          fmt.Sprintf("userWithdraw:%d", row.ID),
	})
	if err != nil {
		return err
	}
	platformNo := ""
	if payRes != nil {
		platformNo = payRes.OrderNo
		if responseOrderID := strings.TrimSpace(payRes.OrderID); responseOrderID != "" && responseOrderID != orderId {
			return fmt.Errorf("haipay payout order mismatch got=%s want=%s", responseOrderID, orderId)
		}
	}
	row.SetTransferPayout(orderId, platformNo, currencyCode, localAmount)
	row.SetStatus(userentity.UserWithdrawStatusTransferring)
	row.SetTransferFailMsg("")
	userwithdrawdao.SyncUserWithdrawLogCache(row)
	return nil
}
