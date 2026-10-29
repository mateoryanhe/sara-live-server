package privacypolicy

import (
	"strings"
	"sync/atomic"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity/user"
)

type sofieCfgSnapshot struct {
	PackageName       string
	PrivacyPolicyUrl  string
	TermsOfServiceUrl string
	CreatorTermsUrl   string
	RoomOwnerTermsUrl string
	VipDescUrl        string
	AboutSiteUrl      string
	SafetyCenterUrl   string
}

var (
	sofieCfgCache         atomic.Value // *sofieCfgSnapshot
	emptySofieCfgSnapshot = &sofieCfgSnapshot{}
)

// ReloadSofieCfgMemory 从 DB 刷新 Sofie 隐私政策内存(数据同步接收等场景).
func ReloadSofieCfgMemory() {
	reloadSofieCfgMemory()
}

func reloadSofieCfgMemory() {
	sofieCfgCache.Store(toSofieCfgSnapshot(cfgdao.LoadSofiePrivacyPolicyCfg()))
}

func getSofieCfgCache() *sofieCfgSnapshot {
	v := sofieCfgCache.Load()
	if v == nil {
		return emptySofieCfgSnapshot
	}
	cfg, ok := v.(*sofieCfgSnapshot)
	if !ok || cfg == nil {
		return emptySofieCfgSnapshot
	}
	return cfg
}

func toSofieCfgSnapshot(row *entity.SofiePrivacyPolicyCfg) *sofieCfgSnapshot {
	if row == nil {
		return emptySofieCfgSnapshot
	}
	return &sofieCfgSnapshot{
		PackageName:       strings.TrimSpace(row.PackageName),
		PrivacyPolicyUrl:  row.PrivacyPolicyUrl,
		TermsOfServiceUrl: row.TermsOfServiceUrl,
		CreatorTermsUrl:   row.CreatorTermsUrl,
		RoomOwnerTermsUrl: row.RoomOwnerTermsUrl,
		VipDescUrl:        row.VipDescUrl,
		AboutSiteUrl:      row.AboutSiteUrl,
		SafetyCenterUrl:   row.SafetyCenterUrl,
	}
}

// MatchSofiePackage 请求包名是否命中 Sofie 隐私配置(双方均非空且相等).
func MatchSofiePackage(packageName string) bool {
	cfgPkg := getSofieCfgCache().PackageName
	reqPkg := strings.TrimSpace(packageName)
	return cfgPkg != "" && reqPkg != "" && cfgPkg == reqPkg
}

// UseSofiePrivacyForPackage Sofie 路径且包名命中时使用 Sofie 配置,否则回落 Sara 隐私政策.
func UseSofiePrivacyForPackage(sofie bool, packageName string) bool {
	return sofie && MatchSofiePackage(packageName)
}

func GetSofiePrivacyPolicyUrl() string {
	return getSofieCfgCache().PrivacyPolicyUrl
}

func GetSofieTermsOfServiceUrl() string {
	return getSofieCfgCache().TermsOfServiceUrl
}

func GetSofieCreatorTermsUrl() string {
	return getSofieCfgCache().CreatorTermsUrl
}

func GetSofieRoomOwnerTermsUrl() string {
	return getSofieCfgCache().RoomOwnerTermsUrl
}

func GetSofieVipDescUrl() string {
	return getSofieCfgCache().VipDescUrl
}

func GetSofieAboutSiteUrl() string {
	return getSofieCfgCache().AboutSiteUrl
}

func GetSofieSafetyCenterUrl() string {
	return getSofieCfgCache().SafetyCenterUrl
}

func GetSofiePackageName() string {
	return getSofieCfgCache().PackageName
}
