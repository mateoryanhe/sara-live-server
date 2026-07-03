package sysdto

import "github.com/gogf/gf/v2/frame/g"

type SysCfgReq struct {
	g.Meta `path:"/cfg" method:"get" summary:"获取系统配置" tags:"系统"`
}

type SysCfgResp struct {
	SysTime                     int64   `json:"sysTime"`
	PaidDanmakuPrice            float64 `json:"paidDanmakuPrice" dc:"直播间付费弹幕价格(钻石)"`
	PrivateRoomFreeWatchSeconds uint32  `json:"privateRoomFreeWatchSeconds" dc:"私密直播间免费观看时长(秒)"`
}
