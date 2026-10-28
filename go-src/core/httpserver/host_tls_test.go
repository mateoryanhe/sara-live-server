package httpserver

import (
	"crypto/tls"
	"testing"
)

func TestNormalizeHostTLSNameWildcard(t *testing.T) {
	got, err := normalizeHostTLSName(" *.SofieLive.net. ")
	if err != nil {
		t.Fatalf("normalizeHostTLSName() error = %v", err)
	}
	if got != "*.sofielive.net" {
		t.Fatalf("got %q, want *.sofielive.net", got)
	}
}

func TestNormalizeHostTLSNameRejectsBadWildcard(t *testing.T) {
	cases := []string{
		"*",
		"*.net",
		"a.*.sofielive.net",
		"*.*.sofielive.net",
		"api.*.net",
	}
	for _, input := range cases {
		if _, err := normalizeHostTLSName(input); err == nil {
			t.Fatalf("normalizeHostTLSName(%q) unexpectedly succeeded", input)
		}
	}
}

func TestLookupHostTLSCertificateExactThenWildcard(t *testing.T) {
	exactCert := &tls.Certificate{}
	wildCert := &tls.Certificate{}
	hostTLSRegistry.Store(&hostTLSRegistrySnapshot{
		certificates: map[string]*tls.Certificate{
			"api.sofielive.net": exactCert,
			"*.sofielive.net":   wildCert,
		},
	})
	t.Cleanup(func() {
		hostTLSRegistry.Store(nil)
	})

	if got := lookupHostTLSCertificate("api.sofielive.net"); got != exactCert {
		t.Fatalf("exact match should win")
	}
	if got := lookupHostTLSCertificate("v1.sofielive.net"); got != wildCert {
		t.Fatalf("wildcard should match v1.sofielive.net")
	}
	if got := lookupHostTLSCertificate("sofielive.net"); got != nil {
		t.Fatalf("apex should not match *.sofielive.net")
	}
	if got := lookupHostTLSCertificate("a.b.sofielive.net"); got != nil {
		t.Fatalf("multi-level subdomain should not match single-level wildcard")
	}
}
