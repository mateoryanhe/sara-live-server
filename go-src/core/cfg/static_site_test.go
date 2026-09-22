package cfg

import (
	"path/filepath"
	"testing"
)

func TestFindThirdPayDomain(t *testing.T) {
	sites := []*StaticSiteCfg{
		{Domain: "static.bigtktool.shop"},
		{Domain: " third-pay.bigtktool.shop,third-pay-backup.bigtktool.shop ", T: true},
	}

	if got, want := findThirdPayDomain(sites), "third-pay.bigtktool.shop"; got != want {
		t.Fatalf("findThirdPayDomain() = %q, want %q", got, want)
	}
}

func TestFindThirdPayDomainWithoutMarker(t *testing.T) {
	if got := findThirdPayDomain([]*StaticSiteCfg{{Domain: "third-pay.bigtktool.shop"}}); got != "" {
		t.Fatalf("findThirdPayDomain() = %q, want empty domain", got)
	}
}

func TestRuntimeStaticSiteOverridesDomainAndPathByPrefix(t *testing.T) {
	runtimeStaticMu.Lock()
	originalSites := runtimeStaticSites
	originalPaths := runtimeStaticPaths
	originalDomains := runtimeDomainSites
	runtimeStaticSites = make(map[string]*StaticSiteCfg)
	runtimeStaticPaths = nil
	runtimeDomainSites = nil
	runtimeStaticMu.Unlock()
	t.Cleanup(func() {
		runtimeStaticMu.Lock()
		runtimeStaticSites = originalSites
		runtimeStaticPaths = originalPaths
		runtimeDomainSites = originalDomains
		runtimeStaticMu.Unlock()
	})

	root := filepath.Join(t.TempDir(), "h5-live")
	RegisterRuntimeStaticSite("/h5-live", "new-h5.example.com", root)
	if got := GetStaticSiteDomain("/h5-live"); got != "new-h5.example.com" {
		t.Fatalf("GetStaticSiteDomain() = %q, want new-h5.example.com", got)
	}
	if got := GetStaticPathRoot("/h5-live"); got != root {
		t.Fatalf("GetStaticPathRoot() = %q, want %q", got, root)
	}

	// 上传资源配置刷新只清自己的临时路径/域名，不应清除站点覆盖。
	ClearRuntimeStaticMappings()
	if got := GetStaticSiteDomain("/h5-live"); got != "new-h5.example.com" {
		t.Fatalf("domain after ClearRuntimeStaticMappings() = %q", got)
	}
}

func TestRuntimeThirdPaySiteKeepsAppDomainMarker(t *testing.T) {
	runtimeStaticMu.Lock()
	originalRuntimeSites := runtimeStaticSites
	originalStaticSites := staticSiteCfgs
	runtimeStaticSites = make(map[string]*StaticSiteCfg)
	staticSiteCfgs = []*StaticSiteCfg{{
		Domain: "old-third-pay.example.com",
		Prefix: "/third-pay",
		Path:   "/old/third-pay",
		Root:   "/old/third-pay",
		T:      true,
	}}
	runtimeStaticMu.Unlock()
	t.Cleanup(func() {
		runtimeStaticMu.Lock()
		runtimeStaticSites = originalRuntimeSites
		staticSiteCfgs = originalStaticSites
		runtimeStaticMu.Unlock()
	})

	RegisterRuntimeStaticSite("/third-pay", "new-third-pay.example.com", "/new/third-pay")
	if got := GetThirdPayDomain(); got != "new-third-pay.example.com" {
		t.Fatalf("GetThirdPayDomain() = %q, want new-third-pay.example.com", got)
	}
}
