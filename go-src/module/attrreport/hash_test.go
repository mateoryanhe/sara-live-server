package attrreport

import "testing"

func TestHashEmail_JohnSmith(t *testing.T) {
	got := HashEmail(" John_Smith@gmail.com ")
	want := "62a14e44f765419d10fea99367361a727c12365e2520f32218d505ed9aa0f62f"
	if got != want {
		t.Fatalf("HashEmail got %s want %s norm=%s", got, want, NormalizeEmail(" John_Smith@gmail.com "))
	}
}

func TestHashPhone_USExample(t *testing.T) {
	got := HashPhone("1", "(650)555-1212")
	want := "e323ec626319ca94ee8bff2e4c87cf613be6ea19919ed1364124e16807ab3176"
	if got != want {
		t.Fatalf("HashPhone got %s want %s norm=%s", got, want, NormalizePhone("1", "(650)555-1212"))
	}
}

func TestNormalizePhone_AlreadyWithCountry(t *testing.T) {
	if got := NormalizePhone("1", "16505551212"); got != "16505551212" {
		t.Fatalf("got %s", got)
	}
}

func TestSplitClientIPs(t *testing.T) {
	v4, v6 := SplitClientIPs("1.2.3.4")
	if v4 != "1.2.3.4" || v6 != "" {
		t.Fatalf("v4=%q v6=%q", v4, v6)
	}
	v4, v6 = SplitClientIPs("2001:db8::1")
	if v4 != "" || v6 != "2001:db8::1" {
		t.Fatalf("v4=%q v6=%q", v4, v6)
	}
}
