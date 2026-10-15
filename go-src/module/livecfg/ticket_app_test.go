package livecfg

import "testing"

func TestAppLiveRoomCallTicketPrice(t *testing.T) {
	liveCfgCache.Store(&liveCfgSnapshot{VideoCallTicketEnabled: false})
	if got := AppLiveRoomCallTicketPrice(12.5); got != 0 {
		t.Fatalf("disabled ticket switch: got %v want 0", got)
	}
	liveCfgCache.Store(&liveCfgSnapshot{VideoCallTicketEnabled: true})
	if got := AppLiveRoomCallTicketPrice(12.5); got != 12.5 {
		t.Fatalf("enabled ticket switch: got %v want 12.5", got)
	}
}
