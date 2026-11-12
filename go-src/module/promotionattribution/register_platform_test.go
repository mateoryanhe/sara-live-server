package promotionattribution

import "testing"

func TestMapClientXPlatformToRegisterPlatform(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"web", RegisterPlatformH5},
		{"WEB", RegisterPlatformH5},
		{"android", RegisterPlatformApp},
		{"ios", RegisterPlatformApp},
		{"IOS", RegisterPlatformApp},
		{"", ""},
		{"unknown", ""},
	}
	for _, c := range cases {
		if got := MapClientXPlatformToRegisterPlatform(c.in); got != c.want {
			t.Fatalf("MapClientXPlatformToRegisterPlatform(%q)=%q want %q", c.in, got, c.want)
		}
	}
}

func TestResolveRegisterPromotionPlatformPrefersHeader(t *testing.T) {
	// android header beats email channel(7)
	if got := ResolveRegisterPromotionPlatform("android", registerPlatformEmailChannel); got != RegisterPlatformApp {
		t.Fatalf("got %q want app", got)
	}
	// web header
	if got := ResolveRegisterPromotionPlatform("web", registerPlatformDeviceChannel); got != RegisterPlatformH5 {
		t.Fatalf("got %q want h5", got)
	}
	// empty header → channel fallback
	if got := ResolveRegisterPromotionPlatform("", registerPlatformEmailChannel); got != RegisterPlatformH5 {
		t.Fatalf("got %q want h5 fallback", got)
	}
	if got := ResolveRegisterPromotionPlatform("", registerPlatformDeviceChannel); got != RegisterPlatformApp {
		t.Fatalf("got %q want app fallback", got)
	}
}
