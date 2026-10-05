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

func TestToLiveCfgSnapshotAudienceListRefreshDefault(t *testing.T) {
	if toLiveCfgSnapshot(nil).AudienceListRefreshSeconds != DefaultAudienceListRefreshSeconds {
		t.Fatal("missing live config must default audience list refresh to 300 seconds")
	}
	if toLiveCfgSnapshot(&liveentity.LiveCfg{}).AudienceListRefreshSeconds != DefaultAudienceListRefreshSeconds {
		t.Fatal("zero stored refresh seconds must default to 300")
	}
}

func TestToLiveCfgSnapshotAudienceListRefreshValue(t *testing.T) {
	if got := toLiveCfgSnapshot(&liveentity.LiveCfg{AudienceListRefreshSeconds: 60}).AudienceListRefreshSeconds; got != 60 {
		t.Fatalf("stored refresh seconds = %d, want 60", got)
	}
}
