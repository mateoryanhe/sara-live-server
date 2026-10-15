package userwithdraw

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"xr-game-server/dao/userinfodao"
	"xr-game-server/dao/userwithdrawdao"
	"xr-game-server/dto/userwithdrawdto"
	userentity "xr-game-server/entity/user"
	"xr-game-server/errercode"
)

func CMSListWithdrawLogs(_ context.Context, req *userwithdrawdto.CMSListWithdrawLogsReq) (*userwithdrawdto.CMSListWithdrawLogsRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	var userId uint64
	if raw := strings.TrimSpace(req.UserId); raw != "" {
		id, err := strconv.ParseUint(raw, 10, 64)
		if err != nil || id == 0 {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		userId = id
	}
	total, rows, err := userwithdrawdao.ListUserWithdrawLogs(&userwithdrawdao.UserWithdrawListFilter{
		UserId:   userId,
		Status:   req.Status,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	list := make([]*userwithdrawdto.CMSWithdrawLogItem, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		nickname := ""
		if user := userinfodao.GetUserInfoByUserId(row.UserId); user != nil {
			nickname = user.Nickname
		}
		list = append(list, &userwithdrawdto.CMSWithdrawLogItem{
			Id:                  strconv.FormatUint(row.ID, 10),
			UserId:              strconv.FormatUint(row.UserId, 10),
			Nickname:            nickname,
			DiamondAmount:       withdrawLogDiamondDisplay(row),
			GoldAmount:          row.GoldAmount,
			WithdrawUsd:         row.WithdrawUsd,
			Status:              row.Status,
			TransferOrderId:     row.TransferOrderId,
			TransferPlatformNo:  row.TransferPlatformNo,
			TransferCurrency:    row.TransferCurrency,
			TransferLocalAmount: row.TransferLocalAmount,
			TransferFailMsg:     row.TransferFailMsg,
			RejectReason:        row.RejectReason,
			CreatedAt:           formatLogTime(row.CreatedAt),
			TransferAt:          formatLogTimePtr(row.TransferAt),
		})
	}
	return &userwithdrawdto.CMSListWithdrawLogsRes{Total: total, List: list}, nil
}

func CMSBatchApproveWithdraw(_ context.Context, req *userwithdrawdto.CMSBatchApproveWithdrawReq) (*userwithdrawdto.CMSBatchApproveWithdrawRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	res := &userwithdrawdto.CMSBatchApproveWithdrawRes{}
	for _, id := range parseIdList(req.Ids) {
		row := userwithdrawdao.GetUserWithdrawLogById(id)
		if row == nil || row.Status != userentity.UserWithdrawStatusPending {
			res.FailCount++
			continue
		}
		row.SetStatus(userentity.UserWithdrawStatusApproved)
		userwithdrawdao.SyncUserWithdrawLogCache(row)
		res.SuccessCount++
	}
	return res, nil
}

func CMSBatchRejectWithdraw(_ context.Context, req *userwithdrawdto.CMSBatchRejectWithdrawReq) (*userwithdrawdto.CMSBatchRejectWithdrawRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	reason := strings.TrimSpace(req.Reason)
	res := &userwithdrawdto.CMSBatchRejectWithdrawRes{}
	for _, id := range parseIdList(req.Ids) {
		row := userwithdrawdao.GetUserWithdrawLogById(id)
		if row == nil || row.Status != userentity.UserWithdrawStatusPending {
			res.FailCount++
			continue
		}
		if err := refundUserWithdraw(row.UserId, row); err != nil {
			res.FailCount++
			continue
		}
		row.SetStatus(userentity.UserWithdrawStatusRejected)
		row.SetRejectReason(reason)
		userwithdrawdao.SyncUserWithdrawLogCache(row)
		res.SuccessCount++
	}
	return res, nil
}

func CMSBatchSubmitPayout(ctx context.Context, req *userwithdrawdto.CMSBatchSubmitPayoutReq) (*userwithdrawdto.CMSBatchSubmitPayoutRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	res := &userwithdrawdto.CMSBatchSubmitPayoutRes{}
	for _, id := range parseIdList(req.Ids) {
		row := userwithdrawdao.GetUserWithdrawLogById(id)
		if row == nil || row.Status != userentity.UserWithdrawStatusApproved {
			res.FailCount++
			continue
		}
		if err := submitUserHaiPayPayout(ctx, row); err != nil {
			row.SetTransferFailMsg(strings.TrimSpace(err.Error()))
			userwithdrawdao.SyncUserWithdrawLogCache(row)
			res.FailCount++
			continue
		}
		res.SuccessCount++
	}
	if res.SuccessCount == 0 && res.FailCount > 0 {
		res.Message = "代付提交失败,请检查HaiPay配置、用户收款信息或汇率"
	} else if res.FailCount > 0 {
		res.Message = fmt.Sprintf("部分成功:成功%d 失败%d", res.SuccessCount, res.FailCount)
	} else {
		res.Message = fmt.Sprintf("已提交代付%d笔,等待HaiPay回调确认", res.SuccessCount)
	}
	return res, nil
}
