package country

import "testing"

func TestHaiPayPayoutCountryCurrency(t *testing.T) {
	if got := HaiPayPayoutCurrency("id"); got != "IDR" {
		t.Fatalf("ID currency=%q", got)
	}
	if got := HaiPayPayoutCountryFromCurrency("idr"); got != "ID" {
		t.Fatalf("IDR country=%q", got)
	}
	if IsHaiPayPayoutCountry("US") {
		t.Fatal("US should not be payout country by default")
	}
	if n := len(ListHaiPayPayoutCountries()); n == 0 {
		t.Fatal("empty payout countries")
	}
}

func TestGetAndFlagIcon(t *testing.T) {
	c, ok := Get("id")
	if !ok {
		t.Fatal("expected ID")
	}
	if c.Code != "ID" || c.NameEn != "Indonesia" || c.NameZh == "" || c.FlagIcon != "id.png" {
		t.Fatalf("unexpected country: %+v", c)
	}
	ci, ok := Get("CI")
	if !ok || ci.NameZh != "科特迪瓦" {
		t.Fatalf("CI zh=%q", ci.NameZh)
	}
	if RelPath("ID", "20260102150405") != "country-flags/20260102150405/id.png" {
		t.Fatalf("rel=%s", RelPath("ID", "20260102150405"))
	}
	if !Exists("EU") {
		t.Fatal("expected EU")
	}
	if Exists("ZZ") {
		t.Fatal("ZZ should not exist")
	}
	if n := len(All()); n != len(all) {
		t.Fatalf("All len=%d want %d", n, len(all))
	}
}
