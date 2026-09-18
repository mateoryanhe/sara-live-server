package wallet

import (
	"context"
	"strconv"
	"strings"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/currencylogdao"
	"xr-game-server/dto/golddto"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

func parseOptionalCoinMerchantUserId(value string) (uint64, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, nil
	}
	userId, err := strconv.ParseUint(value, 10, 64)
	if err != nil || userId == 0 {
		return 0, errercode.CreateCode(errercode.InvalidParam)
	}
	return userId, nil
}

// GetCMSCoinMerchantTransferLogList 直接联表查询币商金币转账流水，不读取用户缓存。
func GetCMSCoinMerchantTransferLogList(
	_ context.Context,
	req *golddto.CMSCoinMerchantTransferLogListReq,
) (*httpserver.CMSQueryResp, error) {
	merchantUserId, err := parseOptionalCoinMerchantUserId(req.CoinMerchantUserId)
	if err != nil {
		return nil, err
	}
	if err = validateCoinMerchantTransferRecordTimeRange(req.StartTime, req.EndTime); err != nil {
		return nil, err
	}
	pageIndex, pageSize := normalizeCoinMerchantTransferRecordPage(req.PageIndex, req.PageSize)
	total, rows, err := currencylogdao.ListCoinMerchantGoldTransferLogs(&currencylogdao.CoinMerchantGoldTransferListFilter{
		MerchantUserId: merchantUserId,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		PageIndex:      pageIndex,
		PageSize:       pageSize,
	})
	if err != nil {
		return nil, err
	}

	list := make([]*golddto.CMSCoinMerchantTransferLogItem, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		item := &golddto.CMSCoinMerchantTransferLogItem{
			Id:                   row.ID,
			CoinMerchantUserId:   row.CoinMerchantUserId,
			CoinMerchantNickname: row.CoinMerchantNickname,
			CoinMerchantAvatar:   upload.ResolveAvatarUrlForUser(row.CoinMerchantUserId, row.CoinMerchantAvatar),
			TargetUserId:         row.TargetUserId,
			TargetNickname:       row.TargetNickname,
			TargetAvatar:         upload.ResolveAvatarUrlForUser(row.TargetUserId, row.TargetAvatar),
			Amount:               row.Amount,
			SenderGoldBefore:     row.SenderGoldBefore,
			SenderGoldAfter:      row.SenderGoldAfter,
			TargetGoldBefore:     row.TargetGoldBefore,
			TargetGoldAfter:      row.TargetGoldAfter,
		}
		if !row.CreatedAt.IsZero() {
			item.CreatedAt = &row.CreatedAt
		}
		list = append(list, item)
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}
