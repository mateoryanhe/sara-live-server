package recharge

import "testing"

func TestHaiPayGlobalCashierAmountIndonesiaIDR(t *testing.T) {
	if got := haiPayFormatGlobalCashierAmount("ID", "IDR", 17717.04); got != "17717" {
		t.Fatalf("format ID+IDR got=%q want=17717", got)
	}
	if got := haiPayRoundGlobalCashierAmount("ID", "IDR", 17717.04); got != 17717 {
		t.Fatalf("round ID+IDR got=%v want=17717", got)
	}
}

func TestHaiPayGlobalCashierAmountIndonesiaUSD(t *testing.T) {
	if got := haiPayFormatGlobalCashierAmount("ID", "USD", 9.99); got != "9.99" {
		t.Fatalf("format ID+USD got=%q want=9.99", got)
	}
}

func TestHaiPayGlobalCashierAmountOtherRegionIDR(t *testing.T) {
	// 全球收银台仅印尼区 IDR 取整；其他地区仍两位小数（若未来目录扩展 IDR）。
	if got := haiPayFormatGlobalCashierAmount("SG", "IDR", 17717.04); got != "17717.04" {
		t.Fatalf("format SG+IDR got=%q want=17717.04", got)
	}
}

func TestHaiPayGlobalCashierAmountDefault(t *testing.T) {
	if got := haiPayFormatGlobalCashierAmount("US", "USD", 10.005); got != "10.01" {
		t.Fatalf("format US+USD got=%q want=10.01", got)
	}
}
