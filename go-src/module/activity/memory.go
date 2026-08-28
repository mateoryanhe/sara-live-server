package activity

import (
	"sync/atomic"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/activitydto"
	activityentity "xr-game-server/entity/activity"
	"xr-game-server/module/upload"
)

// cfgSnapshot 首充活动配置内存快照(App/CMS 查询均读此缓存,保存后 reload)
type cfgSnapshot struct {
	ID                uint64
	Enabled           bool
	IconName          string
	Icon              string
	TitleEn           string
	TitleEs           string
	TitlePt           string
	TitleHi           string
	TitleId           string
	RechargeBtnTextEn string
	RechargeBtnTextEs string
	RechargeBtnTextPt string
	RechargeBtnTextHi string
	RechargeBtnTextId string
	FirstRechargeRatio float64
	Privileges        []*activityentity.FirstRechargeActivityPrivilege
	CreatedAt         string
	UpdatedAt         string
}

var cfgCache atomic.Value

func reloadCfgMemory() {
	row := cfgdao.LoadFirstRechargeActivityCfg()
	privileges := cfgdao.LoadFirstRechargeActivityPrivileges()
	cfgCache.Store(toCfgSnapshot(row, privileges))
}

// ReloadFirstRechargeActivityCache 从 DB 重新加载缓存(供启动/保存后调用)
func ReloadFirstRechargeActivityCache() {
	reloadCfgMemory()
}

func getCfgCache() *cfgSnapshot {
	if cfgCache.Load() == nil {
		reloadCfgMemory()
	}
	v := cfgCache.Load()
	if v == nil {
		return &cfgSnapshot{}
	}
	snap, ok := v.(*cfgSnapshot)
	if !ok || snap == nil {
		return &cfgSnapshot{}
	}
	return snap
}

func toCfgSnapshot(row *activityentity.FirstRechargeActivityCfg, privileges []*activityentity.FirstRechargeActivityPrivilege) *cfgSnapshot {
	if row == nil {
		return &cfgSnapshot{}
	}
	return &cfgSnapshot{
		ID:                row.ID,
		Enabled:           row.Enabled,
		IconName:          row.Icon,
		Icon:              upload.GetUrlByName(row.Icon),
		TitleEn:           row.TitleEn,
		TitleEs:           row.TitleEs,
		TitlePt:           row.TitlePt,
		TitleHi:           row.TitleHi,
		TitleId:           row.TitleId,
		RechargeBtnTextEn: row.RechargeBtnTextEn,
		RechargeBtnTextEs: row.RechargeBtnTextEs,
		RechargeBtnTextPt: row.RechargeBtnTextPt,
		RechargeBtnTextHi: row.RechargeBtnTextHi,
		RechargeBtnTextId: row.RechargeBtnTextId,
		FirstRechargeRatio: normalizeFirstRechargeRatio(row.FirstRechargeRatio),
		Privileges:        privileges,
		CreatedAt:         formatTime(row.CreatedAt),
		UpdatedAt:         formatTime(row.UpdatedAt),
	}
}

func toCfgItemFromSnapshot(snap *cfgSnapshot) *activitydto.FirstRechargeActivityCfgItem {
	if snap == nil || snap.ID == 0 {
		return defaultCfgItem()
	}
	return &activitydto.FirstRechargeActivityCfgItem{
		ID:                formatUintID(snap.ID),
		Enabled:           snap.Enabled,
		IconName:          snap.IconName,
		Icon:              snap.Icon,
		TitleEn:           snap.TitleEn,
		TitleEs:           snap.TitleEs,
		TitlePt:           snap.TitlePt,
		TitleHi:           snap.TitleHi,
		TitleId:           snap.TitleId,
		RechargeBtnTextEn: snap.RechargeBtnTextEn,
		RechargeBtnTextEs: snap.RechargeBtnTextEs,
		RechargeBtnTextPt: snap.RechargeBtnTextPt,
		RechargeBtnTextHi: snap.RechargeBtnTextHi,
		RechargeBtnTextId: snap.RechargeBtnTextId,
		FirstRechargeRatio: snap.FirstRechargeRatio,
		Privileges:        toPrivilegeDTOItems(snap.Privileges),
		CreatedAt:         snap.CreatedAt,
		UpdatedAt:         snap.UpdatedAt,
	}
}

func toAppRes(snap *cfgSnapshot) *activitydto.AppFirstRechargeActivityCfgRes {
	if snap == nil || snap.ID == 0 {
		return &activitydto.AppFirstRechargeActivityCfgRes{
			FirstRechargeRatio: defaultFirstRechargeRatio,
			Privileges:        []*activitydto.AppFirstRechargePrivilegeItem{},
		}
	}
	return &activitydto.AppFirstRechargeActivityCfgRes{
		Enabled:            snap.Enabled,
		Icon:               snap.Icon,
		TitleEn:            snap.TitleEn,
		TitleEs:            snap.TitleEs,
		TitlePt:            snap.TitlePt,
		TitleHi:            snap.TitleHi,
		TitleId:            snap.TitleId,
		RechargeBtnTextEn:  snap.RechargeBtnTextEn,
		RechargeBtnTextEs:  snap.RechargeBtnTextEs,
		RechargeBtnTextPt:  snap.RechargeBtnTextPt,
		RechargeBtnTextHi:  snap.RechargeBtnTextHi,
		RechargeBtnTextId:  snap.RechargeBtnTextId,
		FirstRechargeRatio: snap.FirstRechargeRatio,
		Privileges:         toPrivilegeAppItems(snap.Privileges),
	}
}
