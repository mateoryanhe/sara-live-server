package globalcfgdto

import "github.com/gogf/gf/v2/frame/g"

type GetGlobalCfgReq struct {
	g.Meta `path:"/getGlobalCfg" method:"post,get" summary:"获取全局配置" tags:"全局配置"`
	Module string `json:"module" summary:"模块"`
}
