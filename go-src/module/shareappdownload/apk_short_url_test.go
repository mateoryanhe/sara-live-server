package shareappdownload

import (
	"strings"
	"testing"
)

func TestBuildAppShareCodeThreeLines(t *testing.T) {
	oldCfg := getCfgCache()
	oldURL := GetDownloadUrl()
	t.Cleanup(func() {
		cfgCache.Store(oldCfg)
		storeApkDownloadURL(oldURL)
	})

	cfgCache.Store(&cfgSnapshot{ShareText: "Join Sofie with me!"})
	storeApkDownloadURL("")
	got := BuildAppShareCode("ABC123")
	want := "Join Sofie with me!\n" + InviteCodeLinePrefix + "ABC123"
	if got != want {
		t.Fatalf("without short url got %q want %q", got, want)
	}

	cfgCache.Store(&cfgSnapshot{ShareText: "Hi {inviteCode}\nurl={downloadUrl}"})
	got = BuildAppShareCode("XYZ")
	if strings.Contains(got, "{inviteCode}") || strings.Contains(got, "{downloadUrl}") {
		t.Fatalf("placeholders should be stripped, got %q", got)
	}
	if !strings.HasSuffix(got, "\n"+InviteCodeLinePrefix+"XYZ") && got != InviteCodeLinePrefix+"XYZ" {
		t.Fatalf("invite line should end with labeled code, got %q", got)
	}
}

func TestBuildAppShareCodeInviteOnly(t *testing.T) {
	oldCfg := getCfgCache()
	oldURL := GetDownloadUrl()
	t.Cleanup(func() {
		cfgCache.Store(oldCfg)
		storeApkDownloadURL(oldURL)
	})
	cfgCache.Store(&cfgSnapshot{})
	storeApkDownloadURL("")
	if got := BuildAppShareCode("CODE"); got != InviteCodeLinePrefix+"CODE" {
		t.Fatalf("got %q", got)
	}
}

func TestBuildPreviewShareCode(t *testing.T) {
	oldCfg := getCfgCache()
	oldURL := GetDownloadUrl()
	t.Cleanup(func() {
		cfgCache.Store(oldCfg)
		storeApkDownloadURL(oldURL)
	})
	cfgCache.Store(&cfgSnapshot{ShareText: "Hello"})
	storeApkDownloadURL("")
	got := BuildPreviewShareCode()
	want := "Hello\n" + InviteCodeLinePrefix + PreviewInviteCodePlaceholder
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
