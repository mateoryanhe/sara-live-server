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
	"xr-game-server/dao/guilddao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/incomesettlementdto"
	"xr-game-server/entity/live"
	"xr-game-server/errercode"
	"xr-game-server/module/cmsvis"
	"xr-game-server/module/fxrate"
	"xr-game-server/module/recharge"
	"xr-game-server/module/wallet"
)

func parseIdList(ids []string) []uint64 {
	if len(ids) == 0 {
		return nil
	}
	out := make([]uint64, 0, len(ids))
	seen := make(map[uint64]struct{}, len(ids))
	for _, raw := range ids {
		id, err := strconv.ParseUint(raw, 10, 64)
		if err != nil || id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func guildVisibleSet(ctx context.Context) (map[uint64]struct{}, bool, bool) {
	guildIds, restrict, empty := cmsvis.VisibilityGuildFilter(ctx)
	if empty {
		return nil, true, true
	}
	if !restrict {
		return nil, false, false
	}
	set := make(map[uint64]struct{}, len(guildIds))
	for _, id := range guildIds {
		set[id] = struct{}{}
	}
	return set, true, false
}

// BatchApproveGuildSettlement 批量审核：审核中(0) -> 审核通过(1)
func BatchApproveGuildSettlement(ctx context.Context, req *incomesettlementdto.CMSBatchApproveGuildSettlementReq) (*incomesettlementdto.CMSBatchApproveGuildSettlementRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	ids := parseIdList(req.Ids)
	if len(ids) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	visibleSet, restrict, empty := guildVisibleSet(ctx)
	res := &incomesettlementdto.CMSBatchApproveGuildSettlementRes{}
	if empty {
		res.FailCount = len(ids)
		return res, nil
	}
	for _, id := range ids {
		row := liveroomdao.GetGuildIncomeSettlementLogById(id)
		if row == nil || row.Status != entity.GuildIncomeSettlementStatusPending {
			res.FailCount++
			continue
		}
		if restrict {
			if _, ok := visibleSet[row.GuildId]; !ok {
				res.FailCount++
				continue
			}
		}
		row.SetStatus(entity.GuildIncomeSettlementStatusApproved)
		liveroomdao.PublishGuildIncomeSettlementLog(row)
		res.SuccessCount++
	}
	return res, nil
}

func getGuildSettlementForStatus(ctx context.Context, rawId string, status uint8) (*entity.GuildIncomeSettlementLog, error) {
	id, err := strconv.ParseUint(strings.TrimSpace(rawId), 10, 64)
	if err != nil || id == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	visibleSet, restrict, empty := guildVisibleSet(ctx)
	if empty {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}
	row := liveroomdao.GetGuildIncomeSettlementLogById(id)
	if row == nil || row.Status != status {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if restrict {
		if _, ok := visibleSet[row.GuildId]; !ok {
			return nil, errercode.CreateCode(errercode.NoPermission)
		}
	}
	return row, nil
}

// ReopenGuildSettlementApproval 将审核通过(1)的结算退回审核中(0)。
func ReopenGuildSettlementApproval(ctx context.Context, req *incomesettlementdto.CMSReopenGuildSettlementApprovalReq) (*incomesettlementdto.CMSReopenGuildSettlementApprovalRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row, err := getGuildSettlementForStatus(ctx, req.Id, entity.GuildIncomeSettlementStatusApproved)
	if err != nil {
		return nil, err
	}
	row.SetStatus(entity.GuildIncomeSettlementStatusPending)
	row.SetTransferFailMsg("")
	liveroomdao.PublishGuildIncomeSettlementLog(row)
	xrlog.DetailLog.Infof(ctx, "guild settlement approval reopened settlementId=%d guildId=%d", row.ID, row.GuildId)
	return &incomesettlementdto.CMSReopenGuildSettlementApprovalRes{Success: true}, nil
}

// CopyGuildSettlementPayout 复制已成功的代付记录，生成新的审核中结算流水。
// 新记录使用新的主键，后续提交 HaiPay 时该主键即为新的唯一 orderId。
func CopyGuildSettlementPayout(ctx context.Context, req *incomesettlementdto.CMSCopyGuildSettlementPayoutReq) (*incomesettlementdto.CMSCopyGuildSettlementPayoutRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	source, err := getGuildSettlementForStatus(ctx, req.Id, entity.GuildIncomeSettlementStatusTransferred)
	if err != nil {
		return nil, err
	}
	cloned := entity.NewGuildIncomeSettlementLogWithBreakdown(
		source.GuildId,
		&source.LiveRoomIncomeAmounts,
		source.SettlementSalary,
		source.SettlementShareAmount,
		source.SettlementShareAmountUsd,
		source.SettlementReceivableUsd,
		source.GuildSharePercent,
		&entity.GuildIncomeSettlementBreakdown{
			SettlementRuleType:        source.SettlementRuleType,
			AnchorSocialShareAmount:   source.AnchorSocialShareAmount,
			GuildSocialShareAmount:    source.GuildSocialShareAmount,
			AnchorGameShareAmountGold: source.AnchorGameShareAmountGold,
			GuildGameShareAmountGold:  source.GuildGameShareAmountGold,
		},
	)
	if cloned == nil || cloned.ID == 0 {
		return nil, errercode.CreateCode(errercode.SysError)
	}
	cloned.SetPayoutConversion(
		source.GoldToDiamondRate,
		source.UsdToGoldRate,
		source.GameShareAmountDiamond,
		source.TotalSettlementDiamond,
		source.SettlementReceivableUsd,
	)
	liveroomdao.PublishGuildIncomeSettlementLog(cloned)
	xrlog.DetailLog.Infof(ctx, "guild settlement payout copied sourceSettlementId=%d newSettlementId=%d guildId=%d receivableUsd=%.4f",
		source.ID, cloned.ID, cloned.GuildId, cloned.SettlementReceivableUsd)
	return &incomesettlementdto.CMSCopyGuildSettlementPayoutRes{Id: strconv.FormatUint(cloned.ID, 10)}, nil
}

// UpdateGuildSettlementReceivableUsd 修改审核中(0)结算的最终可收金额。
func UpdateGuildSettlementReceivableUsd(ctx context.Context, req *incomesettlementdto.CMSUpdateGuildSettlementReceivableUsdReq) (*incomesettlementdto.CMSUpdateGuildSettlementReceivableUsdRes, error) {
	if req == nil || math.IsNaN(req.SettlementReceivableUsd) || math.IsInf(req.SettlementReceivableUsd, 0) || req.SettlementReceivableUsd <= 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row, err := getGuildSettlementForStatus(ctx, req.Id, entity.GuildIncomeSettlementStatusPending)
	if err != nil {
		return nil, err
	}
	amount := math.Round(req.SettlementReceivableUsd*10000) / 10000
	if amount <= 0 || amount > 999999999999.9999 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	oldAmount := row.SettlementReceivableUsd
	row.SetSettlementReceivableUsd(amount)
	liveroomdao.PublishGuildIncomeSettlementLog(row)
	xrlog.DetailLog.Infof(ctx, "guild settlement receivable updated settlementId=%d guildId=%d oldUsd=%.4f newUsd=%.4f",
		row.ID, row.GuildId, oldAmount, amount)
	return &incomesettlementdto.CMSUpdateGuildSettlementReceivableUsdRes{Success: true}, nil
}

// BatchTransferGuildSettlement 批量代付:审核通过(1) -> 提交 HaiPay 代付(3),回调成功后再变(2)
func BatchTransferGuildSettlement(ctx context.Context, req *incomesettlementdto.CMSBatchTransferGuildSettlementReq) (*incomesettlementdto.CMSBatchTransferGuildSettlementRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	ids := parseIdList(req.Ids)
	if len(ids) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	visibleSet, restrict, empty := guildVisibleSet(ctx)
	res := &incomesettlementdto.CMSBatchTransferGuildSettlementRes{}
	if empty {
		res.FailCount = len(ids)
		res.Message = "无可见工会权限"
		return res, nil
	}

	for _, id := range ids {
		row := liveroomdao.GetGuildIncomeSettlementLogById(id)
		if row == nil || row.Status != entity.GuildIncomeSettlementStatusApproved {
			res.FailCount++
			continue
		}
		if restrict {
			if _, ok := visibleSet[row.GuildId]; !ok {
				res.FailCount++
				continue
			}
		}
		if err := submitGuildHaiPayPayout(ctx, row); err != nil {
			xrlog.DetailLog.Warningf(ctx, "guild payout submit failed settlementId=%d guildId=%d err=%v",
				row.ID, row.GuildId, err)
			row.SetTransferFailMsg(err.Error())
			liveroomdao.PublishGuildIncomeSettlementLog(row)
			res.FailCount++
			continue
		}
		res.SuccessCount++
	}
	if res.SuccessCount == 0 && res.FailCount > 0 {
		res.Message = "代付提交失败,请检查HaiPay全局配置、工会转账信息或汇率接口"
	} else if res.FailCount > 0 {
		res.Message = fmt.Sprintf("部分成功:成功%d 失败%d", res.SuccessCount, res.FailCount)
	} else {
		res.Message = fmt.Sprintf("已提交代付%d笔,等待HaiPay回调确认", res.SuccessCount)
	}
	return res, nil
}

func submitGuildHaiPayPayout(ctx context.Context, row *entity.GuildIncomeSettlementLog) error {
	if row == nil || row.ID == 0 {
		return fmt.Errorf("empty settlement")
	}
	if row.SettlementRuleType == entity.GuildIncomeSettlementRuleTiered {
		if err := prepareTieredGuildPayout(row); err != nil {
			return err
		}
	}
	if row.SettlementReceivableUsd <= 0 {
		return fmt.Errorf("receivable usd <= 0")
	}
	info := guilddao.GetGuildTransferInfo(row.GuildId)
	if info == nil {
		return fmt.Errorf("guild transfer info missing")
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
		"guild payout fx converted settlementId=%d guildId=%d usdAmount=%.4f currency=%s rate=%v localAmount=%v source=%s cached=%t rateDate=%s",
		row.ID, row.GuildId, row.SettlementReceivableUsd, currency, conversion.Rate, localAmount,
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
		PartnerUserID: strconv.FormatUint(row.GuildId, 10),
		Body:          fmt.Sprintf("guildSettlement:%d", row.ID),
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
	row.SetStatus(entity.GuildIncomeSettlementStatusTransferring)
	liveroomdao.PublishGuildIncomeSettlementLog(row)
	return nil
}

func prepareTieredGuildPayout(row *entity.GuildIncomeSettlementLog) error {
	if row == nil {
		return fmt.Errorf("empty settlement")
	}
	// 周结算已固定审核金额；人工修改后的审核金额也必须原样用于代付。
	if row.SettlementReceivableUsd > 0 {
		return nil
	}
	exchangeCfg := wallet.GetExchangeCfgSnapshot()
	gameDiamond, totalDiamond, usd, err := calculateTieredGuildPayout(row, exchangeCfg)
	if err != nil {
		return err
	}
	row.SetPayoutConversion(exchangeCfg.GoldToDiamondRate, exchangeCfg.UsdToGoldRate, gameDiamond, totalDiamond, usd)
	if !syndb.FlushUntilIdle(15*time.Second) || !liveroomdao.VerifyGuildPayoutConversionPersisted(row) {
		return fmt.Errorf("payout conversion persist failed")
	}
	return nil
}

func calculateTieredGuildPayout(row *entity.GuildIncomeSettlementLog, exchangeCfg wallet.ExchangeCfgSnapshot) (gameDiamond, totalDiamond, usd float64, err error) {
	if row == nil {
		return 0, 0, 0, fmt.Errorf("empty settlement")
	}
	if exchangeCfg.GoldToDiamondRate <= 0 || exchangeCfg.UsdToGoldRate <= 0 {
		return 0, 0, 0, fmt.Errorf("wallet exchange rate invalid")
	}
	gameGold := row.AnchorGameShareAmountGold + row.GuildGameShareAmountGold
	baseDiamond := row.SettlementSalary + row.AnchorSocialShareAmount + row.GuildSocialShareAmount
	gameDiamond, totalDiamond, usd = wallet.CalcSettlementUsdWithSnapshot(baseDiamond, gameGold, exchangeCfg)
	usd = math.Round(usd*10000) / 10000
	if usd <= 0 {
		return 0, 0, 0, fmt.Errorf("tiered settlement receivable usd <= 0")
	}
	return gameDiamond, totalDiamond, usd, nil
}
