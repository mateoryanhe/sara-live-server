package userwithdraw

import (
	"testing"

	"xr-game-server/constants/country"
)

func TestUserWithdrawPayoutFixedAccountType(t *testing.T) {
	for _, code := range []string{"ID", "MY", " id ", "my"} {
		fixed, ok := UserWithdrawPayoutFixedAccountType(code)
		if !ok || fixed != country.HaiPayPayoutAccountTypeEWallet {
			t.Fatalf("UserWithdrawPayoutFixedAccountType(%q) = (%q, %v)", code, fixed, ok)
		}
	}
	if fixed, ok := UserWithdrawPayoutFixedAccountType("US"); ok || fixed != "" {
		t.Fatalf("US should not be fixed, got (%q, %v)", fixed, ok)
	}
}

func TestAppPayoutExtraRequirementsForUserWithdrawIDMY(t *testing.T) {
	for _, cc := range []string{"ID", "MY"} {
		extra := appPayoutExtraRequirementsForUserWithdraw(cc, "IDR", "OVO")
		if extra.IdentifyTypeRequired || extra.AddressRequired || extra.CountryRequired {
			t.Fatalf("%s extra should be empty, got %+v", cc, extra)
		}
	}
}
