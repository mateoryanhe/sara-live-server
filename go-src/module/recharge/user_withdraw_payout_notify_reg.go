package recharge

import "github.com/gogf/gf/v2/net/ghttp"

type UserWithdrawPayoutNotifyHandler func(r *ghttp.Request, logID uint64, orderId string)

var userWithdrawPayoutNotifyHandler UserWithdrawPayoutNotifyHandler

func RegisterUserWithdrawPayoutNotifyHandler(fn UserWithdrawPayoutNotifyHandler) {
	userWithdrawPayoutNotifyHandler = fn
}

func dispatchUserWithdrawPayoutNotify(r *ghttp.Request, logID uint64, orderId string) bool {
	if userWithdrawPayoutNotifyHandler == nil || logID == 0 {
		return false
	}
	userWithdrawPayoutNotifyHandler(r, logID, orderId)
	return true
}
