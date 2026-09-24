package anchornosalarysharecfgdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type AnchorNoSalaryShareCfgListReq struct {
	g.Meta `path:"/anchorNoSalaryShareCfgList" method:"post" summary:"获取无底薪社交流水分佣配置" tags:"无底薪社交流水分佣配置"`
	httpserver.CMSQueryReq
}

type AnchorNoSalaryShareCfgItem struct {
	ID                        string  `json:"id"`
	Level                     uint32  `json:"level"`
	SocialTotalDiamondRevenue float64 `json:"socialTotalDiamondRevenue"`
	AnchorSocialSharePercent  float64 `json:"anchorSocialSharePercent"`
	GuildSocialSharePercent   float64 `json:"guildSocialSharePercent"`
	CreatedAt                 string  `json:"createdAt"`
	UpdatedAt                 string  `json:"updatedAt"`
}

type CreateAnchorNoSalaryShareCfgReq struct {
	g.Meta                    `path:"/createAnchorNoSalaryShareCfg" method:"post" summary:"新增无底薪社交流水分佣配置" tags:"无底薪社交流水分佣配置"`
	Level                     uint32  `json:"level" v:"required|min:1#等级不能为空|等级必须大于0" dc:"等级"`
	SocialTotalDiamondRevenue float64 `json:"socialTotalDiamondRevenue" dc:"社交总钻石流水"`
	AnchorSocialSharePercent  float64 `json:"anchorSocialSharePercent" dc:"主播社交提成比(%)"`
	GuildSocialSharePercent   float64 `json:"guildSocialSharePercent" dc:"工会社交提成比(%)"`
}

type CreateAnchorNoSalaryShareCfgRes struct {
	ID string `json:"id"`
}

type UpdateAnchorNoSalaryShareCfgReq struct {
	g.Meta                    `path:"/updateAnchorNoSalaryShareCfg" method:"post" summary:"修改无底薪社交流水分佣配置" tags:"无底薪社交流水分佣配置"`
	ID                        uint64  `json:"id" v:"required#配置ID不能为空" dc:"配置ID"`
	Level                     uint32  `json:"level" v:"required|min:1#等级不能为空|等级必须大于0" dc:"等级"`
	SocialTotalDiamondRevenue float64 `json:"socialTotalDiamondRevenue" dc:"社交总钻石流水"`
	AnchorSocialSharePercent  float64 `json:"anchorSocialSharePercent" dc:"主播社交提成比(%)"`
	GuildSocialSharePercent   float64 `json:"guildSocialSharePercent" dc:"工会社交提成比(%)"`
}

type UpdateAnchorNoSalaryShareCfgRes struct {
	Success bool `json:"success"`
}

type DeleteAnchorNoSalaryShareCfgReq struct {
	g.Meta `path:"/deleteAnchorNoSalaryShareCfg" method:"post" summary:"删除无底薪社交流水分佣配置" tags:"无底薪社交流水分佣配置"`
	ID     uint64 `json:"id" v:"required#配置ID不能为空" dc:"配置ID"`
}

type DeleteAnchorNoSalaryShareCfgRes struct {
	Success bool `json:"success"`
}
