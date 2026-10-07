package recharge

import (
	"testing"

	rechargeentity "xr-game-server/entity/recharge"
)

func TestAppRechargeCfgListCfgTypeForClient(t *testing.T) {
	if got := appRechargeCfgListCfgTypeForClient(true); got != rechargeentity.RechargeCfgTypeChannel {
		t.Fatalf("h5 got cfgType=%d want channel=%d", got, rechargeentity.RechargeCfgTypeChannel)
	}
	if got := appRechargeCfgListCfgTypeForClient(false); got != rechargeentity.RechargeCfgTypeGoogle {
		t.Fatalf("app got cfgType=%d want google=%d", got, rechargeentity.RechargeCfgTypeGoogle)
	}
}
