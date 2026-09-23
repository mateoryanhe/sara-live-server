package anchorsalarysocialsharecfgdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type AnchorSalarySocialShareCfgListReq struct {
	g.Meta `path:"/anchorSalarySocialShareCfgList" method:"post" summary:"获取有底薪社交流水分佣配置" tags:"有底薪社交流水分佣配置"`
	httpserver.CMSQueryReq
}

type AnchorSalarySocialShareCfgItem struct {
	ID                        string  `json:"id"`
	Level                     uint32  `json:"level"`
	SocialTotalDiamondRevenue float64 `json:"socialTotalDiamondRevenue"`
	AnchorSocialSharePercent  float64 `json:"anchorSocialSharePercent"`
	GuildSocialSharePercent   float64 `json:"guildSocialSharePercent"`
	CreatedAt                 string  `json:"createdAt"`
	UpdatedAt                 string  `json:"updatedAt"`
}

type CreateAnchorSalarySocialShareCfgReq struct {
	g.Meta                    `path:"/createAnchorSalarySocialShareCfg" method:"post" summary:"新增有底薪社交流水分佣配置" tags:"有底薪社交流水分佣配置"`
	Level                     uint32  `json:"level" v:"required|min:1#等级不能为空|等级必须大于0" dc:"等级"`
	SocialTotalDiamondRevenue float64 `json:"socialTotalDiamondRevenue" dc:"社交总钻石流水"`
	AnchorSocialSharePercent  float64 `json:"anchorSocialSharePercent" dc:"主播社交提成比(%)"`
	GuildSocialSharePercent   float64 `json:"guildSocialSharePercent" dc:"工会社交提成比(%)"`
}

type CreateAnchorSalarySocialShareCfgRes struct {
	ID string `json:"id"`
}

type UpdateAnchorSalarySocialShareCfgReq struct {
	g.Meta                    `path:"/updateAnchorSalarySocialShareCfg" method:"post" summary:"修改有底薪社交流水分佣配置" tags:"有底薪社交流水分佣配置"`
	ID                        uint64  `json:"id" v:"required#配置ID不能为空" dc:"配置ID"`
	Level                     uint32  `json:"level" v:"required|min:1#等级不能为空|等级必须大于0" dc:"等级"`
	SocialTotalDiamondRevenue float64 `json:"socialTotalDiamondRevenue" dc:"社交总钻石流水"`
	AnchorSocialSharePercent  float64 `json:"anchorSocialSharePercent" dc:"主播社交提成比(%)"`
	GuildSocialSharePercent   float64 `json:"guildSocialSharePercent" dc:"工会社交提成比(%)"`
}

type UpdateAnchorSalarySocialShareCfgRes struct {
	Success bool `json:"success"`
}

type DeleteAnchorSalarySocialShareCfgReq struct {
	g.Meta `path:"/deleteAnchorSalarySocialShareCfg" method:"post" summary:"删除有底薪社交流水分佣配置" tags:"有底薪社交流水分佣配置"`
	ID     uint64 `json:"id" v:"required#配置ID不能为空" dc:"配置ID"`
}

type DeleteAnchorSalarySocialShareCfgRes struct {
	Success bool `json:"success"`
}
