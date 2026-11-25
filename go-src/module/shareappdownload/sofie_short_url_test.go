package shareappdownload

import "testing"

func TestResolveSofieShareSnapshotSingleRow(t *testing.T) {
	sofieShareCfgCache.Store(&sofieShareRowSnapshot{
		PackageName:   "com.sofie.live",
		GooglePlayUrl: "https://play.google.com/store/apps/details?id=com.sofie.live",
		ShareText:     "Sofie",
	})

	if snap := ResolveSofieShareSnapshot("com.sofie.live"); snap == nil || snap.ShareText != "Sofie" {
		t.Fatalf("expected cfg, got %#v", snap)
	}
	if snap := ResolveSofieShareSnapshot("com.other.app"); snap == nil || snap.ShareText != "Sofie" {
		t.Fatalf("expected same row as fallback, got %#v", snap)
	}
}

func TestBuildShareCodeFromSofieSnapshot(t *testing.T) {
	snap := &sofieShareRowSnapshot{
		PackageName:   "com.sofie.live",
		GooglePlayUrl: "https://play.google.com/store/apps/details?id=com.sofie.live",
		ShareText:     "Join Sofie",
	}
	got := buildShareCodeFromSofieSnapshot(snap, "ABC")
	if got == "" {
		t.Fatal("expected non-empty share code")
	}
	if want := "Join Sofie\nInvite Code: ABC"; got[:len(want)] != want {
		t.Fatalf("unexpected prefix: %q", got)
	}
}