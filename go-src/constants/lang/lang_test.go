package lang

import "testing"

func TestParseArabic(t *testing.T) {
	if got := Parse("ar"); got != LangAR {
		t.Fatalf("Parse(ar) = %q, want %q", got, LangAR)
	}
	if got := Parse("ar-SA,ar;q=0.9,en;q=0.8"); got != LangAR {
		t.Fatalf("Parse(ar-SA...) = %q, want %q", got, LangAR)
	}
}