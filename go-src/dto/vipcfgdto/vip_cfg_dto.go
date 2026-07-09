package vipcfgdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// VipCfgListReq CMS分页查询VIP配置
type VipCfgListReq struct {
	g.Meta `path:"/vipCfgList" method:"post" summary:"获取VIP配置列表" tags:"VIP配置"`
	httpserver.CMSQueryReq
	LevelName    string `json:"levelName" dc:"等级名称(模糊匹配)"`
	StatusFilter int    `json:"statusFilter" dc:"状态过滤(0=全部,1=只看关闭,2=只看开启)"`
}

// VipCfgListRes 列表项
type VipCfgListRes struct {
	ID                   string  `json:"id"`
	Level                uint32  `json:"level"`
	LevelName            string  `json:"levelName"`
	Status               uint8   `json:"status"`
	UpgradeRechargeLimit float64 `json:"upgradeRechargeLimit,string"`
	MinWithdrawAmount    float64 `json:"minWithdrawAmount,string"`
	MaxWithdrawAmount    float64 `json:"maxWithdrawAmount,string"`
	Fee                  float64 `json:"fee,string"`
	Animation            string  `json:"animation" dc:"动画完整URL(列表展示)"`
	AnimationName        string  `json:"animationName" dc:"动画资源文件名(编辑保存用)"`
	CreatedAt            string  `json:"createdAt"`
	UpdatedAt            string  `json:"updatedAt"`
}

// CreateVipCfgReq 创建VIP配置
type CreateVipCfgReq struct {
	g.Meta               `path:"/createVipCfg" method:"post" summary:"创建VIP配置" tags:"VIP配置"`
	Level                uint32  `json:"level" v:"required|min:1#等级不能为空|等级需大于0" dc:"VIP等级"`
	LevelName            string  `json:"levelName" v:"required|length:1,64#等级名称不能为空|等级名称长度需在1到64之间" dc:"等级名称"`
	Status               uint8   `json:"status" v:"in:0,1#状态无效" dc:"状态(0关闭,1开启)"`
	UpgradeRechargeLimit float64 `json:"upgradeRechargeLimit" dc:"升级充值上限(USD,保留4位小数)"`
	MinWithdrawAmount    float64 `json:"minWithdrawAmount" dc:"最低提现金额(USD,保留4位小数)"`
	MaxWithdrawAmount    float64 `json:"maxWithdrawAmount" dc:"最高提现金额(USD,保留4位小数)"`
	Fee                  float64 `json:"fee" dc:"手续费(保留4位小数)"`
	Animation            string  `json:"animation" dc:"进场特效动画资源文件名(mp4)"`
}

type CreateVipCfgRes struct {
	ID string `json:"id"`
}

// UpdateVipCfgReq 更新VIP配置
type UpdateVipCfgReq struct {
	g.Meta               `path:"/updateVipCfg" method:"post" summary:"修改VIP配置" tags:"VIP配置"`
	ID                   uint64  `json:"id" v:"required#ID不能为空" dc:"配置ID"`
	Level                uint32  `json:"level" v:"required|min:1#等级不能为空|等级需大于0" dc:"VIP等级"`
	LevelName            string  `json:"levelName" v:"required|length:1,64#等级名称不能为空|等级名称长度需在1到64之间" dc:"等级名称"`
	Status               uint8   `json:"status" v:"in:0,1#状态无效" dc:"状态(0关闭,1开启)"`
	UpgradeRechargeLimit float64 `json:"upgradeRechargeLimit" dc:"升级充值上限(USD,保留4位小数)"`
	MinWithdrawAmount    float64 `json:"minWithdrawAmount" dc:"最低提现金额(USD,保留4位小数)"`
	MaxWithdrawAmount    float64 `json:"maxWithdrawAmount" dc:"最高提现金额(USD,保留4位小数)"`
	Fee                  float64 `json:"fee" dc:"手续费(保留4位小数)"`
	Animation            string  `json:"animation" dc:"进场特效动画资源文件名(mp4)"`
}

type UpdateVipCfgRes struct {
	Success bool `json:"success"`
}

// DeleteVipCfgReq 删除VIP配置
type DeleteVipCfgReq struct {
	g.Meta `path:"/deleteVipCfg" method:"post" summary:"删除VIP配置" tags:"VIP配置"`
	ID     uint64 `json:"id" v:"required#ID不能为空" dc:"配置ID"`
}

type DeleteVipCfgRes struct {
	Success bool `json:"success"`
}

// AppVipCfgItem App端VIP配置项
type AppVipCfgItem struct {
	ID                   string  `json:"id"`
	Level                uint32  `json:"level"`
	LevelName            string  `json:"levelName"`
	Status               uint8   `json:"status"`
	UpgradeRechargeLimit float64 `json:"upgradeRechargeLimit"`
	MinWithdrawAmount    float64 `json:"minWithdrawAmount"`
	MaxWithdrawAmount    float64 `json:"maxWithdrawAmount"`
	Fee                  float64 `json:"fee"`
	Animation            string  `json:"animation" dc:"进场特效动画完整URL"`
}

// AppVipCfgByLevelReq App端按等级查询VIP配置
type AppVipCfgByLevelReq struct {
	g.Meta `path:"/getVipCfgByLevel" method:"post" summary:"按等级查询VIP配置" tags:"VIP配置"`
	Level  uint32 `json:"level" v:"required|min:1#等级不能为空|等级需大于0" dc:"VIP等级"`
}

type AppVipCfgByLevelRes struct {
	Item *AppVipCfgItem `json:"item"`
}

// AppVipCfgListReq App端查询全部VIP配置
type AppVipCfgListReq struct {
	g.Meta `path:"/vipCfgListForApp" method:"post" summary:"查询全部VIP配置" tags:"VIP配置"`
}

type AppVipCfgListRes struct {
	List []*AppVipCfgItem `json:"list"`
}
