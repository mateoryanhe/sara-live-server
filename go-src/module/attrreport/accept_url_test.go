package attrreport

import "testing"

func TestBuildAcceptURL(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"https://cdntk.showsme.click", "https://cdntk.showsme.click/tj/api/accept.php"},
		{"https://cdntk.showsme.click/", "https://cdntk.showsme.click/tj/api/accept.php"},
		{"https://cdntk.showsme.click/tj/api/accept.php", "https://cdntk.showsme.click/tj/api/accept.php"},
		{"https://cdntk.showsme.click/tj/api/accept.php/", "https://cdntk.showsme.click/tj/api/accept.php"},
	}
	for _, c := range cases {
		if got := buildAcceptURL(c.in); got != c.want {
			t.Fatalf("buildAcceptURL(%q)=%q want %q", c.in, got, c.want)
		}
	}
}
