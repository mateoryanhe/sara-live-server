package staticcachecfgdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type StaticCacheRuleListReq struct {
	g.Meta `path:"/staticCacheRuleList" method:"post" summary:"查询静态网页缓存规则" tags:"静态网页缓存配置"`
	httpserver.CMSQueryReq
	Key string `json:"key" dc:"文件名或备注"`
}

type StaticCacheRuleItem struct {
	ID        string `json:"id"`
	FileName  string `json:"fileName"`
	Remark    string `json:"remark"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type CreateStaticCacheRuleReq struct {
	g.Meta   `path:"/createStaticCacheRule" method:"post" summary:"新增静态网页缓存规则" tags:"静态网页缓存配置"`
	FileName string `json:"fileName" v:"required|max-length:255#文件名不能为空|文件名最长255字符" dc:"不缓存文件名,不含目录"`
	Remark   string `json:"remark" v:"max-length:255#备注最长255字符" dc:"备注"`
}

type CreateStaticCacheRuleRes struct {
	ID string `json:"id"`
}

type UpdateStaticCacheRuleReq struct {
	g.Meta   `path:"/updateStaticCacheRule" method:"post" summary:"修改静态网页缓存规则" tags:"静态网页缓存配置"`
	ID       uint64 `json:"id" v:"required#配置ID不能为空" dc:"配置ID"`
	FileName string `json:"fileName" v:"required|max-length:255#文件名不能为空|文件名最长255字符" dc:"不缓存文件名,不含目录"`
	Remark   string `json:"remark" v:"max-length:255#备注最长255字符" dc:"备注"`
}

type UpdateStaticCacheRuleRes struct {
	Success bool `json:"success"`
}

type DeleteStaticCacheRuleReq struct {
	g.Meta `path:"/deleteStaticCacheRule" method:"post" summary:"删除静态网页缓存规则" tags:"静态网页缓存配置"`
	ID     uint64 `json:"id" v:"required#配置ID不能为空" dc:"配置ID"`
}

type DeleteStaticCacheRuleRes struct {
	Success bool `json:"success"`
}
