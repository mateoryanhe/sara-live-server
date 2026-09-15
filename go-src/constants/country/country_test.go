package country

import (
	"reflect"
	"testing"
)

func TestHaiPayCollectionCatalog(t *testing.T) {
	if got := HaiPayCollectionCurrency("id"); got != "USD" {
		t.Fatalf("ID collection currency=%q", got)
	}
	if got := HaiPayCollectionCurrency("hk"); got != "USD" {
		t.Fatalf("HK collection currency=%q", got)
	}
	if got := HaiPayCollectionCurrency("in"); got != "INR" {
		t.Fatalf("IN collection currency=%q", got)
	}
	if got := HaiPayCollectionCountryFromCurrency("eur"); got != "EU" {
		t.Fatalf("EUR collection country=%q", got)
	}
	if got := HaiPayCollectionCurrency("pl"); got != "EUR" {
		t.Fatalf("PL collection currency=%q", got)
	}
	if got := HaiPayCollectionCountryFromCurrency("pln"); got != "" {
		t.Fatalf("PLN must not be a collection currency, country=%q", got)
	}
	if got := HaiPayCollectionCurrencies("hk"); len(got) != 2 || got[0] != "HKD" || got[1] != "USD" {
		t.Fatalf("HK collection currencies=%v", got)
	}
	options := ListHaiPayCollectionOptions()
	want := []HaiPayCollectionOption{
		{CountryCode: "AE", Currencies: []string{"AED"}},
		{CountryCode: "AT", Currencies: []string{"EUR"}},
		{CountryCode: "BE", Currencies: []string{"EUR"}},
		{CountryCode: "BH", Currencies: []string{"BHD"}},
		{CountryCode: "BR", Currencies: []string{"BRL"}},
		{CountryCode: "EG", Currencies: []string{"EGP"}},
		{CountryCode: "GB", Currencies: []string{"GBP"}},
		{CountryCode: "HK", Currencies: []string{"HKD", "USD"}},
		{CountryCode: "ID", Currencies: []string{"IDR", "USD"}},
		{CountryCode: "IN", Currencies: []string{"INR"}},
		{CountryCode: "JP", Currencies: []string{"JPY", "USD"}},
		{CountryCode: "KR", Currencies: []string{"KRW", "USD"}},
		{CountryCode: "KW", Currencies: []string{"KWD"}},
		{CountryCode: "MY", Currencies: []string{"MYR", "USD"}},
		{CountryCode: "NL", Currencies: []string{"EUR"}},
		{CountryCode: "OM", Currencies: []string{"OMR"}},
		{CountryCode: "PH", Currencies: []string{"PHP"}},
		{CountryCode: "PK", Currencies: []string{"PKR"}},
		{CountryCode: "PL", Currencies: []string{"EUR"}},
		{CountryCode: "QA", Currencies: []string{"QAR"}},
		{CountryCode: "SA", Currencies: []string{"SAR"}},
		{CountryCode: "SG", Currencies: []string{"SGD", "USD"}},
		{CountryCode: "TH", Currencies: []string{"THB", "USD"}},
		{CountryCode: "TR", Currencies: []string{"TRY"}},
		{CountryCode: "TW", Currencies: []string{"TWD", "USD"}},
		{CountryCode: "US", Currencies: []string{"USD"}},
		{CountryCode: "VN", Currencies: []string{"VND", "USD"}},
		{CountryCode: "IT", Currencies: []string{"EUR"}},
		{CountryCode: "EU", Currencies: []string{"EUR"}},
	}
	if !reflect.DeepEqual(options, want) {
		t.Fatalf("collection catalog mismatch\ngot:  %#v\nwant: %#v", options, want)
	}
	if len(options) != len(HaiPayRegionCodes) {
		t.Fatalf("collection options=%d regions=%d", len(options), len(HaiPayRegionCodes))
	}
	options[0].CountryCode = "changed"
	options[7].Currencies[0] = "changed"
	if IsHaiPayRegion("changed") {
		t.Fatal("collection catalog must be copied")
	}
	if got := HaiPayCollectionCurrencies("HK")[0]; got != "HKD" {
		t.Fatalf("collection currency catalog was mutated: %q", got)
	}
}

func TestHaiPayPayoutCountryCurrency(t *testing.T) {
	if got := HaiPayPayoutCurrency("id"); got != "IDR" {
		t.Fatalf("ID currency=%q", got)
	}
	if got := HaiPayPayoutCountryFromCurrency("idr"); got != "ID" {
		t.Fatalf("IDR country=%q", got)
	}
	if !IsHaiPayPayoutCountry("US") || HaiPayPayoutCurrency("US") != "USD" {
		t.Fatal("US should be available as USD payout country")
	}
	if n := len(ListHaiPayPayoutCountries()); n == 0 {
		t.Fatal("empty payout countries")
	}
	if !IsHaiPayPayoutAccountType("IDR", HaiPayPayoutAccountTypeBank) ||
		!IsHaiPayPayoutAccountType("IDR", HaiPayPayoutAccountTypeEWallet) {
		t.Fatal("IDR should support bank account and e-wallet")
	}
	if IsHaiPayPayoutAccountType("BRL", HaiPayPayoutAccountTypeBank) ||
		!IsHaiPayPayoutAccountType("BRL", HaiPayPayoutAccountTypeEWallet) {
		t.Fatal("BRL should only support e-wallet")
	}
	if got := HaiPayPayoutRegion("ID"); got != HaiPayPayoutRegionAsia {
		t.Fatalf("ID region=%q", got)
	}
	if got := HaiPayPayoutRegion("BR"); got != HaiPayPayoutRegionSouthAmerica {
		t.Fatalf("BR region=%q", got)
	}
	if got := HaiPayPayoutRegion("EG"); got != HaiPayPayoutRegionMiddleEast {
		t.Fatalf("EG region=%q", got)
	}
	if got := HaiPayPayoutRegion("ZA"); got != HaiPayPayoutRegionAfrica {
		t.Fatalf("ZA region=%q", got)
	}
	if got := HaiPayPayoutRegion("US"); got != HaiPayPayoutRegionNorthAmerica {
		t.Fatalf("US region=%q", got)
	}
	if IsHaiPayPayoutAccountType("USD", HaiPayPayoutAccountTypeBank) ||
		!IsHaiPayPayoutAccountType("USD", HaiPayPayoutAccountTypeEWallet) {
		t.Fatal("USD should expose e-wallet until ACH fields are supported")
	}
	for _, code := range HaiPayPayoutCountryCodes {
		if got := HaiPayPayoutRegion(code); got == "" {
			t.Fatalf("missing payout region for country %s", code)
		}
	}
	types := ListHaiPayPayoutAccountTypes("IDR")
	types[0] = "changed"
	if IsHaiPayPayoutAccountType("IDR", "changed") {
		t.Fatal("payout account type list must be copied")
	}
	options := ListHaiPayPayoutOptions()
	if len(options) == 0 || len(options) != len(HaiPayPayoutCountryCodes) {
		t.Fatalf("payout options=%d countries=%d", len(options), len(HaiPayPayoutCountryCodes))
	}
	options[0].AccountTypes[0] = "changed"
	if IsHaiPayPayoutAccountType("USD", "changed") {
		t.Fatal("payout catalog must deep-copy account types")
	}
	wallets := ListHaiPayPayoutWallets("usd")
	if len(wallets) != 2 || wallets[0].Code != "VENMO" || wallets[1].Code != "ECASHAPP" {
		t.Fatalf("USD wallets=%v", wallets)
	}
	wallets[0].Code = "changed"
	if got := ListHaiPayPayoutWallets("USD")[0].Code; got != "VENMO" {
		t.Fatalf("payout wallet catalog was mutated: %q", got)
	}
	options[0].Wallets[0].Code = "changed"
	if got := ListHaiPayPayoutOptions()[0].Wallets[0].Code; got != "VENMO" {
		t.Fatalf("payout option wallet catalog was mutated: %q", got)
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
