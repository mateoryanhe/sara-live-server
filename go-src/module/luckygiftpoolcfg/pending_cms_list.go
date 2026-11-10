package luckygiftpoolcfg

import (
	"context"
	"strconv"
	"time"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/luckygiftpoolcfgdto"
	liveentity "xr-game-server/entity/live"
)

func ListLuckyGiftPoolPendingBalance(_ context.Context, req *luckygiftpoolcfgdto.ListLuckyGiftPoolPendingBalanceReq) (*httpserver.CMSQueryResp, error) {
	if req == nil {
		req = &luckygiftpoolcfgdto.ListLuckyGiftPoolPendingBalanceReq{}
	}
	total, rows := cfgdao.ListLuckyGiftPoolPendingBalances(&cfgdao.LuckyGiftPoolPendingBalanceListFilter{
		Status:    req.Status,
		PageIndex: req.PageIndex,
		PageSize:  req.PageSize,
	})
	list := make([]*luckygiftpoolcfgdto.LuckyGiftPoolPendingBalanceItem, 0, len(rows))
	for _, row := range rows {
		if item := toPendingItem(row); item != nil {
			list = append(list, item)
		}
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}

func toPendingItem(row *liveentity.LuckyGiftPoolPendingBalance) *luckygiftpoolcfgdto.LuckyGiftPoolPendingBalanceItem {
	if row == nil {
		return nil
	}
	return &luckygiftpoolcfgdto.LuckyGiftPoolPendingBalanceItem{
		ID:        strconv.FormatUint(row.ID, 10),
		Amount:    row.Amount,
		Status:    row.Status,
		CreatedAt: formatPendingTime(row.CreatedAt),
		UpdatedAt: formatPendingTime(row.UpdatedAt),
	}
}

func formatPendingTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
