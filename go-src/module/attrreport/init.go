package attrreport

import "xr-game-server/dao/cfgdao"

func Init() {
	cfgdao.InitAttrReportCfgDao()
	cfgdao.ReloadAttrReportCfgCache()
	initEvents()
}
