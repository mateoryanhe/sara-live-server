package cfg

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

const IpGeoStr = "ipGeo"

// IpGeoCfg IP 地理定位配置
type IpGeoCfg struct {
	// DbPath GeoLite2-Country.mmdb 绝对路径,留空则不启用
	// 免注册下载: https://cdn.jsdelivr.net/npm/geolite2-country/GeoLite2-Country.mmdb.gz
	DbPath string
}

var IpGeoCfgVar = &IpGeoCfg{}

func initIpGeoCfg() {
	data, _ := g.Cfg().GetWithCmd(gctx.New(), IpGeoStr)
	if err := data.Scan(IpGeoCfgVar); err != nil {
		g.Log().Error(gctx.New(), "无法加载 ipGeo 配置")
	}
}
