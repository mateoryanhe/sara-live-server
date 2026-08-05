package recharge

import (
	"context"
	"sync"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gmlock"
	"github.com/gogf/gf/v2/os/gtimer"
	"xr-game-server/core/xrtimer"
	"xr-game-server/dao/rechargeorderdao"
	"xr-game-server/entity"
)

const rechargeOrderPayTimeout = 30 * time.Minute

var rechargeOrderTimeoutEntries sync.Map // orderId uint64 -> *gtimer.Entry

func initRechargeOrderTimeoutWatch() {
	restorePendingRechargeOrderTimeouts()
}

func restorePendingRechargeOrderTimeouts() {
	orders := rechargeorderdao.LoadPendingOrders()
	if len(orders) == 0 {
		return
	}
	//ctx := gctx.New()
	for _, order := range orders {
		if order == nil || order.ID == 0 {
			continue
		}
		rechargeorderdao.AddOrderToCache(order)
		ScheduleRechargeOrderTimeout(order.ID, order.CreatedAt)
	}
	//g.Log().Infof(ctx, "restore recharge order timeout watch count=%d", len(orders))
}

// ScheduleRechargeOrderTimeout 将待支付订单加入超时检查队列
func ScheduleRechargeOrderTimeout(orderId uint64, createdAt time.Time) {
	if orderId == 0 || createdAt.IsZero() {
		return
	}
	CancelRechargeOrderTimeout(orderId)

	expireAt := createdAt.Add(rechargeOrderPayTimeout)
	delay := time.Until(expireAt)
	if delay <= 0 {
		delay = time.Millisecond
	}
	entry := xrtimer.AddOnce(gctx.New(), delay, func(ctx context.Context) {
		cancelRechargeOrderOnTimeout(ctx, orderId)
	})
	rechargeOrderTimeoutEntries.Store(orderId, entry)
}

// CancelRechargeOrderTimeout 订单已完成/已取消时移除超时检查
func CancelRechargeOrderTimeout(orderId uint64) {
	if orderId == 0 {
		return
	}
	if v, ok := rechargeOrderTimeoutEntries.LoadAndDelete(orderId); ok {
		if entry, ok := v.(*gtimer.Entry); ok && entry != nil {
			entry.Close()
		}
	}
}

func cancelRechargeOrderOnTimeout(ctx context.Context, orderId uint64) {
	rechargeOrderTimeoutEntries.Delete(orderId)

	lockKey := rechargeOrderLockKey(orderId)
	gmlock.Lock(lockKey)
	defer gmlock.Unlock(lockKey)

	order := rechargeorderdao.GetById(orderId)
	if order == nil || order.Status != entity.RechargeOrderStatusPending {
		return
	}
	now := time.Now()
	order.SetStatus(entity.RechargeOrderStatusCancelled)
	order.SetUpdatedAt(now)
	rechargeorderdao.FlushOrderCache(order)
	//g.Log().Infof(ctx, "recharge order pay timeout cancelled orderId=%d userId=%d", orderId, order.UserId)
}
