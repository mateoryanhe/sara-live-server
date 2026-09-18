package wallet

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/currencylogdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/golddto"
	userentity "xr-game-server/entity/user"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

const (
	defaultCoinMerchantTransferRecordPageSize = 20
	maxCoinMerchantTransferRecordPageSize     = 100
	coinMerchantTransferRecordTimeLayout      = "2006-01-02 15:04:05"
)

func normalizeCoinMerchantTransferRecordPage(pageIndex, pageSize int) (int, int) {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = defaultCoinMerchantTransferRecordPageSize
	}
	if pageSize > maxCoinMerchantTransferRecordPageSize {
		pageSize = maxCoinMerchantTransferRecordPageSize
	}
	return pageIndex, pageSize
}

func parseCoinMerchantTransferTargetUserId(value string) (uint64, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, nil
	}
	targetUserId, err := strconv.ParseUint(value, 10, 64)
	if err != nil || targetUserId == 0 {
		return 0, errercode.CreateCode(errercode.InvalidParam)
	}
	return targetUserId, nil
}

func validateCoinMerchantTransferRecordTimeRange(startTime, endTime int64) error {
	if startTime < 0 || endTime < 0 || (startTime > 0 && endTime > 0 && startTime > endTime) {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	return nil
}

func buildCoinMerchantTransferRecordItems(
	rows []*currencylogdao.CoinMerchantGoldTransferListRow,
) []*golddto.AppCoinMerchantGoldTransferRecordItem {
	items := make([]*golddto.AppCoinMerchantGoldTransferRecordItem, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		item := &golddto.AppCoinMerchantGoldTransferRecordItem{
			Id:             strconv.FormatUint(row.ID, 10),
			TargetUserId:   strconv.FormatUint(row.TargetUserId, 10),
			TargetNickname: row.TargetNickname,
			TargetAvatar:   upload.ResolveAvatarUrlForUser(row.TargetUserId, row.TargetAvatar),
			Amount:         row.Amount,
		}
		if !row.CreatedAt.IsZero() {
			item.CreatedAt = row.CreatedAt.Unix()
			item.CreatedAtText = row.CreatedAt.In(time.Local).Format(coinMerchantTransferRecordTimeLayout)
		}
		items = append(items, item)
	}
	return items
}

// GetAppCoinMerchantGoldTransferRecordList 分页查询当前登录币商自己的金币转账记录。
func GetAppCoinMerchantGoldTransferRecordList(
	ctx context.Context,
	req *golddto.AppCoinMerchantGoldTransferRecordListReq,
) (*golddto.AppCoinMerchantGoldTransferRecordListRes, error) {
	merchantUserId := httpserver.GetAuthId(ctx)
	if merchantUserId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	merchant := userinfodao.GetUserInfoByUserId(merchantUserId)
	if merchant == nil || merchant.UserType != userentity.UserTypeCoinMerchant {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}

	targetUserId, err := parseCoinMerchantTransferTargetUserId(req.TargetUserId)
	if err != nil {
		return nil, err
	}
	if err = validateCoinMerchantTransferRecordTimeRange(req.StartTime, req.EndTime); err != nil {
		return nil, err
	}
	pageIndex, pageSize := normalizeCoinMerchantTransferRecordPage(req.PageIndex, req.PageSize)
	total, rows, err := currencylogdao.ListCoinMerchantGoldTransferLogs(&currencylogdao.CoinMerchantGoldTransferListFilter{
		MerchantUserId: merchantUserId,
		TargetUserId:   targetUserId,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		PageIndex:      pageIndex,
		PageSize:       pageSize,
	})
	if err != nil {
		return nil, err
	}
	return &golddto.AppCoinMerchantGoldTransferRecordListRes{
		Total:     total,
		PageIndex: pageIndex,
		PageSize:  pageSize,
		List:      buildCoinMerchantTransferRecordItems(rows),
	}, nil
}
