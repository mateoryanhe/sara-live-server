package cfg

import "testing"

func TestFindThirdPayDomain(t *testing.T) {
	sites := []*StaticSiteCfg{
		{Domain: "static.bigtktool.shop"},
		{Domain: " third-pay.bigtktool.shop,third-pay-backup.bigtktool.shop ", T: true},
	}

	if got, want := findThirdPayDomain(sites), "third-pay.bigtktool.shop"; got != want {
		t.Fatalf("findThirdPayDomain() = %q, want %q", got, want)
	}
}

func TestFindThirdPayDomainWithoutMarker(t *testing.T) {
	if got := findThirdPayDomain([]*StaticSiteCfg{{Domain: "third-pay.bigtktool.shop"}}); got != "" {
		t.Fatalf("findThirdPayDomain() = %q, want empty domain", got)
	}
}
