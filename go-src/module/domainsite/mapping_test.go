package domainsite

import "testing"

func TestNormalizeDomain(t *testing.T) {
	got, err := normalizeDomain(" One.Example.com. ")
	if err != nil {
		t.Fatalf("normalizeDomain() error = %v", err)
	}
	if got != "one.example.com" {
		t.Fatalf("normalizeDomain() = %q, want one.example.com", got)
	}
}

func TestNormalizeDomainRejectsComma(t *testing.T) {
	if _, err := normalizeDomain("one.example.com,two.example.com"); err == nil {
		t.Fatal("comma-separated domains unexpectedly accepted")
	}
}
