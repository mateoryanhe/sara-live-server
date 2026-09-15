package incomesettlement

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"xr-game-server/core/xrlog"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dao/guilddao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/incomesettlementdto"
	"xr-game-server/entity/live"
	"xr-game-server/errercode"
	"xr-game-server/module/cmsvis"
	"xr-game-server/module/fxrate"
	"xr-game-server/module/recharge"
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

// BatchApproveGuildSettlement 批量审核：未审核(0) -> 审核通过(1)
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
		res.SuccessCount++
	}
	return res, nil
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
	cfg := cfgdao.GetHaiPayCfgCached()
	if cfg == nil || !cfg.PayoutEnabled {
		return &incomesettlementdto.CMSBatchTransferGuildSettlementRes{
			Message: "HaiPay代付未启用,请先在CMS配置 payoutEnabled / payoutAppIds",
		}, nil
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
			res.FailCount++
			continue
		}
		res.SuccessCount++
	}
	if res.SuccessCount == 0 && res.FailCount > 0 {
		res.Message = "代付提交失败,请检查工会转账信息/汇率接口/代付appId"
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
	localAmount := conversion.TargetAmount
	if localAmount <= 0 {
		return fmt.Errorf("local amount <= 0")
	}

	orderId := fmt.Sprintf("gis%d_%d", row.ID, time.Now().Unix())
	if len(orderId) > 48 {
		orderId = orderId[:48]
	}
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
		PartnerUserID: strconv.FormatUint(row.GuildId, 10),
		Subject:       "",
		Body:          fmt.Sprintf("guildSettlement:%d", row.ID),
	})
	if err != nil {
		return err
	}
	platformNo := ""
	if payRes != nil {
		platformNo = payRes.OrderNo
		if payRes.OrderID != "" {
			orderId = payRes.OrderID
		}
	}
	row.SetTransferPayout(orderId, platformNo, currency, localAmount)
	row.SetStatus(entity.GuildIncomeSettlementStatusTransferring)
	return nil
}
