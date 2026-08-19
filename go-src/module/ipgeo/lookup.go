package ipgeo

import (
	"net"
	"net/netip"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/oschwald/geoip2-golang/v2"
	"xr-game-server/core/xrlog"
	"xr-game-server/dao/cfgdao"
)

var (
	dbMu sync.RWMutex
	db   *geoip2.Reader
)

// Init 加载 IP 国家数据库(GeoLite2-Country.mmdb,MMDB 格式).
func Init() {
	path := strings.TrimSpace(cfgdao.GetIpGeoDbPath())
	if path == "" {
		xrlog.DetailLog.Warning(gctx.New(), "ipGeo.dbPath 未配置,IP 国家解析不可用")
		return
	}
	reader, err := geoip2.Open(path)
	if err != nil {
		xrlog.DetailLog.Warningf(gctx.New(), "ipGeo 数据库加载失败,path=%s,err=%v", path, err)
		return
	}
	dbMu.Lock()
	if db != nil {
		_ = db.Close()
	}
	db = reader
	dbMu.Unlock()
	xrlog.DetailLog.Infof(gctx.New(), "ipGeo 数据库加载成功,path=%s", path)
}

// Enabled 是否已加载 IP 地理库.
func Enabled() bool {
	dbMu.RLock()
	defer dbMu.RUnlock()
	return db != nil
}

// Lookup 根据 IP 查询国家信息;私网/无效 IP 或库未加载时返回 nil.
func Lookup(ip string) *CountryInfo {
	addr, ok := parseLookupAddr(ip)
	if !ok {
		return nil
	}
	reader := getDB()
	if reader == nil {
		return nil
	}
	record, err := reader.Country(addr)
	if err != nil || record == nil || !record.HasData() {
		return nil
	}
	code := strings.TrimSpace(record.Country.ISOCode)
	if code == "" {
		return nil
	}
	return &CountryInfo{
		Code: code,
		Name: pickCountryName(record.Country.Names),
	}
}

// GetCountryCode 返回 ISO 国家码,无法解析时返回空字符串.
func GetCountryCode(ip string) string {
	info := Lookup(ip)
	if info == nil {
		return ""
	}
	return info.Code
}

// GetCountryName 返回国家名称,无法解析时返回空字符串.
func GetCountryName(ip string) string {
	info := Lookup(ip)
	if info == nil {
		return ""
	}
	return info.Name
}

func getDB() *geoip2.Reader {
	dbMu.RLock()
	defer dbMu.RUnlock()
	return db
}

func parseLookupAddr(raw string) (netip.Addr, bool) {
	ipStr := strings.TrimSpace(raw)
	if ipStr == "" {
		return netip.Addr{}, false
	}
	if host, _, err := net.SplitHostPort(ipStr); err == nil {
		ipStr = host
	}
	addr, err := netip.ParseAddr(ipStr)
	if err != nil || !addr.IsValid() || isNonPublicAddr(addr) {
		return netip.Addr{}, false
	}
	return addr, true
}

func isNonPublicAddr(addr netip.Addr) bool {
	if !addr.IsValid() || addr.IsUnspecified() || addr.IsLoopback() || addr.IsPrivate() || addr.IsLinkLocalUnicast() || addr.IsLinkLocalMulticast() {
		return true
	}
	return false
}

func pickCountryName(names geoip2.Names) string {
	if name := strings.TrimSpace(names.SimplifiedChinese); name != "" {
		return name
	}
	if name := strings.TrimSpace(names.English); name != "" {
		return name
	}
	return ""
}
