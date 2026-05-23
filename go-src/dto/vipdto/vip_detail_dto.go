package vipdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/dto/vipcfgdto"
)

// AppVipDetailReq App端查询当前用户VIP详情
type AppVipDetailReq struct {
	g.Meta `path:"/getVipDetail" method:"post" summary:"查询VIP详情" tags:"VIP"`
}

type AppVipDetailRes struct {
	VipLevel            uint32                   `json:"vipLevel"`
	LevelName           string                   `json:"levelName"`
	TotalRecharge       float64                  `json:"totalRecharge,string" dc:"累计充值(USD,保留4位小数)"`
	IsMaxLevel          bool                     `json:"isMaxLevel"`
	NextLevel           uint32                   `json:"nextLevel" dc:"下一VIP等级(0表示已是最高级或无下一级)"`
	NextUpgradeLimit    float64                  `json:"nextUpgradeLimit,string" dc:"下一级升级门槛(USD,保留4位小数)"`
	RechargeToNextLevel float64                  `json:"rechargeToNextLevel,string" dc:"距离下一级还需充值(USD,保留4位小数)"`
	CurrentCfg          *vipcfgdto.AppVipCfgItem `json:"currentCfg" dc:"当前等级权益配置"`
	NextCfg             *vipcfgdto.AppVipCfgItem `json:"nextCfg" dc:"下一等级权益配置"`
}

// VipLevelPushItem VIP等级变更推送载荷
type VipLevelPushItem struct {
	VipLevel  string `json:"vipLevel" dc:"最新VIP等级"`
	LevelName string `json:"levelName" dc:"等级名称"`
}
