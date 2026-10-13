package userwithdraw

import (
	"strings"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/xrlog"
	"xr-game-server/dao/statdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dao/userwithdrawdao"
	userentity "xr-game-server/entity/user"
	"xr-game-server/module/recharge"
)

func handleUserWithdrawPayoutNotify(r *ghttp.Request, logID uint64, orderId string) {
	ctx := r.Context()
	row := userwithdrawdao.GetUserWithdrawLogById(logID)
	if row == nil {
		xrlog.DetailLog.Warningf(ctx, "haipay user withdraw notify log not found orderId=%s", orderId)
		r.Response.WriteStatus(200)
		r.Response.Write([]byte("OK"))
		return
	}
	if row.Status == userentity.UserWithdrawStatusTransferred {
		r.Response.WriteStatus(200)
		r.Response.Write([]byte("OK"))
		return
	}
	query, err := recharge.HaiPayQueryPayoutOrder(ctx, row.TransferCurrency, row.TransferOrderId, row.TransferPlatformNo, row.TransferLocalAmount)
	if err != nil {
		xrlog.DetailLog.Warningf(ctx, "haipay user withdraw notify query failed orderId=%s err=%v", orderId, err)
		r.Response.WriteStatus(500)
		r.Response.Write([]byte("query fail"))
		return
	}

	switch query.Status {
	case 2:
		now := time.Now()
		if orderNo := strings.TrimSpace(query.OrderNo); orderNo != "" {
			row.SetTransferPayout(row.TransferOrderId, orderNo, row.TransferCurrency, row.TransferLocalAmount)
		}
		row.SetStatus(userentity.UserWithdrawStatusTransferred)
		row.SetTransferAt(&now)
		row.SetTransferFailMsg("")
		if err := userwithdrawdao.UpdateUserWithdrawLog(row); err != nil {
			xrlog.DetailLog.Warningf(ctx, "haipay user withdraw notify persist failed orderId=%s err=%v", orderId, err)
			r.Response.WriteStatus(500)
			r.Response.Write([]byte("persist fail"))
			return
		}
		recordWithdrawSucceeded(row.UserId, row.WithdrawUsd)
	case 3:
		msg := strings.TrimSpace(query.ErrorMsg)
		if msg == "" {
			msg = "payout failed"
		}
		row.SetStatus(userentity.UserWithdrawStatusApproved)
		row.SetTransferFailMsg(msg)
		_ = userwithdrawdao.UpdateUserWithdrawLog(row)
	default:
		xrlog.DetailLog.Infof(ctx, "haipay user withdraw notify query pending status=%d orderId=%s", query.Status, orderId)
	}

	r.Response.WriteStatus(200)
	r.Response.Write([]byte("OK"))
}

func recordWithdrawSucceeded(userId uint64, usd float64) {
	if usd <= 0 || userId == 0 {
		return
	}
	if total := statdao.GetSysStat(); total != nil {
		total.AddTotalWithdraw(usd)
	}
	stat := userinfodao.GetUserCumulativeStatByUserId(userId)
	if stat != nil {
		stat.AddTotalWithdraw(usd)
		userinfodao.PublishUserCumulativeStat(stat)
	}
}
