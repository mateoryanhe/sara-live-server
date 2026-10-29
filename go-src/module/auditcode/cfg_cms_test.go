package auditcode

import (
	"strings"
	"testing"
)

func TestNewAuditCodeLengthAndAlphabet(t *testing.T) {
	for _, n := range []int{4, 8, 32} {
		code, err := newAuditCode(n)
		if err != nil {
			t.Fatalf("n=%d: %v", n, err)
		}
		if len(code) != n {
			t.Fatalf("n=%d got len %d code=%q", n, len(code), code)
		}
		for _, r := range code {
			if !strings.ContainsRune(auditCodeAlphabet, r) {
				t.Fatalf("n=%d invalid rune %q in %q", n, r, code)
			}
		}
	}
}
