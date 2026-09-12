package country

import "testing"

func TestFormatZhEn(t *testing.T) {
	got := FormatZhEn("CI")
	want := "科特迪瓦 / Cote d'Ivoire"
	if got != want {
		t.Fatalf("FormatZhEn(CI)=%q want %q", got, want)
	}
	if FormatZhEn("") != "" {
		t.Fatal("empty")
	}
	if FormatZhEn("历史中文") != "历史中文" {
		t.Fatal("legacy passthrough")
	}
}
