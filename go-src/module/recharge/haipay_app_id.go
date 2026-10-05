package recharge

import (
	"xr-game-server/constants/country"
	"xr-game-server/dao/cfgdao"
)

func lookupHaiPayAppID(businessCode string) (int64, bool) {
	return country.LookupHaiPayAppID(businessCode, cfgdao.HaiPayUseProdAppID())
}

func haiPayCashierAppID() int64 {
	return country.HaiPayCashierAppID(cfgdao.HaiPayUseProdAppID())
}
