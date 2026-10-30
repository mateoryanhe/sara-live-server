package upload

import (
	"strings"
	"testing"
	"time"
)

func TestExtractApkVersionLabel(t *testing.T) {
	tests := []struct {
		file string
		want string
	}{
		{"app-release-1.2.3.apk", "1.2.3"},
		{"SaraLive_v2.0.1_build.apk", "2.0.1"},
		{"channel.apk", ""},
		{"", ""},
		{"foo-10.20.30.40.apk", "10.20.30.40"},
	}
	for _, tc := range tests {
		got := extractApkVersionLabel(tc.file)
		if got != tc.want {
			t.Fatalf("extractApkVersionLabel(%q) = %q, want %q", tc.file, got, tc.want)
		}
	}
}

func TestBuildApkPromotionChannelStoredName_uniquePerCall(t *testing.T) {
	a := BuildApkPromotionChannelStoredName("ch01", "app-1.0.0.apk")
	time.Sleep(2 * time.Millisecond)
	b := BuildApkPromotionChannelStoredName("ch01", "app-1.0.0.apk")
	if a == "" || b == "" {
		t.Fatal("empty stored name")
	}
	if a == b {
		t.Fatalf("expected unique names, got same %q", a)
	}
	if !strings.Contains(a, "ch01_v1.0.0_") {
		t.Fatalf("unexpected pattern: %q", a)
	}
}
