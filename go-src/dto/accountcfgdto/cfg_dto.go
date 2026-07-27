package accountcfgdto

import "github.com/gogf/gf/v2/frame/g"

type GetAccountCfgReq struct {
	g.Meta `path:"/getAccountCfg" method:"post" summary:"查询账号配置" tags:"账号配置"`
}

type AccountCfgItem struct {
	ID                         string `json:"id"`
	CancelAccountByCodeEnabled bool   `json:"cancelAccountByCodeEnabled"`
	CreatedAt                  string `json:"createdAt"`
	UpdatedAt                  string `json:"updatedAt"`
}

type GetAccountCfgRes struct {
	Cfg *AccountCfgItem `json:"cfg"`
}

type SaveAccountCfgReq struct {
	g.Meta                     `path:"/saveAccountCfg" method:"post" summary:"保存账号配置" tags:"账号配置"`
	ID                         uint64 `json:"id" dc:"配置ID,首次保存可为0"`
	CancelAccountByCodeEnabled bool   `json:"cancelAccountByCodeEnabled" dc:"注销码销户开关(官网公开接口)"`
}

type SaveAccountCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
