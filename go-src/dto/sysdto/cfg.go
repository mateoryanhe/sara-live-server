package sysdto

import "github.com/gogf/gf/v2/frame/g"

type SysCfgReq struct {
	g.Meta `path:"/cfg" method:"get" summary:"获取系统配置" tags:"系统"`
}

type SysCfgResp struct {
	SysTime int64 `json:"sysTime"`
}
