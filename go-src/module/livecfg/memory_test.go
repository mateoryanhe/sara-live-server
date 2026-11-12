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

func TestToLiveCfgSnapshotShowcaseLiveRoomDefault(t *testing.T) {
	if toLiveCfgSnapshot(nil).ShowcaseLiveRoomEnabled {
		t.Fatal("missing live config must default showcase live room broadcast to disabled")
	}
	if toLiveCfgSnapshot(&liveentity.LiveCfg{ShowcaseLiveRoomEnabled: true}).ShowcaseLiveRoomEnabled != true {
		t.Fatal("stored enabled showcase switch must remain enabled")
	}
}

func TestToLiveCfgSnapshotAudienceListRefreshValue(t *testing.T) {
	if got := toLiveCfgSnapshot(&liveentity.LiveCfg{AudienceListRefreshSeconds: 60}).AudienceListRefreshSeconds; got != 60 {
		t.Fatalf("stored refresh seconds = %d, want 60", got)
	}
}

func TestToLiveCfgSnapshotMaxAnchorPhotosDefault(t *testing.T) {
	if toLiveCfgSnapshot(nil).MaxAnchorPhotos != DefaultMaxAnchorPhotos {
		t.Fatal("missing live config must default max anchor photos to 100")
	}
	if toLiveCfgSnapshot(&liveentity.LiveCfg{}).MaxAnchorPhotos != DefaultMaxAnchorPhotos {
		t.Fatal("zero stored max anchor photos must default to 100")
	}
}

func TestToLiveCfgSnapshotMaxAnchorPhotosValue(t *testing.T) {
	if got := toLiveCfgSnapshot(&liveentity.LiveCfg{MaxAnchorPhotos: 50}).MaxAnchorPhotos; got != 50 {
		t.Fatalf("stored max anchor photos = %d, want 50", got)
	}
}

func TestToLiveCfgSnapshotMaxAnchorVideosDefault(t *testing.T) {
	if toLiveCfgSnapshot(nil).MaxAnchorVideos != DefaultMaxAnchorVideos {
		t.Fatal("missing live config must default max anchor videos to 100")
	}
	if toLiveCfgSnapshot(&liveentity.LiveCfg{}).MaxAnchorVideos != DefaultMaxAnchorVideos {
		t.Fatal("zero stored max anchor videos must default to 100")
	}
}

func TestToLiveCfgSnapshotMaxAnchorVideoFileSizeMBDefault(t *testing.T) {
	if toLiveCfgSnapshot(nil).MaxAnchorVideoFileSizeMB != DefaultMaxAnchorVideoFileSizeMB {
		t.Fatal("missing live config must default max anchor video file size to 100MB")
	}
	if toLiveCfgSnapshot(&liveentity.LiveCfg{}).MaxAnchorVideoFileSizeMB != DefaultMaxAnchorVideoFileSizeMB {
		t.Fatal("zero max anchor video file size must default to 100MB")
	}
}

func TestToLiveCfgSnapshotMaxAnchorVideoFileSizeMBValue(t *testing.T) {
	if got := toLiveCfgSnapshot(&liveentity.LiveCfg{MaxAnchorVideoFileSizeMB: 50}).MaxAnchorVideoFileSizeMB; got != 50 {
		t.Fatalf("got %d, want 50", got)
	}
}

func TestNormalizeMaxAnchorVideoFileSizeMB(t *testing.T) {
	if NormalizeMaxAnchorVideoFileSizeMB(0) != DefaultMaxAnchorVideoFileSizeMB {
		t.Fatal("zero must normalize to default 100")
	}
	if NormalizeMaxAnchorVideoFileSizeMB(200) != 200 {
		t.Fatal("non-zero must keep value")
	}
}

func TestToLiveCfgSnapshotMaxAnchorVideosValue(t *testing.T) {
	if got := toLiveCfgSnapshot(&liveentity.LiveCfg{MaxAnchorVideos: 50}).MaxAnchorVideos; got != 50 {
		t.Fatalf("stored max anchor videos = %d, want 50", got)
	}
}
