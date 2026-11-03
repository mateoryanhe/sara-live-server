package banner

import "testing"

func TestBannerHidesForPackage(t *testing.T) {
	t.Parallel()
	ex := "com.sofie.live, com.other.app"
	if !bannerHidesForPackage(ex, "com.sofie.live") {
		t.Fatal("expected hide for listed package")
	}
	if !bannerHidesForPackage("com.foo\ncom.bar", "com.bar") {
		t.Fatal("expected hide for newline separated")
	}
	if bannerHidesForPackage(ex, "com.sara.live") {
		t.Fatal("expected show for non-listed package")
	}
	if bannerHidesForPackage(ex, "") {
		t.Fatal("empty package should not hide")
	}
}

func TestNormalizeExcludePackageNames(t *testing.T) {
	t.Parallel()
	got := normalizeExcludePackageNames("  a.b ; c.d \n e.f  ")
	want := "a.b,c.d,e.f"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
