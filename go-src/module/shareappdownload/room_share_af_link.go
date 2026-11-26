package shareappdownload

import (
	"net/url"
	"strconv"
	"strings"
)

func buildLiveRoomAfLink(base, inviteCode string, roomId uint64) string {
	base = strings.TrimSpace(base)
	if base == "" || roomId == 0 {
		return base
	}
	sep := "?"
	if strings.Contains(base, "?") {
		sep = "&"
	}
	q := url.Values{}
	if inviteCode = strings.TrimSpace(inviteCode); inviteCode != "" {
		q.Set("inviteCode", inviteCode)
	}
	q.Set("roomId", strconv.FormatUint(roomId, 10))
	return base + sep + q.Encode()
}