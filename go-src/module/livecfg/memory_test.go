package livecfg

import (
	"testing"

	liveentity "xr-game-server/entity/live"
)

func TestToLiveCfgSnapshotVideoCallTicketDefault(t *testing.T) {
	if !toLiveCfgSnapshot(nil).VideoCallTicketEnabled {
		t.Fatal("missing live config must default video call ticket charging to enabled")
	}
}

func TestToLiveCfgSnapshotVideoCallTicketValue(t *testing.T) {
	if toLiveCfgSnapshot(&liveentity.LiveCfg{VideoCallTicketEnabled: false}).VideoCallTicketEnabled {
		t.Fatal("stored disabled switch must remain disabled")
	}
	if !toLiveCfgSnapshot(&liveentity.LiveCfg{VideoCallTicketEnabled: true}).VideoCallTicketEnabled {
		t.Fatal("stored enabled switch must remain enabled")
	}
}
