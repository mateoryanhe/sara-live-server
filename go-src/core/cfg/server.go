package cfg

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
)

const (
	ServerCfg                      = "server"
	defaultClientMaxBodySize int64 = 8 * 1024 * 1024
)

// XRServerCfg 服务器配置
type XRServerCfg struct {
	//http端口
	HttpPort int
	//静态资源配置路径
	StaticResPath string
	//证书路径
	CertFile string
	KeyFile  string
	//服务器id
	Id int64
	//程序最长退出时间
	ExitTime int
	// Go堆内存软上限(MB),0或未设置时默认400M
	MemoryLimitM int
}

var serverCfg *XRServerCfg

// GetServerCfg 获取服务器配置
func GetServerCfg() *XRServerCfg {
	return serverCfg
}

func initServerCfg() {
	data, _ := g.Cfg().GetWithCmd(gctx.New(), ServerCfg)
	err := data.Scan(&serverCfg)
	if err != nil {
		g.Log().Error(gctx.New(), "无法加载到基础配置文件数据")
		return
	}
	applyMemoryLimit(serverCfg.MemoryLimitM)
}

// GetClientMaxBodySize 读取 server.clientMaxBodySize,与 GoFrame HTTP 请求体上限一致
func GetClientMaxBodySize() int64 {
	raw := strings.TrimSpace(g.Cfg().MustGet(gctx.New(), "server.clientMaxBodySize").String())
	if raw == "" {
		return defaultClientMaxBodySize
	}
	size := gfile.StrToSize(raw)
	if size <= 0 {
		return defaultClientMaxBodySize
	}
	return size
}
