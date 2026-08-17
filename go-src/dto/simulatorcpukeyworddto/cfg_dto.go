package simulatorcpukeyworddto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type SimulatorCpuKeywordListReq struct {
	g.Meta `path:"/simulatorCpuKeywordList" method:"post" summary:"获取模拟器CPU关键词列表" tags:"模拟器CPU拦截"`
	httpserver.CMSQueryReq
	Key string `json:"key" dc:"关键词搜索"`
}

type SimulatorCpuKeywordItem struct {
	ID        string `json:"id"`
	Keyword   string `json:"keyword"`
	Remark    string `json:"remark"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type CreateSimulatorCpuKeywordReq struct {
	g.Meta  `path:"/createSimulatorCpuKeyword" method:"post" summary:"新增模拟器CPU关键词" tags:"模拟器CPU拦截"`
	Keyword string `json:"keyword" v:"required#关键词不能为空" dc:"CPU关键词(模糊匹配)"`
	Remark  string `json:"remark" dc:"备注"`
}

type CreateSimulatorCpuKeywordRes struct {
	ID string `json:"id"`
}

type UpdateSimulatorCpuKeywordReq struct {
	g.Meta  `path:"/updateSimulatorCpuKeyword" method:"post" summary:"修改模拟器CPU关键词" tags:"模拟器CPU拦截"`
	ID      uint64 `json:"id" v:"required#配置ID不能为空" dc:"配置ID"`
	Keyword string `json:"keyword" v:"required#关键词不能为空" dc:"CPU关键词(模糊匹配)"`
	Remark  string `json:"remark" dc:"备注"`
}

type UpdateSimulatorCpuKeywordRes struct {
	Success bool `json:"success"`
}

type DeleteSimulatorCpuKeywordReq struct {
	g.Meta `path:"/deleteSimulatorCpuKeyword" method:"post" summary:"删除模拟器CPU关键词" tags:"模拟器CPU拦截"`
	ID     uint64 `json:"id" v:"required#配置ID不能为空" dc:"配置ID"`
}

type DeleteSimulatorCpuKeywordRes struct {
	Success bool `json:"success"`
}
