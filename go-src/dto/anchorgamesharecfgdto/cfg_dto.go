package anchorgamesharecfgdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type AnchorGameShareCfgListReq struct {
	g.Meta `path:"/anchorGameShareCfgList" method:"post" summary:"查询游戏流水档位分佣配置" tags:"游戏流水档位分佣配置"`
	httpserver.CMSQueryReq
	SalaryType uint32 `json:"salaryType" v:"required|in:1,2#薪资类型不能为空|薪资类型无效" dc:"薪资类型(1有底薪 2无底薪)"`
}

type AnchorGameShareCfgItem struct {
	ID                     string  `json:"id"`
	SalaryType             uint32  `json:"salaryType"`
	Level                  uint32  `json:"level"`
	GameTotalGoldRevenue   float64 `json:"gameTotalGoldRevenue"`
	AnchorGameSharePercent float64 `json:"anchorGameSharePercent"`
	GuildGameSharePercent  float64 `json:"guildGameSharePercent"`
	CreatedAt              string  `json:"createdAt"`
	UpdatedAt              string  `json:"updatedAt"`
}

type CreateAnchorGameShareCfgReq struct {
	g.Meta                 `path:"/createAnchorGameShareCfg" method:"post" summary:"新增游戏流水档位分佣配置" tags:"游戏流水档位分佣配置"`
	SalaryType             uint32  `json:"salaryType" v:"required|in:1,2#薪资类型不能为空|薪资类型无效" dc:"薪资类型(1有底薪 2无底薪)"`
	Level                  uint32  `json:"level" v:"required|min:1#等级不能为空|等级必须大于0" dc:"等级"`
	GameTotalGoldRevenue   float64 `json:"gameTotalGoldRevenue" dc:"游戏总金币流水"`
	AnchorGameSharePercent float64 `json:"anchorGameSharePercent" dc:"主播游戏提成比(%)"`
	GuildGameSharePercent  float64 `json:"guildGameSharePercent" dc:"工会游戏提成比(%)"`
}

type CreateAnchorGameShareCfgRes struct {
	ID string `json:"id"`
}

type UpdateAnchorGameShareCfgReq struct {
	g.Meta                 `path:"/updateAnchorGameShareCfg" method:"post" summary:"修改游戏流水档位分佣配置" tags:"游戏流水档位分佣配置"`
	ID                     uint64  `json:"id" v:"required#配置ID不能为空" dc:"配置ID"`
	Level                  uint32  `json:"level" v:"required|min:1#等级不能为空|等级必须大于0" dc:"等级"`
	GameTotalGoldRevenue   float64 `json:"gameTotalGoldRevenue" dc:"游戏总金币流水"`
	AnchorGameSharePercent float64 `json:"anchorGameSharePercent" dc:"主播游戏提成比(%)"`
	GuildGameSharePercent  float64 `json:"guildGameSharePercent" dc:"工会游戏提成比(%)"`
}

type UpdateAnchorGameShareCfgRes struct {
	Success bool `json:"success"`
}

type DeleteAnchorGameShareCfgReq struct {
	g.Meta `path:"/deleteAnchorGameShareCfg" method:"post" summary:"删除游戏流水档位分佣配置" tags:"游戏流水档位分佣配置"`
	ID     uint64 `json:"id" v:"required#配置ID不能为空" dc:"配置ID"`
}

type DeleteAnchorGameShareCfgRes struct {
	Success bool `json:"success"`
}
