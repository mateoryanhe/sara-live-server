package banner

import (
	"strings"
	"testing"
)

func TestAppendShareParamsToLink_AbsoluteURL(t *testing.T) {
	got := appendShareParamsToLink("https://example.com/h5/share?from=banner", "tok-abc")
	if !strings.HasPrefix(got, "https://example.com/h5/share?") {
		t.Fatalf("unexpected link: %s", got)
	}
	if !strings.Contains(got, "token=tok-abc") {
		t.Fatalf("missing token: %s", got)
	}
	if strings.Contains(got, "shareTx=") || strings.Contains(got, "shareCode=") {
		t.Fatalf("should only append token: %s", got)
	}
}

func TestAppendShareParamsToLink_EmptyToken(t *testing.T) {
	in := "https://example.com/h5"
	if got := appendShareParamsToLink(in, "  "); got != in {
		t.Fatalf("want unchanged, got %q", got)
	}
}
