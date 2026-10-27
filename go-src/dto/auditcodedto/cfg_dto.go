package auditcodedto

import "github.com/gogf/gf/v2/frame/g"

type AuditCodeCfgItem struct {
	ID        string `json:"id"`
	Code      string `json:"code"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type GetAuditCodeCfgReq struct {
	g.Meta `path:"/getAuditCodeCfg" method:"post" summary:"查询当前审核码" tags:"审核码"`
}

type GetAuditCodeCfgRes struct {
	Cfg *AuditCodeCfgItem `json:"cfg"`
}

type GenerateAuditCodeCfgReq struct {
	g.Meta `path:"/generateAuditCodeCfg" method:"post" summary:"生成审核码并覆盖入库" tags:"审核码"`
	Length int    `json:"length" dc:"审核码位数,默认8,范围4-32"`
}

type GenerateAuditCodeCfgRes struct {
	Success bool              `json:"success"`
	Cfg     *AuditCodeCfgItem `json:"cfg"`
}
