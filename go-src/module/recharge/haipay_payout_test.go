package recharge

import (
	"strings"
	"testing"
)

func TestHaiPayNormalizePayoutAmount(t *testing.T) {
	tests := []struct {
		name       string
		currency   string
		amount     float64
		wantAmount float64
		wantText   string
	}{
		{name: "IDR rounds to integer", currency: "IDR", amount: 319839.6, wantAmount: 319840, wantText: "319840"},
		{name: "IDR normalizes code", currency: " idr ", amount: 319839.4, wantAmount: 319839, wantText: "319839"},
		{name: "USD keeps two decimals", currency: "USD", amount: 19.999, wantAmount: 20, wantText: "20.00"},
		{name: "two-decimal currency rounds", currency: "PHP", amount: 579.876, wantAmount: 579.88, wantText: "579.88"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := HaiPayNormalizePayoutAmount(test.currency, test.amount); got != test.wantAmount {
				t.Fatalf("amount=%v want=%v", got, test.wantAmount)
			}
			if got := haiPayFormatPayoutAmount(test.currency, test.amount); got != test.wantText {
				t.Fatalf("text=%q want=%q", got, test.wantText)
			}
		})
	}
}

func TestHaiPayPayoutExtraParams(t *testing.T) {
	t.Run("Brazil PIX", func(t *testing.T) {
		params, err := haiPayPayoutExtraParams(&HaiPayPayoutApplyReq{
			Currency: "BRL", BankCode: "PIX", IdentifyType: "cpf",
		})
		if err != nil || params["identifyType"] != "CPF" {
			t.Fatalf("params=%v err=%v", params, err)
		}
	})

	t.Run("Turkey Papara", func(t *testing.T) {
		params, err := haiPayPayoutExtraParams(&HaiPayPayoutApplyReq{
			Currency: "TRY", BankCode: "PAPARA", IdentifyType: "phone",
		})
		if err != nil || params["identifyType"] != "PHONE" {
			t.Fatalf("params=%v err=%v", params, err)
		}
	})

	t.Run("United States ACH", func(t *testing.T) {
		params, err := haiPayPayoutExtraParams(&HaiPayPayoutApplyReq{
			Currency:     "USD",
			BankCode:     "ACH",
			IdentifyType: "021000021",
			Country:      "us",
			Address1:     "1 Main St",
			Address2:     "New York",
			Address3:     "NY",
			PostalCode:   "10001",
		})
		if err != nil {
			t.Fatal(err)
		}
		if params["country"] != "US" || params["identifyType"] != "021000021" ||
			params["address1"] != "1 Main St" || params["address2"] != "New York" ||
			params["address3"] != "NY" || params["postalCode"] != "10001" {
			t.Fatalf("unexpected ACH params: %v", params)
		}
	})

	t.Run("reject missing ACH address", func(t *testing.T) {
		_, err := haiPayPayoutExtraParams(&HaiPayPayoutApplyReq{
			Currency: "USD", BankCode: "ACH", IdentifyType: "021000021", Country: "US",
		})
		if err == nil || !strings.Contains(err.Error(), "address missing") {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("reject invalid PIX identifier type", func(t *testing.T) {
		_, err := haiPayPayoutExtraParams(&HaiPayPayoutApplyReq{
			Currency: "BRL", BankCode: "PIX", IdentifyType: "passport",
		})
		if err == nil || !strings.Contains(err.Error(), "identifyType") {
			t.Fatalf("err=%v", err)
		}
	})
}
