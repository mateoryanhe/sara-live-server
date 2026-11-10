package simulatorblacklistlogdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type CMSListSimulatorBlacklistLogReq struct {
	g.Meta `path:"/listSimulatorBlacklistLog" method:"post" summary:"CMS模拟器黑名单日志列表" tags:"模拟器黑名单日志"`
	httpserver.CMSQueryReq
	Key       string `json:"key" dc:"设备码/IP/邮箱/CPU模糊"`
	LoginType string `json:"loginType" dc:"登录方式:device/email,空=全部"`
	StartTime int64  `json:"startTime" dc:"拦截时间起(unix秒)"`
	EndTime   int64  `json:"endTime" dc:"拦截时间止(unix秒)"`
}

type CMSSimulatorBlacklistLogItem struct {
	ID          uint64 `json:"id,string"`
	IP          string `json:"ip"`
	DeviceId    string `json:"deviceId"`
	CpuModel    string `json:"cpuModel"`
	DeviceType  string `json:"deviceType"`
	DeviceModel string `json:"deviceModel"`
	OsVersion   string `json:"osVersion"`
	AppVersion  string `json:"appVersion"`
	PackageName string `json:"packageName"`
	LoginType   string `json:"loginType"`
	Email       string `json:"email"`
	AccountId   string `json:"accountId"`
	RegisterAt  string `json:"registerAt"`
	CreatedAt   string `json:"createdAt"`
}
