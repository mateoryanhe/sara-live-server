package userwithdraw

import (
	"strconv"
	"strings"
)

const payoutOrderIDPrefix = "U"

func FormatPayoutOrderID(logID uint64) string {
	return payoutOrderIDPrefix + strconv.FormatUint(logID, 10)
}

func ParsePayoutOrderID(orderID string) (uint64, bool) {
	orderID = strings.TrimSpace(orderID)
	if !strings.HasPrefix(orderID, payoutOrderIDPrefix) {
		return 0, false
	}
	id, err := strconv.ParseUint(strings.TrimPrefix(orderID, payoutOrderIDPrefix), 10, 64)
	if err != nil || id == 0 {
		return 0, false
	}
	return id, true
}
