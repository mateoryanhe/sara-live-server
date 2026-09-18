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
	if got := HaiPayCollectionCurrency("in"); got != "USD" {
		t.Fatalf("IN collection currency=%q", got)
	}
	if got := HaiPayCollectionCountryFromCurrency("eur"); got != "" {
		t.Fatalf("shared EUR must not infer a collection country, got=%q", got)
	}
	if got := HaiPayCollectionCountryFromCurrency("usd"); got != "" {
		t.Fatalf("shared USD must not infer a collection country, got=%q", got)
	}
	if got := HaiPayCollectionCountryFromCurrency("bdt"); got != "BD" {
		t.Fatalf("BDT collection country=%q", got)
	}
	if got := HaiPayCollectionCurrency("pl"); got != "USD" {
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
		{CountryCode: "AT", Currencies: []string{"EUR", "USD"}},
		{CountryCode: "BE", Currencies: []string{"EUR", "USD"}},
		{CountryCode: "IT", Currencies: []string{"EUR", "USD"}},
		{CountryCode: "NL", Currencies: []string{"EUR", "USD"}},
		{CountryCode: "PL", Currencies: []string{"USD"}},
		{CountryCode: "TR", Currencies: []string{"TRY"}},
		{CountryCode: "BR", Currencies: []string{"BRL"}},
		{CountryCode: "HK", Currencies: []string{"HKD", "USD"}},
		{CountryCode: "SG", Currencies: []string{"SGD", "USD"}},
		{CountryCode: "TW", Currencies: []string{"TWD", "USD"}},
		{CountryCode: "JP", Currencies: []string{"JPY", "USD"}},
		{CountryCode: "KR", Currencies: []string{"KRW", "USD"}},
		{CountryCode: "PH", Currencies: []string{"PHP"}},
		{CountryCode: "TH", Currencies: []string{"THB", "USD"}},
		{CountryCode: "VN", Currencies: []string{"VND", "USD"}},
		{CountryCode: "ID", Currencies: []string{"IDR", "USD"}},
		{CountryCode: "IN", Currencies: []string{"INR", "USD"}},
		{CountryCode: "MY", Currencies: []string{"MYR", "USD"}},
		{CountryCode: "PK", Currencies: []string{"PKR", "USD"}},
		{CountryCode: "BD", Currencies: []string{"BDT"}},
		{CountryCode: "EG", Currencies: []string{"EGP", "USD"}},
		{CountryCode: "SA", Currencies: []string{"SAR", "USD"}},
		{CountryCode: "AE", Currencies: []string{"AED", "USD"}},
		{CountryCode: "KW", Currencies: []string{"KWD", "USD"}},
		{CountryCode: "QA", Currencies: []string{"QAR", "USD"}},
		{CountryCode: "OM", Currencies: []string{"OMR", "USD"}},
		{CountryCode: "BH", Currencies: []string{"BHD", "USD"}},
		{CountryCode: "JO", Currencies: []string{"JOD"}},
		{CountryCode: "IQ", Currencies: []string{"IQD"}},
	}
	if !reflect.DeepEqual(options, want) {
		t.Fatalf("collection catalog mismatch\ngot:  %#v\nwant: %#v", options, want)
	}
	if len(options) != len(HaiPayRegionCodes) {
		t.Fatalf("collection options=%d regions=%d", len(options), len(HaiPayRegionCodes))
	}
	options[0].CountryCode = "changed"
	options[8].Currencies[0] = "changed"
	if IsHaiPayRegion("changed") {
		t.Fatal("collection catalog must be copied")
	}
	if got := HaiPayCollectionCurrencies("HK")[0]; got != "HKD" {
		t.Fatalf("collection currency catalog was mutated: %q", got)
	}
}

func TestHaiPayCollectionPaymentMethods(t *testing.T) {
	methodCount := 0
	for countryCode, countryMethods := range haiPayCollectionPaymentMethodsByCountry {
		supportedCurrencies := HaiPayCollectionCurrencies(countryCode)
		for _, method := range countryMethods {
			methodCount++
			if !method.Available || method.CurrencyCode == "" || method.PayType == "" || method.InBankCode == "" ||
				method.MinAmount == "" || method.MaxAmount == "" {
				t.Fatalf("incomplete available method country=%s method=%+v", countryCode, method)
			}
			if !IsHaiPayCollectionPaymentMethod(countryCode, method.CurrencyCode, method.PayType, method.InBankCode) {
				t.Fatalf("catalog method rejected country=%s method=%+v", countryCode, method)
			}
			if !containsCode(supportedCurrencies, method.CurrencyCode) {
				t.Fatalf("payment method currency is not selectable country=%s method=%+v supported=%v", countryCode, method, supportedCurrencies)
			}
		}
	}
	if len(haiPayCollectionPaymentMethodsByCountry) != 28 || methodCount != 142 || len(haiPayGlobalCollectionPaymentMethods) != 3 {
		t.Fatalf("payment method catalog countries=%d methods=%d", len(haiPayCollectionPaymentMethodsByCountry), methodCount)
	}
	for _, option := range ListHaiPayCollectionOptions() {
		if option.CountryCode == "TR" {
			if HasHaiPayCollectionPaymentMethodCatalog(option.CountryCode) {
				t.Fatal("Turkey collection methods are all under maintenance")
			}
			continue
		}
		if !HasHaiPayCollectionPaymentMethodCatalog(option.CountryCode) {
			t.Fatalf("missing collection payment method catalog country=%s", option.CountryCode)
		}
	}

	methods := ListHaiPayCollectionPaymentMethods("hk")
	if len(methods) != 7 {
		t.Fatalf("HK collection payment methods=%v", methods)
	}
	if methods[0].CurrencyCode != "HKD" || methods[0].PayType != "EWALLET" || methods[0].InBankCode != "WXPAY_SCANCODE" {
		t.Fatalf("unexpected HK payment method=%+v", methods[0])
	}
	if !IsHaiPayCollectionPaymentMethod("HK", "HKD", "EWALLET", "WXPAY_SCANCODE") {
		t.Fatal("HKD WXPAY_SCANCODE must be available")
	}
	if IsHaiPayCollectionPaymentMethod("HK", "HKD", "EWALLET", "HK_WXPAY_SCANCODE_USD") {
		t.Fatal("USD payment code must not be accepted for HKD")
	}
	if IsHaiPayCollectionPaymentMethod("HK", "USD", "BANK_TRANSFER", "PAYME_USD") {
		t.Fatal("PAYME_USD is under maintenance and must not be accepted")
	}
	if got, ok := ResolveHaiPayCollectionPaymentMethodCode("ID", "IDR", "QR", "DYNAMIC"); !ok || got != "dynamic" {
		t.Fatalf("ID dynamic canonical code=%q ok=%v", got, ok)
	}
	indonesiaMethods := ListHaiPayCollectionPaymentMethods("ID")
	if len(indonesiaMethods) != 13 {
		t.Fatalf("ID available collection payment methods=%d want=13", len(indonesiaMethods))
	}
	if IsHaiPayCollectionPaymentMethod("ID", "IDR", "VA", "002") ||
		IsHaiPayCollectionPaymentMethod("ID", "USD", "VA", "ID_BRI_USD") {
		t.Fatal("Indonesia methods under maintenance must not be configurable")
	}
	if !IsHaiPayCollectionPaymentMethod("AT", "EUR", "BANK_TRANSFER", "EPS_EUR") ||
		!IsHaiPayCollectionPaymentMethod("KR", "KRW", "BANK_TRANSFER", "KAKAOPAY_KRW") ||
		!IsHaiPayCollectionPaymentMethod("AE", "AED", "VA", "BANK_TRANSFER") ||
		!IsHaiPayCollectionPaymentMethod("ID", "USD", "EWALLET", "APPLE_PAY") {
		t.Fatal("new official collection payment methods must be configurable")
	}
	if IsHaiPayCollectionPaymentMethod("PL", "EUR", "BANK_TRANSFER", "BLIK_USD") ||
		IsHaiPayCollectionPaymentMethod("TR", "TRY", "EWALLET", "PAPARA") ||
		IsHaiPayCollectionPaymentMethod("ID", "IDR", "EWALLET", "APPLE_PAY") {
		t.Fatal("currency-mismatched or maintenance methods must not be configurable")
	}
	methods[0].InBankCode = "changed"
	if got := ListHaiPayCollectionPaymentMethods("HK")[0].InBankCode; got != "WXPAY_SCANCODE" {
		t.Fatalf("payment method catalog was mutated: %q", got)
	}
}

func containsCode(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func TestHaiPayGlobalCashierOptions(t *testing.T) {
	expected := map[string][]string{
		"US": {"USD"}, "MX": {"MXN"},
		"AT": {"EUR"}, "BE": {"EUR"}, "GB": {"GBP"}, "NL": {"EUR"},
		"PL": {"EUR"}, "TR": {"TRY"}, "IT": {"EUR"}, "EU": {"EUR"},
		"BR": {"BRL"},
		"HK": {"HKD", "USD"}, "SG": {"SGD", "USD"}, "TW": {"TWD", "USD"},
		"JP": {"JPY", "USD"}, "KR": {"KRW", "USD"}, "PH": {"PHP"},
		"TH": {"THB", "USD"}, "VN": {"VND", "USD"}, "ID": {"IDR", "USD"},
		"MY": {"MYR", "USD"}, "PK": {"PKR"}, "BD": {"BDT"}, "IN": {"INR"},
		"EG": {"EGP"}, "SA": {"SAR"}, "KW": {"KWD"}, "BH": {"BHD"},
		"AE": {"AED"}, "OM": {"OMR"}, "QA": {"QAR"}, "JO": {"JOD"}, "IQ": {"IQD"},
		"NG": {"NGN"}, "GH": {"GHS"}, "KE": {"KES"}, "ZA": {"ZAR"},
		"CM": {"XAF"}, "TZ": {"TZS"},
	}
	options := ListHaiPayGlobalCashierOptions()
	if len(options) != len(expected) {
		t.Fatalf("global cashier options=%d want=%d", len(options), len(expected))
	}
	continentCounts := make(map[string]int)
	for _, option := range options {
		want, ok := expected[option.CountryCode]
		if !ok {
			t.Fatalf("unexpected global cashier region=%s", option.CountryCode)
		}
		if len(option.Currencies) != len(want) {
			t.Fatalf("region=%s currencies=%v want=%v", option.CountryCode, option.Currencies, want)
		}
		for i := range want {
			if option.Currencies[i] != want[i] {
				t.Fatalf("region=%s currencies=%v want=%v", option.CountryCode, option.Currencies, want)
			}
		}
		continent := HaiPayGlobalCashierContinent(option.CountryCode)
		if !IsHaiPayGlobalCashierRegion(option.CountryCode) || continent == "" {
			t.Fatalf("global cashier region metadata missing region=%s", option.CountryCode)
		}
		continentCounts[continent]++
	}
	expectedContinentCounts := map[string]int{
		HaiPayPayoutRegionNorthAmerica: 2,
		HaiPayPayoutRegionEurope:       8,
		HaiPayPayoutRegionSouthAmerica: 1,
		HaiPayPayoutRegionAsia:         13,
		HaiPayPayoutRegionMiddleEast:   9,
		HaiPayPayoutRegionAfrica:       6,
	}
	if !reflect.DeepEqual(continentCounts, expectedContinentCounts) {
		t.Fatalf("global cashier continent counts=%v want=%v", continentCounts, expectedContinentCounts)
	}
	if got := HaiPayGlobalCashierCurrency("ID"); got != "USD" {
		t.Fatalf("ID default global cashier currency=%q want=USD", got)
	}
	if got := HaiPayGlobalCashierCountryFromCurrency("IDR"); got != "ID" {
		t.Fatalf("IDR region=%q want=ID", got)
	}
	if got := HaiPayGlobalCashierCountryFromCurrency("USD"); got != "" {
		t.Fatalf("shared USD must require explicit region, got=%q", got)
	}
	if got := HaiPayGlobalCashierCountryFromCurrency("EUR"); got != "" {
		t.Fatalf("shared EUR must require explicit region, got=%q", got)
	}
	options[0].Currencies[0] = "changed"
	if got := ListHaiPayGlobalCashierOptions()[0].Currencies[0]; got != "USD" {
		t.Fatalf("global cashier catalog was mutated: %q", got)
	}
	if !IsHaiPayGlobalCashierRegion("US") || !IsHaiPayGlobalCashierRegion("AT") {
		t.Fatal("existing global cashier regions must remain available")
	}
	if got := HaiPayGlobalCashierContinent("MX"); got != HaiPayPayoutRegionNorthAmerica {
		t.Fatalf("MX continent=%q", got)
	}
	if got := HaiPayGlobalCashierContinent("NG"); got != HaiPayPayoutRegionAfrica {
		t.Fatalf("NG continent=%q", got)
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
	if got := HaiPayPayoutRegion("EU"); got != HaiPayPayoutRegionEurope {
		t.Fatalf("EU region=%q", got)
	}
	if got := HaiPayPayoutRegion("US"); got != HaiPayPayoutRegionNorthAmerica {
		t.Fatalf("US region=%q", got)
	}
	if !IsHaiPayPayoutAccountType("USD", HaiPayPayoutAccountTypeBank) ||
		!IsHaiPayPayoutAccountType("USD", HaiPayPayoutAccountTypeEWallet) {
		t.Fatal("USD should expose bank account and e-wallet")
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
	methodCount := 0
	methodKeys := make(map[string]struct{})
	for _, option := range options {
		if len(option.Methods) == 0 {
			t.Fatalf("payout option %s has no available methods", option.CountryCode)
		}
		for _, method := range option.Methods {
			if method.AccountType == "" || method.BankCode == "" || method.Limit == "" || method.Description == "" {
				t.Fatalf("incomplete payout method for %s: %+v", option.Currency, method)
			}
			key := option.Currency + "|" + method.AccountType + "|" + method.BankCode
			if _, exists := methodKeys[key]; exists {
				t.Fatalf("duplicate payout method %s", key)
			}
			methodKeys[key] = struct{}{}
			methodCount++
		}
	}
	if methodCount != 478 {
		t.Fatalf("available payout methods=%d, want 478", methodCount)
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
	methods := ListHaiPayPayoutMethods("USD")
	if len(methods) != 3 {
		t.Fatalf("USD methods=%d, want 3", len(methods))
	}
	if !IsHaiPayPayoutMethod("USD", HaiPayPayoutAccountTypeBank, "ACH") ||
		!IsHaiPayPayoutMethod("USD", HaiPayPayoutAccountTypeEWallet, "VENMO") ||
		IsHaiPayPayoutMethod("USD", HaiPayPayoutAccountTypeBank, "VENMO") {
		t.Fatal("USD payout method pairs are incorrect")
	}
	if method, ok := FindHaiPayPayoutMethod("usd", "bank_account", "ach"); !ok || method.Limit != "10-200000" {
		t.Fatalf("USD ACH method=%+v ok=%v", method, ok)
	}
	methods[0].BankCode = "changed"
	if got := ListHaiPayPayoutMethods("USD")[0].BankCode; got != "VENMO" {
		t.Fatalf("payout method catalog was mutated: %q", got)
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
