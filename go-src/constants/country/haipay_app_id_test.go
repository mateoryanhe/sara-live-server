package country

import "testing"

func TestLookupHaiPayAppID(t *testing.T) {
	want := map[string]int64{
		"IDR": 25238, "PHP": 25239, "USD": 25240, "KRW": 25241, "EUR": 25242,
		"TWD": 25243, "VND": 25244, "THB": 25245, "SAR": 25246, "KWD": 25247,
		"OMR": 25248, "QAR": 25249, "AED": 25250, "BHD": 25251, "USDT": 25252,
		"BRL": 25253, "HKD": 25254, "INR": 25255, "MYR": 25256, "TRY": 25257,
		"PKR": 25258, "SGD": 25259, "RUB": 25260, "EGP": 25261, "JPY": 25262,
		"GBP": 25263, "PLN": 25264, "CAD": 25265, "JOD": 25266, "IQD": 25267,
		"BDT": 25268, "USDC": 25269, "MXN": 25270, "NGN": 25271, "CASHIER": 25272,
	}
	for code, wantAppID := range want {
		got, ok := LookupHaiPayAppID(code)
		if !ok || got != wantAppID {
			t.Fatalf("LookupHaiPayAppID(%q) = (%d, %t), want (%d, true)", code, got, ok, wantAppID)
		}
	}
	if got, ok := LookupHaiPayAppID(" usd "); !ok || got != 25240 {
		t.Fatalf("LookupHaiPayAppID should normalize code, got (%d, %t)", got, ok)
	}
	if got, ok := LookupHaiPayAppID("UNKNOWN"); ok || got != 0 {
		t.Fatalf("LookupHaiPayAppID(UNKNOWN) = (%d, %t), want (0, false)", got, ok)
	}
}
