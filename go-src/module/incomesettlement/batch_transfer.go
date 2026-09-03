package incomesettlement

import (
	"context"
	"strconv"

	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/incomesettlementdto"
	"xr-game-server/entity/live"
	"xr-game-server/errercode"
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

// BatchApproveGuildSettlement 批量审核：未审核(0) -> 审核通过(1)
func BatchApproveGuildSettlement(_ context.Context, req *incomesettlementdto.CMSBatchApproveGuildSettlementReq) (*incomesettlementdto.CMSBatchApproveGuildSettlementRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	ids := parseIdList(req.Ids)
	if len(ids) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	res := &incomesettlementdto.CMSBatchApproveGuildSettlementRes{}
	for _, id := range ids {
		row := liveroomdao.GetGuildIncomeSettlementLogById(id)
		if row == nil || row.Status != entity.GuildIncomeSettlementStatusPending {
			res.FailCount++
			continue
		}
		row.SetStatus(entity.GuildIncomeSettlementStatusApproved)
		res.SuccessCount++
	}
	return res, nil
}

// BatchTransferGuildSettlement 批量转账预留(转账API未接入)
func BatchTransferGuildSettlement(_ context.Context, req *incomesettlementdto.CMSBatchTransferGuildSettlementReq) (*incomesettlementdto.CMSBatchTransferGuildSettlementRes, error) {
	if req == nil || len(parseIdList(req.Ids)) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	return &incomesettlementdto.CMSBatchTransferGuildSettlementRes{
		Reserved: true,
		Message:  "转账功能预留中，转账API尚未接入",
	}, nil
}
