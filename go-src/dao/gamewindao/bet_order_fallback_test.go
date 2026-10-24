package gamewindao

import (
	"strings"
	"testing"
)

// 与 winBetOrderFallbackSQL 同规则的 Go 版,仅用于单测核对推导逻辑.
func betOrderIDFallback(orderID string) string {
	orderID = strings.TrimSpace(orderID)
	if orderID == "" {
		return ""
	}
	first := strings.IndexByte(orderID, '-')
	if first <= 0 || first >= len(orderID)-1 {
		return ""
	}
	rest := orderID[first+1:]
	second := strings.IndexByte(rest, '-')
	if second <= 0 {
		return ""
	}
	a := orderID[:first]
	b := rest[:second]
	suffix := rest[second+1:]
	if a == "" || b == "" || a == b {
		return ""
	}
	return b + "-" + b + "-" + suffix
}

func TestBetOrderIDFallback(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in   string
		want string
	}{
		{
			in:   "7668196330310750209-7668196330310750208-04",
			want: "7668196330310750208-7668196330310750208-04",
		},
		{
			in:   "7666082560276971521-7666082560142753792-01",
			want: "7666082560142753792-7666082560142753792-01",
		},
		{in: "7659836252302102528-7659836252302102528-11", want: ""},
		{in: "", want: ""},
		{in: "plain", want: ""},
		{in: "only-one", want: ""},
	}
	for _, tc := range cases {
		if got := betOrderIDFallback(tc.in); got != tc.want {
			t.Fatalf("betOrderIDFallback(%q)=%q want %q", tc.in, got, tc.want)
		}
	}
}

func TestWinBetOrderFallbackSQLMentionsOrderId(t *testing.T) {
	t.Parallel()
	if !strings.Contains(winBetOrderFallbackSQL, "w.order_id") {
		t.Fatalf("winBetOrderFallbackSQL missing w.order_id: %s", winBetOrderFallbackSQL)
	}
}
