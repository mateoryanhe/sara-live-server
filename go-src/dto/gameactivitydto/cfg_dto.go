package gameactivitydto

import "github.com/gogf/gf/v2/frame/g"

type GetGameActivityCfgReq struct {
	g.Meta `path:"/getGameActivityCfg" method:"post" summary:"查询游戏活动配置" tags:"游戏活动配置"`
}

type GameActivityCfgItem struct {
	ID          string `json:"id"`
	ButtonIcon  string `json:"buttonIcon" dc:"游戏按钮icon URL"`
	GameListUrl string `json:"gameListUrl" dc:"游戏列表URL"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type GetGameActivityCfgRes struct {
	Cfg *GameActivityCfgItem `json:"cfg"`
}

type SaveGameActivityCfgReq struct {
	g.Meta      `path:"/saveGameActivityCfg" method:"post" summary:"保存游戏活动配置" tags:"游戏活动配置"`
	ID          uint64 `json:"id" dc:"配置ID,新建传0"`
	ButtonIcon  string `json:"buttonIcon" v:"max-length:512#按钮icon最长512字符" dc:"游戏按钮icon URL"`
	GameListUrl string `json:"gameListUrl" v:"max-length:1024#游戏列表URL最长1024字符" dc:"游戏列表URL"`
}

type SaveGameActivityCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
