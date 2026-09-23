package incomesettlement

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"xr-game-server/core/syndb"
	"xr-game-server/core/xrlog"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/incomesettlementdto"
	live "xr-game-server/entity/live"
	"xr-game-server/errercode"
	"xr-game-server/module/cmsvis"
	"xr-game-server/module/fxrate"
	"xr-game-server/module/recharge"
	"xr-game-server/module/wallet"
)

func platformAnchorVisibleSet(ctx context.Context) (map[uint64]struct{}, bool, bool) {
	ids, restrict, empty := cmsvis.PlatformAnchorVisibilityFilter(ctx)
	if empty {
		return nil, true, true
	}
	if !restrict {
		return nil, false, false
	}
	set := make(map[uint64]struct{}, len(ids))
	for _, id := range ids {
		set[id] = struct{}{}
	}
	return set, true, false
}

// BatchApproveAnchorSettlement 批量审核平台主播直接代付单。
func BatchApproveAnchorSettlement(ctx context.Context, req *incomesettlementdto.CMSBatchApproveAnchorSettlementReq) (*incomesettlementdto.CMSBatchApproveAnchorSettlementRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	ids := parseIdList(req.Ids)
	if len(ids) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	visibleSet, restrict, empty := platformAnchorVisibleSet(ctx)
	res := &incomesettlementdto.CMSBatchApproveAnchorSettlementRes{}
	if empty {
		res.FailCount = len(ids)
		return res, nil
	}
	for _, id := range ids {
		row := liveroomdao.GetAnchorIncomeSettlementLogById(id)
		if row == nil || !row.DirectPayout || row.Status != live.AnchorIncomeSettlementStatusPending {
			res.FailCount++
			continue
		}
		if restrict {
			if _, ok := visibleSet[row.RoomId]; !ok {
				res.FailCount++
				continue
			}
		}
		row.SetStatus(live.AnchorIncomeSettlementStatusApproved)
		liveroomdao.PublishAnchorIncomeSettlementLog(row)
		res.SuccessCount++
	}
	return res, nil
}

// BatchTransferAnchorSettlement 批量提交平台主播 HaiPay 代付。
func BatchTransferAnchorSettlement(ctx context.Context, req *incomesettlementdto.CMSBatchTransferAnchorSettlementReq) (*incomesettlementdto.CMSBatchTransferAnchorSettlementRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	ids := parseIdList(req.Ids)
	if len(ids) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	visibleSet, restrict, empty := platformAnchorVisibleSet(ctx)
	res := &incomesettlementdto.CMSBatchTransferAnchorSettlementRes{}
	if empty {
		res.FailCount = len(ids)
		res.Message = "无可见平台主播权限"
		return res, nil
	}
	for _, id := range ids {
		row := liveroomdao.GetAnchorIncomeSettlementLogById(id)
		if row == nil || !row.DirectPayout || row.Status != live.AnchorIncomeSettlementStatusApproved {
			res.FailCount++
			continue
		}
		if restrict {
			if _, ok := visibleSet[row.RoomId]; !ok {
				res.FailCount++
				continue
			}
		}
		if err := submitAnchorHaiPayPayout(ctx, row); err != nil {
			xrlog.DetailLog.Warningf(ctx, "platform anchor payout submit failed settlementId=%d anchorId=%d err=%v", row.ID, row.RoomId, err)
			row.SetTransferFailMsg(err.Error())
			liveroomdao.PublishAnchorIncomeSettlementLog(row)
			res.FailCount++
			continue
		}
		res.SuccessCount++
	}
	if res.SuccessCount == 0 && res.FailCount > 0 {
		res.Message = "代付提交失败,请检查HaiPay全局配置、主播转账信息或汇率接口"
	} else if res.FailCount > 0 {
		res.Message = fmt.Sprintf("部分成功:成功%d 失败%d", res.SuccessCount, res.FailCount)
	} else {
		res.Message = fmt.Sprintf("已提交代付%d笔,等待HaiPay回调确认", res.SuccessCount)
	}
	return res, nil
}

func submitAnchorHaiPayPayout(ctx context.Context, row *live.AnchorIncomeSettlementLog) error {
	if row == nil || row.ID == 0 || row.RoomId == 0 || !row.DirectPayout {
		return fmt.Errorf("empty platform anchor settlement")
	}
	if err := prepareTieredAnchorPayout(row); err != nil {
		return err
	}
	info := liveroomdao.GetAnchorTransferInfo(row.RoomId)
	if info == nil {
		return fmt.Errorf("anchor transfer info missing")
	}
	currency := strings.ToUpper(strings.TrimSpace(info.Currency))
	if currency == "" {
		return fmt.Errorf("transfer currency empty")
	}
	conversion, err := fxrate.ConvertUSD(ctx, row.SettlementReceivableUsd, currency)
	if err != nil {
		return err
	}
	localAmount := recharge.HaiPayNormalizePayoutAmount(currency, conversion.TargetAmount)
	if localAmount <= 0 {
		return fmt.Errorf("local amount <= 0")
	}
	xrlog.DetailLog.Infof(ctx,
		"platform anchor payout fx converted settlementId=%d anchorId=%d usdAmount=%.4f currency=%s rate=%v localAmount=%v source=%s cached=%t rateDate=%s",
		row.ID, row.RoomId, row.SettlementReceivableUsd, currency, conversion.Rate, localAmount,
		conversion.Source, conversion.Cached, conversion.RateDate)

	orderId := strconv.FormatUint(row.ID, 10)
	payRes, err := recharge.HaiPayApplyPayout(ctx, &recharge.HaiPayPayoutApplyReq{
		OrderID:       orderId,
		Currency:      currency,
		Amount:        localAmount,
		AccountType:   info.AccountType,
		BankCode:      info.BankCode,
		AccountNo:     info.AccountNo,
		Name:          info.PayeeName,
		Phone:         info.Phone,
		Email:         info.Email,
		Country:       info.CountryCode,
		IdentifyType:  info.IdentifyType,
		Address1:      info.Address1,
		Address2:      info.Address2,
		Address3:      info.Address3,
		PostalCode:    info.PostalCode,
		PartnerUserID: strconv.FormatUint(row.RoomId, 10),
		Body:          fmt.Sprintf("anchorSettlement:%d", row.ID),
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
	row.SetTransferPayout(orderId, platformNo, currency, localAmount)
	row.SetStatus(live.AnchorIncomeSettlementStatusTransferring)
	liveroomdao.PublishAnchorIncomeSettlementLog(row)
	return nil
}

func prepareTieredAnchorPayout(row *live.AnchorIncomeSettlementLog) error {
	if row == nil || row.SettlementRuleType != live.AnchorIncomeSettlementRuleTiered {
		return fmt.Errorf("platform anchor settlement rule invalid")
	}
	// 周结算已固定审核金额；仅为历史上金额为空的结算单补算。
	if row.SettlementReceivableUsd > 0 {
		return nil
	}
	exchangeCfg := wallet.GetExchangeCfgSnapshot()
	gameDiamond, totalDiamond, usd, err := calculateTieredAnchorPayout(row, exchangeCfg)
	if err != nil {
		return err
	}
	row.SetPayoutConversion(exchangeCfg.GoldToDiamondRate, exchangeCfg.UsdToGoldRate, gameDiamond, totalDiamond, usd)
	if !syndb.FlushUntilIdle(15*time.Second) || !liveroomdao.VerifyAnchorPayoutConversionPersisted(row) {
		return fmt.Errorf("payout conversion persist failed")
	}
	return nil
}

func calculateTieredAnchorPayout(row *live.AnchorIncomeSettlementLog, exchangeCfg wallet.ExchangeCfgSnapshot) (gameDiamond, totalDiamond, usd float64, err error) {
	if row == nil {
		return 0, 0, 0, fmt.Errorf("empty settlement")
	}
	if exchangeCfg.GoldToDiamondRate <= 0 || exchangeCfg.UsdToGoldRate <= 0 {
		return 0, 0, 0, fmt.Errorf("wallet exchange rate invalid")
	}
	baseDiamond := row.SettlementSalary + row.AnchorSocialShareAmount
	gameDiamond, totalDiamond, usd = wallet.CalcSettlementUsdWithSnapshot(baseDiamond, row.AnchorGameShareAmountGold, exchangeCfg)
	usd = math.Round(usd*10000) / 10000
	if usd <= 0 {
		return 0, 0, 0, fmt.Errorf("platform anchor settlement receivable usd <= 0")
	}
	return gameDiamond, totalDiamond, usd, nil
}
