package recharge

import (
	"reflect"
	"testing"

	"xr-game-server/constants/country"
)

func TestHaiPayCoinMerchantCollectionCountriesAlignWithNormalCollection(t *testing.T) {
	normalOptions := country.ListHaiPayGlobalCashierOptions()
	merchantOptions := haiPayCoinMerchantCollectionOptions()
	if len(merchantOptions) != len(normalOptions) {
		t.Fatalf("coin merchant countries=%d normal countries=%d", len(merchantOptions), len(normalOptions))
	}

	for index := range normalOptions {
		if merchantOptions[index].CountryCode != normalOptions[index].CountryCode {
			t.Fatalf(
				"country order differs at %d: coin merchant=%s normal=%s",
				index, merchantOptions[index].CountryCode, normalOptions[index].CountryCode,
			)
		}
		if continent := country.HaiPayGlobalCashierContinent(merchantOptions[index].CountryCode); continent == "" {
			t.Fatalf("country %s has no continent", merchantOptions[index].CountryCode)
		}
	}

	if got := haiPayCoinMerchantCollectionCurrencies("ID"); !reflect.DeepEqual(got, []string{"IDR", "USD"}) {
		t.Fatalf("ID coin merchant currencies=%v", got)
	}
	if got := haiPayCoinMerchantCollectionCurrencies("US"); !reflect.DeepEqual(got, []string{"USD"}) {
		t.Fatalf("US coin merchant currencies=%v", got)
	}
	if got := haiPayCoinMerchantCollectionCurrency("AT"); got != "USD" {
		t.Fatalf("AT default coin merchant currency=%q", got)
	}
	if got := haiPayCoinMerchantCollectionCurrency("GB"); got != "GBP" {
		t.Fatalf("GB default coin merchant currency=%q", got)
	}
}

func TestHaiPayCoinMerchantNewUSDRegionUsesGlobalMethods(t *testing.T) {
	methods := country.ListHaiPayCollectionPaymentMethods("US")
	if len(methods) != 3 {
		t.Fatalf("US collection methods=%d want=3", len(methods))
	}
	if code, ok := country.ResolveHaiPayCollectionPaymentMethodCode("US", "USD", "BANK_TRANSFER", "CREDIT_CARD"); !ok || code != "CREDIT_CARD" {
		t.Fatalf("US credit card method code=%q ok=%v", code, ok)
	}
	if methods := country.ListHaiPayCollectionPaymentMethods("GB"); len(methods) != 0 {
		t.Fatalf("GB must not inherit USD methods: %v", methods)
	}
}
