package shareappdownload

import "testing"

func TestBuildLiveRoomAfLink(t *testing.T) {
	got := buildLiveRoomAfLink("https://af.example/abc", "INV9", 1006615)
	if got != "https://af.example/abc?inviteCode=INV9&roomId=1006615" {
		t.Fatalf("unexpected link: %q", got)
	}
	got = buildLiveRoomAfLink("https://af.example/abc?pid=x", "INV9", 1006615)
	if got != "https://af.example/abc?pid=x&inviteCode=INV9&roomId=1006615" {
		t.Fatalf("unexpected link with query: %q", got)
	}
}